-- Create task status table
CREATE TABLE IF NOT EXISTS task_management.task_status (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL,
    client_name VARCHAR(100) NOT NULL,
    client_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    status_description TEXT,
    updated_by VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint
    CONSTRAINT fk_task_status_task 
        FOREIGN KEY (task_id) 
        REFERENCES task_management.tasks(id) 
        ON DELETE CASCADE,
        
    -- Check constraint for status values
    CONSTRAINT chk_valid_status 
        CHECK (status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED'))
);

COMMENT ON TABLE task_management.task_status IS 'Stores historical status changes for tasks';
COMMENT ON COLUMN task_management.task_status.task_id IS 'Reference to the task';
COMMENT ON COLUMN task_management.task_status.client_name IS 'Name of the client associated with the task';
COMMENT ON COLUMN task_management.task_status.client_id IS 'Unique identifier for the client';
COMMENT ON COLUMN task_management.task_status.status IS 'Current status of the task';
COMMENT ON COLUMN task_management.task_status.status_description IS 'Optional description or reason for status change';
COMMENT ON COLUMN task_management.task_status.updated_by IS 'User who updated the status';

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_task_status_client 
ON task_management.task_status(client_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_task_status_task 
ON task_management.task_status(task_id, created_at DESC);