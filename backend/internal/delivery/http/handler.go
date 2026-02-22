package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/rh-mithu/rizon/backend/config"
	"github.com/rh-mithu/rizon/backend/internal/delivery/middleware"
	"github.com/rh-mithu/rizon/backend/internal/infrastructure/repository/mail"
	"github.com/rh-mithu/rizon/backend/internal/service"
	"log/slog"
)

type Handler struct {
	l *slog.Logger
}

func ProvideHandler(c *config.Config) chi.Router {
	// Middleware
	authMiddleware := middleware.AuthMiddleware(c.JWTSecret)
	emailRepo := mail.ProvideEmailRepository(c)
	authService := service.NewAuthService(emailRepo)
	authHandler := NewAuthHandler(authService)
	return NewRouter(authMiddleware, authHandler)
}
