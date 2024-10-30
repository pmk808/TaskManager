package interfaces

import (
	"context"
)

// TaskQueryRepository defines the methods for querying tasks
type TaskQueryRepository interface {
	// GetActiveTasks retrieves active tasks for a specific client
	GetActiveTasks(ctx context.Context, clientName string, clientID string) ([]TaskDTO, error)

	// GetTaskStatusHistory retrieves status history for a specific client
	GetTaskStatusHistory(ctx context.Context, clientName string, clientID string) ([]TaskStatusDTO, error)
}

// TaskDTO represents a task query result
type TaskDTO struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Age         int     `json:"age"`
	Address     string  `json:"address"`
	PhoneNumber string  `json:"phone_number"`
	Department  string  `json:"department"`
	Position    string  `json:"position"`
	Salary      float64 `json:"salary"`
	ClientName  string  `json:"client_name"`
	ClientID    string  `json:"client_id"`
	IsActive    bool    `json:"is_active"`
}

// TaskStatusDTO represents a task status history entry
type TaskStatusDTO struct {
	TaskID            int    `json:"task_id"`
	ClientName        string `json:"client_name"`
	ClientID          string `json:"client_id"`
	Status            string `json:"status"`
	StatusDescription string `json:"status_description"`
	UpdatedBy         string `json:"updated_by"`
	CreatedAt         string `json:"created_at"`
}
