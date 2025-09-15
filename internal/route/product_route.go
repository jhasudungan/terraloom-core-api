package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jhasudungan/terraloom-core-api/internal/handler"
)

func SetupProductRoutes(productHandler *handler.ProductHandler, router *gin.Engine) *gin.Engine {

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Product routes
		products := v1.Group("/products")
		{
			products.GET("", productHandler.GetProducts)
		}

		product := v1.Group("/product")
		{
			product.GET("/:id", productHandler.GetProductDetail)
		}
	}

	return router
}
