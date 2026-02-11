package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/rh-mithu/rizon/backend/config"
	"github.com/rh-mithu/rizon/backend/internal/delivery/middleware"
	"log/slog"
)

type Handler struct {
	l *slog.Logger
}

func ProvideHandler(c *config.Config) chi.Router {
	// Middleware
	authMiddleware := middleware.AuthMiddleware(c.JWTSecret)
	return NewRouter(authMiddleware)
}
