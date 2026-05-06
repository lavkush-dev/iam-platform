package interfaces

import (
	"context"

	"iam-platform/internal/models"
)

type RoleRepository interface {
	Create(ctx context.Context, role *models.Role) error
	GetByID(ctx context.Context, id string) (*models.Role, error)
	GetByName(ctx context.Context, name string) (*models.Role, error)
}
