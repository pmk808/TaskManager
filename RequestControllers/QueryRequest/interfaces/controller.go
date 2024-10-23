package interfaces

import "github.com/gin-gonic/gin"

type QueryApiController interface {
    RegisterRoutes(router *gin.RouterGroup)
    // Add query methods later for dashboard
}