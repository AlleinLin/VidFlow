package video_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/video-platform/go/internal/domain/video"
	videoService "github.com/video-platform/go/internal/service/video"
	apperrors "github.com/video-platform/go/pkg/errors"
)

type MockVideoRepository struct {
	mock.Mock
}

func (m *MockVideoRepository) Create(ctx context.Context, v *video.Video) error {
	args := m.Called(ctx, v)
	return args.Error(0)
}

func (m *MockVideoRepository) GetByID(ctx context.Context, id int64) (*video.Video, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*video.Video), args.Error(1)
}

func (m *MockVideoRepository) GetByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*video.Video, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*video.Video), args.Get(1).(int64), args.Error(2)
}

func (m *MockVideoRepository) List(ctx context.Context, filter *video.VideoFilter, page, pageSize int) ([]*video.Video, int64, error) {
	args := m.Called(ctx, filter, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*video.Video), args.Get(1).(int64), args.Error(2)
}

func (m *MockVideoRepository) Update(ctx context.Context, v *video.Video) error {
	args := m.Called(ctx, v)
	return args.Error(0)
}

func (m *MockVideoRepository) UpdateStatus(ctx context.Context, id int64, status video.VideoStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockVideoRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVideoRepository) IncrementViewCount(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVideoRepository) IncrementLikeCount(ctx context.Context, id int64, delta int) error {
	args := m.Called(ctx, id, delta)
	return args.Error(0)
}

func (m *MockVideoRepository) IncrementCommentCount(ctx context.Context, id int64, delta int) error {
	args := m.Called(ctx, id, delta)
	return args.Error(0)
}

func (m *MockVideoRepository) GetHotVideos(ctx context.Context, limit int) ([]*video.Video, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*video.Video), args.Error(1)
}

func (m *MockVideoRepository) Search(ctx context.Context, keyword string, page, pageSize int) ([]*video.Video, int64, error) {
	args := m.Called(ctx, keyword, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*video.Video), args.Get(1).(int64), args.Error(2)
}

func (m *MockVideoRepository) CreateTranscodeTask(ctx context.Context, task *video.TranscodeTask) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockVideoRepository) GetTranscodeTasksByVideoID(ctx context.Context, videoID int64) ([]*video.TranscodeTask, error) {
	args := m.Called(ctx, videoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*video.TranscodeTask), args.Error(1)
}

func (m *MockVideoRepository) UpdateTranscodeTask(ctx context.Context, task *video.TranscodeTask) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockVideoRepository) GetVideoDuration(ctx context.Context, id int64) (int, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int), args.Error(1)
}

func TestGetVideo_Success(t *testing.T) {
	mockRepo := new(MockVideoRepository)
	svc := videoService.NewService(mockRepo)

	expectedVideo := &video.Video{
		ID:          1,
		UserID:      1,
		Title:       "Test Video",
		Description: "Test description",
		Status:      video.StatusPublished,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedVideo, nil)

	v, err := svc.GetVideo(context.Background(), 1, 0)

	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, int64(1), v.ID)
	assert.Equal(t, "Test Video", v.Title)
	mockRepo.AssertExpectations(t)
}

func TestListVideos_Success(t *testing.T) {
	mockRepo := new(MockVideoRepository)
	svc := videoService.NewService(mockRepo)

	videos := []*video.Video{
		{ID: 1, Title: "Video 1", Status: video.StatusPublished},
		{ID: 2, Title: "Video 2", Status: video.StatusPublished},
	}

	filter := &video.VideoFilter{Status: video.StatusPublished}
	mockRepo.On("List", mock.Anything, filter, 1, 20).Return(videos, int64(2), nil)

	resp, err := svc.ListVideos(context.Background(), filter, 1, 20)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Videos, 2)
	assert.Equal(t, int64(2), resp.Total)
	mockRepo.AssertExpectations(t)
}

