-- Create function to generate dummy data
CREATE OR REPLACE FUNCTION task_management.insert_dummy_data()
RETURNS void AS $$
DECLARE
    client1_id UUID := gen_random_uuid();
    client2_id UUID := gen_random_uuid();
    client3_id UUID := gen_random_uuid();
    task_id INTEGER;
BEGIN
    -- Insert dummy tasks for Client 1
    INSERT INTO task_management.tasks 
        (name, email, age, address, phone_number, department, position, salary, hire_date, is_active, client_name, client_id)
    VALUES
        ('John Doe', 'john@client1.com', 30, '123 Client1 St', '1234567890', 'IT', 'Developer', 75000, '2023-01-01', TRUE, 'Client One Corp', client1_id),
        ('Jane Smith', 'jane@client1.com', 28, '456 Client1 St', '0987654321', 'HR', 'Manager', 85000, '2023-02-01', FALSE, 'Client One Corp', client1_id);

    -- Insert dummy tasks for Client 2
    INSERT INTO task_management.tasks 
        (name, email, age, address, phone_number, department, position, salary, hire_date, is_active, client_name, client_id)
    VALUES
        ('Bob Wilson', 'bob@client2.com', 35, '789 Client2 St', '1122334455', 'Sales', 'Executive', 90000, '2023-03-01', TRUE, 'Client Two LLC', client2_id),
        ('Alice Brown', 'alice@client2.com', 32, '321 Client2 St', '5544332211', 'Marketing', 'Director', 95000, '2023-04-01', TRUE, 'Client Two LLC', client2_id);

    -- Insert dummy tasks for Client 3
    INSERT INTO task_management.tasks 
        (name, email, age, address, phone_number, department, position, salary, hire_date, is_active, client_name, client_id)
    VALUES
        ('Charlie Davis', 'charlie@client3.com', 40, '654 Client3 St', '6677889900', 'Operations', 'Manager', 88000, '2023-05-01', TRUE, 'Client Three Inc', client3_id);

    -- Insert status history for each task
    FOR task_id IN (SELECT id FROM task_management.tasks) LOOP
        INSERT INTO task_management.task_status 
            (task_id, client_name, client_id, status, status_description, updated_by)
        VALUES
            (task_id, 
             (SELECT client_name FROM task_management.tasks WHERE id = task_id),
             (SELECT client_id FROM task_management.tasks WHERE id = task_id),
             'PENDING',
             'Initial status',
             'system');

        -- Add some tasks with multiple status changes
        IF task_id % 2 = 0 THEN
            INSERT INTO task_management.task_status 
                (task_id, client_name, client_id, status, status_description, updated_by)
            VALUES
                (task_id,
                 (SELECT client_name FROM task_management.tasks WHERE id = task_id),
                 (SELECT client_id FROM task_management.tasks WHERE id = task_id),
                 'IN_PROGRESS',
                 'Work started',
                 'system');
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Execute the function to insert dummy data
SELECT task_management.insert_dummy_data();

-- Drop the function after use
DROP FUNCTION task_management.insert_dummy_data();