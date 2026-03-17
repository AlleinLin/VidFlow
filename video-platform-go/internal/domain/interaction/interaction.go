package interaction

import (
	"time"
)

type CommentStatus string

const (
	CommentStatusVisible CommentStatus = "visible"
	CommentStatusHidden  CommentStatus = "hidden"
	CommentStatusDeleted CommentStatus = "deleted"
)

type Comment struct {
	ID        int64         `json:"id" db:"id"`
	VideoID   int64         `json:"video_id" db:"video_id"`
	UserID    int64         `json:"user_id" db:"user_id"`
	ParentID  *int64        `json:"parent_id,omitempty" db:"parent_id"`
	RootID    *int64        `json:"root_id,omitempty" db:"root_id"`
	Content   string        `json:"content" db:"content"`
	LikeCount int           `json:"like_count" db:"like_count"`
	Status    CommentStatus `json:"status" db:"status"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	User      *CommentUser  `json:"user,omitempty"`
}

type CommentUser struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type DanmakuStyle string

const (
	DanmakuStyleScroll DanmakuStyle = "scroll"
	DanmakuStyleTop    DanmakuStyle = "top"
	DanmakuStyleBottom DanmakuStyle = "bottom"
)

type Danmaku struct {
	ID              int64        `json:"id" db:"id"`
	VideoID         int64        `json:"video_id" db:"video_id"`
	UserID          int64        `json:"user_id" db:"user_id"`
	Content         string       `json:"content" db:"content"`
	PositionSeconds float64      `json:"position_seconds" db:"position_seconds"`
	Style           DanmakuStyle `json:"style" db:"style"`
	Color           string       `json:"color" db:"color"`
	FontSize        int          `json:"font_size" db:"font_size"`
	CreatedAt       time.Time    `json:"created_at" db:"created_at"`
}

type Like struct {
	UserID    int64     `json:"user_id" db:"user_id"`
	VideoID   int64     `json:"video_id" db:"video_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Favorite struct {
	UserID    int64     `json:"user_id" db:"user_id"`
	VideoID   int64     `json:"video_id" db:"video_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateCommentRequest struct {
	VideoID  int64  `json:"video_id" validate:"required"`
	ParentID *int64 `json:"parent_id,omitempty"`
	Content  string `json:"content" validate:"required,min=1,max=5000"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=5000"`
}

type CreateDanmakuRequest struct {
	VideoID         int64        `json:"video_id" validate:"required"`
	Content         string       `json:"content" validate:"required,min=1,max=100"`
	PositionSeconds float64      `json:"position_seconds" validate:"required,min=0"`
	Style           DanmakuStyle `json:"style" validate:"required,oneof=scroll top bottom"`
	Color           string       `json:"color" validate:"required,hexcolor"`
	FontSize        int          `json:"font_size" validate:"required,min=12,max=36"`
}

type CommentListResponse struct {
	Comments   []*Comment `json:"comments"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
}

type DanmakuListResponse struct {
	Danmakus []*Danmaku `json:"danmakus"`
}

type LikeStatus struct {
	IsLiked     bool `json:"is_liked"`
	IsFavorited bool `json:"is_favorited"`
	LikeCount   int64 `json:"like_count"`
}

func (c *Comment) IsRoot() bool {
	return c.ParentID == nil
}

func (c *Comment) IsReply() bool {
	return c.ParentID != nil
}

func (c *Comment) CanEdit(userID int64, isAdmin bool) bool {
	return c.UserID == userID || isAdmin
}

func (c *Comment) CanDelete(userID int64, isAdmin bool) bool {
	return c.UserID == userID || isAdmin
}
