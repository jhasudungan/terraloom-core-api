package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/handler"
	"github.com/jhasudungan/terraloom-core-api/internal/middlewares"
	"github.com/jhasudungan/terraloom-core-api/internal/repository"
	"github.com/jhasudungan/terraloom-core-api/internal/route"
	"github.com/jhasudungan/terraloom-core-api/internal/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	// Service port
	srvPort := os.Getenv("PORT")

	// Mode
	env := os.Getenv("ENV")

	// JWT Secret
	jwtSecret := os.Getenv("JWT_SECRET")

	// Prepare DB
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	idGenerator := common.NewIDGenerator()

	// Initialize repository
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	orderItemRepo := repository.NewOrderItemRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	accountRepo := repository.NewAccountRepository(db)

	// Initalize service
	jwtService := service.NewJwtService(jwtSecret)

	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(
		orderRepo,
		productRepo,
		orderItemRepo,
		paymentRepo,
		accountRepo,
		idGenerator)
	accountService := service.NewAccountService(jwtService, accountRepo)
	paymentService := service.NewPaymentService(orderRepo, paymentRepo)

	// Initalize handler
	errorHandler := handler.NewErrorHandler()
	productHandler := handler.NewProductHandler(productService, errorHandler)
	orderHandler := handler.NewOrderHandler(orderService, errorHandler)
	accountHandler := handler.NewAccountHandler(accountService, orderService, errorHandler)
	paymentHandler := handler.NewPaymentHandler(paymentService, errorHandler)

	// Initialize middleware
	authMiddleware := middlewares.NewAuthMiddleware(jwtService, errorHandler)

	// Setup routes
	router := gin.New()

	router = route.SetupProductRoutes(productHandler, router)
	router = route.SetupOrderRoutes(orderHandler, authMiddleware, router)
	router = route.SetupAuthRoutes(accountHandler, authMiddleware, router)
	router = route.SetupAccountRoutes(accountHandler, orderHandler, authMiddleware, router)
	router = route.SetuPaymentRoutes(paymentHandler, authMiddleware, router)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + srvPort,
		Handler: router}

	if env == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Start in go routine
	go func() {

		logrus.WithField("port", srvPort).Info("Starting HTTP server")

		err = srv.ListenAndServe()

		if err != nil {
			logrus.WithError(err).Fatal("Failed to start server")
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// Create a deadline for the shutdown (Gracefull Shutdown)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatal("Server forced to shutdown")
	}

	logrus.Info("Server shutdown complete")
}
