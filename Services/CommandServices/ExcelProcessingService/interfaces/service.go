package interfaces

import "taskmanager/schemas"

type ExcelProcessingService interface {
    ReadTasksFromSpreadsheet(filePath string) ([]schemas.Task, error)
}
