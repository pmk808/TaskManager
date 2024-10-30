package schemas

import "time"

// TaskModel represents the internal task data structure
type TaskModel struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Email       string    `db:"email"`
	Age         int       `db:"age"`
	Address     string    `db:"address"`
	PhoneNumber string    `db:"phone_number"`
	Department  string    `db:"department"`
	Position    string    `db:"position"`
	Salary      float64   `db:"salary"`
	HireDate    time.Time `db:"hire_date"`
	IsActive    bool      `db:"is_active"`
	ClientName  string    `db:"client_name"`
	ClientID    string    `db:"client_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// MapFromDTO converts a TaskImportDTO to a TaskModel
func (m *TaskModel) MapFromDTO(dto TaskImportDTO) {
	m.Name = dto.Name
	m.Email = dto.Email
	m.Age = dto.Age
	m.Address = dto.Address
	m.PhoneNumber = dto.PhoneNumber
	m.Department = dto.Department
	m.Position = dto.Position
	m.Salary = dto.Salary
	m.HireDate = dto.HireDate
	m.IsActive = true // default value for new tasks
}