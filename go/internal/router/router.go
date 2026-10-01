package router

import (
	"net/http"

	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/apperror"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/handler"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/middleware"
	"github.com/rsfzxx/test-teknologi-operator-prima/go/internal/response"
)

const apiV1 = "/api/v1"

type Dependencies struct {
	Health *handler.HealthHandler
	Bank   *handler.BankHandler
}

func New(d Dependencies) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", d.Health.Check)
	mux.HandleFunc("/health", methodNotAllowed(http.MethodGet))

	// Bank
	mux.HandleFunc("GET "+apiV1+"/banks", d.Bank.List)
	mux.HandleFunc("GET "+apiV1+"/banks/{id}", d.Bank.GetByID)
	mux.HandleFunc("PATCH "+apiV1+"/banks/{id}/status", d.Bank.UpdateStatus)
	mux.HandleFunc(apiV1+"/banks", methodNotAllowed(http.MethodGet))
	mux.HandleFunc(apiV1+"/banks/{id}", methodNotAllowed(http.MethodGet))
	mux.HandleFunc(apiV1+"/banks/{id}/status", methodNotAllowed(http.MethodPatch))

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		response.Error(w, apperror.NotFound("Endpoint not found"))
	})

	return middleware.Chain(mux,
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
	)
}

func methodNotAllowed(allow string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Allow", allow)
		response.Error(w, apperror.MethodNotAllowed())
	}
}