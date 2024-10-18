package interfaces

import "taskmanager/schemas"

type TaskRepository interface {
	CreateTableIfNotExists() error
	CreateTask(task *schemas.Task) error
	BulkCreateTasks(tasks []schemas.Task) error
	// Add other repository methods as needed
}
