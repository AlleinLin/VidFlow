package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/video-platform/go/internal/domain/video"
	apperrors "github.com/video-platform/go/pkg/errors"
)

type VideoRepository interface {
	Create(ctx context.Context, v *video.Video) error
	GetByID(ctx context.Context, id int64) (*video.Video, error)
	GetByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*video.Video, int64, error)
	List(ctx context.Context, filter *video.VideoFilter, page, pageSize int) ([]*video.Video, int64, error)
	Update(ctx context.Context, v *video.Video) error
	UpdateStatus(ctx context.Context, id int64, status video.VideoStatus) error
	Delete(ctx context.Context, id int64) error
	IncrementViewCount(ctx context.Context, id int64) error
	IncrementLikeCount(ctx context.Context, id int64, delta int) error
	IncrementCommentCount(ctx context.Context, id int64, delta int) error
	GetHotVideos(ctx context.Context, limit int) ([]*video.Video, error)
	GetVideoDuration(ctx context.Context, id int64) (int, error)
	Search(ctx context.Context, keyword string, page, pageSize int) ([]*video.Video, int64, error)
	
	CreateTranscodeTask(ctx context.Context, task *video.TranscodeTask) error
	GetTranscodeTasksByVideoID(ctx context.Context, videoID int64) ([]*video.TranscodeTask, error)
	UpdateTranscodeTask(ctx context.Context, task *video.TranscodeTask) error
}

type videoRepository struct {
	pool *pgxpool.Pool
}

func NewVideoRepository(pool *pgxpool.Pool) VideoRepository {
	return &videoRepository{pool: pool}
}

func (r *videoRepository) Create(ctx context.Context, v *video.Video) error {
	query := `
		INSERT INTO videos (user_id, title, description, status, visibility, original_filename, storage_key, category_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	
	now := time.Now()
	err := r.pool.QueryRow(ctx, query,
		v.UserID, v.Title, v.Description, v.Status, v.Visibility,
		v.OriginalFilename, v.StorageKey, v.CategoryID,
		now, now,
	).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create video: %w", err)
	}
	
	return nil
}

func (r *videoRepository) GetByID(ctx context.Context, id int64) (*video.Video, error) {
	query := `
		SELECT id, user_id, title, description, status, visibility, duration_seconds, original_filename,
			   storage_key, thumbnail_url, category_id, view_count, like_count, comment_count, published_at, created_at, updated_at
		FROM videos
		WHERE id = $1
	`
	
	var v video.Video
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&v.ID, &v.UserID, &v.Title, &v.Description, &v.Status, &v.Visibility,
		&v.DurationSeconds, &v.OriginalFilename, &v.StorageKey, &v.ThumbnailURL,
		&v.CategoryID, &v.ViewCount, &v.LikeCount, &v.CommentCount,
		&v.PublishedAt, &v.CreatedAt, &v.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrVideoNotFound
		}
		return nil, fmt.Errorf("failed to get video: %w", err)
	}
	
	return &v, nil
}

func (r *videoRepository) GetByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*video.Video, int64, error) {
	offset := (page - 1) * pageSize
	
	countQuery := `SELECT COUNT(*) FROM videos WHERE user_id = $1`
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user videos: %w", err)
	}
	
	query := `
		SELECT id, user_id, title, description, status, visibility, duration_seconds, original_filename,
			   storage_key, thumbnail_url, category_id, view_count, like_count, comment_count, published_at, created_at, updated_at
		FROM videos
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.pool.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user videos: %w", err)
	}
	defer rows.Close()
	
	var videos []*video.Video
	for rows.Next() {
		var v video.Video
		if err := rows.Scan(
			&v.ID, &v.UserID, &v.Title, &v.Description, &v.Status, &v.Visibility,
			&v.DurationSeconds, &v.OriginalFilename, &v.StorageKey, &v.ThumbnailURL,
			&v.CategoryID, &v.ViewCount, &v.LikeCount, &v.CommentCount,
			&v.PublishedAt, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan video: %w", err)
		}
		videos = append(videos, &v)
	}
	
	return videos, total, nil
}

