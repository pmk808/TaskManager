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

	// Create table if it doesn't exist
	if err := repo.CreateTableIfNotExists(); err != nil {
		return nil, fmt.Errorf("failed to initialize table: %w", err)
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
        phone_number VARCHAR(15) NOT NULL,
        department VARCHAR(50) NOT NULL,
        position VARCHAR(50) NOT NULL,
        salary DECIMAL(10, 2) NOT NULL,
        hire_date DATE NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
    `

	r.logger.Info("Creating tasks table if it doesn't exist")

	_, err := r.db.Exec(query)
	if err != nil {
		r.logger.WithError(err).Error("Failed to create tasks table")
		return fmt.Errorf("failed to create table: %w", err)
	}

	r.logger.Info("Tasks table created or verified successfully")
	return nil
}

func (r *taskCommandRepository) BulkCreateTasks(tasks []schemas.TaskImportEntry) error {
	r.logger.WithField("task_count", len(tasks)).Info("Starting bulk task creation")

	// Ensure table exists before proceeding
	if err := r.CreateTableIfNotExists(); err != nil {
		return fmt.Errorf("failed to verify table existence: %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		r.logger.WithError(err).Error("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create a temporary table for bulk insert
	createTempTableQuery := `
        CREATE TEMP TABLE temp_tasks (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            email VARCHAR(100) NOT NULL,
            age INT NOT NULL,
            address TEXT NOT NULL,
            phone_number VARCHAR(15) NOT NULL,
            department VARCHAR(50) NOT NULL,
            position VARCHAR(50) NOT NULL,
            salary DECIMAL(10, 2) NOT NULL,
            hire_date DATE NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `
	_, err = tx.Exec(createTempTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create temp table: %w", err)
	}

	// Prepare the INSERT statement
	insertStmt := `
        INSERT INTO temp_tasks (
            name, email, age, address, phone_number, 
            department, position, salary, hire_date
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	stmt, err := tx.Prepare(insertStmt)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert all tasks into the temporary table
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

	// Insert from temporary table to actual table
	insertFromTempQuery := `
        INSERT INTO tasks (
            name, email, age, address, phone_number,
            department, position, salary, hire_date
        )
        SELECT 
            name, email, age, address, phone_number,
            department, position, salary, hire_date
        FROM temp_tasks
    `
	_, err = tx.Exec(insertFromTempQuery)
	if err != nil {
		return fmt.Errorf("failed to insert from temp table: %w", err)
	}

	if err := tx.Commit(); err != nil {
		r.logger.WithError(err).Error("Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.WithField("task_count", len(tasks)).Info("Bulk task creation completed successfully")
	return nil
}
