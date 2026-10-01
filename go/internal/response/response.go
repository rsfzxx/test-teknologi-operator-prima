package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/apperror"
)

const requestIDHeader = "X-Request-ID"

type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type successBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type errorBody struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Errors  []apperror.FieldError `json:"errors,omitempty"`
}

func JSON(w http.ResponseWriter, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		slog.Error("failed to marshal response", "error", err, "request_id", w.Header().Get(requestIDHeader))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"success":false,"message":"Internal server error"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func Success(w http.ResponseWriter, status int, message string, data any) {
	JSON(w, status, successBody{Success: true, Message: message, Data: data})
}

func Paginated(w http.ResponseWriter, message string, data any, meta Meta) {
	JSON(w, http.StatusOK, successBody{Success: true, Message: message, Data: data, Meta: &meta})
}

func Error(w http.ResponseWriter, err error) {
	var appErr *apperror.Error
	if !errors.As(err, &appErr) {
		appErr = apperror.Internal(err)
	}

	if appErr.Status >= http.StatusInternalServerError {
		slog.Error("request failed",
			"error", appErr.Error(),
			"request_id", w.Header().Get(requestIDHeader),
		)
	}

	JSON(w, appErr.Status, errorBody{
		Success: false,
		Message: appErr.Message,
		Errors:  appErr.Details,
	})
}

func NewMeta(page, perPage int, total int64) Meta {
	totalPages := 0
	if perPage > 0 {
		totalPages = int((total + int64(perPage) - 1) / int64(perPage))
	}
	return Meta{Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}
}