package middleware

import (
	"net/http"
	"strings"
	
	"github.com/yeahmerey/go-auth-service/internal/services"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		
		//check if token is in the blacklist
		if services.IsTokenBlacklisted(token) {
			http.Error(w, "Token has been revoked", http.StatusUnauthorized)
			return
		}
		
		// standard token validation
		if _, err := services.ValidateToken(token); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		next(w, r)
	}
}