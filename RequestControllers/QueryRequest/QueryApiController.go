package QueryRequest

import (
	"net/http"
	"taskmanager/RequestControllers/QueryRequest/dto"
	controllerInterfaces "taskmanager/RequestControllers/QueryRequest/interfaces"
	serviceInterfaces "taskmanager/Services/QueryServices/TaskQueryService/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type queryApiController struct {
	queryService serviceInterfaces.TaskQueryService
	logger       *logrus.Logger
}

func NewQueryApiController(
	queryService serviceInterfaces.TaskQueryService,
	logger *logrus.Logger,
) controllerInterfaces.QueryApiController {
	return &queryApiController{
		queryService: queryService,
		logger:       logger,
	}
}

func (c *queryApiController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/tasks/active", c.GetActiveTasks)
	router.POST("/tasks/history", c.GetTaskStatusHistory)
}

// GetActiveTasks godoc
// @Summary Get active tasks for a client
// @Description Retrieves all active tasks for a specific client
// @Tags queries
// @Accept json
// @Produce json
// @Param request body dto.ClientQueryRequest true "Client details"
// @Success 200 {object} interfaces.TasksResponseDTO
// @Failure 400 {object} interfaces.TasksResponseDTO
// @Failure 500 {object} interfaces.TasksResponseDTO
// @Router /api/queries/tasks/active [post]
func (c *queryApiController) GetActiveTasks(ctx *gin.Context) {
	var request dto.ClientQueryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.logger.WithError(err).Error("Invalid request body for GetActiveTasks")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"errors":  []string{err.Error()},
		})
		return
	}

	response, err := c.queryService.GetActiveTasks(ctx, request.ClientName, request.ClientID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to get active tasks")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve active tasks",
			"errors":  []string{err.Error()},
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetTaskStatusHistory godoc
// @Summary Get task status history for a client
// @Description Retrieves the status history of all tasks for a specific client
// @Tags queries
// @Accept json
// @Produce json
// @Param request body dto.ClientQueryRequest true "Client details"
// @Success 200 {object} interfaces.StatusHistoryResponseDTO
// @Failure 400 {object} interfaces.StatusHistoryResponseDTO
// @Failure 500 {object} interfaces.StatusHistoryResponseDTO
// @Router /api/queries/tasks/history [post]
func (c *queryApiController) GetTaskStatusHistory(ctx *gin.Context) {
	var request dto.ClientQueryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.logger.WithError(err).Error("Invalid request body for GetTaskStatusHistory")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"errors":  []string{err.Error()},
		})
		return
	}

	response, err := c.queryService.GetTaskStatusHistory(ctx, request.ClientName, request.ClientID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to get task status history")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve task status history",
			"errors":  []string{err.Error()},
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
