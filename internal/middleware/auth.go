package middleware

import (
	"net/http"
	"strings"

	"iam-platform/internal/utils"
	"iam-platform/pkg/jwt"
)

type AuthMiddleware struct {
	jwt *jwt.Manager
}

func NewAuthMiddleware(jwt *jwt.Manager) *AuthMiddleware {
	return &AuthMiddleware{jwt}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" {
			next.ServeHTTP(w, r)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		claims, err := m.jwt.Validate(token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := utils.SetUserID(r.Context(), claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
