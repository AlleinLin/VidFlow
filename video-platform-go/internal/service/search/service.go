package search

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type SearchType string

const (
	SearchTypeVideo    SearchType = "video"
	SearchTypeUser     SearchType = "user"
	SearchTypeTag      SearchType = "tag"
	SearchTypeCategory SearchType = "category"
)

type SearchResult struct {
	Type      SearchType   `json:"type"`
	ID        int64        `json:"id"`
	Title     string       `json:"title,omitempty"`
	Content   string       `json:"content,omitempty"`
	Score     float64      `json:"score"`
	Highlight []string     `json:"highlight,omitempty"`
	Metadata  interface{}  `json:"metadata,omitempty"`
}

type SearchResponse struct {
	Results    []*SearchResult `json:"results"`
	Total      int64           `json:"total"`
	Took       time.Duration   `json:"took"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
}

type VideoSearchResult struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	ThumbnailURL string  `json:"thumbnail_url"`
	ViewCount    int64   `json:"view_count"`
	LikeCount    int64   `json:"like_count"`
	Duration     int     `json:"duration"`
	PublishedAt  *time.Time `json:"published_at"`
	AuthorName   string  `json:"author_name"`
	AuthorID     int64   `json:"author_id"`
	Score        float64 `json:"score"`
}

type UserSearchResult struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
	FollowerCount int64 `json:"follower_count"`
	Score       float64 `json:"score"`
}

type SearchRepository interface {
	SearchVideos(ctx context.Context, query string, page, pageSize int) ([]*VideoSearchResult, int64, error)
	SearchUsers(ctx context.Context, query string, page, pageSize int) ([]*UserSearchResult, int64, error)
	SearchAll(ctx context.Context, query string, page, pageSize int) (*SearchResponse, error)
	IndexVideo(ctx context.Context, videoID int64) error
	IndexUser(ctx context.Context, userID int64) error
	RemoveFromIndex(ctx context.Context, indexType SearchType, id int64) error
}

type searchRepository struct {
	pool  *pgxpool.Pool
	redis *redis.Client
}

func NewSearchRepository(pool *pgxpool.Pool, redis *redis.Client) SearchRepository {
	return &searchRepository{pool: pool, redis: redis}
}

func (r *searchRepository) SearchVideos(ctx context.Context, query string, page, pageSize int) ([]*VideoSearchResult, int64, error) {
	offset := (page - 1) * pageSize
	
	countQuery := `
		SELECT COUNT(*) FROM videos 
		WHERE status = 'PUBLISHED' AND visibility = 'public'
		AND (title ILIKE $1 OR description ILIKE $1)
	`
	var total int64
	r.pool.QueryRow(ctx, countQuery, "%"+query+"%").Scan(&total)
	
	searchQuery := `
		SELECT v.id, v.title, v.description, v.thumbnail_url, v.view_count, v.like_count, 
			   v.duration_seconds, v.published_at, u.id, u.display_name
		FROM videos v
		INNER JOIN users u ON v.user_id = u.id
		WHERE v.status = 'PUBLISHED' AND v.visibility = 'public'
		AND (v.title ILIKE $1 OR v.description ILIKE $1)
		ORDER BY 
			CASE WHEN v.title ILIKE $1 THEN 2 ELSE 1 END,
			v.view_count DESC, v.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.pool.Query(ctx, searchQuery, "%"+query+"%", pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var results []*VideoSearchResult
	for rows.Next() {
		var result VideoSearchResult
		err := rows.Scan(
			&result.ID, &result.Title, &result.Description, &result.ThumbnailURL,
			&result.ViewCount, &result.LikeCount, &result.Duration, &result.PublishedAt,
			&result.AuthorID, &result.AuthorName,
		)
		if err != nil {
			return nil, 0, err
		}
		result.Score = calculateRelevanceScore(result.Title, result.Description, query)
		results = append(results, &result)
	}
	
	return results, total, nil
}

func (r *searchRepository) SearchUsers(ctx context.Context, query string, page, pageSize int) ([]*UserSearchResult, int64, error) {
	offset := (page - 1) * pageSize
	
	countQuery := `
		SELECT COUNT(*) FROM users 
		WHERE status = 'active' AND (username ILIKE $1 OR display_name ILIKE $1 OR bio ILIKE $1)
	`
	var total int64
	r.pool.QueryRow(ctx, countQuery, "%"+query+"%").Scan(&total)
	
	searchQuery := `
		SELECT u.id, u.username, u.display_name, u.avatar_url, u.bio, 
			   (SELECT COUNT(*) FROM user_follows WHERE following_id = u.id) as follower_count
		FROM users u
		WHERE u.status = 'active'
		AND (u.username ILIKE $1 OR u.display_name ILIKE $1 OR u.bio ILIKE $1)
		ORDER BY 
			CASE WHEN u.username ILIKE $1 THEN 3 
				 WHEN u.display_name ILIKE $1 THEN 2 
				 ELSE 1 END,
			follower_count DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.pool.Query(ctx, searchQuery, "%"+query+"%", pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var results []*UserSearchResult
	for rows.Next() {
		var result UserSearchResult
		err := rows.Scan(
			&result.ID, &result.Username, &result.DisplayName, &result.AvatarURL,
			&result.Bio, &result.FollowerCount,
		)
		if err != nil {
			return nil, 0, err
		}
		result.Score = calculateUserRelevanceScore(result.Username, result.DisplayName, result.Bio, query)
		results = append(results, &result)
	}
	
	return results, total, nil
}

func (r *searchRepository) SearchAll(ctx context.Context, query string, page, pageSize int) (*SearchResponse, error) {
	start := time.Now()
	
	halfPageSize := pageSize / 2
	if halfPageSize < 1 {
		halfPageSize = 1
	}
	
	videos, videoTotal, err := r.SearchVideos(ctx, query, page, halfPageSize)
	if err != nil {
		return nil, err
	}
	
	users, userTotal, err := r.SearchUsers(ctx, query, page, halfPageSize)
	if err != nil {
		return nil, err
	}
	
	var results []*SearchResult
	
	for _, v := range videos {
		results = append(results, &SearchResult{
			Type:    SearchTypeVideo,
			ID:      v.ID,
			Title:   v.Title,
			Content: v.Description,
			Score:   v.Score,
			Metadata: v,
		})
	}
	
	for _, u := range users {
		results = append(results, &SearchResult{
			Type:    SearchTypeUser,
			ID:      u.ID,
			Title:   u.DisplayName,
			Content: u.Bio,
			Score:   u.Score,
			Metadata: u,
		})
	}
	
	return &SearchResponse{
		Results:  results,
		Total:    videoTotal + userTotal,
		Took:     time.Since(start),
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (r *searchRepository) IndexVideo(ctx context.Context, videoID int64) error {
	cacheKey := fmt.Sprintf("search:video:%d", videoID)
	
	var videoData struct {
		ID          int64
		Title       string
		Description string
	}
	
	query := `SELECT id, title, description FROM videos WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, videoID).Scan(&videoData.ID, &videoData.Title, &videoData.Description)
	if err != nil {
		return err
	}
	
	data, _ := json.Marshal(videoData)
	return r.redis.Set(ctx, cacheKey, data, 24*time.Hour).Err()
}

