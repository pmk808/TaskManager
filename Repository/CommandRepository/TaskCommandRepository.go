package CommandRepository

import (
	"database/sql"
	"fmt"
	"time"

	"taskmanager/Repository/CommandRepository/interfaces"
	"taskmanager/RequestControllers/httpSetup/config"
	"taskmanager/Services/CommandServices/ImportTaskService/schemas"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type taskCommandRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewTaskCommandRepository(cfg *config.DatabaseConfig, logger *logrus.Logger) (interfaces.TaskCommandRepository, error) {
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

	repo := &taskCommandRepository{
		db:     db,
		logger: logger,
	}

	logger.Info("Task command repository initialized successfully")
	return repo, nil
}

// BulkCreateTasks handles bulk insertion of tasks
func (r *taskCommandRepository) BulkCreateTasks(tasks []schemas.TaskImportEntry) error {
	r.logger.WithField("task_count", len(tasks)).Info("Starting bulk task creation")

	tx, err := r.db.Begin()
	if err != nil {
		r.logger.WithError(err).Error("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare statement for inserting tasks
	stmt, err := tx.Prepare(`
        INSERT INTO task_management.tasks (
            name, email, age, address, phone_number, 
            department, position, salary, hire_date
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert all tasks
	for i, task := range tasks {
		_, err = stmt.Exec(
			task.Name,
			task.Email,
			task.Age,
			task.Address,
			task.PhoneNumber,
			task.Department,
			task.Position,
			task.Salary,
			task.HireDate,
		)
		if err != nil {
			return fmt.Errorf("failed to insert task at row %d: %w", i+1, err)
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.WithError(err).Error("Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.WithField("task_count", len(tasks)).Info("Bulk task creation completed successfully")
	return nil
}
