package middleware

import (
	"4-order-api/config"
	"4-order-api/internal/user"
	"4-order-api/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserContextKey   contextKey = "user"
	UserIDContextKey contextKey = "user_id"
)

type AuthMiddleware struct {
	config         *config.Config
	userRepository *user.Repository
}

func NewAuthMiddleware(config *config.Config, userRepository *user.Repository) *AuthMiddleware {
	return &AuthMiddleware{
		config:         config,
		userRepository: userRepository,
	}
}

func GetUserFromContext(ctx context.Context) *user.User {
	user, ok := ctx.Value(UserContextKey).(*user.User)
	if !ok {
		return nil
	}
	return user
}

func GetUserIDFromContext(ctx context.Context) uint {
	userID, ok := ctx.Value(UserIDContextKey).(uint)
	if !ok {
		return 0
	}
	return userID
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		jwtService := jwt.NewJwt(m.config.Auth.Secret)
		userID, err := jwtService.ParseToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
