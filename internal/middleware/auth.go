package middleware

import (
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		token := strings.TrimSpace(parts[1])

		// 🔐 TODO: Validate the token here (JWT or custom logic)
		if token != "your-secret-token" { // ← replace with actual validation
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// ✅ Token is valid
		next.ServeHTTP(w, r)
	})
}