func (r *videoRepository) List(ctx context.Context, filter *video.VideoFilter, page, pageSize int) ([]*video.Video, int64, error) {
	offset := (page - 1) * pageSize
	
	whereClause := "WHERE 1=1"
	args := make([]interface{}, 0)
	argNum := 1
	
	if filter.UserID > 0 {
		whereClause += fmt.Sprintf(" AND user_id = $%d", argNum)
		args = append(args, filter.UserID)
		argNum++
	}
	
	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, filter.Status)
		argNum++
	}
	
	if filter.Visibility != "" {
		whereClause += fmt.Sprintf(" AND visibility = $%d", argNum)
		args = append(args, filter.Visibility)
		argNum++
	}
	
	if filter.CategoryID > 0 {
		whereClause += fmt.Sprintf(" AND category_id = $%d", argNum)
		args = append(args, filter.CategoryID)
		argNum++
	}
	
	if filter.Keyword != "" {
		whereClause += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argNum, argNum+1)
		args = append(args, "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
		argNum += 2
	}
	
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM videos %s", whereClause)
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count videos: %w", err)
	}
	
	orderClause := "ORDER BY created_at DESC"
	if filter.SortBy != "" {
		orderDir := "ASC"
		if filter.SortDesc {
			orderDir = "DESC"
		}
		orderClause = fmt.Sprintf("ORDER BY %s %s", filter.SortBy, orderDir)
	}
	
	query := fmt.Sprintf(`
		SELECT id, user_id, title, description, status, visibility, duration_seconds, original_filename,
			   storage_key, thumbnail_url, category_id, view_count, like_count, comment_count, published_at, created_at, updated_at
		FROM videos
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderClause, argNum, argNum+1)
	
	args = append(args, pageSize, offset)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list videos: %w", err)
	}
	defer rows.Close()
	
	var videos []*video.Video
	for rows.Next() {
		var v video.Video
		if err := rows.Scan(
			&v.ID, &v.UserID, &v.Title, &v.Description, &v.Status, &v.Visibility,
			&v.DurationSeconds, &v.OriginalFilename, &v.StorageKey, &v.ThumbnailURL,
			&v.CategoryID, &v.ViewCount, &v.LikeCount, &v.CommentCount,
			&v.PublishedAt, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan video: %w", err)
		}
		videos = append(videos, &v)
	}
	
	return videos, total, nil
}

func (r *videoRepository) Update(ctx context.Context, v *video.Video) error {
	query := `
		UPDATE videos
		SET title = $2, description = $3, visibility = $4, category_id = $5, thumbnail_url = $6, updated_at = $7
		WHERE id = $1
	`
	
	result, err := r.pool.Exec(ctx, query,
		v.ID, v.Title, v.Description, v.Visibility, v.CategoryID, v.ThumbnailURL, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to update video: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return apperrors.ErrVideoNotFound
	}
	
	return nil
}

func (r *videoRepository) UpdateStatus(ctx context.Context, id int64, status video.VideoStatus) error {
	now := time.Now()
	var publishedAt *time.Time
	if status == video.StatusPublished {
		publishedAt = &now
	}
	
	query := `
		UPDATE videos
		SET status = $2, published_at = $3, updated_at = $4
		WHERE id = $1
	`
	
	result, err := r.pool.Exec(ctx, query, id, status, publishedAt, now)
	if err != nil {
		return fmt.Errorf("failed to update video status: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return apperrors.ErrVideoNotFound
	}
	
	return nil
}

func (r *videoRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM videos WHERE id = $1`
	
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return apperrors.ErrVideoNotFound
	}
	
	return nil
}

func (r *videoRepository) IncrementViewCount(ctx context.Context, id int64) error {
	query := `UPDATE videos SET view_count = view_count + 1 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *videoRepository) IncrementLikeCount(ctx context.Context, id int64, delta int) error {
	if delta >= 0 {
		query := `UPDATE videos SET like_count = like_count + $1 WHERE id = $2`
		_, err := r.pool.Exec(ctx, query, delta, id)
		return err
	}
	query := `UPDATE videos SET like_count = GREATEST(0, like_count + $1) WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, delta, id)
	return err
}

