package interfaces

import (
	"context"

	"iam-platform/internal/models"
)

type PermissionRepository interface {
	Create(ctx context.Context, p *models.Permission) error
	GetByID(ctx context.Context, id string) (*models.Permission, error)
	GetByName(ctx context.Context, name string) (*models.Permission, error)
}
