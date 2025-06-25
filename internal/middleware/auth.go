package middleware

import (
	"context"
	"net/http"
	"payslip-generation-system/utils"
	"strings"
)

type key string

const (
	UserIDKey key = "user_id"
	RoleKey   key = "role"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or malformed token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userID, role, err := utils.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Store in context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) string {
	val := r.Context().Value(UserIDKey)
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

func GetUserRole(r *http.Request) string {
	val := r.Context().Value(RoleKey)
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}
