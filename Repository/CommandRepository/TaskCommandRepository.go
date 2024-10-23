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

func (r *taskCommandRepository) CreateTableIfNotExists() error {
	query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) NOT NULL,
        age INT NOT NULL,
        address TEXT NOT NULL,
        phone_number VARCHAR(15),
        department VARCHAR(50),
        position VARCHAR(50),
        salary DECIMAL(10, 2),
        hire_date DATE
    );`

	_, err := r.db.Exec(query)
	if err != nil {
		r.logger.WithError(err).Error("Failed to create tasks table")
		return fmt.Errorf("failed to create table: %w", err)
	}

	r.logger.Info("Tasks table created or already exists")
	return nil
}

func (r *taskCommandRepository) BulkCreateTasks(tasks []schemas.TaskImportEntry) error {
	r.logger.WithField("entry_count", len(tasks)).Info("Starting bulk task creation")

	tx, err := r.db.Begin()
	if err != nil {
		r.logger.WithError(err).Error("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Using COPY for better performance with bulk inserts
	stmt, err := tx.Prepare(`
        COPY tasks (name, email, age, address, phone_number, department, position, salary, hire_date)
        FROM STDIN WITH (FORMAT csv)
    `)
	if err != nil {
		r.logger.WithError(err).Error("Failed to prepare statement")
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, task := range tasks {
		_, err := stmt.Exec(
			task.Name, task.Email, task.Age, task.Address,
			task.PhoneNumber, task.Department, task.Position,
			task.Salary, task.HireDate,
		)
		if err != nil {
			r.logger.WithError(err).Error("Failed to insert task")
			return fmt.Errorf("failed to insert task: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.WithError(err).Error("Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.WithField("task_count", len(tasks)).Info("Bulk task creation completed successfully")
	return nil
}
