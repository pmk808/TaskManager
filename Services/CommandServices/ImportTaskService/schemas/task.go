// Services/CommandServices/ImportTaskService/schemas/task.go
package schemas

import (
	"time"
)

// TaskImportEntry represents a single row of task data to be imported
type TaskImportEntry struct {
	Name        string    `json:"name" db:"name" validate:"required,max=100"`
	Email       string    `json:"email" db:"email" validate:"required,email,max=100"`
	Age         int       `json:"age" db:"age" validate:"required,gt=0,lt=150"`
	Address     string    `json:"address" db:"address" validate:"required"`
	PhoneNumber string    `json:"phone_number" db:"phone_number" validate:"required,max=15"`
	Department  string    `json:"department" db:"department" validate:"required,max=50"`
	Position    string    `json:"position" db:"position" validate:"required,max=50"`
	Salary      float64   `json:"salary" db:"salary" validate:"required,gte=0"`
	HireDate    time.Time `json:"hire_date" db:"hire_date" validate:"required"`
}

// TaskImportResponse represents the response after importing tasks
type TaskImportResponse struct {
	Success      bool             `json:"success"`
	Message      string           `json:"message"`
	ImportedAt   time.Time        `json:"imported_at"`
	TotalEntries int              `json:"total_entries"`
	Errors       []string         `json:"errors,omitempty"`
	Stats        *TaskImportStats `json:"stats,omitempty"`
}

// TaskImportStats represents statistics about the import process
type TaskImportStats struct {
	TotalProcessed int       `json:"total_processed"`
	SuccessCount   int       `json:"success_count"`
	ErrorCount     int       `json:"error_count"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	DurationMS     int64     `json:"duration_ms"` // Changed from time.Duration to int64
}
