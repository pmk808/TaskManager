package interfaces

import (
	"taskmanager/Services/CommandServices/ImportTaskService/schemas"
)

type TaskCommandRepository interface {
	CreateTableIfNotExists() error
	BulkCreateTasks(tasks []schemas.TaskImportEntry) error
}
