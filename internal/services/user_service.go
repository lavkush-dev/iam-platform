package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"iam-platform/internal/dto"
	"iam-platform/internal/models"
	"iam-platform/internal/repository/interfaces"
	"iam-platform/internal/utils"
)

type UserService struct {
	repo interfaces.UserRepository
}

func NewUserService(r interfaces.UserRepository) *UserService {
	return &UserService{r}
}

func (s *UserService) CreateUser(
	ctx context.Context,
	req dto.CreateUserRequest,
) error {

	// 1. basic validation (cheap, early fail)
	if req.Email == "" || req.Password == "" {
		return errors.New("email and password are required")
	}

	// 2. hash password (must handle error)
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. generate ID (no panic)
	id, err := uuid.NewV7() // Use uuid v7 instead of uuid v4 by default
	if err != nil {
		return fmt.Errorf("failed to generate user id: %w", err)
	}

	user := &models.User{
		ID:           id.String(),
		Email:        req.Email,
		PasswordHash: hash,
	}

	// 4. persist user with context
	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}