func TestDeleteVideo_Success(t *testing.T) {
	mockRepo := new(MockVideoRepository)
	svc := videoService.NewService(mockRepo)

	existingVideo := &video.Video{
		ID:     1,
		UserID: 1,
		Status: video.StatusUploading,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingVideo, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := svc.DeleteVideo(context.Background(), 1, 1, false)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteVideo_Unauthorized(t *testing.T) {
	mockRepo := new(MockVideoRepository)
	svc := videoService.NewService(mockRepo)

	existingVideo := &video.Video{
		ID:     1,
		UserID: 2,
		Status: video.StatusUploading,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingVideo, nil)

	err := svc.DeleteVideo(context.Background(), 1, 1, false)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, apperrors.ErrForbidden))
	mockRepo.AssertExpectations(t)
}

func TestUploadVideo_Success(t *testing.T) {
	mockRepo := new(MockVideoRepository)
	svc := videoService.NewService(mockRepo)

	req := &video.VideoUploadRequest{
		Title:       "Test Video",
		Description: "Test description",
		Visibility:  video.VisibilityPublic,
		CategoryID:  1,
		Filename:    "test.mp4",
		FileSize:    1024000,
	}

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	resp, err := svc.UploadVideo(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestPublishVideo_Success(t *testing.T) {
	mockRepo := new(MockVideoRepository)
	svc := videoService.NewService(mockRepo)

	existingVideo := &video.Video{
		ID:     1,
		UserID: 1,
		Status: video.StatusAuditing,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingVideo, nil)
	mockRepo.On("UpdateStatus", mock.Anything, int64(1), video.StatusPublished).Return(nil)

	err := svc.PublishVideo(context.Background(), 1, 1, false)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRecordView_Success(t *testing.T) {
	mockRepo := new(MockVideoRepository)
	svc := videoService.NewService(mockRepo)

	mockRepo.On("IncrementViewCount", mock.Anything, int64(1)).Return(nil)

	err := svc.RecordView(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetHotVideos_Success(t *testing.T) {
	mockRepo := new(MockVideoRepository)
	svc := videoService.NewService(mockRepo)

	videos := []*video.Video{
		{ID: 1, Title: "Hot Video 1", Status: video.StatusPublished},
		{ID: 2, Title: "Hot Video 2", Status: video.StatusPublished},
	}

	mockRepo.On("GetHotVideos", mock.Anything, 10).Return(videos, nil)

	result, err := svc.GetHotVideos(context.Background(), 10)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestSearchVideos_Success(t *testing.T) {
	mockRepo := new(MockVideoRepository)
	svc := videoService.NewService(mockRepo)

	videos := []*video.Video{
		{ID: 1, Title: "Test Video", Status: video.StatusPublished},
	}

	mockRepo.On("Search", mock.Anything, "test", 1, 20).Return(videos, int64(1), nil)

	resp, err := svc.SearchVideos(context.Background(), "test", 1, 20)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Videos, 1)
	assert.Equal(t, int64(1), resp.Total)
	mockRepo.AssertExpectations(t)
}

func TestGetUserVideos_Success(t *testing.T) {
	mockRepo := new(MockVideoRepository)
	svc := videoService.NewService(mockRepo)

	videos := []*video.Video{
		{ID: 1, Title: "User Video 1", Status: video.StatusPublished, UserID: 1},
		{ID: 2, Title: "User Video 2", Status: video.StatusPublished, UserID: 1},
	}

	mockRepo.On("GetByUserID", mock.Anything, int64(1), 1, 20).Return(videos, int64(2), nil)

	resp, err := svc.GetUserVideos(context.Background(), 1, 1, 20)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Videos, 2)
	assert.Equal(t, int64(2), resp.Total)
	mockRepo.AssertExpectations(t)
}

func TestUpdateVideo_Success(t *testing.T) {
	mockRepo := new(MockVideoRepository)
	svc := videoService.NewService(mockRepo)

	existingVideo := &video.Video{
		ID:          1,
		UserID:      1,
		Title:       "Old Title",
		Description: "Old Description",
		Status:      video.StatusUploading,
	}

	newTitle := "New Title"
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingVideo, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	req := &video.VideoUpdateRequest{
		Title: &newTitle,
	}

	v, err := svc.UpdateVideo(context.Background(), 1, 1, req, false)

	assert.NoError(t, err)
	assert.NotNil(t, v)
	mockRepo.AssertExpectations(t)
}
