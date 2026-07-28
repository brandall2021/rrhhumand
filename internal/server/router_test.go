package server

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRouterSetup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	routes := router.Routes()
	if len(routes) == 0 {
		t.Error("expected routes to be registered")
	}
}
