package services

import (
	"fmt"
	// "os"
	"path/filepath"
	"time"

	"taskmanager/interfaces"
	"taskmanager/schemas"

	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

type TaskService struct {
	repo      interfaces.TaskRepository
	validator interfaces.Validator
	logger    *logrus.Logger
	directory string
}

func (s *TaskService) ImportData() error {
	start := time.Now()
	s.logger.Info("Starting data import process")

	tasks, err := s.readTasksFromSpreadsheet()
	if err != nil {
		s.logger.WithError(err).Error("Failed to read tasks from spreadsheet")
		return fmt.Errorf("failed to read tasks from spreadsheet: %w", err)
	}

	s.logger.WithField("task_count", len(tasks)).Info("Tasks read from spreadsheet")

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
	s.logger.WithFields(logrus.Fields{
		"duration":   duration,
		"task_count": len(tasks),
	}).Info("Data import process completed successfully")
	return nil
}

func (s *TaskService) readTasksFromSpreadsheet() ([]schemas.Task, error) {
	files, err := filepath.Glob(filepath.Join(s.directory, "*.xlsx"))
	if err != nil {
		return nil, fmt.Errorf("failed to find Excel files: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no Excel files found in directory: %s", s.directory)
	}

	// Use the first Excel file found
	filePath := files[0]
	s.logger.WithField("file", filePath).Info("Reading spreadsheet")

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1") // Assumes data is in Sheet1
	if err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"file":  filePath,
		"sheet": "Sheet1",
	}).Info("Reading spreadsheet")

	var tasks []schemas.Task
	for i, row := range rows {
		if i == 0 {
			continue // Skip header row
		}
		if len(row) < 10 {
			return nil, fmt.Errorf("row %d has insufficient columns", i+1)
		}

		task := schemas.Task{
			Name:        row[0],
			Email:       row[1],
			Age:         parseInt(row[2]),
			Address:     row[3],
			PhoneNumber: row[4],
			Department:  row[5],
			Position:    row[6],
			Salary:      parseFloat(row[7]),
			HireDate:    parseDate(row[8]),
		}
		tasks = append(tasks, task)
	}
	s.logger.WithField("task_count", len(tasks)).Info("Finished reading tasks from spreadsheet")
	return tasks, nil
}

// Helper functions for parsing (implement these)
func parseInt(s string) int {
	// Implement parsing string to int
}

func parseFloat(s string) float64 {
	// Implement parsing string to float64
}

func parseDate(s string) time.Time {
	// Implement parsing string to time.Time
}
