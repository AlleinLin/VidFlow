package interaction

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/video-platform/go/internal/domain/interaction"
	apperrors "github.com/video-platform/go/pkg/errors"
)

type InteractionRepository interface {
	CreateComment(ctx context.Context, c *interaction.Comment) error
	GetCommentByID(ctx context.Context, id int64) (*interaction.Comment, error)
	GetCommentsByVideoID(ctx context.Context, videoID int64, page, pageSize int) ([]*interaction.Comment, int64, error)
	GetReplies(ctx context.Context, rootID int64, page, pageSize int) ([]*interaction.Comment, int64, error)
	UpdateComment(ctx context.Context, id int64, content string) error
	DeleteComment(ctx context.Context, id int64) error
	
	CreateLike(ctx context.Context, userID, videoID int64) error
	DeleteLike(ctx context.Context, userID, videoID int64) error
	IsLiked(ctx context.Context, userID, videoID int64) (bool, error)
	GetLikeCount(ctx context.Context, videoID int64) (int64, error)
	
	CreateFavorite(ctx context.Context, userID, videoID int64) error
	DeleteFavorite(ctx context.Context, userID, videoID int64) error
	IsFavorited(ctx context.Context, userID, videoID int64) (bool, error)
	GetFavorites(ctx context.Context, userID int64, page, pageSize int) ([]int64, int64, error)
	
	CreateDanmaku(ctx context.Context, d *interaction.Danmaku) error
	GetDanmakusByVideoID(ctx context.Context, videoID int64, startTime, endTime float64) ([]*interaction.Danmaku, error)
}

type interactionRepository struct {
	pool *pgxpool.Pool
}

func NewInteractionRepository(pool *pgxpool.Pool) InteractionRepository {
	return &interactionRepository{pool: pool}
}

func (r *interactionRepository) CreateComment(ctx context.Context, c *interaction.Comment) error {
	query := `
		INSERT INTO comments (video_id, user_id, parent_id, root_id, content, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	
	now := time.Now()
	return r.pool.QueryRow(ctx, query,
		c.VideoID, c.UserID, c.ParentID, c.RootID, c.Content, c.Status, now,
	).Scan(&c.ID, &c.CreatedAt)
}

func (r *interactionRepository) GetCommentByID(ctx context.Context, id int64) (*interaction.Comment, error) {
	query := `
		SELECT id, video_id, user_id, parent_id, root_id, content, like_count, status, created_at
		FROM comments WHERE id = $1
	`
	
	var c interaction.Comment
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.VideoID, &c.UserID, &c.ParentID, &c.RootID,
		&c.Content, &c.LikeCount, &c.Status, &c.CreatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrCommentNotFound
		}
		return nil, err
	}
	
	return &c, nil
}

func (r *interactionRepository) GetCommentsByVideoID(ctx context.Context, videoID int64, page, pageSize int) ([]*interaction.Comment, int64, error) {
	countQuery := `SELECT COUNT(*) FROM comments WHERE video_id = $1 AND parent_id IS NULL AND status = 'visible'`
	var total int64
	r.pool.QueryRow(ctx, countQuery, videoID).Scan(&total)
	
	offset := (page - 1) * pageSize
	query := `
		SELECT c.id, c.video_id, c.user_id, c.parent_id, c.root_id, c.content, c.like_count, c.status, c.created_at,
			   u.id, u.username, u.display_name, u.avatar_url
		FROM comments c
		INNER JOIN users u ON c.user_id = u.id
		WHERE c.video_id = $1 AND c.parent_id IS NULL AND c.status = 'visible'
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.pool.Query(ctx, query, videoID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var comments []*interaction.Comment
	for rows.Next() {
		var c interaction.Comment
		c.User = &interaction.CommentUser{}
		if err := rows.Scan(
			&c.ID, &c.VideoID, &c.UserID, &c.ParentID, &c.RootID,
			&c.Content, &c.LikeCount, &c.Status, &c.CreatedAt,
			&c.User.ID, &c.User.Username, &c.User.DisplayName, &c.User.AvatarURL,
		); err != nil {
			return nil, 0, err
		}
		comments = append(comments, &c)
	}
	
	return comments, total, nil
}

func (r *interactionRepository) GetReplies(ctx context.Context, rootID int64, page, pageSize int) ([]*interaction.Comment, int64, error) {
	countQuery := `SELECT COUNT(*) FROM comments WHERE root_id = $1 AND status = 'visible'`
	var total int64
	r.pool.QueryRow(ctx, countQuery, rootID).Scan(&total)
	
	offset := (page - 1) * pageSize
	query := `
		SELECT c.id, c.video_id, c.user_id, c.parent_id, c.root_id, c.content, c.like_count, c.status, c.created_at,
			   u.id, u.username, u.display_name, u.avatar_url
		FROM comments c
		INNER JOIN users u ON c.user_id = u.id
		WHERE c.root_id = $1 AND c.status = 'visible'
		ORDER BY c.created_at ASC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.pool.Query(ctx, query, rootID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var comments []*interaction.Comment
	for rows.Next() {
		var c interaction.Comment
		c.User = &interaction.CommentUser{}
		if err := rows.Scan(
			&c.ID, &c.VideoID, &c.UserID, &c.ParentID, &c.RootID,
			&c.Content, &c.LikeCount, &c.Status, &c.CreatedAt,
			&c.User.ID, &c.User.Username, &c.User.DisplayName, &c.User.AvatarURL,
		); err != nil {
			return nil, 0, err
		}
		comments = append(comments, &c)
	}
	
	return comments, total, nil
}

func (r *interactionRepository) UpdateComment(ctx context.Context, id int64, content string) error {
	query := `UPDATE comments SET content = $2 WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id, content)
	if err != nil {
		return err
	}
	
	if result.RowsAffected() == 0 {
		return apperrors.ErrCommentNotFound
	}
	
	return nil
}