func (r *videoRepository) IncrementCommentCount(ctx context.Context, id int64, delta int) error {
	if delta >= 0 {
		query := `UPDATE videos SET comment_count = comment_count + $1 WHERE id = $2`
		_, err := r.pool.Exec(ctx, query, delta, id)
		return err
	}
	query := `UPDATE videos SET comment_count = GREATEST(0, comment_count + $1) WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, delta, id)
	return err
}

func (r *videoRepository) GetHotVideos(ctx context.Context, limit int) ([]*video.Video, error) {
	query := `
		SELECT id, user_id, title, description, status, visibility, duration_seconds, original_filename,
			   storage_key, thumbnail_url, category_id, view_count, like_count, comment_count, published_at, created_at, updated_at
		FROM videos
		WHERE status = 'published' AND visibility = 'public'
		ORDER BY (view_count * 0.3 + like_count * 0.5 + comment_count * 0.2) DESC, published_at DESC
		LIMIT $1
	`
	
	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get hot videos: %w", err)
	}
	defer rows.Close()
	
	var videos []*video.Video
	for rows.Next() {
		var v video.Video
		if err := rows.Scan(
			&v.ID, &v.UserID, &v.Title, &v.Description, &v.Status, &v.Visibility,
			&v.DurationSeconds, &v.OriginalFilename, &v.StorageKey, &v.ThumbnailURL,
			&v.CategoryID, &v.ViewCount, &v.LikeCount, &v.CommentCount,
			&v.PublishedAt, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan hot video: %w", err)
		}
		videos = append(videos, &v)
	}
	
	return videos, nil
}

func (r *videoRepository) GetVideoDuration(ctx context.Context, id int64) (int, error) {
	query := `SELECT COALESCE(duration_seconds, 0) FROM videos WHERE id = $1`
	var duration int
	err := r.pool.QueryRow(ctx, query, id).Scan(&duration)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, apperrors.ErrVideoNotFound
		}
		return 0, fmt.Errorf("failed to get video duration: %w", err)
	}
	return duration, nil
}

func (r *videoRepository) Search(ctx context.Context, keyword string, page, pageSize int) ([]*video.Video, int64, error) {
	offset := (page - 1) * pageSize
	
	countQuery := `
		SELECT COUNT(*) FROM videos 
		WHERE status = 'published' AND visibility = 'public'
		AND (title ILIKE $1 OR description ILIKE $1)
	`
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, "%"+keyword+"%").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}
	
	query := `
		SELECT id, user_id, title, description, status, visibility, duration_seconds, original_filename,
			   storage_key, thumbnail_url, category_id, view_count, like_count, comment_count, published_at, created_at, updated_at
		FROM videos
		WHERE status = 'published' AND visibility = 'public'
		AND (title ILIKE $1 OR description ILIKE $1)
		ORDER BY 
			CASE WHEN title ILIKE $2 THEN 0 ELSE 1 END,
			published_at DESC
		LIMIT $3 OFFSET $4
	`
	
	rows, err := r.pool.Query(ctx, query, "%"+keyword+"%", keyword, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search videos: %w", err)
	}
	defer rows.Close()
	
	var videos []*video.Video
	for rows.Next() {
		var v video.Video
		if err := rows.Scan(
			&v.ID, &v.UserID, &v.Title, &v.Description, &v.Status, &v.Visibility,
			&v.DurationSeconds, &v.OriginalFilename, &v.StorageKey, &v.ThumbnailURL,
			&v.CategoryID, &v.ViewCount, &v.LikeCount, &v.CommentCount,
			&v.PublishedAt, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan search result: %w", err)
		}
		videos = append(videos, &v)
	}
	
	return videos, total, nil
}

func (r *videoRepository) CreateTranscodeTask(ctx context.Context, task *video.TranscodeTask) error {
	query := `
		INSERT INTO video_transcode_tasks (video_id, resolution, format, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	
	err := r.pool.QueryRow(ctx, query,
		task.VideoID, task.Resolution, task.Format, task.Status, time.Now(),
	).Scan(&task.ID, &task.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create transcode task: %w", err)
	}
	
	return nil
}

func (r *videoRepository) GetTranscodeTasksByVideoID(ctx context.Context, videoID int64) ([]*video.TranscodeTask, error) {
	query := `
		SELECT id, video_id, resolution, format, status, output_storage_key, error_message, started_at, completed_at, created_at
		FROM video_transcode_tasks
		WHERE video_id = $1
		ORDER BY created_at
	`
	
	rows, err := r.pool.Query(ctx, query, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transcode tasks: %w", err)
	}
	defer rows.Close()
	
	var tasks []*video.TranscodeTask
	for rows.Next() {
		var task video.TranscodeTask
		if err := rows.Scan(
			&task.ID, &task.VideoID, &task.Resolution, &task.Format, &task.Status,
			&task.OutputStorageKey, &task.ErrorMessage, &task.StartedAt, &task.CompletedAt, &task.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transcode task: %w", err)
		}
		tasks = append(tasks, &task)
	}
	
	return tasks, nil
}

func (r *videoRepository) UpdateTranscodeTask(ctx context.Context, task *video.TranscodeTask) error {
	query := `
		UPDATE video_transcode_tasks
		SET status = $2, output_storage_key = $3, error_message = $4, started_at = $5, completed_at = $6
		WHERE id = $1
	`
	
	_, err := r.pool.Exec(ctx, query,
		task.ID, task.Status, task.OutputStorageKey, task.ErrorMessage, task.StartedAt, task.CompletedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update transcode task: %w", err)
	}
	
	return nil
}
