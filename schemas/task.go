package schemas

import (
	"time"
)

type Task struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Email       string    `json:"email" db:"email"`
	Age         int       `json:"age" db:"age"`
	Address     string    `json:"address" db:"address"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Department  string    `json:"department" db:"department"`
	Position    string    `json:"position" db:"position"`
	Salary      float64   `json:"salary" db:"salary"`
	HireDate    time.Time `json:"hire_date" db:"hire_date"`
}
