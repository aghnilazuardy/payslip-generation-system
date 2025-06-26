package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"payslip-generation-system/internal/helper"
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
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusUnauthorized, "missing or malformed token", nil, nil))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userID, role, err := utils.ParseToken(tokenStr)
		if err != nil {
			json.NewEncoder(w).Encode(helper.WriteJSONResponse(w, http.StatusUnauthorized, "unauthorized", nil, nil))
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
