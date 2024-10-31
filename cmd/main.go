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
	cmdRepoInterfaces "taskmanager/Repository/CommandRepository/interfaces"
	"taskmanager/Repository/QueryRepository"
	queryRepoInterfaces "taskmanager/Repository/QueryRepository/interfaces"
	"taskmanager/RequestControllers/AuthRequest"
	authInterfaces "taskmanager/RequestControllers/AuthRequest/interfaces"
	"taskmanager/RequestControllers/CommandRequest"
	"taskmanager/RequestControllers/QueryRequest"
	"taskmanager/RequestControllers/httpSetup"
	"taskmanager/RequestControllers/httpSetup/config"
	"taskmanager/RequestControllers/httpSetup/jwt"
	"taskmanager/RequestControllers/httpSetup/logger"
	"taskmanager/Services/CommandServices/ImportTaskService"
	commandServiceInterfaces "taskmanager/Services/CommandServices/ImportTaskService/interfaces"
	"taskmanager/Services/CommandServices/ImportTaskService/validation"
	"taskmanager/Services/QueryServices/TaskQueryService"
	queryServiceInterfaces "taskmanager/Services/QueryServices/TaskQueryService/interfaces"

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

	// Initialize JWT Manager
	jwtManager := jwt.NewJWTManager(cfg.JWT.SecretKey, cfg.JWT.ExpiryHours)

	// Initialize repositories
	commandRepo, queryRepo, err := initializeRepositories(cfg, logger)
	if err != nil {
		return nil, err
	}

	// Initialize services
	commandService, queryService, err := initializeServices(cfg, logger, commandRepo, queryRepo)
	if err != nil {
		return nil, err
	}

	// Initialize auth controller
	authController := AuthRequest.NewAuthController(jwtManager, logger)

	// Initialize controllers and router
	router, err := initializeControllers(cfg, logger, commandService, queryService, authController, jwtManager)
	if err != nil {
		return nil, err
	}

	logger.Info("Application dependencies initialized successfully")
	return &appDependencies{
		router: router,
	}, nil
}

func initializeRepositories(cfg *config.Config, logger *logrus.Logger) (
	cmdRepo cmdRepoInterfaces.TaskCommandRepository,
	queryRepo queryRepoInterfaces.TaskQueryRepository,
	err error,
) {
	logger.Info("Initializing repositories")

	cmdRepo, err = CommandRepository.NewTaskCommandRepository(&cfg.Database, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize command repository: %w", err)
	}

	queryRepo, err = QueryRepository.NewTaskQueryRepository(&cfg.Database, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize query repository: %w", err)
	}

	return cmdRepo, queryRepo, nil
}

func initializeServices(
	cfg *config.Config,
	logger *logrus.Logger,
	commandRepo cmdRepoInterfaces.TaskCommandRepository,
	queryRepo queryRepoInterfaces.TaskQueryRepository,
) (
	commandService commandServiceInterfaces.ImportService,
	queryService queryServiceInterfaces.TaskQueryService,
	err error,
) {
	logger.Info("Initializing services")

	// Initialize validator
	dataValidator := validation.NewDataValidator(logger)

	// Initialize services
	commandService = ImportTaskService.NewImportService(
		commandRepo,
		dataValidator,
		logger,
		cfg.Import.Directory,
	)

	queryService = TaskQueryService.NewTaskQueryService(
		queryRepo,
		logger,
	)

	return commandService, queryService, nil
}

func initializeControllers(
	cfg *config.Config,
	logger *logrus.Logger,
	commandService commandServiceInterfaces.ImportService,
	queryService queryServiceInterfaces.TaskQueryService,
	authController authInterfaces.AuthController,
	jwtManager *jwt.JWTManager,
) (*gin.Engine, error) {
	logger.Info("Initializing controllers")

	// Initialize controllers
	commandController := CommandRequest.NewCommandApiController(commandService, logger)
	queryController := QueryRequest.NewQueryApiController(queryService, logger)

	// Setup HTTP router
	routerConfig := httpSetup.RouterConfig{
		CommandController: commandController,
		QueryController:   queryController,
		AuthController:    authController,
		JWTManager:        jwtManager,
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
