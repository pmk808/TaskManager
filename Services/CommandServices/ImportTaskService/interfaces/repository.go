package interfaces

import "taskmanager/schemas"

type TaskRepository interface {
    BulkCreateTasks(tasks []schemas.Task) error
}
