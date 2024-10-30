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
	repoInterfaces "taskmanager/Repository/CommandRepository/interfaces"
	"taskmanager/RequestControllers/CommandRequest"
	"taskmanager/RequestControllers/httpSetup"
	"taskmanager/RequestControllers/httpSetup/config"
	"taskmanager/RequestControllers/httpSetup/logger"
	"taskmanager/Services/CommandServices/ImportTaskService"
	serviceInterfaces "taskmanager/Services/CommandServices/ImportTaskService/interfaces"
	"taskmanager/Services/CommandServices/ImportTaskService/validation"

	_ "taskmanager/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Application constants
const (
	shutdownTimeout = 5 * time.Second
	readTimeout     = 10 * time.Second
	writeTimeout    = 30 * time.Second
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
	startServerWithGracefulShutdown(app, cfg, appLogger)
}

type appDependencies struct {
	router *gin.Engine
}

func initializeApp(cfg *config.Config, logger *logrus.Logger) (*appDependencies, error) {
	logger.Info("Initializing application dependencies")

	// Initialize repositories
	commandRepo, err := initializeRepositories(cfg, logger)
	if err != nil {
		return nil, err
	}

	// Initialize services
	services, err := initializeServices(cfg, logger, commandRepo)
	if err != nil {
		return nil, err
	}

	// Initialize controllers
	router, err := initializeControllers(cfg, logger, services)
	if err != nil {
		return nil, err
	}

	logger.Info("Application dependencies initialized successfully")
	return &appDependencies{
		router: router,
	}, nil
}

func initializeRepositories(cfg *config.Config, logger *logrus.Logger) (repoInterfaces.TaskCommandRepository, error) {
	logger.Info("Initializing repositories")

	commandRepo, err := CommandRepository.NewTaskCommandRepository(&cfg.Database, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize command repository: %w", err)
	}

	return commandRepo, nil
}

func initializeServices(
	cfg *config.Config,
	logger *logrus.Logger,
	commandRepo repoInterfaces.TaskCommandRepository,
) (serviceInterfaces.ImportService, error) {
	logger.Info("Initializing services")

	// Initialize validator
	dataValidator := validation.NewDataValidator(logger)

	// Initialize import service
	importService := ImportTaskService.NewImportService(
		commandRepo,
		dataValidator,
		logger,
		cfg.Import.Directory,
	)

	return importService, nil
}

func initializeControllers(
	cfg *config.Config,
	logger *logrus.Logger,
	importService serviceInterfaces.ImportService,
) (*gin.Engine, error) {
	logger.Info("Initializing controllers")

	// Initialize controller
	commandController := CommandRequest.NewCommandApiController(importService, logger)

	// Setup HTTP router
	routerConfig := httpSetup.RouterConfig{
		CommandController: commandController,
		Logger:            logger,
	}
	router := httpSetup.SetupRouter(routerConfig)

	return router, nil
}

func startServerWithGracefulShutdown(app *appDependencies, cfg *config.Config, logger *logrus.Logger) {
	serverAddress := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.WithField("address", serverAddress).Info("Starting server")

	server := &http.Server{
		Addr:         serverAddress,
		Handler:      app.router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// Server shutdown channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		logger.Info("Server is ready to handle requests")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for shutdown signal
	<-quit
	logger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	logger.Info("Server exited gracefully")
}
