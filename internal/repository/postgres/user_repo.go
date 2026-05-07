package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"iam-platform/internal/models"
	"iam-platform/internal/repository/interfaces"
)

type UserRepository struct {
	db *sql.DB
}

// compile-time check
var _ interfaces.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO users (id, email, password_hash) VALUES ($1,$2,$3)",
		u.ID, u.Email, u.PasswordHash,
	)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT id, email, password_hash FROM users WHERE email=$1",
		email,
	)

	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash)
	return &u, err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT id, email, created_at FROM users WHERE id=$1",
		id,
	)

	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.CreatedAt)
	return &u, err
}

func (r *UserRepository) UpdatePasswordHash(
	ctx context.Context,
	userID string,
	newHash string,
) error {

	query := `
		UPDATE users
		SET password_hash = $1,
		    updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		newHash,
		userID,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to update password hash: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf(
			"failed to get rows affected: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
