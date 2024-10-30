package TaskQueryService

import (
	"context"
	"fmt"
	repoInterfaces "taskmanager/Repository/QueryRepository/interfaces"
	serviceInterfaces "taskmanager/Services/QueryServices/interfaces"
	"taskmanager/Services/QueryServices/validation"

	"github.com/sirupsen/logrus"
)

type taskQueryService struct {
	repo      repoInterfaces.TaskQueryRepository
	validator *validation.QueryValidator
	logger    *logrus.Logger
}

// NewTaskQueryService creates a new instance of TaskQueryService
func NewTaskQueryService(
	repo repoInterfaces.TaskQueryRepository,
	logger *logrus.Logger,
) serviceInterfaces.TaskQueryService {
	return &taskQueryService{
		repo:      repo,
		validator: validation.NewQueryValidator(),
		logger:    logger,
	}
}

func (s *taskQueryService) GetActiveTasks(
	ctx context.Context,
	clientName string,
	clientID string,
) (*serviceInterfaces.TasksResponseDTO, error) {
	// Validate input parameters
	if err := s.validator.ValidateClientParams(clientName, clientID); err != nil {
		return &serviceInterfaces.TasksResponseDTO{
			Success: false,
			Message: fmt.Sprintf("Invalid parameters: %v", err),
		}, nil
	}

	// Get tasks from repository
	tasks, err := s.repo.GetActiveTasks(ctx, clientName, clientID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get active tasks")
		return &serviceInterfaces.TasksResponseDTO{
			Success: false,
			Message: "Failed to retrieve active tasks",
		}, err
	}

	// Map repository data to DTOs
	var taskDTOs []serviceInterfaces.TaskDetailDTO
	for _, task := range tasks {
		taskDTOs = append(taskDTOs, serviceInterfaces.TaskDetailDTO{
			ID:         task.ID,
			Name:       task.Name,
			Email:      task.Email,
			Department: task.Department,
			Position:   task.Position,
			IsActive:   task.IsActive,
			ClientName: task.ClientName,
			ClientID:   task.ClientID,
		})
	}

	return &serviceInterfaces.TasksResponseDTO{
		Success:    true,
		Message:    "Successfully retrieved active tasks",
		Tasks:      taskDTOs,
		TotalCount: len(taskDTOs),
	}, nil
}

func (s *taskQueryService) GetTaskStatusHistory(
	ctx context.Context,
	clientName string,
	clientID string,
) (*serviceInterfaces.StatusHistoryResponseDTO, error) {
	// Validate input parameters
	if err := s.validator.ValidateClientParams(clientName, clientID); err != nil {
		return &serviceInterfaces.StatusHistoryResponseDTO{
			Success: false,
			Message: fmt.Sprintf("Invalid parameters: %v", err),
		}, nil
	}

	// Get status history from repository
	history, err := s.repo.GetTaskStatusHistory(ctx, clientName, clientID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get task status history")
		return &serviceInterfaces.StatusHistoryResponseDTO{
			Success: false,
			Message: "Failed to retrieve status history",
		}, err
	}

	// Map repository data to DTOs
	var historyDTOs []serviceInterfaces.StatusDetailDTO
	for _, status := range history {
		historyDTOs = append(historyDTOs, serviceInterfaces.StatusDetailDTO{
			TaskID:            status.TaskID,
			Status:            status.Status,
			StatusDescription: status.StatusDescription,
			UpdatedBy:         status.UpdatedBy,
			CreatedAt:         status.CreatedAt,
		})
	}

	return &serviceInterfaces.StatusHistoryResponseDTO{
		Success:    true,
		Message:    "Successfully retrieved status history",
		History:    historyDTOs,
		TotalCount: len(historyDTOs),
	}, nil
}
