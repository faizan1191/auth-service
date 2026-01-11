package router

import (
	"net/http"

	"github.com/faizan1191/auth-service/internal/auth"
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *auth.Handler) *gin.Engine {
	r := gin.Default() // Create Gin router

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// group all routes so all start with auth/
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", handler.Signup)
		authGroup.POST("/login", handler.Login)
	}

	return r
}
