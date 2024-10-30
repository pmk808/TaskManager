package interfaces

import "github.com/gin-gonic/gin"

type QueryApiController interface {
    RegisterRoutes(router *gin.RouterGroup)
    GetActiveTasks(c *gin.Context)
    GetTaskStatusHistory(c *gin.Context)
}