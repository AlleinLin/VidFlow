package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/video-platform/go/internal/service/search"
	"github.com/video-platform/go/pkg/response"
)

type SearchHandler struct {
	service search.SearchService
}

func NewSearchHandler(service search.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

func (h *SearchHandler) RegisterRoutes(r chi.Router) {
	r.Route("/search", func(r chi.Router) {
		r.Get("/", h.Search)
		r.Get("/videos", h.SearchVideos)
		r.Get("/users", h.SearchUsers)
		r.Get("/suggestions", h.GetSuggestions)
	})
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		response.BadRequest(w, "Search query is required")
		return
	}
	
	searchType := search.SearchType(r.URL.Query().Get("type"))
	if searchType == "" {
		searchType = search.SearchTypeVideo
	}
	
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	result, err := h.service.Search(r.Context(), query, searchType, page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, result)
}

func (h *SearchHandler) SearchVideos(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		response.BadRequest(w, "Search query is required")
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
	
	results, total, err := h.service.SearchVideos(r.Context(), query, page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]interface{}{
		"results":   results,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *SearchHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		response.BadRequest(w, "Search query is required")
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
	
	results, total, err := h.service.SearchUsers(r.Context(), query, page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]interface{}{
		"results":   results,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *SearchHandler) GetSuggestions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		response.BadRequest(w, "Search query is required")
		return
	}
	
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 20 {
		limit = 10
	}
	
	suggestions, err := h.service.GetSearchSuggestions(r.Context(), query, limit)
	if err != nil {
		response.Error(w, err)
		return
	}
	
	response.Success(w, map[string]interface{}{
		"suggestions": suggestions,
	})
}
