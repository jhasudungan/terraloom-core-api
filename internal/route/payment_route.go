package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jhasudungan/terraloom-core-api/internal/handler"
)

func SetuPaymentRoutes(paymentHandler *handler.PaymentHandler, authMiddleware gin.HandlerFunc, router *gin.Engine) *gin.Engine {

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Product routes
		order := v1.Group("/payment")
		order.Use(authMiddleware)
		{
			order.POST("/submit", paymentHandler.SubmitPayment)
		}

	}

	return router
}
