package router

import (
	"net/http"

	"github.com/faizan1191/auth-service/internal/auth"
	"github.com/faizan1191/auth-service/internal/middleware"
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
		authGroup.POST("/refresh", handler.Refresh)
		authGroup.POST("/logout", handler.Logout)
	}

	// protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"user_id": c.GetString("user_id"),
				"email":   c.GetString("email"),
			})
		})
	}

	return r
}
