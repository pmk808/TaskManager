package httpSetup

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
	logger *logrus.Logger
	port   int
}

func NewServer(logger *logrus.Logger, port int) *Server {
	router := gin.Default()

	// Add swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Server{
		router: router,
		logger: logger,
		port:   port,
	}
}

func (s *Server) Router() *gin.Engine {
	return s.router
}

func (s *Server) Start() error {
	serverAddress := fmt.Sprintf(":%d", s.port)
	s.logger.Infof("Starting server on %s", serverAddress)

	return http.ListenAndServe(serverAddress, s.router)
}
