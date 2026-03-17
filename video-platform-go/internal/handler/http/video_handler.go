package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/video-platform/go/internal/domain/video"
	videoService "github.com/video-platform/go/internal/service/video"
	"github.com/video-platform/go/pkg/jwt"
	"github.com/video-platform/go/pkg/response"
)

type VideoHandler struct {
	service  videoService.Service
	validate *validator.Validate
}

func NewVideoHandler(service videoService.Service) *VideoHandler {
	return &VideoHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *VideoHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/videos", func(r chi.Router) {
		r.Get("/", h.ListVideos)
		r.Get("/hot", h.GetHotVideos)
		r.Get("/search", h.SearchVideos)
		r.With(authMiddleware).Post("/", h.UploadVideo)
		
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetVideo)
			r.With(authMiddleware).Put("/", h.UpdateVideo)
			r.With(authMiddleware).Delete("/", h.DeleteVideo)
			r.With(authMiddleware).Post("/publish", h.PublishVideo)
			r.Post("/view", h.RecordView)
		})
		
		r.With(authMiddleware).Get("/my", h.GetMyVideos)
		r.Get("/user/{userId}", h.GetUserVideos)
	})
}

func (h *VideoHandler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	var req video.VideoUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	
	if err := h.validate.Struct(&req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	
	uploadResp, err := h.service.UploadVideo(r.Context(), claims.UserID, &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Created(w, uploadResp)
}

func (h *VideoHandler) GetVideo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid video ID")
		return
	}
	
	var viewerID int64
	if claims, ok := r.Context().Value("userClaims").(*jwt.Claims); ok {
		viewerID = claims.UserID
	}
	
	v, err := h.service.GetVideo(r.Context(), id, viewerID)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, v)
}

func (h *VideoHandler) ListVideos(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	filter := &video.VideoFilter{
		Status:     video.StatusPublished,
		Visibility: video.VisibilityPublic,
	}
	
	if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
		if cid, err := strconv.Atoi(categoryID); err == nil {
			filter.CategoryID = cid
		}
	}
	
	if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
		filter.SortBy = sortBy
	}
	filter.SortDesc = r.URL.Query().Get("sort_order") == "desc"
	
	resp, err := h.service.ListVideos(r.Context(), filter, page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, resp)
}

func (h *VideoHandler) UpdateVideo(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid video ID")
		return
	}
	
	var req video.VideoUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	
	v, err := h.service.UpdateVideo(r.Context(), id, claims.UserID, &req, false)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, v)
}

func (h *VideoHandler) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid video ID")
		return
	}
	
	if err := h.service.DeleteVideo(r.Context(), id, claims.UserID, false); err != nil {
		response.Error(w, err)
		return
	}
	
	response.NoContent(w)
}

func (h *VideoHandler) PublishVideo(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid video ID")
		return
	}
	
	if err := h.service.PublishVideo(r.Context(), id, claims.UserID, false); err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]string{"message": "Video published successfully"})
}

func (h *VideoHandler) GetHotVideos(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	videos, err := h.service.GetHotVideos(r.Context(), limit)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, videos)
}

func (h *VideoHandler) SearchVideos(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	if keyword == "" {
		response.BadRequest(w, "Search keyword is required")
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
	
	resp, err := h.service.SearchVideos(r.Context(), keyword, page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, resp)
}

func (h *VideoHandler) GetUserVideos(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
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
	
	resp, err := h.service.GetUserVideos(r.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, resp)
}

func (h *VideoHandler) GetMyVideos(w http.ResponseWriter, r *http.Request) {
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
	
	resp, err := h.service.GetUserVideos(r.Context(), claims.UserID, page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, resp)
}

func (h *VideoHandler) RecordView(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid video ID")
		return
	}
	
	if err := h.service.RecordView(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]string{"message": "View recorded"})
}
