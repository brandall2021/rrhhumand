package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/auth"
	"github.com/rrhhumand/api/pkg/response"
)

type PermissionChecker interface {
	HasPermission(ctx context.Context, userID, companyID, resource, action string) (bool, error)
}

func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "Invalid authorization format")
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("company_id", claims.CompanyID)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID, exists := c.Get("company_id")
		if !exists || companyID == "" {
			response.Forbidden(c, "Company context required")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RBACMiddleware(checker PermissionChecker, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, _ := c.Get("roles")
		roleList, ok := roles.([]string)
		if !ok {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		for _, role := range roleList {
			if role == "SUPER_ADMIN" || role == "COMPANY_ADMIN" {
				c.Next()
				return
			}
		}

		userID, _ := c.Get("user_id")
		companyID, _ := c.Get("company_id")

		userIDStr, _ := userID.(string)
		companyIDStr, _ := companyID.(string)

		if checker != nil && userIDStr != "" && companyIDStr != "" {
			hasPermission, err := checker.HasPermission(c.Request.Context(), userIDStr, companyIDStr, resource, action)
			if err == nil && hasPermission {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Insufficient permissions")
		c.Abort()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
