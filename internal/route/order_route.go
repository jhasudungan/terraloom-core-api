package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jhasudungan/terraloom-core-api/internal/handler"
)

func SetupOrderRoutes(orderHandler *handler.OrderHandler, authMiddleware gin.HandlerFunc, router *gin.Engine) *gin.Engine {

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Product routes
		order := v1.Group("/order")
		order.Use(authMiddleware)
		{
			order.POST("/submit", orderHandler.SubmitOrder)
			order.POST("/cancel", orderHandler.CancelOrder)
			order.GET("/detail/:orderReference", orderHandler.GetOrderDetail)
		}

	}

	return router
}
