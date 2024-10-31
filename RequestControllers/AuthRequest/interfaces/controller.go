package interfaces

import "github.com/gin-gonic/gin"

type AuthController interface {
	RegisterRoutes(router *gin.RouterGroup)
	GenerateToken(c *gin.Context)
}
