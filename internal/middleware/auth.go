package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DanilWaliev/wishlist-api/internal/services"
)

// ключ для context
type contextKey string

const userIDKey contextKey = "user_id"

// middleware
type AuthMiddleware struct {
	authService *services.AuthService
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// основной метод проверки доступа
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeError(w, http.StatusUnauthorized, "authorization header is required")
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			writeError(w, http.StatusUnauthorized, "invalid authorization format")
			return
		}

		token := strings.TrimPrefix(authHeader, prefix)
		token = strings.TrimSpace(token)

		if token == "" {
			writeError(w, http.StatusUnauthorized, "empty token")
			return
		}

		user, err := m.authService.AuthorizeByToken(r.Context(), token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// кладём user_id в context
		ctx := context.WithValue(r.Context(), userIDKey, user.ID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// функция для получения user_id из context
func GetUserID(ctx context.Context) (uint32, bool) {
	id, ok := ctx.Value(userIDKey).(uint32)
	return id, ok
}

// вспомогательный метод для ошибок
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
