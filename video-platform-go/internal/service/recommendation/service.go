package recommendation

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RecommendationRepository interface {
	GetUserWatchHistory(ctx context.Context, userID int64, limit int) ([]*WatchRecord, error)
	GetUserLikes(ctx context.Context, userID int64, limit int) ([]int64, error)
	GetUserFavorites(ctx context.Context, userID int64, limit int) ([]int64, error)
	GetVideoCategories(ctx context.Context, videoIDs []int64) (map[int64][]int, error)
	GetCategoryVideos(ctx context.Context, categoryIDs []int, limit int) ([]int64, error)
	GetHotVideos(ctx context.Context, limit int) ([]int64, error)
	GetSimilarUsers(ctx context.Context, userID int64, limit int) ([]int64, error)
	GetUserFollowings(ctx context.Context, userID int64, limit int) ([]int64, error)
	GetFollowingVideos(ctx context.Context, userIDs []int64, limit int) ([]int64, error)
	GetVideoTags(ctx context.Context, videoIDs []int64) (map[int64][]string, error)
	GetTagVideos(ctx context.Context, tags []string, limit int) ([]int64, error)
}

type WatchRecord struct {
	VideoID       int64
	WatchDuration int64
	WatchProgress float64
	WatchedAt     time.Time
}

type recommendationRepository struct {
	pool *pgxpool.Pool
}

func NewRecommendationRepository(pool *pgxpool.Pool) RecommendationRepository {
	return &recommendationRepository{pool: pool}
}

func (r *recommendationRepository) GetUserWatchHistory(ctx context.Context, userID int64, limit int) ([]*WatchRecord, error) {
	query := `
		SELECT video_id, watch_duration, watch_progress, watched_at
		FROM watch_history
		WHERE user_id = $1
		ORDER BY watched_at DESC
		LIMIT $2
	`
	
	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var records []*WatchRecord
	for rows.Next() {
		var record WatchRecord
		if err := rows.Scan(&record.VideoID, &record.WatchDuration, &record.WatchProgress, &record.WatchedAt); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	
	return records, nil
}

func (r *recommendationRepository) GetUserLikes(ctx context.Context, userID int64, limit int) ([]int64, error) {
	query := `SELECT video_id FROM likes WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`
	
	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var videoIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		videoIDs = append(videoIDs, id)
	}
	
	return videoIDs, nil
}

func (r *recommendationRepository) GetUserFavorites(ctx context.Context, userID int64, limit int) ([]int64, error) {
	query := `SELECT video_id FROM favorites WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`
	
	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var videoIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		videoIDs = append(videoIDs, id)
	}
	
	return videoIDs, nil
}

func (r *recommendationRepository) GetVideoCategories(ctx context.Context, videoIDs []int64) (map[int64][]int, error) {
	if len(videoIDs) == 0 {
		return make(map[int64][]int), nil
	}
	
	query := `SELECT id, category_id FROM videos WHERE id = ANY($1)`
	
	rows, err := r.pool.Query(ctx, query, videoIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	result := make(map[int64][]int)
	for rows.Next() {
		var videoID int64
		var categoryID int
		if err := rows.Scan(&videoID, &categoryID); err != nil {
			return nil, err
		}
		result[videoID] = append(result[videoID], categoryID)
	}
	
	return result, nil
}

func (r *recommendationRepository) GetCategoryVideos(ctx context.Context, categoryIDs []int, limit int) ([]int64, error) {
	if len(categoryIDs) == 0 {
		return nil, nil
	}
	
	query := `
		SELECT id FROM videos 
		WHERE category_id = ANY($1) AND status = 'PUBLISHED'
		ORDER BY view_count DESC, created_at DESC
		LIMIT $2
	`
	
	rows, err := r.pool.Query(ctx, query, categoryIDs, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var videoIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		videoIDs = append(videoIDs, id)
	}
	
	return videoIDs, nil
}

func (r *recommendationRepository) GetHotVideos(ctx context.Context, limit int) ([]int64, error) {
	query := `
		SELECT id FROM videos 
		WHERE status = 'PUBLISHED' AND visibility = 'public'
		ORDER BY (view_count * 1 + like_count * 5 + comment_count * 10) DESC, created_at DESC
		LIMIT $1
	`
	
	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var videoIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		videoIDs = append(videoIDs, id)
	}
	
	return videoIDs, nil
}

