package httpSetup

import (
	"context"
	"net/http"
	"os"
	authInterfaces "taskmanager/RequestControllers/AuthRequest/interfaces"
	cmdControllerInterfaces "taskmanager/RequestControllers/CommandRequest/interfaces"
	queryControllerInterfaces "taskmanager/RequestControllers/QueryRequest/interfaces"
	"taskmanager/RequestControllers/httpSetup/jwt"
	"taskmanager/RequestControllers/httpSetup/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	requestTimeout = 30 * time.Second
)

type RouterConfig struct {
	CommandController cmdControllerInterfaces.CommandApiController
	QueryController   queryControllerInterfaces.QueryApiController
	Logger            *logrus.Logger
	AuthController    authInterfaces.AuthController
	JWTManager        *jwt.JWTManager
}

func InitializeGin(logger *logrus.Logger) {
	env := os.Getenv("GIN_MODE")
	if env == "" {
		env = "debug"
	}

	switch env {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
		logger.Info("Running Gin in debug mode")
	}

	gin.DefaultWriter = logger.Writer()
	gin.DefaultErrorWriter = logger.Writer()
}

func SetupRouter(config RouterConfig) *gin.Engine {
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(requestLoggerMiddleware(config.Logger))

	// API routes
	api := router.Group("/api")

	// Auth routes (no JWT required)
	auth := api.Group("/auth")
	config.AuthController.RegisterRoutes(auth)

	// Query routes (with JWT)
	queries := api.Group("/queries")
	queries.Use(middleware.JWTAuthMiddleware(config.JWTManager))
	config.QueryController.RegisterRoutes(queries)

	// Command routes
	commands := api.Group("/commands")
	config.CommandController.RegisterRoutes(commands)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func requestLoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		logger.WithFields(logrus.Fields{
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"status":    c.Writer.Status(),
			"duration":  duration,
			"client_ip": c.ClientIP(),
		}).Info("Request processed")
	}
}

func timeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			c.Next()
			done <- struct{}{}
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
				"success": false,
				"message": "Request timeout",
			})
			return
		}
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
