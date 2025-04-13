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

const UserContextKey contextKey = "user"

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

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Отсутствует токен авторизации", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
			return
		}

		jwtService := jwt.NewJwt(m.config.Auth.Secret)
		valid, jwtData := jwtService.Parce(tokenParts[1])
		if !valid || jwtData == nil {
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		currentUser, err := m.userRepository.FindByPhone(jwtData.Phone)
		if err != nil {
			http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, currentUser)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
