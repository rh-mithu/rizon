package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	authMiddleware func(http.Handler) http.Handler,
) chi.Router {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(authMiddleware)

	// Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	})
	return r
}
