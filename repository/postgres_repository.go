package repository

import (
	"database/sql"
	"fmt"
	"time"

	"taskmanager/config"
	"taskmanager/interfaces"
	"taskmanager/schemas"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type PostgresRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPostgresRepository(cfg *config.DatabaseConfig, logger *logrus.Logger) (interfaces.TaskRepository, error) {
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

	repo := &PostgresRepository{db: db, logger: logger}

	if err := repo.CreateTableIfNotExists(); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	logger.Info("PostgreSQL repository initialized successfully")
	return repo, nil
}

func (r *PostgresRepository) BulkCreateTasks(tasks []schemas.Task) error {
	r.logger.WithField("task_count", len(tasks)).Info("Starting bulk task creation")

	tx, err := r.db.Begin()
	if err != nil {
		r.logger.WithError(err).Error("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
        INSERT INTO tasks (name, email, age, address, phone_number, department, position, salary, hire_date)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `)
	if err != nil {
		r.logger.WithError(err).Error("Failed to prepare statement")
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for i, task := range tasks {
		_, err := stmt.Exec(
			task.Name, task.Email, task.Age, task.Address, task.PhoneNumber,
			task.Department, task.Position, task.Salary, task.HireDate,
		)
		if err != nil {
			r.logger.WithFields(logrus.Fields{
				"row":   i + 1,
				"error": err,
			}).Error("Failed to insert task")
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
