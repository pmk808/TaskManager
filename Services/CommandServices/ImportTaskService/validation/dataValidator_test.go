package validation

import (
	"taskmanager/Services/CommandServices/ImportTaskService/schemas"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDataValidator(t *testing.T) {
	logger := logrus.New()
	validator := NewDataValidator(logger)

	tests := []struct {
		name    string
		entry   schemas.TaskImportDTO // Changed from TaskImportEntry
		wantErr bool
	}{
		{
			name: "Valid entry",
			entry: schemas.TaskImportDTO{ // Changed struct name
				Name:        "John Doe",
				Email:       "john@example.com",
				Age:         30,
				Address:     "123 Main St",
				PhoneNumber: "123-456-7890",
				Department:  "Engineering",
				Position:    "Software Engineer",
				Salary:      75000,
				HireDate:    time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Missing name",
			entry: schemas.TaskImportDTO{ // Changed struct name
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
			entry: schemas.TaskImportDTO{ // Changed struct name
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
		{
			name: "Invalid email",
			entry: schemas.TaskImportDTO{ // Changed struct name
				Name:        "John Doe",
				Email:       "invalid-email",
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
			name: "Invalid phone number",
			entry: schemas.TaskImportDTO{ // Changed struct name
				Name:        "John Doe",
				Email:       "john@example.com",
				Age:         30,
				Address:     "123 Main St",
				PhoneNumber: "123", // Too short
				Department:  "IT",
				Position:    "Developer",
				Salary:      50000,
				HireDate:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Negative salary",
			entry: schemas.TaskImportDTO{ // Changed struct name
				Name:        "John Doe",
				Email:       "john@example.com",
				Age:         30,
				Address:     "123 Main St",
				PhoneNumber: "123-456-7890",
				Department:  "IT",
				Position:    "Developer",
				Salary:      -1000,
				HireDate:    time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateEntry(&tt.entry)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDataValidatorBatch(t *testing.T) {
	logger := logrus.New()
	validator := NewDataValidator(logger)

	validEntry := schemas.TaskImportDTO{ // Changed struct name
		Name:        "John Doe",
		Email:       "john@example.com",
		Age:         30,
		Address:     "123 Main St",
		PhoneNumber: "123-456-7890",
		Department:  "Engineering",
		Position:    "Software Engineer",
		Salary:      75000,
		HireDate:    time.Now(),
	}

	invalidEntry := schemas.TaskImportDTO{ // Changed struct name
		Name:        "", // Invalid: empty name
		Email:       "john@example.com",
		Age:         30,
		Address:     "123 Main St",
		PhoneNumber: "123-456-7890",
		Department:  "Engineering",
		Position:    "Software Engineer",
		Salary:      75000,
		HireDate:    time.Now(),
	}

	tests := []struct {
		name    string
		entries []schemas.TaskImportDTO // Changed type
		wantErr bool
	}{
		{
			name:    "Empty batch",
			entries: []schemas.TaskImportDTO{}, // Changed type
			wantErr: false,
		},
		{
			name:    "Single valid entry",
			entries: []schemas.TaskImportDTO{validEntry}, // Changed type
			wantErr: false,
		},
		{
			name:    "Multiple valid entries",
			entries: []schemas.TaskImportDTO{validEntry, validEntry}, // Changed type
			wantErr: false,
		},
		{
			name:    "Contains invalid entry",
			entries: []schemas.TaskImportDTO{validEntry, invalidEntry}, // Changed type
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateBatch(tt.entries)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
