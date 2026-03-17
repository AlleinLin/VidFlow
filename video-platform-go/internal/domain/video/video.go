package video

import (
	"time"
)

type VideoStatus string

const (
	StatusUploading   VideoStatus = "UPLOADING"
	StatusTranscoding VideoStatus = "TRANSCODING"
	StatusAuditing    VideoStatus = "AUDITING"
	StatusPublished   VideoStatus = "PUBLISHED"
	StatusRejected   VideoStatus = "REJECTED"
	StatusDeleted     VideoStatus = "DELETED"
	StatusHidden      VideoStatus = "HIDDEN"
)

type Visibility string

const (
	VisibilityPublic      Visibility = "public"
	VisibilityFollowers  Visibility = "followers_only"
	VisibilityPrivate    Visibility = "private"
)

type Resolution string

const (
	Resolution240p  Resolution = "240p"
	Resolution480p  Resolution = "480p"
	Resolution720p  Resolution = "720p"
	Resolution1080p Resolution = "1080p"
	Resolution4K    Resolution = "4k"
)

type ResolutionConfig struct {
	Name        string
	Width       int
	Height      int
	BitrateH264 int
	BitrateH265 int
}

var ResolutionConfigs = map[string]ResolutionConfig{
	"240p":  {"240p", 426, 240, 400, 200},
	"480p":  {"480p", 854, 480, 1000, 500},
	"720p":  {"720p", 1280, 720, 2500, 1200},
	"1080p": {"1080p", 1920, 1080, 5000, 2500},
	"4k":    {"4k", 3840, 2160, 15000, 8000},
}

type Video struct {
	ID              int64        `json:"id" db:"id"`
	UserID          int64        `json:"user_id" db:"user_id"`
	Title           string       `json:"title" db:"title"`
	Description     string       `json:"description" db:"description"`
	Status          VideoStatus  `json:"status" db:"status"`
	Visibility      Visibility   `json:"visibility" db:"visibility"`
	DurationSeconds int          `json:"duration_seconds" db:"duration_seconds"`
	OriginalFilename string       `json:"original_filename" db:"original_filename"`
	StorageKey      string       `json:"storage_key" db:"storage_key"`
	ThumbnailURL    string       `json:"thumbnail_url" db:"thumbnail_url"`
	CategoryID      int          `json:"category_id" db:"category_id"`
	ViewCount       int64        `json:"view_count" db:"view_count"`
	LikeCount       int64        `json:"like_count" db:"like_count"`
	CommentCount    int          `json:"comment_count" db:"comment_count"`
	PublishedAt     *time.Time  `json:"published_at,omitempty" db:"published_at"`
	CreatedAt       time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at" db:"updated_at"`
	Tags            []string     `json:"tags,omitempty"`
}

type TranscodeTask struct {
	ID               int64       `json:"id" db:"id"`
	VideoID          int64       `json:"video_id" db:"video_id"`
	Resolution       string      `json:"resolution" db:"resolution"`
	Format           string      `json:"format" db:"format"`
	Status           string      `json:"status" db:"status"`
	OutputStorageKey string      `json:"output_storage_key" db:"output_storage_key"`
	ErrorMessage     string      `json:"error_message,omitempty" db:"error_message"`
	StartedAt        *time.Time  `json:"started_at,omitempty" db:"started_at"`
	CompletedAt      *time.Time  `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
}

type Tag struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	VideoCount int    `json:"video_count" db:"video_count"`
}

type Category struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	ParentID    *int   `json:"parent_id,omitempty" db:"parent_id"`
}

type VideoUploadRequest struct {
	Title       string   `json:"title" validate:"required,min=1,max=200"`
	Description string   `json:"description" validate:"max=5000"`
	Visibility  Visibility `json:"visibility" validate:"required,oneof=public followers_only private"`
	CategoryID  int      `json:"category_id" validate:"required"`
	Tags        []string `json:"tags" validate:"max=10,dive,min=1,max=50"`
	Filename    string   `json:"filename" validate:"required"`
	FileSize    int64    `json:"file_size" validate:"required,min=1"`
}

type VideoUploadResponse struct {
	VideoID       int64  `json:"video_id"`
	UploadURL     string `json:"upload_url"`
	UploadToken   string `json:"upload_token"`
	ExpiresAt     int64  `json:"expires_at"`
}

type VideoUpdateRequest struct {
	Title       *string    `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=5000"`
	Visibility  *Visibility `json:"visibility,omitempty" validate:"omitempty,oneof=public followers_only private"`
	CategoryID  *int       `json:"category_id,omitempty"`
	Tags        []string   `json:"tags,omitempty" validate:"max=10,dive,min=1,max=50"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
}

type VideoFilter struct {
	UserID     int64       `json:"user_id,omitempty"`
	Status     VideoStatus `json:"status,omitempty"`
	Visibility Visibility  `json:"visibility,omitempty"`
	CategoryID int         `json:"category_id,omitempty"`
	Tags       []string    `json:"tags,omitempty"`
	Keyword    string      `json:"keyword,omitempty"`
	SortBy     string      `json:"sort_by,omitempty"`
	SortDesc   bool        `json:"sort_desc,omitempty"`
}

type VideoListResponse struct {
	Videos     []*Video `json:"videos"`
	Total      int64    `json:"total"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalPages int      `json:"total_pages"`
}

func (v *Video) IsPublished() bool {
	return v.Status == StatusPublished
}

func (v *Video) CanView(viewerID int64, isFollowing bool) bool {
	if v.Status != StatusPublished {
		return false
	}
	
	switch v.Visibility {
	case VisibilityPublic:
		return true
	case VisibilityFollowers:
		return v.UserID == viewerID || isFollowing
	case VisibilityPrivate:
		return v.UserID == viewerID
	default:
		return false
	}
}

func (v *Video) CanEdit(userID int64, isAdmin bool) bool {
	return v.UserID == userID || isAdmin
}

func (v *Video) CanDelete(userID int64, isAdmin bool) bool {
	return v.UserID == userID || isAdmin
}

func ValidTransitions() map[VideoStatus][]VideoStatus {
	return map[VideoStatus][]VideoStatus{
		StatusUploading:   {StatusTranscoding, StatusDeleted},
		StatusTranscoding: {StatusAuditing, StatusPublished, StatusRejected},
		StatusAuditing:    {StatusPublished, StatusRejected},
		StatusPublished:   {StatusHidden, StatusDeleted},
		StatusRejected:    {StatusDeleted},
		StatusHidden:      {StatusPublished, StatusDeleted},
	}
}

func (v *Video) CanTransitionTo(newStatus VideoStatus) bool {
	transitions := ValidTransitions()
	allowed, exists := transitions[v.Status]
	if !exists {
		return false
	}
	
	for _, s := range allowed {
		if s == newStatus {
			return true
		}
	}
	return false
}
