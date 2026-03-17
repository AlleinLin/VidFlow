package interaction_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/video-platform/go/internal/domain/interaction"
	interactionService "github.com/video-platform/go/internal/service/interaction"
)

type MockInteractionRepository struct {
	mock.Mock
}

func (m *MockInteractionRepository) CreateComment(ctx context.Context, c *interaction.Comment) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockInteractionRepository) GetCommentByID(ctx context.Context, id int64) (*interaction.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interaction.Comment), args.Error(1)
}

func (m *MockInteractionRepository) GetCommentsByVideoID(ctx context.Context, videoID int64, page, pageSize int) ([]*interaction.Comment, int64, error) {
	args := m.Called(ctx, videoID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*interaction.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *MockInteractionRepository) GetReplies(ctx context.Context, rootID int64, page, pageSize int) ([]*interaction.Comment, int64, error) {
	args := m.Called(ctx, rootID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*interaction.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *MockInteractionRepository) UpdateComment(ctx context.Context, id int64, content string) error {
	args := m.Called(ctx, id, content)
	return args.Error(0)
}

func (m *MockInteractionRepository) DeleteComment(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInteractionRepository) CreateLike(ctx context.Context, userID, videoID int64) error {
	args := m.Called(ctx, userID, videoID)
	return args.Error(0)
}

func (m *MockInteractionRepository) DeleteLike(ctx context.Context, userID, videoID int64) error {
	args := m.Called(ctx, userID, videoID)
	return args.Error(0)
}

func (m *MockInteractionRepository) IsLiked(ctx context.Context, userID, videoID int64) (bool, error) {
	args := m.Called(ctx, userID, videoID)
	return args.Bool(0), args.Error(1)
}

func (m *MockInteractionRepository) GetLikeCount(ctx context.Context, videoID int64) (int64, error) {
	args := m.Called(ctx, videoID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockInteractionRepository) CreateFavorite(ctx context.Context, userID, videoID int64) error {
	args := m.Called(ctx, userID, videoID)
	return args.Error(0)
}

func (m *MockInteractionRepository) DeleteFavorite(ctx context.Context, userID, videoID int64) error {
	args := m.Called(ctx, userID, videoID)
	return args.Error(0)
}

func (m *MockInteractionRepository) IsFavorited(ctx context.Context, userID, videoID int64) (bool, error) {
	args := m.Called(ctx, userID, videoID)
	return args.Bool(0), args.Error(1)
}

func (m *MockInteractionRepository) GetFavorites(ctx context.Context, userID int64, page, pageSize int) ([]int64, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]int64), args.Get(1).(int64), args.Error(2)
}

func (m *MockInteractionRepository) CreateDanmaku(ctx context.Context, d *interaction.Danmaku) error {
	args := m.Called(ctx, d)
	return args.Error(0)
}

func (m *MockInteractionRepository) GetDanmakusByVideoID(ctx context.Context, videoID int64, startTime, endTime float64) ([]*interaction.Danmaku, error) {
	args := m.Called(ctx, videoID, startTime, endTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interaction.Danmaku), args.Error(1)
}

type MockVideoRepository struct {
	mock.Mock
}

func (m *MockVideoRepository) IncrementCommentCount(ctx context.Context, id int64, delta int) error {
	args := m.Called(ctx, id, delta)
	return args.Error(0)
}

func (m *MockVideoRepository) IncrementLikeCount(ctx context.Context, id int64, delta int) error {
	args := m.Called(ctx, id, delta)
	return args.Error(0)
}

func TestCreateComment_Success(t *testing.T) {
	mockRepo := new(MockInteractionRepository)
	mockVideoRepo := new(MockVideoRepository)
	svc := interactionService.NewInteractionService(mockRepo, mockVideoRepo)

	req := &interaction.CreateCommentRequest{
		VideoID: 1,
		Content: "Test comment",
	}

	mockRepo.On("CreateComment", mock.Anything, mock.Anything).Return(nil)
	mockVideoRepo.On("IncrementCommentCount", mock.Anything, int64(1), 1).Return(nil)

	comment, err := svc.CreateComment(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, comment)
	mockRepo.AssertExpectations(t)
}

func TestLikeVideo_Success(t *testing.T) {
	mockRepo := new(MockInteractionRepository)
	mockVideoRepo := new(MockVideoRepository)
	svc := interactionService.NewInteractionService(mockRepo, mockVideoRepo)

	mockRepo.On("IsLiked", mock.Anything, int64(1), int64(1)).Return(false, nil)
	mockRepo.On("CreateLike", mock.Anything, int64(1), int64(1)).Return(nil)
	mockVideoRepo.On("IncrementLikeCount", mock.Anything, int64(1), 1).Return(nil)

	err := svc.LikeVideo(context.Background(), 1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLikeVideo_AlreadyLiked(t *testing.T) {
	mockRepo := new(MockInteractionRepository)
	mockVideoRepo := new(MockVideoRepository)
	svc := interactionService.NewInteractionService(mockRepo, mockVideoRepo)

	mockRepo.On("IsLiked", mock.Anything, int64(1), int64(1)).Return(true, nil)

	err := svc.LikeVideo(context.Background(), 1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetLikeStatus_Success(t *testing.T) {
	mockRepo := new(MockInteractionRepository)
	mockVideoRepo := new(MockVideoRepository)
	svc := interactionService.NewInteractionService(mockRepo, mockVideoRepo)

	mockRepo.On("IsLiked", mock.Anything, int64(1), int64(1)).Return(true, nil)
	mockRepo.On("IsFavorited", mock.Anything, int64(1), int64(1)).Return(false, nil)
	mockRepo.On("GetLikeCount", mock.Anything, int64(1)).Return(int64(100), nil)

	status, err := svc.GetLikeStatus(context.Background(), 1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.True(t, status.IsLiked)
	assert.False(t, status.IsFavorited)
	mockRepo.AssertExpectations(t)
}

func TestUnlikeVideo_Success(t *testing.T) {
	mockRepo := new(MockInteractionRepository)
	mockVideoRepo := new(MockVideoRepository)
	svc := interactionService.NewInteractionService(mockRepo, mockVideoRepo)

	mockRepo.On("IsLiked", mock.Anything, int64(1), int64(1)).Return(true, nil)
	mockRepo.On("DeleteLike", mock.Anything, int64(1), int64(1)).Return(nil)
	mockVideoRepo.On("IncrementLikeCount", mock.Anything, int64(1), -1).Return(nil)

	err := svc.UnlikeVideo(context.Background(), 1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestFavoriteVideo_Success(t *testing.T) {
	mockRepo := new(MockInteractionRepository)
	mockVideoRepo := new(MockVideoRepository)
	svc := interactionService.NewInteractionService(mockRepo, mockVideoRepo)

	mockRepo.On("CreateFavorite", mock.Anything, int64(1), int64(1)).Return(nil)

	err := svc.FavoriteVideo(context.Background(), 1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateDanmaku_Success(t *testing.T) {
	mockRepo := new(MockInteractionRepository)
	mockVideoRepo := new(MockVideoRepository)
	svc := interactionService.NewInteractionService(mockRepo, mockVideoRepo)

	req := &interaction.CreateDanmakuRequest{
		VideoID:         1,
		Content:         "Test danmaku",
		PositionSeconds: 10.5,
		Style:           interaction.DanmakuStyleScroll,
		Color:           "#FFFFFF",
		FontSize:        24,
	}

	mockRepo.On("CreateDanmaku", mock.Anything, mock.Anything).Return(nil)

	danmaku, err := svc.CreateDanmaku(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, danmaku)
	mockRepo.AssertExpectations(t)
}