func (r *interactionRepository) DeleteComment(ctx context.Context, id int64) error {
	query := `UPDATE comments SET status = 'deleted' WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	
	if result.RowsAffected() == 0 {
		return apperrors.ErrCommentNotFound
	}
	
	return nil
}

func (r *interactionRepository) CreateLike(ctx context.Context, userID, videoID int64) error {
	query := `INSERT INTO likes (user_id, video_id, created_at) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, query, userID, videoID, time.Now())
	return err
}

func (r *interactionRepository) DeleteLike(ctx context.Context, userID, videoID int64) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND video_id = $2`
	_, err := r.pool.Exec(ctx, query, userID, videoID)
	return err
}

func (r *interactionRepository) IsLiked(ctx context.Context, userID, videoID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND video_id = $2)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, userID, videoID).Scan(&exists)
	return exists, err
}

func (r *interactionRepository) GetLikeCount(ctx context.Context, videoID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM likes WHERE video_id = $1`
	var count int64
	err := r.pool.QueryRow(ctx, query, videoID).Scan(&count)
	return count, err
}

func (r *interactionRepository) CreateFavorite(ctx context.Context, userID, videoID int64) error {
	query := `INSERT INTO favorites (user_id, video_id, created_at) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, query, userID, videoID, time.Now())
	return err
}

func (r *interactionRepository) DeleteFavorite(ctx context.Context, userID, videoID int64) error {
	query := `DELETE FROM favorites WHERE user_id = $1 AND video_id = $2`
	_, err := r.pool.Exec(ctx, query, userID, videoID)
	return err
}

func (r *interactionRepository) IsFavorited(ctx context.Context, userID, videoID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND video_id = $2)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, userID, videoID).Scan(&exists)
	return exists, err
}

func (r *interactionRepository) GetFavorites(ctx context.Context, userID int64, page, pageSize int) ([]int64, int64, error) {
	countQuery := `SELECT COUNT(*) FROM favorites WHERE user_id = $1`
	var total int64
	r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	
	offset := (page - 1) * pageSize
	query := `SELECT video_id FROM favorites WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := r.pool.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var videoIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		videoIDs = append(videoIDs, id)
	}
	
	return videoIDs, total, nil
}

