package httpSetup

import (
	"taskmanager/RequestControllers/CommandRequest/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterConfig struct {
	CommandController interfaces.CommandApiController
	Logger            *logrus.Logger
}

// InitializeGin sets up Gin with proper mode and middleware
func InitializeGin(logger *logrus.Logger) {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	// Setup Gin logging and recovery middleware
	gin.DefaultWriter = logger.Writer()
	gin.DefaultErrorWriter = logger.Writer()
}

// SetupSwagger initializes swagger documentation
func SetupSwagger(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// SetupMiddleware adds common middleware to router
func SetupMiddleware(router *gin.Engine, logger *logrus.Logger) {
	// Add logging middleware
	router.Use(gin.LoggerWithWriter(logger.Writer()))
	router.Use(gin.Recovery())

	// Add any other common middleware here
	router.Use(corsMiddleware())
}

// corsMiddleware handles CORS setup
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func SetupRouter(config RouterConfig) *gin.Engine {
	// Initialize Gin
	router := gin.Default()

	// Setup API groups
	apiGroup := router.Group("/api")
	commandsGroup := apiGroup.Group("/commands")

	// Register command routes
	config.CommandController.RegisterRoutes(commandsGroup)

	// Setup swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
