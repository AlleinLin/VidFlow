package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/video-platform/go/internal/domain/user"
	apperrors "github.com/video-platform/go/pkg/errors"
)

type UserRepository interface {
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

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, display_name, role, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	
	now := time.Now()
	err := r.pool.QueryRow(ctx, query,
		u.Username, u.Email, u.PasswordHash, u.DisplayName, u.Role, u.Status, now, now,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*user.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar_url, bio, role, status,
			   follower_count, following_count, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`
	
	var u user.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.DisplayName, &u.AvatarURL, &u.Bio,
		&u.Role, &u.Status, &u.FollowerCount, &u.FollowingCount,
		&u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	
	return &u, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar_url, bio, role, status,
			   follower_count, following_count, created_at, updated_at, last_login_at
		FROM users
		WHERE username = $1
	`
	
	var u user.User
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.DisplayName, &u.AvatarURL, &u.Bio,
		&u.Role, &u.Status, &u.FollowerCount, &u.FollowingCount,
		&u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	
	return &u, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar_url, bio, role, status,
			   follower_count, following_count, created_at, updated_at, last_login_at
		FROM users
		WHERE email = $1
	`
	
	var u user.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.DisplayName, &u.AvatarURL, &u.Bio,
		&u.Role, &u.Status, &u.FollowerCount, &u.FollowingCount,
		&u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return &u, nil
}

func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users
		SET display_name = $2, avatar_url = $3, bio = $4, updated_at = $5
		WHERE id = $1
	`
	
	result, err := r.pool.Exec(ctx, query, u.ID, u.DisplayName, u.AvatarURL, u.Bio, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return apperrors.ErrUserNotFound
	}
	
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return apperrors.ErrUserNotFound
	}
	
	return nil
}

func (r *userRepository) List(ctx context.Context, filter *user.UserFilter, page, pageSize int) ([]*user.User, int64, error) {
	whereClause := "WHERE 1=1"
	args := make([]interface{}, 0)
	argNum := 1
	
	if filter.Username != "" {
		whereClause += fmt.Sprintf(" AND username LIKE $%d", argNum)
		args = append(args, filter.Username+"%")
		argNum++
	}
	
	if filter.Email != "" {
		whereClause += fmt.Sprintf(" AND email = $%d", argNum)
		args = append(args, filter.Email)
		argNum++
	}
	
	if filter.Role != "" {
		whereClause += fmt.Sprintf(" AND role = $%d", argNum)
		args = append(args, filter.Role)
		argNum++
	}
	
	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, filter.Status)
		argNum++
	}
	
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}
	
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, username, email, password_hash, display_name, avatar_url, bio, role, status,
			   follower_count, following_count, created_at, updated_at, last_login_at
		FROM users
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)
	
	args = append(args, pageSize, offset)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()
	
	var users []*user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.DisplayName, &u.AvatarURL, &u.Bio,
			&u.Role, &u.Status, &u.FollowerCount, &u.FollowingCount,
			&u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &u)
	}
	
	return users, total, nil
}

func (r *userRepository) CreateFollow(ctx context.Context, followerID, followingID int64) error {
	query := `
		INSERT INTO user_follows (follower_id, following_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (follower_id, following_id) DO NOTHING
	`
	
	_, err := r.pool.Exec(ctx, query, followerID, followingID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create follow: %w", err)
	}
	
	return nil
}

func (r *userRepository) DeleteFollow(ctx context.Context, followerID, followingID int64) error {
	query := `DELETE FROM user_follows WHERE follower_id = $1 AND following_id = $2`
	
	_, err := r.pool.Exec(ctx, query, followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to delete follow: %w", err)
	}
	
	return nil
}

func (r *userRepository) GetFollowers(ctx context.Context, userID int64, page, pageSize int) ([]*user.User, int64, error) {
	countQuery := `SELECT COUNT(*) FROM user_follows WHERE following_id = $1`
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count followers: %w", err)
	}
	
	offset := (page - 1) * pageSize
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.display_name, u.avatar_url, u.bio, u.role, u.status,
			   u.follower_count, u.following_count, u.created_at, u.updated_at, u.last_login_at
		FROM users u
		INNER JOIN user_follows uf ON u.id = uf.follower_id
		WHERE uf.following_id = $1
		ORDER BY uf.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.pool.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get followers: %w", err)
	}
	defer rows.Close()
	
	var users []*user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.DisplayName, &u.AvatarURL, &u.Bio,
			&u.Role, &u.Status, &u.FollowerCount, &u.FollowingCount,
			&u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan follower: %w", err)
		}
		users = append(users, &u)
	}
	
	return users, total, nil
}

func (r *userRepository) GetFollowing(ctx context.Context, userID int64, page, pageSize int) ([]*user.User, int64, error) {
	countQuery := `SELECT COUNT(*) FROM user_follows WHERE follower_id = $1`
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count following: %w", err)
	}
	
	offset := (page - 1) * pageSize
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.display_name, u.avatar_url, u.bio, u.role, u.status,
			   u.follower_count, u.following_count, u.created_at, u.updated_at, u.last_login_at
		FROM users u
		INNER JOIN user_follows uf ON u.id = uf.following_id
		WHERE uf.follower_id = $1
		ORDER BY uf.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.pool.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get following: %w", err)
	}
	defer rows.Close()
	
	var users []*user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.DisplayName, &u.AvatarURL, &u.Bio,
			&u.Role, &u.Status, &u.FollowerCount, &u.FollowingCount,
			&u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan following: %w", err)
		}
		users = append(users, &u)
	}
	
	return users, total, nil
}

func (r *userRepository) IsFollowing(ctx context.Context, followerID, followingID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_follows WHERE follower_id = $1 AND following_id = $2)`
	
	var exists bool
	err := r.pool.QueryRow(ctx, query, followerID, followingID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check is following: %w", err)
	}
	
	return exists, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID int64) error {
	query := `UPDATE users SET last_login_at = $2 WHERE id = $1`
	
	_, err := r.pool.Exec(ctx, query, userID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	
	return nil
}

func (r *userRepository) UpdateFollowerCount(ctx context.Context, userID int64, delta int) error {
	query := `UPDATE users SET follower_count = follower_count + $2 WHERE id = $1`
	
	_, err := r.pool.Exec(ctx, query, userID, delta)
	if err != nil {
		return fmt.Errorf("failed to update follower count: %w", err)
	}
	
	return nil
}

func (r *userRepository) UpdateFollowingCount(ctx context.Context, userID int64, delta int) error {
	query := `UPDATE users SET following_count = following_count + $2 WHERE id = $1`
	
	_, err := r.pool.Exec(ctx, query, userID, delta)
	if err != nil {
		return fmt.Errorf("failed to update following count: %w", err)
	}
	
	return nil
}