func (r *recommendationRepository) GetSimilarUsers(ctx context.Context, userID int64, limit int) ([]int64, error) {
	query := `
		SELECT DISTINCT l2.user_id
		FROM likes l1
		INNER JOIN likes l2 ON l1.video_id = l2.video_id AND l1.user_id != l2.user_id
		WHERE l1.user_id = $1
		GROUP BY l2.user_id
		ORDER BY COUNT(*) DESC
		LIMIT $2
	`
	
	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var userIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}
	
	return userIDs, nil
}

func (r *recommendationRepository) GetUserFollowings(ctx context.Context, userID int64, limit int) ([]int64, error) {
	query := `SELECT following_id FROM user_follows WHERE follower_id = $1 ORDER BY created_at DESC LIMIT $2`
	
	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var userIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}
	
	return userIDs, nil
}

func (r *recommendationRepository) GetFollowingVideos(ctx context.Context, userIDs []int64, limit int) ([]int64, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	
	query := `
		SELECT id FROM videos 
		WHERE user_id = ANY($1) AND status = 'PUBLISHED' AND visibility = 'public'
		ORDER BY created_at DESC
		LIMIT $2
	`
	
	rows, err := r.pool.Query(ctx, query, userIDs, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var videoIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		videoIDs = append(videoIDs, id)
	}
	
	return videoIDs, nil
}

func (r *recommendationRepository) GetVideoTags(ctx context.Context, videoIDs []int64) (map[int64][]string, error) {
	return make(map[int64][]string), nil
}

func (r *recommendationRepository) GetTagVideos(ctx context.Context, tags []string, limit int) ([]int64, error) {
	return nil, nil
}

type RecommendationService interface {
	GetPersonalizedRecommendations(ctx context.Context, userID int64, limit int) ([]int64, error)
	GetHotRecommendations(ctx context.Context, limit int) ([]int64, error)
	GetSimilarVideos(ctx context.Context, videoID int64, limit int) ([]int64, error)
	GetFollowingFeed(ctx context.Context, userID int64, limit int) ([]int64, error)
	RefreshUserPreferences(ctx context.Context, userID int64) error
}

type recommendationService struct {
	repo RecommendationRepository
	cache map[int64]*UserPreferences
	mu    sync.RWMutex
}

type UserPreferences struct {
	UserID         int64
	CategoryScores map[int]float64
	TagScores      map[string]float64
	LastUpdated    time.Time
}

func NewRecommendationService(repo RecommendationRepository) RecommendationService {
	return &recommendationService{
		repo:  repo,
		cache: make(map[int64]*UserPreferences),
	}
}

func (s *recommendationService) GetPersonalizedRecommendations(ctx context.Context, userID int64, limit int) ([]int64, error) {
	s.mu.RLock()
	prefs, exists := s.cache[userID]
	s.mu.RUnlock()
	
	if !exists || time.Since(prefs.LastUpdated) > time.Hour {
		if err := s.RefreshUserPreferences(ctx, userID); err != nil {
			return s.GetHotRecommendations(ctx, limit)
		}
		s.mu.RLock()
		prefs = s.cache[userID]
		s.mu.RUnlock()
	}
	
	watchedVideos, err := s.repo.GetUserWatchHistory(ctx, userID, 100)
	if err != nil {
		return nil, err
	}
	
	watchedSet := make(map[int64]bool)
	for _, record := range watchedVideos {
		watchedSet[record.VideoID] = true
	}
	
	likedVideos, _ := s.repo.GetUserLikes(ctx, userID, 50)
	for _, id := range likedVideos {
		watchedSet[id] = true
	}
	
	favoriteVideos, _ := s.repo.GetUserFavorites(ctx, userID, 50)
	for _, id := range favoriteVideos {
		watchedSet[id] = true
	}
	
	categoryIDs := make([]int, 0, len(prefs.CategoryScores))
	for id := range prefs.CategoryScores {
		categoryIDs = append(categoryIDs, id)
	}
	
	categoryVideos, err := s.repo.GetCategoryVideos(ctx, categoryIDs, limit*3)
	if err != nil {
		return nil, err
	}
	
	scoredVideos := make(map[int64]float64)
	for _, videoID := range categoryVideos {
		if watchedSet[videoID] {
			continue
		}
		scoredVideos[videoID] += 1.0
	}
	
	similarUsers, _ := s.repo.GetSimilarUsers(ctx, userID, 20)
	if len(similarUsers) > 0 {
		for _, uid := range similarUsers {
			similarLikes, _ := s.repo.GetUserLikes(ctx, uid, 10)
			for _, videoID := range similarLikes {
				if !watchedSet[videoID] {
					scoredVideos[videoID] += 0.5
				}
			}
		}
	}
	
	type scoredVideo struct {
		VideoID int64
		Score   float64
	}
	
	var sortedVideos []scoredVideo
	for videoID, score := range scoredVideos {
		sortedVideos = append(sortedVideos, scoredVideo{VideoID: videoID, Score: score})
	}
	
	sort.Slice(sortedVideos, func(i, j int) bool {
		return sortedVideos[i].Score > sortedVideos[j].Score
	})
	
	result := make([]int64, 0, limit)
	for i := 0; i < len(sortedVideos) && len(result) < limit; i++ {
		result = append(result, sortedVideos[i].VideoID)
	}
	
	if len(result) < limit {
		hotVideos, _ := s.repo.GetHotVideos(ctx, limit-len(result))
		for _, videoID := range hotVideos {
			if !watchedSet[videoID] {
				result = append(result, videoID)
			}
		}
	}
	
	return result, nil
}

