package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/video-platform/go/internal/service/playback"
	"github.com/video-platform/go/pkg/jwt"
	"github.com/video-platform/go/pkg/response"
)

type PlaybackHandler struct {
	service  playback.PlaybackService
	validate *validator.Validate
}

func NewPlaybackHandler(service playback.PlaybackService) *PlaybackHandler {
	return &PlaybackHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *PlaybackHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/playback", func(r chi.Router) {
		r.With(authMiddleware).Post("/progress", h.UpdateProgress)
		r.With(authMiddleware).Get("/progress/{videoId}", h.GetProgress)
		r.With(authMiddleware).Get("/history", h.GetWatchHistory)
		r.With(authMiddleware).Get("/continue-watching", h.GetContinueWatching)
		r.With(authMiddleware).Delete("/history/{videoId}", h.DeleteWatchHistory)
		r.With(authMiddleware).Delete("/history", h.ClearWatchHistory)
	})
}

type UpdateProgressRequest struct {
	VideoID       int64   `json:"video_id" validate:"required"`
	Position      float64 `json:"position" validate:"required,min=0"`
	Duration      float64 `json:"duration" validate:"required,min=0"`
	WatchDuration int64   `json:"watch_duration" validate:"required,min=0"`
}

func (h *PlaybackHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	var req UpdateProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	
	if err := h.validate.Struct(&req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	
	if err := h.service.UpdateProgress(r.Context(), claims.UserID, req.VideoID, req.Position, req.Duration, req.WatchDuration); err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]string{"message": "Progress updated successfully"})
}

func (h *PlaybackHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	videoIDStr := chi.URLParam(r, "videoId")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid video ID")
		return
	}
	
	history, err := h.service.GetProgress(r.Context(), claims.UserID, videoID)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	if history == nil {
		response.Success(w, nil)
		return
	}
	
	response.Success(w, history)
}

func (h *PlaybackHandler) GetWatchHistory(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	histories, total, err := h.service.GetWatchHistory(r.Context(), claims.UserID, page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]interface{}{
		"histories": histories,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *PlaybackHandler) GetContinueWatching(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 10
	}
	
	histories, err := h.service.GetContinueWatching(r.Context(), claims.UserID, limit)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, histories)
}

func (h *PlaybackHandler) DeleteWatchHistory(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	videoIDStr := chi.URLParam(r, "videoId")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid video ID")
		return
	}
	
	if err := h.service.DeleteWatchHistory(r.Context(), claims.UserID, videoID); err != nil {
		response.Error(w, err)
		return
	}
	
	response.NoContent(w)
}

func (h *PlaybackHandler) ClearWatchHistory(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	if err := h.service.ClearWatchHistory(r.Context(), claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	
	response.NoContent(w)
}
