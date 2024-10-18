package handlers

import (
	"net/http"
	"taskmanager/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TaskHandler struct {
	service interfaces.TaskService
	logger  *logrus.Logger
}

func NewTaskHandler(service interfaces.TaskService, logger *logrus.Logger) *TaskHandler {
	return &TaskHandler{
		service: service,
		logger:  logger,
	}
}

func (h *TaskHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/import", h.ImportData)
	// Register other routes
}

// ImportData godoc
// @Summary Import data from spreadsheet
// @Description Triggers the data import process from a specified directory
// @Tags import
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /import [post]
func (h *TaskHandler) ImportData(c *gin.Context) {
	h.logger.Info("Received request to import data")

	err := h.service.ImportData()
	if err != nil {
		h.logger.WithError(err).Error("Failed to import data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data import completed successfully"})
}