func (s *recommendationService) GetHotRecommendations(ctx context.Context, limit int) ([]int64, error) {
	return s.repo.GetHotVideos(ctx, limit)
}

func (s *recommendationService) GetSimilarVideos(ctx context.Context, videoID int64, limit int) ([]int64, error) {
	videoCategories, err := s.repo.GetVideoCategories(ctx, []int64{videoID})
	if err != nil {
		return nil, err
	}
	
	categoryIDs := videoCategories[videoID]
	if len(categoryIDs) == 0 {
		return s.GetHotRecommendations(ctx, limit)
	}
	
	similarVideos, err := s.repo.GetCategoryVideos(ctx, categoryIDs, limit+1)
	if err != nil {
		return nil, err
	}
	
	var result []int64
	for _, id := range similarVideos {
		if id != videoID {
			result = append(result, id)
		}
		if len(result) >= limit {
			break
		}
	}
	
	return result, nil
}

func (s *recommendationService) GetFollowingFeed(ctx context.Context, userID int64, limit int) ([]int64, error) {
	followings, err := s.repo.GetUserFollowings(ctx, userID, 100)
	if err != nil {
		return nil, err
	}
	
	if len(followings) == 0 {
		return s.GetHotRecommendations(ctx, limit)
	}
	
	return s.repo.GetFollowingVideos(ctx, followings, limit)
}

func (s *recommendationService) RefreshUserPreferences(ctx context.Context, userID int64) error {
	prefs := &UserPreferences{
		UserID:         userID,
		CategoryScores: make(map[int]float64),
		TagScores:      make(map[string]float64),
		LastUpdated:    time.Now(),
	}
	
	watchHistory, err := s.repo.GetUserWatchHistory(ctx, userID, 100)
	if err != nil {
		return err
	}
	
	videoIDs := make([]int64, 0, len(watchHistory))
	for _, record := range watchHistory {
		videoIDs = append(videoIDs, record.VideoID)
	}
	
	if len(videoIDs) > 0 {
		categories, err := s.repo.GetVideoCategories(ctx, videoIDs)
		if err != nil {
			return err
		}
		
		for _, record := range watchHistory {
			weight := math.Min(record.WatchProgress, 1.0)
			if cats, ok := categories[record.VideoID]; ok {
				for _, catID := range cats {
					prefs.CategoryScores[catID] += weight
				}
			}
		}
	}
	
	likedVideos, _ := s.repo.GetUserLikes(ctx, userID, 50)
	if len(likedVideos) > 0 {
		categories, _ := s.repo.GetVideoCategories(ctx, likedVideos)
		for _, videoID := range likedVideos {
			if cats, ok := categories[videoID]; ok {
				for _, catID := range cats {
					prefs.CategoryScores[catID] += 2.0
				}
			}
		}
	}
	
	favoriteVideos, _ := s.repo.GetUserFavorites(ctx, userID, 50)
	if len(favoriteVideos) > 0 {
		categories, _ := s.repo.GetVideoCategories(ctx, favoriteVideos)
		for _, videoID := range favoriteVideos {
			if cats, ok := categories[videoID]; ok {
				for _, catID := range cats {
					prefs.CategoryScores[catID] += 3.0
				}
			}
		}
	}
	
	s.mu.Lock()
	s.cache[userID] = prefs
	s.mu.Unlock()
	
	return nil
}
