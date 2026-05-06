package postgres

import (
	"context"
	"database/sql"
	"iam-platform/internal/repository/interfaces"
)

type UserRoleRepository struct {
	db *sql.DB
}

// compile-time check
var _ interfaces.UserRoleRepository = (*UserRoleRepository)(nil)

func NewUserRoleRepository(db *sql.DB) *UserRoleRepository {
	return &UserRoleRepository{db}
}

func (r *UserRoleRepository) AssignRole(ctx context.Context, userID, roleID string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		userID, roleID,
	)
	return err
}

func (r *UserRoleRepository) GetRolesByUserID(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT role_id FROM user_roles WHERE user_id=$1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string

	for rows.Next() {
		var roleID string
		if err := rows.Scan(&roleID); err != nil {
			return nil, err
		}
		roles = append(roles, roleID)
	}

	return roles, nil
}
