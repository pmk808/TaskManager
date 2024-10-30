package QueryRepository

import (
	"context"
	"database/sql"
	"testing"

	"taskmanager/Repository/QueryRepository/interfaces"
	"taskmanager/RequestControllers/httpSetup/config"

	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	clientOneUUID = "9ebcc92c-e186-41b3-834b-f75ab3f110ae"
	clientTwoUUID = "34fb4178-bee7-4c5d-b13c-7a4ac405d56d"
)

func setupTestDB(t *testing.T) (*sql.DB, *logrus.Logger) {
	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Load test configuration
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "peemak", // Use your actual test DB password
		DBName:   "taskmanager",
		SSLMode:  "disable",
	}

	// Initialize repository
	repo, err := NewTaskQueryRepository(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Get the underlying database connection
	db := repo.(*taskQueryRepository).db

	return db, logger
}

func TestGetActiveTasks(t *testing.T) {
	db, logger := setupTestDB(t)
	defer db.Close()

	repo := &taskQueryRepository{
		db:     db,
		logger: logger,
	}

	tests := []struct {
		name       string
		clientName string
		clientID   string
		wantErr    bool
		setup      func(t *testing.T, db *sql.DB)
		verify     func(t *testing.T, tasks []interfaces.TaskDTO, err error)
	}{
		{
			name:       "Valid client with active tasks",
			clientName: "Client One Corp",
			clientID:   clientOneUUID,
			wantErr:    false,
			setup: func(t *testing.T, db *sql.DB) {
			},
			verify: func(t *testing.T, tasks []interfaces.TaskDTO, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, tasks)
				for _, task := range tasks {
					assert.True(t, task.IsActive)
					assert.Equal(t, "Client One Corp", task.ClientName)
				}
			},
		},
		{
			name:       "Non-existent client",
			clientName: "Non Existent Corp",
			clientID:   "123e4567-e89b-12d3-a456-426614174999",
			wantErr:    false,
			verify: func(t *testing.T, tasks []interfaces.TaskDTO, err error) {
				assert.NoError(t, err)
				assert.Empty(t, tasks)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, db)
			}

			tasks, err := repo.GetActiveTasks(context.Background(), tt.clientName, tt.clientID)
			if tt.verify != nil {
				tt.verify(t, tasks, err)
			} else {
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestGetTaskStatusHistory(t *testing.T) {
	db, logger := setupTestDB(t)
	defer db.Close()

	repo := &taskQueryRepository{
		db:     db,
		logger: logger,
	}

	tests := []struct {
		name       string
		clientName string
		clientID   string
		wantErr    bool
		setup      func(t *testing.T, db *sql.DB)
		verify     func(t *testing.T, history []interfaces.TaskStatusDTO, err error)
	}{
		{
			name:       "Valid client with status history",
			clientName: "Client One Corp",
			clientID:   clientOneUUID,
			wantErr:    false,
			verify: func(t *testing.T, history []interfaces.TaskStatusDTO, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, history)
				for _, status := range history {
					assert.Equal(t, "Client One Corp", status.ClientName)
					assert.NotEmpty(t, status.Status)
				}
			},
		},
		{
			name:       "Non-existent client",
			clientName: "Non Existent Corp",
			clientID:   "123e4567-e89b-12d3-a456-426614174999",
			wantErr:    false,
			verify: func(t *testing.T, history []interfaces.TaskStatusDTO, err error) {
				assert.NoError(t, err)
				assert.Empty(t, history)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, db)
			}

			history, err := repo.GetTaskStatusHistory(context.Background(), tt.clientName, tt.clientID)
			if tt.verify != nil {
				tt.verify(t, history, err)
			} else {
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}
