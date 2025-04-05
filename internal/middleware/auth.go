package middleware

import (
	"net/http"
	"strings"
	
	"github.com/yeahmerey/go-auth-service/internal/services"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		
		// Проверка, не находится ли токен в черном списке
		if services.IsTokenBlacklisted(token) {
			http.Error(w, "Token has been revoked", http.StatusUnauthorized)
			return
		}
		
		// Стандартная проверка токена
		if _, err := services.ValidateToken(token); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		next(w, r)
	}
}