package postgres

import (
	"context"
	"database/sql"

	"iam-platform/internal/models"
	"iam-platform/internal/repository/interfaces"
)

type PermissionRepository struct {
	db *sql.DB
}

// compile-time check
var _ interfaces.PermissionRepository = (*PermissionRepository)(nil)

func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{db}
}

func (r *PermissionRepository) Create(ctx context.Context, p *models.Permission) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO permissions (id, name) VALUES ($1, $2)",
		p.ID, p.Name,
	)
	return err
}

func (r *PermissionRepository) GetByID(ctx context.Context, id string) (*models.Permission, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT id, name FROM permissions WHERE id=$1",
		id,
	)

	var p models.Permission
	err := row.Scan(&p.ID, &p.Name)
	return &p, err
}

func (r *PermissionRepository) GetByName(ctx context.Context, name string) (*models.Permission, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT id, name FROM permissions WHERE name=$1",
		name,
	)

	var p models.Permission
	err := row.Scan(&p.ID, &p.Name)
	return &p, err
}
