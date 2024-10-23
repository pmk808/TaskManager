// cmd/main.go
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

	"taskmanager/Repository/CommandRepository"
	"taskmanager/RequestControllers/CommandRequest"
	"taskmanager/RequestControllers/httpSetup"
	"taskmanager/RequestControllers/httpSetup/config"
	"taskmanager/RequestControllers/httpSetup/logger"
	"taskmanager/Services/CommandServices/ImportTaskService"
	"taskmanager/Services/CommandServices/ImportTaskService/validation"

	_ "taskmanager/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize configuration
	cfg, err := config.InitializeConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	appLogger := logger.InitializeLogger()
	appLogger.Info("Application starting")

	// Initialize business logic dependencies
	app, err := initializeApp(cfg, appLogger)
	if err != nil {
		appLogger.WithError(err).Fatal("Failed to initialize application")
	}

	// Start server and handle graceful shutdown
	startServerWithGracefulShutdown(app, cfg.Server.Port, appLogger)
}

type appDependencies struct {
	router *gin.Engine
}

func initializeApp(cfg *config.Config, logger *logrus.Logger) (*appDependencies, error) {
	// Initialize repository
	commandRepo, err := CommandRepository.NewTaskCommandRepository(&cfg.Database, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize command repository: %w", err)
	}

	// Initialize validator
	dataValidator := validation.NewDataValidator(logger)

	// Initialize service
	importService := ImportTaskService.NewImportService(
		commandRepo,
		dataValidator,
		logger,
		cfg.Import.Directory,
	)

	// Initialize controller
	commandController := CommandRequest.NewCommandApiController(importService, logger)

	// Setup HTTP router
	routerConfig := httpSetup.RouterConfig{
		CommandController: commandController,
		Logger:            logger,
	}
	router := httpSetup.SetupRouter(routerConfig)

	return &appDependencies{
		router: router,
	}, nil
}

func startServerWithGracefulShutdown(app *appDependencies, port int, logger *logrus.Logger) {
	serverAddress := fmt.Sprintf(":%d", port)
	logger.Infof("Starting server on %s", serverAddress)

	server := &http.Server{
		Addr:    serverAddress,
		Handler: app.router,
	}

	// Server shutdown channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	logger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exiting")
}
