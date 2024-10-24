package interfaces

import (
	"taskmanager/Services/CommandServices/ImportTaskService/schemas"
)

type TaskRepository interface {
	CreateTableIfNotExists() error
	BulkCreateTasks(tasks []schemas.TaskImportEntry) error
}
