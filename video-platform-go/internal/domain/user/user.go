package user

import (
	"time"
)

type UserRole string

const (
	RoleUser      UserRole = "user"
	RoleCreator   UserRole = "creator"
	RoleModerator UserRole = "moderator"
	RoleAdmin     UserRole = "admin"
)

type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusSuspended UserStatus = "suspended"
	StatusBanned    UserStatus = "banned"
)

type User struct {
	ID             int64      `json:"id" db:"id"`
	Username       string     `json:"username" db:"username"`
	Email          string     `json:"email" db:"email"`
	PasswordHash   string     `json:"-" db:"password_hash"`
	DisplayName    string     `json:"display_name" db:"display_name"`
	AvatarURL      string     `json:"avatar_url" db:"avatar_url"`
	Bio            string     `json:"bio" db:"bio"`
	Role           UserRole   `json:"role" db:"role"`
	Status         UserStatus `json:"status" db:"status"`
	FollowerCount  int64      `json:"follower_count" db:"follower_count"`
	FollowingCount int64      `json:"following_count" db:"following_count"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

type UserFollow struct {
	FollowerID  int64     `json:"follower_id" db:"follower_id"`
	FollowingID int64     `json:"following_id" db:"following_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type RegisterRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,max=72"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	User         *User  `json:"user"`
}

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name,omitempty" validate:"omitempty,min=1,max=100"`
	Bio         *string `json:"bio,omitempty" validate:"omitempty,max=1000"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
}

type UserFilter struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     UserRole `json:"role,omitempty"`
	Status   UserStatus `json:"status,omitempty"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsModerator() bool {
	return u.Role == RoleModerator || u.Role == RoleAdmin
}

func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

func (u *User) CanUploadVideo() bool {
	return u.IsActive() && (u.Role == RoleCreator || u.Role == RoleModerator || u.Role == RoleAdmin)
}
