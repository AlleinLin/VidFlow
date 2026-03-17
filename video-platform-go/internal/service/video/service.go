package video

import (
	"context"
	"fmt"
	"time"

	"github.com/video-platform/go/internal/domain/video"
	apperrors "github.com/video-platform/go/pkg/errors"
)

type Repository interface {
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
	Search(ctx context.Context, keyword string, page, pageSize int) ([]*video.Video, int64, error)
}

type Service interface {
	UploadVideo(ctx context.Context, userID int64, req *video.VideoUploadRequest) (*video.VideoUploadResponse, error)
	GetVideo(ctx context.Context, id int64, viewerID int64) (*video.Video, error)
	GetUserVideos(ctx context.Context, userID int64, page, pageSize int) (*video.VideoListResponse, error)
	ListVideos(ctx context.Context, filter *video.VideoFilter, page, pageSize int) (*video.VideoListResponse, error)
	UpdateVideo(ctx context.Context, id int64, userID int64, req *video.VideoUpdateRequest, isAdmin bool) (*video.Video, error)
	DeleteVideo(ctx context.Context, id int64, userID int64, isAdmin bool) error
	PublishVideo(ctx context.Context, id int64, userID int64, isAdmin bool) error
	GetHotVideos(ctx context.Context, limit int) ([]*video.Video, error)
	SearchVideos(ctx context.Context, keyword string, page, pageSize int) (*video.VideoListResponse, error)
	RecordView(ctx context.Context, videoID int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) UploadVideo(ctx context.Context, userID int64, req *video.VideoUploadRequest) (*video.VideoUploadResponse, error) {
	v := &video.Video{
		UserID:           userID,
		Title:            req.Title,
		Description:      req.Description,
		Status:           video.StatusUploading,
		Visibility:       req.Visibility,
		OriginalFilename: req.Filename,
		CategoryID:       req.CategoryID,
	}

	if err := s.repo.Create(ctx, v); err != nil {
		return nil, fmt.Errorf("failed to create video record: %w", err)
	}

	uploadToken := generateUploadToken(v.ID)
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	return &video.VideoUploadResponse{
		VideoID:     v.ID,
		UploadURL:   fmt.Sprintf("/api/v1/videos/%d/upload", v.ID),
		UploadToken: uploadToken,
		ExpiresAt:   expiresAt,
	}, nil
}

func (s *service) GetVideo(ctx context.Context, id int64, viewerID int64) (*video.Video, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !v.IsPublished() {
		if viewerID == 0 || (viewerID != v.UserID) {
			return nil, apperrors.ErrVideoNotFound
		}
	}

	return v, nil
}

func (s *service) GetUserVideos(ctx context.Context, userID int64, page, pageSize int) (*video.VideoListResponse, error) {
	videos, total, err := s.repo.GetByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &video.VideoListResponse{
		Videos:     videos,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *service) ListVideos(ctx context.Context, filter *video.VideoFilter, page, pageSize int) (*video.VideoListResponse, error) {
	videos, total, err := s.repo.List(ctx, filter, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &video.VideoListResponse{
		Videos:     videos,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *service) UpdateVideo(ctx context.Context, id int64, userID int64, req *video.VideoUpdateRequest, isAdmin bool) (*video.Video, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !v.CanEdit(userID, isAdmin) {
		return nil, apperrors.ErrForbidden
	}

	if req.Title != nil {
		v.Title = *req.Title
	}
	if req.Description != nil {
		v.Description = *req.Description
	}
	if req.Visibility != nil {
		v.Visibility = *req.Visibility
	}
	if req.CategoryID != nil {
		v.CategoryID = *req.CategoryID
	}
	if req.ThumbnailURL != nil {
		v.ThumbnailURL = *req.ThumbnailURL
	}

	if err := s.repo.Update(ctx, v); err != nil {
		return nil, err
	}

	return v, nil
}

func (s *service) DeleteVideo(ctx context.Context, id int64, userID int64, isAdmin bool) error {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !v.CanDelete(userID, isAdmin) {
		return apperrors.ErrForbidden
	}

	return s.repo.Delete(ctx, id)
}

func (s *service) PublishVideo(ctx context.Context, id int64, userID int64, isAdmin bool) error {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !v.CanEdit(userID, isAdmin) {
		return apperrors.ErrForbidden
	}

	if !v.CanTransitionTo(video.StatusPublished) {
		return apperrors.ErrInvalidVideoStatus
	}

	return s.repo.UpdateStatus(ctx, id, video.StatusPublished)
}

func (s *service) GetHotVideos(ctx context.Context, limit int) ([]*video.Video, error) {
	return s.repo.GetHotVideos(ctx, limit)
}

func (s *service) SearchVideos(ctx context.Context, keyword string, page, pageSize int) (*video.VideoListResponse, error) {
	videos, total, err := s.repo.Search(ctx, keyword, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &video.VideoListResponse{
		Videos:     videos,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *service) RecordView(ctx context.Context, videoID int64) error {
	return s.repo.IncrementViewCount(ctx, videoID)
}

func generateUploadToken(videoID int64) string {
	return fmt.Sprintf("upload_%d_%d", videoID, time.Now().UnixNano())
}
