package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rh-mithu/rizon/backend/internal/apperrors"
	"github.com/rh-mithu/rizon/backend/internal/delivery/http/response"
	"github.com/rh-mithu/rizon/backend/internal/dto"
	"github.com/rh-mithu/rizon/backend/internal/service"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/request-link", h.SendEmail)
		})
	})
}

func (h *AuthHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperrors.NewBadRequest("invalid request body"))
		return
	}
	err := h.authService.SendEmail(r.Context(), req.Email)
	if err != nil {
		response.Error(w, err)
	}
	response.JSON(w, http.StatusOK, dto.AuthResponse{Message: "success"})
}
