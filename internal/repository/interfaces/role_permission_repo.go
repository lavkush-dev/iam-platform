package interfaces

import "context"

type RolePermissionRepository interface {
	AssignPermission(ctx context.Context, roleID, permissionID string) error
	GetPermissionsByRoleID(ctx context.Context, roleID string) ([]string, error)
}
