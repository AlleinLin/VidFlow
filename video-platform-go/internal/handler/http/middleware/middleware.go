package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/video-platform/go/internal/infrastructure/metrics"
	"github.com/video-platform/go/pkg/jwt"
	"github.com/video-platform/go/pkg/logger"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		
		defer func() {
			logger.LogRequest(r.Context(), r.Method, r.URL.Path, ww.Status(), time.Since(start), nil)
		}()
		
		next.ServeHTTP(ww, r)
	})
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		
		defer func() {
			duration := time.Since(start).Seconds()
			metrics.RecordHTTPRequest(r.Method, r.URL.Path, string(rune(ww.Status())), duration)
		}()
		
		next.ServeHTTP(ww, r)
	})
}

func Recovery(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}

func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			
			allowed := false
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}
			
			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			}
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

func RateLimit(requestsPerMinute int) func(http.Handler) http.Handler {
	type client struct {
		count     int
		expiresAt time.Time
	}
	
	var mu sync.Mutex
	clients := make(map[string]*client)
	
	go func() {
		for range time.Tick(time.Minute) {
			mu.Lock()
			for ip, c := range clients {
				if time.Now().After(c.expiresAt) {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getRealIP(r)
			
			mu.Lock()
			if _, exists := clients[ip]; !exists {
				clients[ip] = &client{
					count:     0,
					expiresAt: time.Now().Add(time.Minute),
				}
			}
			
			if clients[ip].count >= requestsPerMinute {
				mu.Unlock()
				writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}
			
			clients[ip].count++
			mu.Unlock()
			
			next.ServeHTTP(w, r)
		})
	}
}

func RateLimitWithKey(requestsPerMinute int, keyFunc func(r *http.Request) string) func(http.Handler) http.Handler {
	type client struct {
		count     int
		expiresAt time.Time
	}
	
	var mu sync.Mutex
	clients := make(map[string]*client)
	
	go func() {
		for range time.Tick(time.Minute) {
			mu.Lock()
			for key, c := range clients {
				if time.Now().After(c.expiresAt) {
					delete(clients, key)
				}
			}
			mu.Unlock()
		}
	}()
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)
			if key == "" {
				key = getRealIP(r)
			}
			
			mu.Lock()
			if _, exists := clients[key]; !exists {
				clients[key] = &client{
					count:     0,
					expiresAt: time.Now().Add(time.Minute),
				}
			}
			
			if clients[key].count >= requestsPerMinute {
				mu.Unlock()
				writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}
			
			clients[key].count++
			mu.Unlock()
			
			next.ServeHTTP(w, r)
		})
	}
}

func getRealIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}

func RequestID(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}

func RealIP(next http.Handler) http.Handler {
	return middleware.RealIP(next)
}

func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, timeout, "Request timeout")
	}
}

func Compress(level int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware.Compress(level)(next)
	}
}

func StripSlashes(next http.Handler) http.Handler {
	return middleware.StripSlashes(next)
}

func Heartbeat(endpoint string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware.Heartbeat(endpoint)(next)
	}
}

func AuthMiddleware(jwtManager *jwt.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeJSONError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}
			
			tokenString, err := jwt.ExtractTokenFromHeader(authHeader)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}
			
			claims, err := jwtManager.ValidateToken(r.Context(), tokenString)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}
			
			ctx := context.WithValue(r.Context(), "userClaims", claims)
			ctx = context.WithValue(ctx, "userID", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func OptionalAuthMiddleware(jwtManager *jwt.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}
			
			tokenString, err := jwt.ExtractTokenFromHeader(authHeader)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			
			claims, err := jwtManager.ValidateToken(r.Context(), tokenString)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			
			ctx := context.WithValue(r.Context(), "userClaims", claims)
			ctx = context.WithValue(ctx, "userID", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		
		if claims.Role != "admin" {
			writeJSONError(w, http.StatusForbidden, "admin access required")
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