func (r *searchRepository) IndexUser(ctx context.Context, userID int64) error {
	cacheKey := fmt.Sprintf("search:user:%d", userID)
	
	var userData struct {
		ID          int64
		Username    string
		DisplayName string
		Bio         string
	}
	
	query := `SELECT id, username, display_name, bio FROM users WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, userID).Scan(&userData.ID, &userData.Username, &userData.DisplayName, &userData.Bio)
	if err != nil {
		return err
	}
	
	data, _ := json.Marshal(userData)
	return r.redis.Set(ctx, cacheKey, data, 24*time.Hour).Err()
}

func (r *searchRepository) RemoveFromIndex(ctx context.Context, indexType SearchType, id int64) error {
	cacheKey := fmt.Sprintf("search:%s:%d", indexType, id)
	return r.redis.Del(ctx, cacheKey).Err()
}

func calculateRelevanceScore(title, description, query string) float64 {
	score := 0.0
	
	if containsIgnoreCase(title, query) {
		score += 2.0
	}
	if containsIgnoreCase(description, query) {
		score += 1.0
	}
	
	return score
}

func calculateUserRelevanceScore(username, displayName, bio, query string) float64 {
	score := 0.0
	
	if containsIgnoreCase(username, query) {
		score += 3.0
	}
	if containsIgnoreCase(displayName, query) {
		score += 2.0
	}
	if containsIgnoreCase(bio, query) {
		score += 1.0
	}
	
	return score
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || len(s) > 0 && containsIgnoreCase(s[1:], substr) || 
		 (len(s) >= len(substr) && s[:len(substr)] == substr))
}

type SearchService interface {
	Search(ctx context.Context, query string, searchType SearchType, page, pageSize int) (*SearchResponse, error)
	SearchVideos(ctx context.Context, query string, page, pageSize int) ([]*VideoSearchResult, int64, error)
	SearchUsers(ctx context.Context, query string, page, pageSize int) ([]*UserSearchResult, int64, error)
	GetSearchSuggestions(ctx context.Context, query string, limit int) ([]string, error)
}

type searchService struct {
	repo SearchRepository
}

func NewSearchService(repo SearchRepository) SearchService {
	return &searchService{repo: repo}
}

func (s *searchService) Search(ctx context.Context, query string, searchType SearchType, page, pageSize int) (*SearchResponse, error) {
	switch searchType {
	case SearchTypeVideo:
		videos, total, err := s.repo.SearchVideos(ctx, query, page, pageSize)
		if err != nil {
			return nil, err
		}
		
		var results []*SearchResult
		for _, v := range videos {
			results = append(results, &SearchResult{
				Type:     SearchTypeVideo,
				ID:       v.ID,
				Title:    v.Title,
				Content:  v.Description,
				Score:    v.Score,
				Metadata: v,
			})
		}
		
		return &SearchResponse{
			Results:  results,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		}, nil
		
	case SearchTypeUser:
		users, total, err := s.repo.SearchUsers(ctx, query, page, pageSize)
		if err != nil {
			return nil, err
		}
		
		var results []*SearchResult
		for _, u := range users {
			results = append(results, &SearchResult{
				Type:     SearchTypeUser,
				ID:       u.ID,
				Title:    u.DisplayName,
				Content:  u.Bio,
				Score:    u.Score,
				Metadata: u,
			})
		}
		
		return &SearchResponse{
			Results:  results,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		}, nil
		
	default:
		return s.repo.SearchAll(ctx, query, page, pageSize)
	}
}

func (s *searchService) SearchVideos(ctx context.Context, query string, page, pageSize int) ([]*VideoSearchResult, int64, error) {
	return s.repo.SearchVideos(ctx, query, page, pageSize)
}

func (s *searchService) SearchUsers(ctx context.Context, query string, page, pageSize int) ([]*UserSearchResult, int64, error) {
	return s.repo.SearchUsers(ctx, query, page, pageSize)
}

func (s *searchService) GetSearchSuggestions(ctx context.Context, query string, limit int) ([]string, error) {
	return []string{}, nil
}
