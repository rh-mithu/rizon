package service

import (
	"context"
	"github.com/rh-mithu/rizon/backend/internal/entity/repository"
)

type AuthService struct {
	repo repository.EmailRepository
}

func NewAuthService(repo repository.EmailRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) SendEmail(ctx context.Context, email string) error {
	return s.repo.SendEmail(ctx, email)
}
