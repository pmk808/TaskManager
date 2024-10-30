package interfaces

import (
	"context"
	"time"
)

// TaskQueryService defines the interface for querying tasks
type TaskQueryService interface {
	// GetActiveTasks retrieves active tasks for a specific client
	GetActiveTasks(ctx context.Context, clientName string, clientID string) (*TasksResponseDTO, error)

	// GetTaskStatusHistory retrieves status history for a specific client
	GetTaskStatusHistory(ctx context.Context, clientName string, clientID string) (*StatusHistoryResponseDTO, error)
}

// TasksResponseDTO represents the response for active tasks query
type TasksResponseDTO struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	Tasks      []TaskDetailDTO `json:"tasks,omitempty"`
	TotalCount int             `json:"total_count"`
}

// StatusHistoryResponseDTO represents the response for status history query
type StatusHistoryResponseDTO struct {
	Success    bool              `json:"success"`
	Message    string            `json:"message"`
	History    []StatusDetailDTO `json:"history,omitempty"`
	TotalCount int               `json:"total_count"`
}

// TaskDetailDTO represents detailed task information
type TaskDetailDTO struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department"`
	Position   string `json:"position"`
	IsActive   bool   `json:"is_active"`
	ClientName string `json:"client_name"`
	ClientID   string `json:"client_id"`
}

// StatusDetailDTO represents detailed status information
type StatusDetailDTO struct {
	TaskID            int       `json:"task_id"`
	Status            string    `json:"status"`
	StatusDescription string    `json:"status_description"`
	UpdatedBy         string    `json:"updated_by"`
	CreatedAt         time.Time `json:"created_at"`
}
