package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/video-platform/go/internal/domain/user"
	apperrors "github.com/video-platform/go/pkg/errors"
	"github.com/video-platform/go/pkg/hash"
	"github.com/video-platform/go/pkg/jwt"
)

type Repository interface {
	Create(ctx context.Context, u *user.User) error
	GetByID(ctx context.Context, id int64) (*user.User, error)
	GetByUsername(ctx context.Context, username string) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	Update(ctx context.Context, u *user.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter *user.UserFilter, page, pageSize int) ([]*user.User, int64, error)
	
	CreateFollow(ctx context.Context, followerID, followingID int64) error
	DeleteFollow(ctx context.Context, followerID, followingID int64) error
	GetFollowers(ctx context.Context, userID int64, page, pageSize int) ([]*user.User, int64, error)
	GetFollowing(ctx context.Context, userID int64, page, pageSize int) ([]*user.User, int64, error)
	IsFollowing(ctx context.Context, followerID, followingID int64) (bool, error)
	
	UpdateLastLogin(ctx context.Context, userID int64) error
	UpdateFollowerCount(ctx context.Context, userID int64, delta int) error
	UpdateFollowingCount(ctx context.Context, userID int64, delta int) error
}

type Service interface {
	Register(ctx context.Context, req *user.RegisterRequest) (*user.User, error)
	Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*user.LoginResponse, error)
	GetProfile(ctx context.Context, userID int64) (*user.User, error)
	UpdateProfile(ctx context.Context, userID int64, req *user.UpdateProfileRequest) (*user.User, error)
	DeleteUser(ctx context.Context, userID int64) error
	GetUserByID(ctx context.Context, id int64) (*user.User, error)
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
	ListUsers(ctx context.Context, filter *user.UserFilter, page, pageSize int) ([]*user.User, int64, error)
	
	Follow(ctx context.Context, followerID, followingID int64) error
	Unfollow(ctx context.Context, followerID, followingID int64) error
	GetFollowers(ctx context.Context, userID int64, page, pageSize int) ([]*user.User, int64, error)
	GetFollowing(ctx context.Context, userID int64, page, pageSize int) ([]*user.User, int64, error)
	IsFollowing(ctx context.Context, followerID, followingID int64) (bool, error)
}

type service struct {
	repo       Repository
	jwtManager *jwt.JWTManager
}

func NewService(repo Repository, jwtManager *jwt.JWTManager) Service {
	return &service{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *service) Register(ctx context.Context, req *user.RegisterRequest) (*user.User, error) {
	existingUser, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, apperrors.ErrUserNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, apperrors.ErrUserAlreadyExists
	}
	
	existingEmail, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, apperrors.ErrUserNotFound) {
		return nil, err
	}
	if existingEmail != nil {
		return nil, apperrors.ErrUserAlreadyExists
	}
	
	passwordHash, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	
	u := &user.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		DisplayName:  req.DisplayName,
		Role:         user.RoleUser,
		Status:       user.StatusActive,
	}
	
	if err := s.repo.Create(ctx, u); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, apperrors.ErrUserAlreadyExists
		}
		return nil, err
	}
	
	return u, nil
}

func (s *service) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	u, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, err
	}
	
	if !u.IsActive() {
		return nil, apperrors.ErrForbidden
	}
	
	if !hash.CheckPassword(req.Password, u.PasswordHash) {
		return nil, apperrors.ErrInvalidCredentials
	}
	
	tokenPair, err := s.jwtManager.GenerateTokenPair(ctx, u.ID, u.Username, string(u.Role))
	if err != nil {
		return nil, err
	}
	
	go func() {
		_ = s.repo.UpdateLastLogin(context.Background(), u.ID)
	}()
	
	return &user.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    tokenPair.TokenType,
		User:         u,
	}, nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*user.LoginResponse, error) {
	claims, err := s.jwtManager.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	
	u, err := s.repo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	
	if !u.IsActive() {
		return nil, apperrors.ErrForbidden
	}
	
	tokenPair, err := s.jwtManager.GenerateTokenPair(ctx, u.ID, u.Username, string(u.Role))
	if err != nil {
		return nil, err
	}
	
	return &user.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    tokenPair.TokenType,
		User:         u,
	}, nil
}

func (s *service) GetProfile(ctx context.Context, userID int64) (*user.User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *service) UpdateProfile(ctx context.Context, userID int64, req *user.UpdateProfileRequest) (*user.User, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	if req.DisplayName != nil {
		u.DisplayName = *req.DisplayName
	}
	if req.Bio != nil {
		u.Bio = *req.Bio
	}
	if req.AvatarURL != nil {
		u.AvatarURL = *req.AvatarURL
	}
	
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}
	
	return u, nil
}

func (s *service) DeleteUser(ctx context.Context, userID int64) error {
	return s.repo.Delete(ctx, userID)
}

func (s *service) GetUserByID(ctx context.Context, id int64) (*user.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *service) ListUsers(ctx context.Context, filter *user.UserFilter, page, pageSize int) ([]*user.User, int64, error) {
	return s.repo.List(ctx, filter, page, pageSize)
}

func (s *service) Follow(ctx context.Context, followerID, followingID int64) error {
	if followerID == followingID {
		return apperrors.ErrBadRequest
	}
	
	isFollowing, err := s.repo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if isFollowing {
		return nil
	}
	
	if err := s.repo.CreateFollow(ctx, followerID, followingID); err != nil {
		return err
	}
	
	go func() {
		_ = s.repo.UpdateFollowerCount(context.Background(), followingID, 1)
		_ = s.repo.UpdateFollowingCount(context.Background(), followerID, 1)
	}()
	
	return nil
}

func (s *service) Unfollow(ctx context.Context, followerID, followingID int64) error {
	isFollowing, err := s.repo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if !isFollowing {
		return nil
	}
	
	if err := s.repo.DeleteFollow(ctx, followerID, followingID); err != nil {
		return err
	}
	
	go func() {
		_ = s.repo.UpdateFollowerCount(context.Background(), followingID, -1)
		_ = s.repo.UpdateFollowingCount(context.Background(), followerID, -1)
	}()
	
	return nil
}

func (s *service) GetFollowers(ctx context.Context, userID int64, page, pageSize int) ([]*user.User, int64, error) {
	return s.repo.GetFollowers(ctx, userID, page, pageSize)
}

func (s *service) GetFollowing(ctx context.Context, userID int64, page, pageSize int) ([]*user.User, int64, error) {
	return s.repo.GetFollowing(ctx, userID, page, pageSize)
}

func (s *service) IsFollowing(ctx context.Context, followerID, followingID int64) (bool, error) {
	return s.repo.IsFollowing(ctx, followerID, followingID)
}
