package interfaces

import "github.com/gin-gonic/gin"

type CommandApiController interface {
    RegisterRoutes(router *gin.RouterGroup)
    ImportTasks(c *gin.Context)
}