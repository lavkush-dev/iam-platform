package services

import (
	"context"
	"errors"

	"iam-platform/internal/dto"
	"iam-platform/internal/repository/interfaces"
	"iam-platform/internal/utils"
	"iam-platform/pkg/jwt"
)

type AuthService struct {
	userRepo interfaces.UserRepository
	jwt      *jwt.Manager
}

func NewAuthService(repo interfaces.UserRepository, jwt *jwt.Manager) *AuthService {
	return &AuthService{repo, jwt}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	token, err := s.jwt.Generate(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
