package interfaces

import "taskmanager/schemas"

type Validator interface {
	ValidateTasks(tasks []schemas.Task) error
}
