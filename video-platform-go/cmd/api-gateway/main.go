package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/video-platform/go/internal/config"
	httpHandler "github.com/video-platform/go/internal/handler/http"
	"github.com/video-platform/go/internal/infrastructure/cache"
	"github.com/video-platform/go/internal/infrastructure/messaging"
	"github.com/video-platform/go/internal/infrastructure/metrics"
	"github.com/video-platform/go/internal/infrastructure/tracing"
	"github.com/video-platform/go/internal/repository/postgres"
	auditService "github.com/video-platform/go/internal/service/audit"
	cdnService "github.com/video-platform/go/internal/service/cdn"
	interactionService "github.com/video-platform/go/internal/service/interaction"
	notificationService "github.com/video-platform/go/internal/service/notification"
	paymentService "github.com/video-platform/go/internal/service/payment"
	playbackService "github.com/video-platform/go/internal/service/playback"
	recommendationService "github.com/video-platform/go/internal/service/recommendation"
	searchService "github.com/video-platform/go/internal/service/search"
	subscriptionService "github.com/video-platform/go/internal/service/subscription"
	"github.com/video-platform/go/internal/service/transcode"
	"github.com/video-platform/go/internal/service/user"
	videoService "github.com/video-platform/go/internal/service/video"
	apperrors "github.com/video-platform/go/pkg/errors"
	"github.com/video-platform/go/pkg/jwt"
	"github.com/video-platform/go/pkg/logger"
	"github.com/video-platform/go/pkg/response"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logInstance, err := logger.Init(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	_ = logInstance

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	pool, err := postgres.NewPostgresPool(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to create postgres pool: %v", err)
	}
	defer pool.Close()

	var redisClient *redis.Client
	var redisCache *cache.RedisCache
	
	redisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
	
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		redisCache, _ = cache.NewRedisCache(&cfg.Redis)
	}
	_ = redisCache

	if cfg.Tracing.Enabled {
		tracerProvider, err := tracing.NewTracerProvider(ctx, "video-platform", "1.0.0", cfg.Tracing.Endpoint)
		if err != nil {
			log.Printf("Warning: Failed to initialize tracing: %v", err)
		} else {
			defer tracerProvider.Shutdown(ctx)
		}
	}

	jwtManager, err := jwt.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
		cfg.JWT.Issuer,
	)
	if err != nil {
		log.Fatalf("Failed to create JWT manager: %v", err)
	}

	userRepo := postgres.NewUserRepository(pool)
	userSvc := user.NewService(userRepo, jwtManager)

	videoRepo := postgres.NewVideoRepository(pool)
	videoSvc := videoService.NewService(videoRepo)

	interactionRepo := interactionService.NewInteractionRepository(pool)
	interactionSvc := interactionService.NewInteractionService(interactionRepo, videoRepo)

	playbackRepo := playbackService.NewPlaybackRepository(pool)
	playbackSvc := playbackService.NewPlaybackService(playbackRepo, videoRepo)

	recommendationRepo := recommendationService.NewRecommendationRepository(pool)
	recommendationSvc := recommendationService.NewRecommendationService(recommendationRepo)

	searchRepo := searchService.NewSearchRepository(pool, redisClient)
	searchSvc := searchService.NewSearchService(searchRepo)

	notificationRepo := notificationService.NewNotificationRepository(pool, redisClient)
	notificationSvc := notificationService.NewNotificationService(notificationRepo, redisClient)
	_ = notificationSvc

	paymentRepo := paymentService.NewPaymentRepository(pool)
	paymentSvc := paymentService.NewPaymentService(paymentRepo, nil)
	_ = paymentSvc

	subscriptionRepo := subscriptionService.NewSubscriptionRepository(pool)
	subscriptionSvc := subscriptionService.NewSubscriptionService(subscriptionRepo)
	_ = subscriptionSvc

	auditRepo := auditService.NewAuditRepository(pool)
	auditSvc := auditService.NewAuditService(auditRepo, nil)
	_ = auditSvc

	cdnRepo := cdnService.NewCDNRepository(pool)
	cdnConfig := &cdnService.CDNConfig{
		Provider:    cdnService.CDNProviderCustom,
		Domain:      "cdn.videoplatform.com",
		EnableHTTPS: true,
	}
	cdnSvc := cdnService.NewCDNService(cdnRepo, cdnConfig)

	transcodeSvc := transcode.NewService(transcode.TranscodeConfig{
		OutputDir:       "./transcoded",
		FFmpegPath:      "ffmpeg",
		FFprobePath:     "ffprobe",
		MaxWorkers:      4,
		SegmentDuration: 6,
	})
	_ = transcodeSvc

	var kafkaProducer *messaging.KafkaProducer
	if len(cfg.Kafka.Brokers) > 0 {
		kafkaProducer, err = messaging.NewKafkaProducer(&cfg.Kafka)
		if err != nil {
			log.Printf("Warning: Failed to create Kafka producer: %v", err)
		} else {
			defer kafkaProducer.Close()
		}
	}
	_ = kafkaProducer

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID", "X-Client-Version"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	authMiddleware := AuthMiddleware(jwtManager)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"status": "not ready", "error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	})

	if cfg.Metrics.Enabled {
		r.Handle(cfg.Metrics.Path, promhttp.Handler())
	}

	r.Route("/api/v1", func(r chi.Router) {
		userHandler := httpHandler.NewUserHandler(userSvc, jwtManager)
		userHandler.RegisterRoutes(r, authMiddleware)

		videoHandler := httpHandler.NewVideoHandler(videoSvc)
		videoHandler.RegisterRoutes(r, authMiddleware)

		interactionHandler := httpHandler.NewInteractionHandler(interactionSvc)
		interactionHandler.RegisterRoutes(r, authMiddleware)

		playbackHandler := httpHandler.NewPlaybackHandler(playbackSvc)
		playbackHandler.RegisterRoutes(r, authMiddleware)

		recommendationHandler := httpHandler.NewRecommendationHandler(recommendationSvc)
		recommendationHandler.RegisterRoutes(r, authMiddleware)

		searchHandler := httpHandler.NewSearchHandler(searchSvc)
		searchHandler.RegisterRoutes(r)

		r.Get("/notifications", func(w http.ResponseWriter, r *http.Request) {
			response.Success(w, map[string]string{"message": "Notification service ready"})
		})

		r.Get("/subscriptions", func(w http.ResponseWriter, r *http.Request) {
			response.Success(w, map[string]string{"message": "Subscription service ready"})
		})

		r.Get("/payments", func(w http.ResponseWriter, r *http.Request) {
			response.Success(w, map[string]string{"message": "Payment service ready"})
		})

		r.Get("/cdn/stats", func(w http.ResponseWriter, r *http.Request) {
			stats, err := cdnSvc.GetStatistics(ctx, time.Now().Add(-24*time.Hour), time.Now())
			if err != nil {
				response.Error(w, err)
				return
			}
			response.Success(w, stats)
		})
	})

	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.NotFound(w, "Route not found")
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		response.Error(w, apperrors.NewBadRequestError("Method not allowed", nil))
	})

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		logger.Info(ctx, "Starting API Gateway server", "addr", server.Addr)
		logger.Info(ctx, "Available endpoints:")
		logger.Info(ctx, "  - GET  /health - Health check")
		logger.Info(ctx, "  - GET  /ready  - Readiness check")
		logger.Info(ctx, "  - GET  /metrics - Prometheus metrics")
		logger.Info(ctx, "  - POST /api/v1/auth/register - User registration")
		logger.Info(ctx, "  - POST /api/v1/auth/login - User login")
		logger.Info(ctx, "  - GET  /api/v1/videos - List videos")
		logger.Info(ctx, "  - POST /api/v1/videos - Upload video (auth required)")
		logger.Info(ctx, "  - GET  /api/v1/search - Search videos and users")
		logger.Info(ctx, "  - GET  /api/v1/recommendations - Get recommendations")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-sigChan
	logger.Info(ctx, "Received shutdown signal, gracefully shutting down...")

	metrics.SetActiveUsers(0)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(shutdownCtx, "Server shutdown error", "error", err)
	}

	if redisClient != nil {
		redisClient.Close()
	}

	logger.Info(ctx, "API Gateway stopped successfully")
}

func AuthMiddleware(jwtManager *jwt.JWTManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "Missing authorization header")
				return
			}

			tokenString, err := jwt.ExtractTokenFromHeader(authHeader)
			if err != nil {
				response.Unauthorized(w, "Invalid authorization header format")
				return
			}

			claims, err := jwtManager.ValidateToken(r.Context(), tokenString)
			if err != nil {
				if apperrors.IsAppError(err) {
					response.Error(w, err)
				} else {
					response.Unauthorized(w, "Invalid token")
				}
				return
			}

			ctx := context.WithValue(r.Context(), "userClaims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
