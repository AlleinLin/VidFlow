package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/video-platform/go/internal/domain/user"
	userService "github.com/video-platform/go/internal/service/user"
	apperrors "github.com/video-platform/go/pkg/errors"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, filter *user.UserFilter, page, pageSize int) ([]*user.User, int64, error) {
	args := m.Called(ctx, filter, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*user.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) CreateFollow(ctx context.Context, followerID, followingID int64) error {
	args := m.Called(ctx, followerID, followingID)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteFollow(ctx context.Context, followerID, followingID int64) error {
	args := m.Called(ctx, followerID, followingID)
	return args.Error(0)
}

func (m *MockUserRepository) GetFollowers(ctx context.Context, userID int64, page, pageSize int) ([]*user.User, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*user.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) GetFollowing(ctx context.Context, userID int64, page, pageSize int) ([]*user.User, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*user.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) IsFollowing(ctx context.Context, followerID, followingID int64) (bool, error) {
	args := m.Called(ctx, followerID, followingID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateFollowerCount(ctx context.Context, userID int64, delta int) error {
	args := m.Called(ctx, userID, delta)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateFollowingCount(ctx context.Context, userID int64, delta int) error {
	args := m.Called(ctx, userID, delta)
	return args.Error(0)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := userService.NewService(mockRepo, nil)

	req := &user.RegisterRequest{
		Username:    "testuser",
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "Test User",
	}

	mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(nil, apperrors.ErrUserNotFound)
	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, apperrors.ErrUserNotFound)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	u, err := svc.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, "testuser", u.Username)
	mockRepo.AssertExpectations(t)
}

func TestRegister_DuplicateUsername(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := userService.NewService(mockRepo, nil)

	req := &user.RegisterRequest{
		Username:    "existinguser",
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "Test User",
	}

	existingUser := &user.User{ID: 1, Username: "existinguser"}
	mockRepo.On("GetByUsername", mock.Anything, "existinguser").Return(existingUser, nil)

	u, err := svc.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, u)
	assert.True(t, errors.Is(err, apperrors.ErrUserAlreadyExists))
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := userService.NewService(mockRepo, nil)

	expectedUser := &user.User{
		ID:          1,
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedUser, nil)

	u, err := svc.GetUserByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, int64(1), u.ID)
	assert.Equal(t, "testuser", u.Username)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := userService.NewService(mockRepo, nil)

	mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, apperrors.ErrUserNotFound)

	u, err := svc.GetUserByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, u)
	assert.True(t, errors.Is(err, apperrors.ErrUserNotFound))
	mockRepo.AssertExpectations(t)
}

func TestUpdateProfile_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := userService.NewService(mockRepo, nil)

	existingUser := &user.User{
		ID:          1,
		Username:    "testuser",
		DisplayName: "Old Name",
	}

	newName := "New Name"
	newBio := "New bio"
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingUser, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	req := &user.UpdateProfileRequest{
		DisplayName: &newName,
		Bio:         &newBio,
	}

	u, err := svc.UpdateProfile(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	mockRepo.AssertExpectations(t)
}

func TestFollow_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := userService.NewService(mockRepo, nil)

	mockRepo.On("IsFollowing", mock.Anything, int64(1), int64(2)).Return(false, nil)
	mockRepo.On("CreateFollow", mock.Anything, int64(1), int64(2)).Return(nil)

	err := svc.Follow(context.Background(), 1, 2)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUnfollow_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := userService.NewService(mockRepo, nil)

	mockRepo.On("IsFollowing", mock.Anything, int64(1), int64(2)).Return(true, nil)
	mockRepo.On("DeleteFollow", mock.Anything, int64(1), int64(2)).Return(nil)

	err := svc.Unfollow(context.Background(), 1, 2)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestFollow_SelfFollow(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := userService.NewService(mockRepo, nil)

	err := svc.Follow(context.Background(), 1, 1)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, apperrors.ErrBadRequest))
}
