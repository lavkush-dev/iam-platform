package services

import (
	"context"

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

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) error {
	hash, _ := utils.HashPassword(req.Password)

	user := &models.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		PasswordHash: hash,
	}

	return s.repo.Create(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}
