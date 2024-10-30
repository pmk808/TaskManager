-- Clear existing test data
TRUNCATE task_management.tasks CASCADE;
TRUNCATE task_management.task_status CASCADE;

-- Declare variables for client UUIDs
DO $$ 
DECLARE
    client_1_uuid UUID := gen_random_uuid();
    client_2_uuid UUID := gen_random_uuid();
BEGIN
    -- Insert test tasks
    INSERT INTO task_management.tasks (
        name, email, age, address, phone_number, 
        department, position, salary, hire_date, 
        is_active, client_name, client_id
    ) VALUES 
        ('John Doe', 'john@test.com', 30, '123 Test St', '1234567890',
         'IT', 'Developer', 75000, '2023-01-01', true, 
         'Client One Corp', client_1_uuid),
        ('Jane Smith', 'jane@test.com', 28, '456 Test Ave', '0987654321',
         'HR', 'Manager', 85000, '2023-02-01', false,
         'Client One Corp', client_1_uuid),
        ('Bob Wilson', 'bob@test.com', 35, '789 Test Rd', '1122334455',
         'Sales', 'Executive', 90000, '2023-03-01', true,
         'Client Two LLC', client_2_uuid);

    -- Insert test status history
    INSERT INTO task_management.task_status (
        task_id, client_name, client_id, status, 
        status_description, updated_by
    ) 
    SELECT 
        id,
        client_name,
        client_id,
        'PENDING',
        'Initial status',
        'system'
    FROM task_management.tasks;

    -- Add some status changes
    INSERT INTO task_management.task_status (
        task_id, client_name, client_id, status, 
        status_description, updated_by
    )
    SELECT 
        id,
        client_name,
        client_id,
        'IN_PROGRESS',
        'Work started',
        'system'
    FROM task_management.tasks
    WHERE is_active = true;
END $$;