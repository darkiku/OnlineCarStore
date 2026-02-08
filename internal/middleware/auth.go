package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/teamserik/online-car-store/internal/auth"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(UserContextKey).(*auth.Claims)
		if claims.Role != "admin" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
