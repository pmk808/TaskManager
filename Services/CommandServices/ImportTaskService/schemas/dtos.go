package schemas

import "time"

// ImportTaskRequestDTO represents the incoming request for task import
type ImportTaskRequestDTO struct {
	DirectoryPath string `json:"directory_path,omitempty"`
}

// ImportTaskResponseDTO represents the response after importing tasks
type ImportTaskResponseDTO struct {
	Success      bool            `json:"success"`
	Message      string          `json:"message"`
	ImportedAt   time.Time       `json:"imported_at"`
	TotalEntries int             `json:"total_entries"`
	Errors       []string        `json:"errors,omitempty"`
	Stats        *ImportStatsDTO `json:"stats,omitempty"`
}

// ImportStatsDTO represents statistics about the import process
type ImportStatsDTO struct {
	TotalProcessed int       `json:"total_processed"`
	SuccessCount   int       `json:"success_count"`
	ErrorCount     int       `json:"error_count"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	DurationMS     int64     `json:"duration_ms"`
}

// TaskImportDTO represents a single row of task data to be imported
type TaskImportDTO struct {
	Name        string    `json:"name" validate:"required,max=100"`
	Email       string    `json:"email" validate:"required,email,max=100"`
	Age         int       `json:"age" validate:"required,gt=0,lt=150"`
	Address     string    `json:"address" validate:"required"`
	PhoneNumber string    `json:"phone_number" validate:"required,max=15"`
	Department  string    `json:"department" validate:"required,max=50"`
	Position    string    `json:"position" validate:"required,max=50"`
	Salary      float64   `json:"salary" validate:"required,gte=0"`
	HireDate    time.Time `json:"hire_date" validate:"required"`
}
