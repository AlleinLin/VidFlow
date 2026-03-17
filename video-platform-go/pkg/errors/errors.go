package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

var (
	ErrBadRequest = &AppError{
		Code:       "BAD_REQUEST",
		Message:    "Invalid request",
		StatusCode: http.StatusBadRequest,
	}
	
	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "Authentication required",
		StatusCode: http.StatusUnauthorized,
	}
	
	ErrForbidden = &AppError{
		Code:       "FORBIDDEN",
		Message:    "Access denied",
		StatusCode: http.StatusForbidden,
	}
	
	ErrNotFound = &AppError{
		Code:       "NOT_FOUND",
		Message:    "Resource not found",
		StatusCode: http.StatusNotFound,
	}
	
	ErrConflict = &AppError{
		Code:       "CONFLICT",
		Message:    "Resource conflict",
		StatusCode: http.StatusConflict,
	}
	
	ErrInternal = &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    "Internal server error",
		StatusCode: http.StatusInternalServerError,
	}
	
	ErrServiceUnavailable = &AppError{
		Code:       "SERVICE_UNAVAILABLE",
		Message:    "Service temporarily unavailable",
		StatusCode: http.StatusServiceUnavailable,
	}
)

var (
	ErrUserNotFound       = NewAppError("USER_NOT_FOUND", "User not found", http.StatusNotFound, nil)
	ErrUserAlreadyExists  = NewAppError("USER_ALREADY_EXISTS", "User already exists", http.StatusConflict, nil)
	ErrInvalidCredentials = NewAppError("INVALID_CREDENTIALS", "Invalid username or password", http.StatusUnauthorized, nil)
	ErrInvalidToken       = NewAppError("INVALID_TOKEN", "Invalid or expired token", http.StatusUnauthorized, nil)
	ErrTokenExpired       = NewAppError("TOKEN_EXPIRED", "Token has expired", http.StatusUnauthorized, nil)
	
	ErrVideoNotFound      = NewAppError("VIDEO_NOT_FOUND", "Video not found", http.StatusNotFound, nil)
	ErrVideoNotPublished  = NewAppError("VIDEO_NOT_PUBLISHED", "Video is not published", http.StatusForbidden, nil)
	ErrVideoAlreadyExists = NewAppError("VIDEO_ALREADY_EXISTS", "Video already exists", http.StatusConflict, nil)
	ErrInvalidVideoStatus = NewAppError("INVALID_VIDEO_STATUS", "Invalid video status for this operation", http.StatusBadRequest, nil)
	
	ErrCommentNotFound    = NewAppError("COMMENT_NOT_FOUND", "Comment not found", http.StatusNotFound, nil)
	ErrCommentTooLong     = NewAppError("COMMENT_TOO_LONG", "Comment exceeds maximum length", http.StatusBadRequest, nil)
	
	ErrUploadFailed       = NewAppError("UPLOAD_FAILED", "Failed to upload file", http.StatusInternalServerError, nil)
	ErrTranscodeFailed    = NewAppError("TRANSCODE_FAILED", "Failed to transcode video", http.StatusInternalServerError, nil)
	ErrInvalidFileType    = NewAppError("INVALID_FILE_TYPE", "Invalid file type", http.StatusBadRequest, nil)
	ErrFileTooLarge       = NewAppError("FILE_TOO_LARGE", "File size exceeds limit", http.StatusBadRequest, nil)
	
	ErrRateLimitExceeded  = NewAppError("RATE_LIMIT_EXCEEDED", "Rate limit exceeded", http.StatusTooManyRequests, nil)
)

func Wrap(err error, appErr *AppError) *AppError {
	return &AppError{
		Code:       appErr.Code,
		Message:    appErr.Message,
		StatusCode: appErr.StatusCode,
		Err:        err,
	}
}

func NewBadRequestError(message string, err error) *AppError {
	if message == "" {
		message = "Bad request"
	}
	return &AppError{
		Code:       "BAD_REQUEST",
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Err:        err,
	}
}

func NewNotFoundError(message string, err error) *AppError {
	if message == "" {
		message = "Not found"
	}
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    message,
		StatusCode: http.StatusNotFound,
		Err:        err,
	}
}

func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return ErrInternal
}
