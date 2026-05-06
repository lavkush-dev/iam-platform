package postgres

import (
	"context"
	"database/sql"

	"iam-platform/internal/models"
	"iam-platform/internal/repository/interfaces"
)

type RoleRepository struct {
	db *sql.DB
}

// compile-time check
var _ interfaces.RoleRepository = (*RoleRepository)(nil)

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db}
}

func (r *RoleRepository) Create(ctx context.Context, role *models.Role) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO roles (id, name) VALUES ($1, $2)",
		role.ID, role.Name,
	)
	return err
}

func (r *RoleRepository) GetByID(ctx context.Context, id string) (*models.Role, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT id, name FROM roles WHERE id=$1",
		id,
	)

	var role models.Role
	err := row.Scan(&role.ID, &role.Name)
	return &role, err
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT id, name FROM roles WHERE name=$1",
		name,
	)

	var role models.Role
	err := row.Scan(&role.ID, &role.Name)
	return &role, err
}
