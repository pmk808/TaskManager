package validation

import (
	"io"
	"testing"
	"time"

	"taskmanager/schemas"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestValidateTask(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard) // Discard logs in tests

	validator := NewValidator(logger)

	tests := []struct {
		name    string
		task    schemas.Task
		wantErr bool
	}{
		{
			name: "Valid task",
			task: schemas.Task{
				Name:        "John Doe",
				Email:       "john@example.com",
				Age:         30,
				Address:     "123 Main St",
				PhoneNumber: "123-456-7890",
				Department:  "IT",
				Position:    "Developer",
				Salary:      50000,
				HireDate:    time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Missing name",
			task: schemas.Task{
				Email:       "john@example.com",
				Age:         30,
				Address:     "123 Main St",
				PhoneNumber: "123-456-7890",
				Department:  "IT",
				Position:    "Developer",
				Salary:      50000,
				HireDate:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Zero age",
			task: schemas.Task{
				Name:        "John Doe",
				Email:       "john@example.com",
				Age:         0,
				Address:     "123 Main St",
				PhoneNumber: "123-456-7890",
				Department:  "IT",
				Position:    "Developer",
				Salary:      50000,
				HireDate:    time.Now(),
			},
			wantErr: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTask(&tt.task)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTasks(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard) // Discard logs in tests

	validator := NewValidator(logger)

	validTask := schemas.Task{
		Name:        "John Doe",
		Email:       "john@example.com",
		Age:         30,
		Address:     "123 Main St",
		PhoneNumber: "123-456-7890",
		Department:  "IT",
		Position:    "Developer",
		Salary:      50000,
		HireDate:    time.Now(),
	}

	invalidTask := schemas.Task{
		// Missing name
		Email:       "jane@example.com",
		Age:         25,
		Address:     "456 Elm St",
		PhoneNumber: "987-654-3210",
		Department:  "HR",
		Position:    "Manager",
		Salary:      60000,
		HireDate:    time.Now(),
	}

	tests := []struct {
		name    string
		tasks   []schemas.Task
		wantErr bool
	}{
		{
			name:    "All valid tasks",
			tasks:   []schemas.Task{validTask, validTask},
			wantErr: false,
		},
		{
			name:    "One invalid task",
			tasks:   []schemas.Task{validTask, invalidTask},
			wantErr: true,
		},
		{
			name:    "Empty task list",
			tasks:   []schemas.Task{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTasks(tt.tasks)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
