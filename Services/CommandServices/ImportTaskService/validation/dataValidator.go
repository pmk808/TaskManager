package validation

import (
	"fmt"
	"regexp"
	"strings"
	"taskmanager/Services/CommandServices/ImportTaskService/schemas"
	"taskmanager/Services/CommandServices/ImportTaskService/validation/interfaces"

	"github.com/sirupsen/logrus"
)

type dataValidator struct {
	logger *logrus.Logger
}

func NewDataValidator(logger *logrus.Logger) interfaces.Validator {
	return &dataValidator{
		logger: logger,
	}
}

func (v *dataValidator) ValidateEntry(entry *schemas.TaskImportDTO) error {
	if err := v.checkName(entry.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if err := v.checkEmail(entry.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	if err := v.checkAge(entry.Age); err != nil {
		return fmt.Errorf("invalid age: %w", err)
	}

	if err := v.checkAddress(entry.Address); err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	if err := v.checkPhone(entry.PhoneNumber); err != nil {
		return fmt.Errorf("invalid phone number: %w", err)
	}

	if err := v.checkDepartment(entry.Department); err != nil {
		return fmt.Errorf("invalid department: %w", err)
	}

	if err := v.checkPosition(entry.Position); err != nil {
		return fmt.Errorf("invalid position: %w", err)
	}

	if err := v.checkSalary(entry.Salary); err != nil {
		return fmt.Errorf("invalid salary: %w", err)
	}

	return nil
}

func (v *dataValidator) ValidateBatch(entries []schemas.TaskImportDTO) error {
	v.logger.WithField("entry_count", len(entries)).Info("Starting batch validation")

	for i, entry := range entries {
		if err := v.ValidateEntry(&entry); err != nil {
			v.logger.WithFields(logrus.Fields{
				"row":   i + 1,
				"error": err,
			}).Error("Validation failed for entry")
			return fmt.Errorf("validation failed for row %d: %w", i+1, err)
		}
	}

	v.logger.Info("All entries passed validation")
	return nil
}

func (v *dataValidator) checkName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("name cannot be longer than 100 characters")
	}
	return nil
}

func (v *dataValidator) checkEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (v *dataValidator) checkAge(age int) error {
	if age <= 0 {
		return fmt.Errorf("age must be greater than 0")
	}
	if age > 150 {
		return fmt.Errorf("age cannot be greater than 150")
	}
	return nil
}

func (v *dataValidator) checkAddress(address string) error {
	address = strings.TrimSpace(address)
	if address == "" {
		return fmt.Errorf("address cannot be empty")
	}
	return nil
}

func (v *dataValidator) checkPhone(phoneNumber string) error {
	phoneNumber = strings.TrimSpace(phoneNumber)
	if phoneNumber == "" {
		return fmt.Errorf("phone number cannot be empty")
	}

	// Remove common phone number formatting characters
	cleanPhone := strings.Map(func(r rune) rune {
		switch r {
		case ' ', '-', '(', ')', '+':
			return -1
		default:
			return r
		}
	}, phoneNumber)

	// Check if the cleaned number contains only digits
	if !regexp.MustCompile(`^\d{7,15}$`).MatchString(cleanPhone) {
		return fmt.Errorf("phone number must contain 7-15 digits")
	}

	return nil
}
func (v *dataValidator) checkDepartment(department string) error {
	department = strings.TrimSpace(department)
	if department == "" {
		return fmt.Errorf("department cannot be empty")
	}
	if len(department) > 50 {
		return fmt.Errorf("department name cannot be longer than 50 characters")
	}
	return nil
}

func (v *dataValidator) checkPosition(position string) error {
	position = strings.TrimSpace(position)
	if position == "" {
		return fmt.Errorf("position cannot be empty")
	}
	if len(position) > 50 {
		return fmt.Errorf("position name cannot be longer than 50 characters")
	}
	return nil
}

func (v *dataValidator) checkSalary(salary float64) error {
	if salary < 0 {
		return fmt.Errorf("salary cannot be negative")
	}
	return nil
}