func (r *interactionRepository) CreateDanmaku(ctx context.Context, d *interaction.Danmaku) error {
	query := `
		INSERT INTO danmakus (video_id, user_id, content, position_seconds, style, color, font_size, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	
	return r.pool.QueryRow(ctx, query,
		d.VideoID, d.UserID, d.Content, d.PositionSeconds, d.Style, d.Color, d.FontSize, time.Now(),
	).Scan(&d.ID, &d.CreatedAt)
}

func (r *interactionRepository) GetDanmakusByVideoID(ctx context.Context, videoID int64, startTime, endTime float64) ([]*interaction.Danmaku, error) {
	query := `
		SELECT id, video_id, user_id, content, position_seconds, style, color, font_size, created_at
		FROM danmakus
		WHERE video_id = $1 AND position_seconds >= $2 AND position_seconds <= $3
		ORDER BY position_seconds
	`
	
	rows, err := r.pool.Query(ctx, query, videoID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var danmakus []*interaction.Danmaku
	for rows.Next() {
		var d interaction.Danmaku
		if err := rows.Scan(
			&d.ID, &d.VideoID, &d.UserID, &d.Content, &d.PositionSeconds,
			&d.Style, &d.Color, &d.FontSize, &d.CreatedAt,
		); err != nil {
			return nil, err
		}
		danmakus = append(danmakus, &d)
	}
	
	return danmakus, nil
}

type InteractionService interface {
	CreateComment(ctx context.Context, userID int64, req *interaction.CreateCommentRequest) (*interaction.Comment, error)
	GetComments(ctx context.Context, videoID int64, page, pageSize int) (*interaction.CommentListResponse, error)
	GetReplies(ctx context.Context, rootID int64, page, pageSize int) (*interaction.CommentListResponse, error)
	UpdateComment(ctx context.Context, id int64, userID int64, content string, isAdmin bool) error
	DeleteComment(ctx context.Context, id int64, userID int64, isAdmin bool) error
	
	LikeVideo(ctx context.Context, userID, videoID int64) error
	UnlikeVideo(ctx context.Context, userID, videoID int64) error
	GetLikeStatus(ctx context.Context, userID, videoID int64) (*interaction.LikeStatus, error)
	
	FavoriteVideo(ctx context.Context, userID, videoID int64) error
	UnfavoriteVideo(ctx context.Context, userID, videoID int64) error
	
	CreateDanmaku(ctx context.Context, userID int64, req *interaction.CreateDanmakuRequest) (*interaction.Danmaku, error)
	GetDanmakus(ctx context.Context, videoID int64, startTime, endTime float64) (*interaction.DanmakuListResponse, error)
}

type interactionService struct {
	repo       InteractionRepository
	videoRepo  VideoRepository
}

type VideoRepository interface {
	IncrementCommentCount(ctx context.Context, id int64, delta int) error
	IncrementLikeCount(ctx context.Context, id int64, delta int) error
}

func NewInteractionService(repo InteractionRepository, videoRepo VideoRepository) InteractionService {
	return &interactionService{repo: repo, videoRepo: videoRepo}
}

func (s *interactionService) CreateComment(ctx context.Context, userID int64, req *interaction.CreateCommentRequest) (*interaction.Comment, error) {
	c := &interaction.Comment{
		VideoID: req.VideoID,
		UserID:  userID,
		Content: req.Content,
		Status:  interaction.CommentStatusVisible,
	}
	
	if req.ParentID != nil {
		c.ParentID = req.ParentID
		parent, err := s.repo.GetCommentByID(ctx, *req.ParentID)
		if err != nil {
			return nil, err
		}
		if parent.RootID != nil {
			c.RootID = parent.RootID
		} else {
			rootID := parent.ID
			c.RootID = &rootID
		}
	}
	
	if err := s.repo.CreateComment(ctx, c); err != nil {
		return nil, err
	}
	
	go func() {
		_ = s.videoRepo.IncrementCommentCount(context.Background(), req.VideoID, 1)
	}()
	
	return c, nil
}

func (s *interactionService) GetComments(ctx context.Context, videoID int64, page, pageSize int) (*interaction.CommentListResponse, error) {
	comments, total, err := s.repo.GetCommentsByVideoID(ctx, videoID, page, pageSize)
	if err != nil {
		return nil, err
	}
	
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	
	return &interaction.CommentListResponse{
		Comments:   comments,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *interactionService) GetReplies(ctx context.Context, rootID int64, page, pageSize int) (*interaction.CommentListResponse, error) {
	comments, total, err := s.repo.GetReplies(ctx, rootID, page, pageSize)
	if err != nil {
		return nil, err
	}
	
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	
	return &interaction.CommentListResponse{
		Comments:   comments,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *interactionService) UpdateComment(ctx context.Context, id int64, userID int64, content string, isAdmin bool) error {
	c, err := s.repo.GetCommentByID(ctx, id)
	if err != nil {
		return err
	}
	
	if !c.CanEdit(userID, isAdmin) {
		return apperrors.ErrForbidden
	}
	
	return s.repo.UpdateComment(ctx, id, content)
}

func (s *interactionService) DeleteComment(ctx context.Context, id int64, userID int64, isAdmin bool) error {
	c, err := s.repo.GetCommentByID(ctx, id)
	if err != nil {
		return err
	}
	
	if !c.CanDelete(userID, isAdmin) {
		return apperrors.ErrForbidden
	}
	
	if err := s.repo.DeleteComment(ctx, id); err != nil {
		return err
	}
	
	go func() {
		_ = s.videoRepo.IncrementCommentCount(context.Background(), c.VideoID, -1)
	}()
	
	return nil
}

func (s *interactionService) LikeVideo(ctx context.Context, userID, videoID int64) error {
	isLiked, err := s.repo.IsLiked(ctx, userID, videoID)
	if err != nil {
		return err
	}
	
	if isLiked {
		return nil
	}
	
	if err := s.repo.CreateLike(ctx, userID, videoID); err != nil {
		return err
	}
	
	go func() {
		_ = s.videoRepo.IncrementLikeCount(context.Background(), videoID, 1)
	}()
	
	return nil
}

func (s *interactionService) UnlikeVideo(ctx context.Context, userID, videoID int64) error {
	isLiked, err := s.repo.IsLiked(ctx, userID, videoID)
	if err != nil {
		return err
	}
	
	if !isLiked {
		return nil
	}
	
	if err := s.repo.DeleteLike(ctx, userID, videoID); err != nil {
		return err
	}
	
	go func() {
		_ = s.videoRepo.IncrementLikeCount(context.Background(), videoID, -1)
	}()
	
	return nil
}

func (s *interactionService) GetLikeStatus(ctx context.Context, userID, videoID int64) (*interaction.LikeStatus, error) {
	isLiked, err := s.repo.IsLiked(ctx, userID, videoID)
	if err != nil {
		return nil, err
	}
	
	isFavorited, err := s.repo.IsFavorited(ctx, userID, videoID)
	if err != nil {
		return nil, err
	}
	
	likeCount, err := s.repo.GetLikeCount(ctx, videoID)
	if err != nil {
		return nil, err
	}
	
	return &interaction.LikeStatus{
		IsLiked:     isLiked,
		IsFavorited: isFavorited,
		LikeCount:   likeCount,
	}, nil
}

func (s *interactionService) FavoriteVideo(ctx context.Context, userID, videoID int64) error {
	return s.repo.CreateFavorite(ctx, userID, videoID)
}

func (s *interactionService) UnfavoriteVideo(ctx context.Context, userID, videoID int64) error {
	return s.repo.DeleteFavorite(ctx, userID, videoID)
}

func (s *interactionService) CreateDanmaku(ctx context.Context, userID int64, req *interaction.CreateDanmakuRequest) (*interaction.Danmaku, error) {
	d := &interaction.Danmaku{
		VideoID:         req.VideoID,
		UserID:          userID,
		Content:         req.Content,
		PositionSeconds: req.PositionSeconds,
		Style:           req.Style,
		Color:           req.Color,
		FontSize:        req.FontSize,
	}
	
	if err := s.repo.CreateDanmaku(ctx, d); err != nil {
		return nil, err
	}
	
	return d, nil
}

func (s *interactionService) GetDanmakus(ctx context.Context, videoID int64, startTime, endTime float64) (*interaction.DanmakuListResponse, error) {
	danmakus, err := s.repo.GetDanmakusByVideoID(ctx, videoID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	
	return &interaction.DanmakuListResponse{Danmakus: danmakus}, nil
}
