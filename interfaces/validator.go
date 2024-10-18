package interfaces

import "taskmanager/schemas"

type Validator interface {
	ValidateTask(task *schemas.Task) error
	ValidateTasks(tasks []schemas.Task) error
}
