package interfaces

import "taskmanager/schemas"

type TaskValidator interface {
    ValidateTask(task *schemas.Task) error
    ValidateTasks(tasks []schemas.Task) error
}