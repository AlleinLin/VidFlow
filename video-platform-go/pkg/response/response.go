package response

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/video-platform/go/pkg/errors"
)

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
	TraceID   string      `json:"trace_id,omitempty"`
}

type PagedData struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

func Success(w http.ResponseWriter, data interface{}) {
	resp := Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
	writeJSON(w, http.StatusOK, resp)
}

func SuccessWithStatus(w http.ResponseWriter, statusCode int, data interface{}) {
	resp := Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
	writeJSON(w, statusCode, resp)
}

func Created(w http.ResponseWriter, data interface{}) {
	resp := Response{
		Code:      0,
		Message:   "created",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
	writeJSON(w, http.StatusCreated, resp)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Paged(w http.ResponseWriter, items interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	
	data := PagedData{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
	
	Success(w, data)
}

func Error(w http.ResponseWriter, err error) {
	appErr := errors.GetAppError(err)
	
	resp := Response{
		Code:      appErr.StatusCode,
		Message:   appErr.Message,
		Timestamp: time.Now().Unix(),
	}
	
	writeJSON(w, appErr.StatusCode, resp)
}

func ErrorWithMessage(w http.ResponseWriter, statusCode int, message string) {
	resp := Response{
		Code:      statusCode,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
	writeJSON(w, statusCode, resp)
}

func ValidationError(w http.ResponseWriter, message string) {
	resp := Response{
		Code:      http.StatusBadRequest,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
	writeJSON(w, http.StatusBadRequest, resp)
}

func Unauthorized(w http.ResponseWriter, message string) {
	if message == "" {
		message = "unauthorized"
	}
	resp := Response{
		Code:      http.StatusUnauthorized,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
	writeJSON(w, http.StatusUnauthorized, resp)
}

func Forbidden(w http.ResponseWriter, message string) {
	if message == "" {
		message = "forbidden"
	}
	resp := Response{
		Code:      http.StatusForbidden,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
	writeJSON(w, http.StatusForbidden, resp)
}

func NotFound(w http.ResponseWriter, message string) {
	if message == "" {
		message = "not found"
	}
	resp := Response{
		Code:      http.StatusNotFound,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
	writeJSON(w, http.StatusNotFound, resp)
}

func BadRequest(w http.ResponseWriter, message string) {
	if message == "" {
		message = "bad request"
	}
	resp := Response{
		Code:      http.StatusBadRequest,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
	writeJSON(w, http.StatusBadRequest, resp)
}

func InternalError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "internal server error"
	}
	resp := Response{
		Code:      http.StatusInternalServerError,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
	writeJSON(w, http.StatusInternalServerError, resp)
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
