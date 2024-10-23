package interfaces

import (
	"taskmanager/Services/CommandServices/ImportTaskService/schemas"
)

type Validator interface {
	ValidateEntry(entry *schemas.TaskImportEntry) error
	ValidateBatch(entries []schemas.TaskImportEntry) error
}
