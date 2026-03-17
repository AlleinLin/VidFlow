package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/video-platform/go/internal/domain/interaction"
	interactionService "github.com/video-platform/go/internal/service/interaction"
	"github.com/video-platform/go/pkg/jwt"
	"github.com/video-platform/go/pkg/response"
)

type InteractionHandler struct {
	service  interactionService.InteractionService
	validate *validator.Validate
}

func NewInteractionHandler(service interactionService.InteractionService) *InteractionHandler {
	return &InteractionHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *InteractionHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/interactions", func(r chi.Router) {
		r.Route("/comments", func(r chi.Router) {
			r.Get("/video/{videoId}", h.GetComments)
			r.With(authMiddleware).Post("/", h.CreateComment)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/replies", h.GetReplies)
				r.With(authMiddleware).Put("/", h.UpdateComment)
				r.With(authMiddleware).Delete("/", h.DeleteComment)
			})
		})
		
		r.Route("/likes", func(r chi.Router) {
			r.With(authMiddleware).Post("/video/{videoId}", h.LikeVideo)
			r.With(authMiddleware).Delete("/video/{videoId}", h.UnlikeVideo)
			r.Get("/video/{videoId}/status", h.GetLikeStatus)
		})
		
		r.Route("/favorites", func(r chi.Router) {
			r.With(authMiddleware).Post("/video/{videoId}", h.FavoriteVideo)
			r.With(authMiddleware).Delete("/video/{videoId}", h.UnfavoriteVideo)
		})
		
		r.Route("/danmakus", func(r chi.Router) {
			r.Get("/video/{videoId}", h.GetDanmakus)
			r.With(authMiddleware).Post("/", h.CreateDanmaku)
		})
	})
}

func (h *InteractionHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	var req interaction.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	
	if err := h.validate.Struct(&req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	
	comment, err := h.service.CreateComment(r.Context(), claims.UserID, &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Created(w, comment)
}

func (h *InteractionHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	videoIDStr := chi.URLParam(r, "videoId")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid video ID")
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
	
	resp, err := h.service.GetComments(r.Context(), videoID, page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, resp)
}

func (h *InteractionHandler) GetReplies(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	rootID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid comment ID")
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
	
	resp, err := h.service.GetReplies(r.Context(), rootID, page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, resp)
}

func (h *InteractionHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid comment ID")
		return
	}
	
	var req interaction.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	
	if err := h.validate.Struct(&req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	
	if err := h.service.UpdateComment(r.Context(), id, claims.UserID, req.Content, false); err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]string{"message": "Comment updated successfully"})
}

func (h *InteractionHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid comment ID")
		return
	}
	
	if err := h.service.DeleteComment(r.Context(), id, claims.UserID, false); err != nil {
		response.Error(w, err)
		return
	}
	
	response.NoContent(w)
}

func (h *InteractionHandler) LikeVideo(w http.ResponseWriter, r *http.Request) {
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
	
	if err := h.service.LikeVideo(r.Context(), claims.UserID, videoID); err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]string{"message": "Video liked successfully"})
}

func (h *InteractionHandler) UnlikeVideo(w http.ResponseWriter, r *http.Request) {
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
	
	if err := h.service.UnlikeVideo(r.Context(), claims.UserID, videoID); err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]string{"message": "Video unliked successfully"})
}

func (h *InteractionHandler) GetLikeStatus(w http.ResponseWriter, r *http.Request) {
	videoIDStr := chi.URLParam(r, "videoId")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid video ID")
		return
	}
	
	var userID int64
	if claims, ok := r.Context().Value("userClaims").(*jwt.Claims); ok {
		userID = claims.UserID
	}
	
	status, err := h.service.GetLikeStatus(r.Context(), userID, videoID)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, status)
}

func (h *InteractionHandler) FavoriteVideo(w http.ResponseWriter, r *http.Request) {
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
	
	if err := h.service.FavoriteVideo(r.Context(), claims.UserID, videoID); err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]string{"message": "Video favorited successfully"})
}

func (h *InteractionHandler) UnfavoriteVideo(w http.ResponseWriter, r *http.Request) {
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
	
	if err := h.service.UnfavoriteVideo(r.Context(), claims.UserID, videoID); err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]string{"message": "Video unfavorited successfully"})
}

func (h *InteractionHandler) CreateDanmaku(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	var req interaction.CreateDanmakuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	
	if err := h.validate.Struct(&req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	
	danmaku, err := h.service.CreateDanmaku(r.Context(), claims.UserID, &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Created(w, danmaku)
}

func (h *InteractionHandler) GetDanmakus(w http.ResponseWriter, r *http.Request) {
	videoIDStr := chi.URLParam(r, "videoId")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid video ID")
		return
	}
	
	startTime, _ := strconv.ParseFloat(r.URL.Query().Get("start"), 64)
	endTime, _ := strconv.ParseFloat(r.URL.Query().Get("end"), 64)
	
	if endTime == 0 {
		endTime = 999999
	}
	
	resp, err := h.service.GetDanmakus(r.Context(), videoID, startTime, endTime)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, resp)
}
