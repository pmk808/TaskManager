package QueryRepository

import (
	"context"
	"database/sql"
	"fmt"
	"taskmanager/Repository/QueryRepository/interfaces"
	"taskmanager/RequestControllers/httpSetup/config"
	"time"

	"github.com/sirupsen/logrus"
)

type taskQueryRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewTaskQueryRepository creates a new instance of TaskQueryRepository
func NewTaskQueryRepository(cfg *config.DatabaseConfig, logger *logrus.Logger) (interfaces.TaskQueryRepository, error) {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &taskQueryRepository{
		db:     db,
		logger: logger,
	}, nil
}

// GetActiveTasks retrieves active tasks for a specific client
func (r *taskQueryRepository) GetActiveTasks(ctx context.Context, clientName string, clientID string) ([]interfaces.TaskDTO, error) {
	r.logger.WithFields(logrus.Fields{
		"client_name": clientName,
		"client_id":   clientID,
	}).Debug("Querying active tasks")
	query := `
		SELECT 
			id, name, email, age, address, phone_number,
			department, position, salary, client_name,
			client_id, is_active
		FROM task_management.tasks
		WHERE client_name = $1 
		AND client_id = $2
		AND is_active = true
	`

	r.logger.WithFields(logrus.Fields{
		"query":  query,
		"params": []interface{}{clientName, clientID},
	}).Debug("Executing query")

	rows, err := r.db.QueryContext(ctx, query, clientName, clientID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to query active tasks")
		return nil, fmt.Errorf("failed to query active tasks: %w", err)
	}
	defer rows.Close()

	var tasks []interfaces.TaskDTO
	for rows.Next() {
		var task interfaces.TaskDTO
		err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Email,
			&task.Age,
			&task.Address,
			&task.PhoneNumber,
			&task.Department,
			&task.Position,
			&task.Salary,
			&task.ClientName,
			&task.ClientID,
			&task.IsActive,
		)
		if err != nil {
			r.logger.WithError(err).Error("Failed to scan task row")
			return nil, fmt.Errorf("failed to scan task row: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		r.logger.WithError(err).Error("Error during row iteration")
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}
	r.logger.WithField("task_count", len(tasks)).Info("Retrieved active tasks")
	return tasks, nil
}

// GetTaskStatusHistory retrieves status history for a specific client
func (r *taskQueryRepository) GetTaskStatusHistory(ctx context.Context, clientName string, clientID string) ([]interfaces.TaskStatusDTO, error) {
	r.logger.WithFields(logrus.Fields{
		"client_name": clientName,
		"client_id":   clientID,
	}).Debug("Querying task status history")
	query := `
		SELECT 
			task_id, client_name, client_id, status,
			status_description, updated_by, 
			created_at
		FROM task_management.task_status
		WHERE client_name = $1 
		AND client_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, clientName, clientID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to query task status history")
		return nil, fmt.Errorf("failed to query task status history: %w", err)
	}
	defer rows.Close()

	var statusHistory []interfaces.TaskStatusDTO
	for rows.Next() {
		var status interfaces.TaskStatusDTO
		err := rows.Scan(
			&status.TaskID,
			&status.ClientName,
			&status.ClientID,
			&status.Status,
			&status.StatusDescription,
			&status.UpdatedBy,
			&status.CreatedAt,
		)
		if err != nil {
			r.logger.WithError(err).Error("Failed to scan status row")
			return nil, fmt.Errorf("failed to scan status row: %w", err)
		}
		statusHistory = append(statusHistory, status)
	}

	if err = rows.Err(); err != nil {
		r.logger.WithError(err).Error("Error during row iteration")
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}
	r.logger.WithField("history_count", len(statusHistory)).Info("Retrieved status history")
	return statusHistory, nil
}
