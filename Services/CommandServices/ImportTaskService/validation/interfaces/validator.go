package interfaces

import (
	"taskmanager/Services/CommandServices/ImportTaskService/schemas"
)

type Validator interface {
	ValidateEntry(entry *schemas.TaskImportDTO) error
	ValidateBatch(entries []schemas.TaskImportDTO) error
}
