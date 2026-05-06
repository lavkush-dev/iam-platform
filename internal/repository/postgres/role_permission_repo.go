package postgres

import (
	"context"
	"database/sql"
	"iam-platform/internal/repository/interfaces"
)

type RolePermissionRepository struct {
	db *sql.DB
}

// compile-time check
var _ interfaces.RolePermissionRepository = (*RolePermissionRepository)(nil)

func NewRolePermissionRepository(db *sql.DB) *RolePermissionRepository {
	return &RolePermissionRepository{db}
}

func (r *RolePermissionRepository) AssignPermission(ctx context.Context, roleID, permissionID string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		roleID, permissionID,
	)
	return err
}

func (r *RolePermissionRepository) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT permission_id FROM role_permissions WHERE role_id=$1",
		roleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string

	for rows.Next() {
		var pid string
		if err := rows.Scan(&pid); err != nil {
			return nil, err
		}
		permissions = append(permissions, pid)
	}

	return permissions, nil
}
