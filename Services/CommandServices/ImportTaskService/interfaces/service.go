package interfaces

import (
	"taskmanager/Services/CommandServices/ImportTaskService/schemas"
)

type ImportService interface {
	Import() (*schemas.TaskImportResponse, error)
}
