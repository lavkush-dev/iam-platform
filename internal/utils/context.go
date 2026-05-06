package utils

import "context"

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RolesKey  contextKey = "roles"
)

// ---- USER ID ----

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func GetUserID(ctx context.Context) string {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return ""
	}
	return val.(string)
}
