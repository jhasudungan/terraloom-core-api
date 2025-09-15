package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jhasudungan/terraloom-core-api/internal/handler"
)

func SetupAccountRoutes(accountHandler *handler.AccountHandler, orderHandler *handler.OrderHandler, authMiddleware gin.HandlerFunc, router *gin.Engine) *gin.Engine {

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Product routes
		account := v1.Group("/account")
		account.Use(authMiddleware)
		{
			account.GET("/detail", accountHandler.GetAccountDetail)
			account.GET("/orders", accountHandler.GetAccountOrders)
			account.PUT("/update", accountHandler.UpdateAccount)
			account.PUT("/update/password", accountHandler.UpdatePassword)
		}
	}

	return router
}
