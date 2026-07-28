package tenant

import (
	"github.com/gin-gonic/gin"
)

func GetCompanyID(c *gin.Context) string {
	val, exists := c.Get("company_id")
	if !exists {
		return ""
	}
	id, ok := val.(string)
	if !ok {
		return ""
	}
	return id
}

func GetUserID(c *gin.Context) string {
	val, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	id, ok := val.(string)
	if !ok {
		return ""
	}
	return id
}

func GetEmail(c *gin.Context) string {
	val, exists := c.Get("email")
	if !exists {
		return ""
	}
	email, ok := val.(string)
	if !ok {
		return ""
	}
	return email
}

func GetRoles(c *gin.Context) []string {
	val, exists := c.Get("roles")
	if !exists {
		return nil
	}
	roles, ok := val.([]string)
	if !ok {
		return nil
	}
	return roles
}
