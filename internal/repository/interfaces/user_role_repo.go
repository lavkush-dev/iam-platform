package interfaces

import "context"

type UserRoleRepository interface {
	AssignRole(ctx context.Context, userID, roleID string) error
	GetRolesByUserID(ctx context.Context, userID string) ([]string, error)
}
