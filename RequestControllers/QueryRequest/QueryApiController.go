package QueryRequest

import (
	"net/http"
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
	router.GET("/tasks/active", c.GetActiveTasks)
	router.GET("/tasks/history", c.GetTaskStatusHistory)
}

// GetActiveTasks godoc
// @Summary Get active tasks for a client
// @Description Retrieves all active tasks for a specific client
// @Tags queries
// @Produce json
// @Security Bearer
// @Success 200 {object} interfaces.TasksResponseDTO
// @Failure 400 {object} interfaces.TasksResponseDTO
// @Failure 500 {object} interfaces.TasksResponseDTO
// @Router /api/queries/tasks/active [post]
func (c *queryApiController) GetActiveTasks(ctx *gin.Context) {
	clientName, _ := ctx.Get("client_name")
	clientID, _ := ctx.Get("client_id")

	response, err := c.queryService.GetActiveTasks(
		ctx,
		clientName.(string),
		clientID.(string),
	)
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
// @Summary Get task status history
// @Description Retrieves task status history for a client
// @Tags queries
// @Produce json
// @Security Bearer
// @Success 200 {object} interfaces.StatusHistoryResponseDTO
// @Failure 400 {object} interfaces.StatusHistoryResponseDTO
// @Failure 500 {object} interfaces.StatusHistoryResponseDTO
// @Router /api/queries/tasks/history [post]
func (c *queryApiController) GetTaskStatusHistory(ctx *gin.Context) {
	// Get client info from context (set by JWT middleware)
	clientName, _ := ctx.Get("client_name")
	clientID, _ := ctx.Get("client_id")

	response, err := c.queryService.GetTaskStatusHistory(
		ctx,
		clientName.(string),
		clientID.(string),
	)
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
