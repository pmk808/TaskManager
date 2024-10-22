package commandrequest

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CommandApiController struct {
	importTaskService ImportTaskService
	logger           *logrus.Logger
}

// ImportTaskService interface defined within the same package
type ImportTaskService interface {
	ImportData() error
}

func NewCommandApiController(importTaskService ImportTaskService, logger *logrus.Logger) *CommandApiController {
	return &CommandApiController{
		importTaskService: importTaskService,
		logger:           logger,
	}
}

func (c *CommandApiController) RegisterRoutes(router *gin.Engine) {
	commandGroup := router.Group("/api/v1/commands")
	{
		commandGroup.POST("/import-tasks", c.importTasks)
	}
}

// @Summary Import tasks from spreadsheet
// @Description Triggers the task import process from a specified directory
// @Tags commands
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/commands/import-tasks [post]
func (c *CommandApiController) importTasks(ctx *gin.Context) {
	c.logger.Info("Received request to import tasks")

	if err := c.importTaskService.ImportData(); err != nil {
		c.logger.WithError(err).Error("Failed to import tasks")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import tasks"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tasks import completed successfully"})
}