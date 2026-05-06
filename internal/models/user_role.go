package models

type UserRole struct {
	UserID string `db:"user_id"`
	RoleID string `db:"role_id"`
}
