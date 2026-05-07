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

	valid, rehashRequired, newHash, err :=
		utils.CheckPassword(
			req.Password,
			user.PasswordHash,
		)

	if err != nil {
		return "", err
	}

	if !valid {
		return "", errors.New("invalid credentials")
	}

	// migrate bcrypt -> argon2id
	if rehashRequired {

		err = s.userRepo.UpdatePasswordHash(
			ctx,
			user.ID,
			newHash,
		)

		if err != nil {
			return "", err
		}
	}

	token, err := s.jwt.Generate(user.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}
