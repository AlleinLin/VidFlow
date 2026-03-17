package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/video-platform/go/internal/service/recommendation"
	"github.com/video-platform/go/pkg/jwt"
	"github.com/video-platform/go/pkg/response"
)

type RecommendationHandler struct {
	service recommendation.RecommendationService
}

func NewRecommendationHandler(service recommendation.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{service: service}
}

func (h *RecommendationHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/recommendations", func(r chi.Router) {
		r.Get("/hot", h.GetHotRecommendations)
		r.With(authMiddleware).Get("/personalized", h.GetPersonalizedRecommendations)
		r.Get("/similar/{videoId}", h.GetSimilarVideos)
		r.With(authMiddleware).Get("/following", h.GetFollowingFeed)
	})
}

func (h *RecommendationHandler) GetHotRecommendations(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	videoIDs, err := h.service.GetHotRecommendations(r.Context(), limit)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]interface{}{
		"video_ids": videoIDs,
		"type":      "hot",
	})
}

func (h *RecommendationHandler) GetPersonalizedRecommendations(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	videoIDs, err := h.service.GetPersonalizedRecommendations(r.Context(), claims.UserID, limit)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]interface{}{
		"video_ids": videoIDs,
		"type":      "personalized",
	})
}

func (h *RecommendationHandler) GetSimilarVideos(w http.ResponseWriter, r *http.Request) {
	videoIDStr := chi.URLParam(r, "videoId")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid video ID")
		return
	}
	
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}
	
	videoIDs, err := h.service.GetSimilarVideos(r.Context(), videoID, limit)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]interface{}{
		"video_ids": videoIDs,
		"type":      "similar",
	})
}

func (h *RecommendationHandler) GetFollowingFeed(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	videoIDs, err := h.service.GetFollowingFeed(r.Context(), claims.UserID, limit)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]interface{}{
		"video_ids": videoIDs,
		"type":      "following",
	})
}
