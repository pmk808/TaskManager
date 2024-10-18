package main

import (
	"fmt" // Import fmt for string formatting
	"log"
	"net/http"
	"os"

	"taskmanager/config"
	"taskmanager/handlers"
	"taskmanager/repository"
	"taskmanager/services"
	"taskmanager/validation"

	_ "taskmanager/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	logger.Info("Application starting")

	// Initialize repository
	repo, err := repository.NewPostgresRepository(&cfg.Database, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize repository")
	}

	// Initialize validation
	validator := validation.NewValidator(logger)

	// Initialize service
	service := services.NewTaskService(repo, validator, logger, cfg.Import.Directory)

	// Initialize handler
	handler := handlers.NewTaskHandler(service, logger)

	// Initialize Gin router
	router := gin.Default()

	// Register routes
	handler.RegisterRoutes(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	serverAddress := fmt.Sprintf(":%d", cfg.Server.Port) // Create the server address
	logger.Infof("Starting server on %s", serverAddress)
	if err := http.ListenAndServe(serverAddress, router); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
	logger.Info("Application initialized successfully")
}
