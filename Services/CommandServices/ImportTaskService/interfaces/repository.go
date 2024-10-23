package interfaces

import (
    // "context"
    "taskmanager/schemas"
)

type TaskRepository interface {
    CreateTableIfNotExists() error
    BulkCreateTasks(tasks []schemas.Task) error
}