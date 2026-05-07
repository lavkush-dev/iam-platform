package interfaces

import (
	"context"

	"iam-platform/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	UpdatePasswordHash(ctx context.Context, userID string, newHash string) error
}
