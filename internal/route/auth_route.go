package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jhasudungan/terraloom-core-api/internal/handler"
)

func SetupAuthRoutes(accountHandler *handler.AccountHandler, authMiddleware gin.HandlerFunc, router *gin.Engine) *gin.Engine {

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Product routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", accountHandler.Register)
			auth.POST("/login", accountHandler.Login)
		}
	}

	return router
}
