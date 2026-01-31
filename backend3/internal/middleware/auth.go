package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/katerinasolntsye/fulleng/internal/service"
)

type contextKey string

const UserIdKey contextKey = "userId"
const UserEmailKey contextKey = "userEmail"

type AuthMiddleware struct {
	jwtService *service.JWTService
}

func NewAuthMiddleware(jwtService *service.JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

// RequireAuth - middleware для проверки JWT токена
func (am *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Извлекаем токен из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Authorization header required"})
			return
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]

		// Валидируем токен
		claims, err := am.jwtService.ValidateToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid or expired token"})
			return
		}

		// Проверяем, что это access токен
		if claims.Type != service.AccessToken {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token type"})
			return
		}

		// Добавляем userId и email в context
		ctx := context.WithValue(r.Context(), UserIdKey, claims.UserId)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

		// Передаем запрос дальше с обновленным context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
