package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	authMiddleware func(http.Handler) http.Handler,
	authHandler *AuthHandler,
) chi.Router {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)

	// Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	})
	authHandler.RegisterRoutes(r)
	_ = chi.Walk(r, func(method string, route string, handler http.Handler,
		middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	})
	return r
}
