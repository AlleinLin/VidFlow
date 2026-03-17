package playback

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WatchHistory struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	VideoID        int64     `json:"video_id"`
	WatchDuration  int64     `json:"watch_duration"`
	WatchProgress  float64   `json:"watch_progress"`
	LastPosition   float64   `json:"last_position"`
	Completed      bool      `json:"completed"`
	WatchedAt      time.Time `json:"watched_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PlaybackProgress struct {
	VideoID       int64   `json:"video_id"`
	Position      float64 `json:"position"`
	Duration      float64 `json:"duration"`
	Progress      float64 `json:"progress"`
	WatchDuration int64   `json:"watch_duration"`
}

type PlaybackRepository interface {
	GetWatchHistory(ctx context.Context, userID int64, page, pageSize int) ([]*WatchHistory, int64, error)
	GetWatchHistoryByVideo(ctx context.Context, userID, videoID int64) (*WatchHistory, error)
	UpsertWatchHistory(ctx context.Context, history *WatchHistory) error
	DeleteWatchHistory(ctx context.Context, userID, videoID int64) error
	ClearWatchHistory(ctx context.Context, userID int64) error
	GetContinueWatching(ctx context.Context, userID int64, limit int) ([]*WatchHistory, error)
}

type playbackRepository struct {
	pool *pgxpool.Pool
}

func NewPlaybackRepository(pool *pgxpool.Pool) PlaybackRepository {
	return &playbackRepository{pool: pool}
}

func (r *playbackRepository) GetWatchHistory(ctx context.Context, userID int64, page, pageSize int) ([]*WatchHistory, int64, error) {
	countQuery := `SELECT COUNT(*) FROM watch_history WHERE user_id = $1`
	var total int64
	r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	
	offset := (page - 1) * pageSize
	query := `
		SELECT id, user_id, video_id, watch_duration, watch_progress, last_position, completed, watched_at, created_at, updated_at
		FROM watch_history
		WHERE user_id = $1
		ORDER BY watched_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.pool.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var histories []*WatchHistory
	for rows.Next() {
		var h WatchHistory
		if err := rows.Scan(
			&h.ID, &h.UserID, &h.VideoID, &h.WatchDuration, &h.WatchProgress,
			&h.LastPosition, &h.Completed, &h.WatchedAt, &h.CreatedAt, &h.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		histories = append(histories, &h)
	}
	
	return histories, total, nil
}

