package validation

import (
	"fmt"
	"reflect"
	"taskmanager/schemas"

	"github.com/sirupsen/logrus"
)

type Validator struct {
	logger *logrus.Logger
}

func NewValidator(logger *logrus.Logger) *Validator {
	return &Validator{logger: logger}
}

func (v *Validator) ValidateTask(task *schemas.Task) error {
	val := reflect.ValueOf(*task)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if field.Interface() == reflect.Zero(field.Type()).Interface() {
			v.logger.WithFields(logrus.Fields{
				"field": fieldType.Name,
				"task":  task,
			}).Warn("Field has zero value")
			return fmt.Errorf("field %s is required but has a zero value", fieldType.Name)
		}
	}

	return nil
}

func (v *Validator) ValidateTasks(tasks []schemas.Task) error {
	v.logger.WithField("task_count", len(tasks)).Info("Starting task validation")

	for i, task := range tasks {
		if err := v.ValidateTask(&task); err != nil {
			v.logger.WithFields(logrus.Fields{
				"row":   i + 1,
				"error": err,
			}).Error("Validation failed for task")
			return fmt.Errorf("validation failed for row %d: %w", i+1, err)
		}
	}

	v.logger.Info("All tasks passed validation")
	return nil
}
