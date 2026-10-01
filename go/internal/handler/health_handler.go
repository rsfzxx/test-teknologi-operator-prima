package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/response"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	db Pinger
}

func NewHealthHandler(db Pinger) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		response.JSON(w, http.StatusServiceUnavailable, map[string]any{
			"success": false,
			"message": "Service unavailable",
			"data":    map[string]string{"status": "unhealthy", "database": "down"},
		})
		return
	}

	response.Success(w, http.StatusOK, "Service is healthy", map[string]string{
		"status":   "healthy",
		"database": "up",
	})
}