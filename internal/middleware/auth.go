package middleware

import (
	"net/http"
	"payslip-generation-system/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or malformed token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userID, role, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Store in context
		c.Set("user_id", userID)
		c.Set("role", role)
		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	if val, exists := c.Get("user_id"); exists {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func GetUserRole(c *gin.Context) string {
	if val, exists := c.Get("role"); exists {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}
