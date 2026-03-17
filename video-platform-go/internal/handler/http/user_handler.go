package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/video-platform/go/internal/domain/user"
	userService "github.com/video-platform/go/internal/service/user"
	"github.com/video-platform/go/pkg/jwt"
	"github.com/video-platform/go/pkg/response"
)

type UserHandler struct {
	service    userService.Service
	jwtManager *jwt.JWTManager
	validate   *validator.Validate
}

func NewUserHandler(service userService.Service, jwtManager *jwt.JWTManager) *UserHandler {
	return &UserHandler{
		service:    service,
		jwtManager: jwtManager,
		validate:   validator.New(),
	}
}

func (h *UserHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/refresh", h.RefreshToken)
		r.With(authMiddleware).Post("/logout", h.Logout)
	})
	
	r.Route("/users", func(r chi.Router) {
		r.With(authMiddleware).Get("/me", h.GetCurrentUser)
		r.With(authMiddleware).Put("/me", h.UpdateProfile)
		r.With(authMiddleware).Delete("/me", h.DeleteUser)
		r.Get("/{id}", h.GetUserByID)
		r.With(authMiddleware).Post("/{id}/follow", h.FollowUser)
		r.With(authMiddleware).Delete("/{id}/follow", h.UnfollowUser)
		r.Get("/{id}/followers", h.GetFollowers)
		r.Get("/{id}/following", h.GetFollowing)
	})
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req user.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	
	if err := h.validate.Struct(&req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	
	u, err := h.service.Register(r.Context(), &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Created(w, u)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req user.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	
	if err := h.validate.Struct(&req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	
	loginResp, err := h.service.Login(r.Context(), &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, loginResp)
}

func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	
	if req.RefreshToken == "" {
		response.BadRequest(w, "Refresh token is required")
		return
	}
	
	loginResp, err := h.service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, loginResp)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "Logged out successfully"})
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	u, err := h.service.GetProfile(r.Context(), claims.UserID)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, u)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	var req user.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	
	u, err := h.service.UpdateProfile(r.Context(), claims.UserID, &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, u)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	if err := h.service.DeleteUser(r.Context(), claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	
	response.NoContent(w)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}
	
	u, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	publicProfile := map[string]interface{}{
		"id":              u.ID,
		"username":        u.Username,
		"display_name":    u.DisplayName,
		"avatar_url":      u.AvatarURL,
		"bio":             u.Bio,
		"follower_count":  u.FollowerCount,
		"following_count": u.FollowingCount,
		"created_at":      u.CreatedAt,
	}
	
	response.Success(w, publicProfile)
}

func (h *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	idStr := chi.URLParam(r, "id")
	followingID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}
	
	if err := h.service.Follow(r.Context(), claims.UserID, followingID); err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]string{"message": "Followed successfully"})
}

func (h *UserHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("userClaims").(*jwt.Claims)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}
	
	idStr := chi.URLParam(r, "id")
	followingID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}
	
	if err := h.service.Unfollow(r.Context(), claims.UserID, followingID); err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]string{"message": "Unfollowed successfully"})
}

func (h *UserHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
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
	
	users, total, err := h.service.GetFollowers(r.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	publicUsers := make([]map[string]interface{}, len(users))
	for i, u := range users {
		publicUsers[i] = map[string]interface{}{
			"id":              u.ID,
			"username":        u.Username,
			"display_name":    u.DisplayName,
			"avatar_url":      u.AvatarURL,
			"follower_count":  u.FollowerCount,
			"following_count": u.FollowingCount,
		}
	}
	
	response.Success(w, map[string]interface{}{
		"users":     publicUsers,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *UserHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
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
	
	users, total, err := h.service.GetFollowing(r.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	publicUsers := make([]map[string]interface{}, len(users))
	for i, u := range users {
		publicUsers[i] = map[string]interface{}{
			"id":              u.ID,
			"username":        u.Username,
			"display_name":    u.DisplayName,
			"avatar_url":      u.AvatarURL,
			"follower_count":  u.FollowerCount,
			"following_count": u.FollowingCount,
		}
	}
	
	response.Success(w, map[string]interface{}{
		"users":     publicUsers,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
