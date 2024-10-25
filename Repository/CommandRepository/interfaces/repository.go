package interfaces

import (
	"taskmanager/Services/CommandServices/ImportTaskService/schemas"
)

type TaskCommandRepository interface {
	BulkCreateTasks(tasks []schemas.TaskImportEntry) error
}
