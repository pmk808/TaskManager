// RequestControllers/CommandRequest/CommandApiController.go
package CommandRequest

import (
	"net/http"
	"taskmanager/Services/CommandServices/ImportTaskService/interfaces"
	"taskmanager/Services/CommandServices/ImportTaskService/schemas"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type commandApiController struct {
	importService interfaces.ImportService
	logger        *logrus.Logger
}

func NewCommandApiController(
	importService interfaces.ImportService,
	logger *logrus.Logger,
) *commandApiController {
	return &commandApiController{
		importService: importService,
		logger:        logger,
	}
}

func (c *commandApiController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/import", c.ImportTasks)
}

// ImportTasks godoc
// @Summary Import tasks from CSV
// @Description Triggers the task import process from a specified directory
// @Tags commands
// @Accept json
// @Produce json
// @Success 200 {object} schemas.TaskImportResponse "Successful import response"
// @Failure 500 {object} schemas.TaskImportResponse "Error import response"
// @Router /api/commands/import [post]
func (c *commandApiController) ImportTasks(ctx *gin.Context) {
	c.logger.Info("Received request to import tasks")

	response, err := c.importService.Import()
	if err != nil {
		c.logger.WithError(err).Error("Failed to import tasks")
		if response != nil {
			// Use the structured response even in error case
			ctx.JSON(http.StatusInternalServerError, response)
		} else {
			// Fallback if no response was received
			ctx.JSON(http.StatusInternalServerError, schemas.TaskImportResponse{
				Success: false,
				Message: "Failed to import tasks",
				Errors:  []string{err.Error()},
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, response)
}
