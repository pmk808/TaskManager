package services

import (
	"fmt"
	"path/filepath"
	"strconv"
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

func NewTaskService(repo interfaces.TaskRepository, validator interfaces.Validator, logger *logrus.Logger, directory string) *TaskService {
	return &TaskService{
		repo:      repo,
		validator: validator,
		logger:    logger,
		directory: directory,
	}
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
		// Log the entire row content
		s.logger.WithField("row", row).Infof("Processing row %d", i+1)

		if i == 0 {
			continue // Skip header row
		}
		if len(row) < 9 {
			s.logger.WithField("row", row).Errorf("Row %d has insufficient columns", i+1)
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

// Helper functions for parsing
func parseInt(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		logrus.WithField("value", s).WithError(err).Error("Failed to parse int")
		return 0 // Handle as appropriate
	}
	return value
}

func parseFloat(s string) float64 {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logrus.WithField("value", s).WithError(err).Error("Failed to parse float")
		return 0.0 // Handle as appropriate
	}
	return value
}

func parseDate(s string) time.Time {
	value, err := time.Parse("01-02-06", s) // Adjust the layout as needed
	if err != nil {
		logrus.WithField("value", s).WithError(err).Error("Failed to parse date")
		return time.Time{} // Handle as appropriate
	}
	return value
}
