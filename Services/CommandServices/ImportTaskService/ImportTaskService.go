package services

import (
	"fmt"
	"taskmanager/Services/CommandServices/ImportTaskService/interfaces"
	"taskmanager/schemas"
	"time"

	"github.com/sirupsen/logrus"
)

type ImportTaskService struct {
	repo      interfaces.TaskRepository
	validator interfaces.Validator
	logger    *logrus.Logger
}

func NewImportTaskService(repo interfaces.TaskRepository, validator interfaces.Validator, logger *logrus.Logger) *ImportTaskService {
	return &ImportTaskService{
		repo:      repo,
		validator: validator,
		logger:    logger,
	}
}

func (s *ImportTaskService) ImportData(tasks []schemas.Task) error {
	start := time.Now()
	s.logger.Info("Starting data import process")

	if err := s.validator.ValidateTasks(tasks); err != nil {
		s.logger.WithError(err).Error("Validation failed")
		return fmt.Errorf("validation failed: %w", err)
	}

	s.logger.Info("All tasks passed validation")

	if err := s.repo.BulkCreateTasks(tasks); err != nil {
		s.logger.WithError(err).Error("Failed to import tasks")
		return fmt.Errorf("failed to import tasks: %w", err)
	}

	duration := time.Since(start)
	s.logger.WithField("duration", duration).Info("Data import completed")
	return nil
}