func (r *playbackRepository) GetWatchHistoryByVideo(ctx context.Context, userID, videoID int64) (*WatchHistory, error) {
	query := `
		SELECT id, user_id, video_id, watch_duration, watch_progress, last_position, completed, watched_at, created_at, updated_at
		FROM watch_history
		WHERE user_id = $1 AND video_id = $2
	`
	
	var h WatchHistory
	err := r.pool.QueryRow(ctx, query, userID, videoID).Scan(
		&h.ID, &h.UserID, &h.VideoID, &h.WatchDuration, &h.WatchProgress,
		&h.LastPosition, &h.Completed, &h.WatchedAt, &h.CreatedAt, &h.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	return &h, nil
}

func (r *playbackRepository) UpsertWatchHistory(ctx context.Context, history *WatchHistory) error {
	query := `
		INSERT INTO watch_history (user_id, video_id, watch_duration, watch_progress, last_position, completed, watched_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (user_id, video_id) DO UPDATE SET
			watch_duration = EXCLUDED.watch_duration,
			watch_progress = EXCLUDED.watch_progress,
			last_position = EXCLUDED.last_position,
			completed = EXCLUDED.completed,
			watched_at = EXCLUDED.watched_at,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	
	return r.pool.QueryRow(ctx, query,
		history.UserID, history.VideoID, history.WatchDuration, history.WatchProgress,
		history.LastPosition, history.Completed, history.WatchedAt,
	).Scan(&history.ID, &history.CreatedAt, &history.UpdatedAt)
}

func (r *playbackRepository) DeleteWatchHistory(ctx context.Context, userID, videoID int64) error {
	query := `DELETE FROM watch_history WHERE user_id = $1 AND video_id = $2`
	_, err := r.pool.Exec(ctx, query, userID, videoID)
	return err
}

func (r *playbackRepository) ClearWatchHistory(ctx context.Context, userID int64) error {
	query := `DELETE FROM watch_history WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}

func (r *playbackRepository) GetContinueWatching(ctx context.Context, userID int64, limit int) ([]*WatchHistory, error) {
	query := `
		SELECT id, user_id, video_id, watch_duration, watch_progress, last_position, completed, watched_at, created_at, updated_at
		FROM watch_history
		WHERE user_id = $1 AND completed = false AND watch_progress > 0.05
		ORDER BY watched_at DESC
		LIMIT $2
	`
	
	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var histories []*WatchHistory
	for rows.Next() {
		var h WatchHistory
		if err := rows.Scan(
			&h.ID, &h.UserID, &h.VideoID, &h.WatchDuration, &h.WatchProgress,
			&h.LastPosition, &h.Completed, &h.WatchedAt, &h.CreatedAt, &h.UpdatedAt,
		); err != nil {
			return nil, err
		}
		histories = append(histories, &h)
	}
	
	return histories, nil
}

type VideoRepository interface {
	IncrementViewCount(ctx context.Context, id int64) error
	GetVideoDuration(ctx context.Context, id int64) (int, error)
}

type PlaybackService interface {
	UpdateProgress(ctx context.Context, userID, videoID int64, position, duration float64, watchDuration int64) error
	GetProgress(ctx context.Context, userID, videoID int64) (*WatchHistory, error)
	GetWatchHistory(ctx context.Context, userID int64, page, pageSize int) ([]*WatchHistory, int64, error)
	GetContinueWatching(ctx context.Context, userID int64, limit int) ([]*WatchHistory, error)
	DeleteWatchHistory(ctx context.Context, userID, videoID int64) error
	ClearWatchHistory(ctx context.Context, userID int64) error
}

type playbackService struct {
	repo       PlaybackRepository
	videoRepo  VideoRepository
}

func NewPlaybackService(repo PlaybackRepository, videoRepo VideoRepository) PlaybackService {
	return &playbackService{
		repo:      repo,
		videoRepo: videoRepo,
	}
}

func (s *playbackService) UpdateProgress(ctx context.Context, userID, videoID int64, position, duration float64, watchDuration int64) error {
	var progress float64
	if duration > 0 {
		progress = position / duration
	}
	
	completed := progress >= 0.95
	
	history := &WatchHistory{
		UserID:        userID,
		VideoID:       videoID,
		WatchDuration: watchDuration,
		WatchProgress: progress,
		LastPosition:  position,
		Completed:     completed,
		WatchedAt:     time.Now(),
	}
	
	if err := s.repo.UpsertWatchHistory(ctx, history); err != nil {
		return fmt.Errorf("failed to update watch history: %w", err)
	}
	
	if progress < 0.1 {
		if err := s.videoRepo.IncrementViewCount(ctx, videoID); err != nil {
			return fmt.Errorf("failed to increment view count: %w", err)
		}
	}
	
	return nil
}

func (s *playbackService) GetProgress(ctx context.Context, userID, videoID int64) (*WatchHistory, error) {
	return s.repo.GetWatchHistoryByVideo(ctx, userID, videoID)
}

func (s *playbackService) GetWatchHistory(ctx context.Context, userID int64, page, pageSize int) ([]*WatchHistory, int64, error) {
	return s.repo.GetWatchHistory(ctx, userID, page, pageSize)
}

func (s *playbackService) GetContinueWatching(ctx context.Context, userID int64, limit int) ([]*WatchHistory, error) {
	return s.repo.GetContinueWatching(ctx, userID, limit)
}

func (s *playbackService) DeleteWatchHistory(ctx context.Context, userID, videoID int64) error {
	return s.repo.DeleteWatchHistory(ctx, userID, videoID)
}

func (s *playbackService) ClearWatchHistory(ctx context.Context, userID int64) error {
	return s.repo.ClearWatchHistory(ctx, userID)
}
