# Task Manager Application

## Project Structure
```
TaskManager/
├── cmd/
│   └── main.go                 # Application entry point
├── Deployment/
│   ├── scripts/
│   │   └── setup_db.sh        # Database initialization script
│   └── sql/
│       ├── 01_create_database.sql
│       ├── 02_create_schema.sql
│       └── 03_create_table.sql
├── Repository/
│   ├── CommandRepository/
│   │   ├── interfaces/
│   │   │   └── repository.go
│   │   └── TaskCommandRepository.go
│   └── QueryRepository/       # (Future implementation)
├── RequestControllers/
│   ├── httpSetup/
│   │   ├── config/
│   │   └── logger/
│   ├── CommandRequest/
│   │   ├── interfaces/
│   │   └── CommandApiController.go
│   └── QueryRequest/          # (Future implementation)
├── Services/
│   └── CommandServices/
│       └── ImportTaskService/
│           ├── interfaces/
│           ├── validation/
│           ├── schemas/
│           └── ImportTaskService.go
├── config.yaml                # Application configuration
└── start.sh                   # Application startup script
```

## Key Changes from Previous Structure

### 1. CQRS Implementation
- Separated Command (write) and Query (read) operations
- Command operations: Import tasks from CSV
- Query operations: (Prepared for future implementation)

### 2. Separation of Database Initialization
**Before:**
- Database table creation was handled in repository layer
- Mixed concerns between data operations and schema management

**After:**
- Moved database initialization to deployment scripts
- Clear separation between database setup and application logic
- Using SQL scripts for schema and table management

### 3. Improved Configuration
- Centralized configuration in config.yaml
- Configuration handling moved to httpSetup/config
- Environment-specific settings support

## Components

### 1. Database Scripts
Located in `deployment/sql/`:
```sql
-- 01_create_database.sql
CREATE DATABASE taskmanager ...

-- 02_create_schema.sql
CREATE SCHEMA task_management ...

-- 03_create_table.sql
CREATE TABLE task_management.tasks ...
```

### 2. Repository Layer
- Focuses purely on data operations
- No schema management responsibilities
- Clear interface definitions

```go
type TaskCommandRepository interface {
    BulkCreateTasks(tasks []schemas.TaskImportEntry) error
}
```

### 3. Service Layer
- Business logic implementation
- CSV processing and validation
- Error handling and logging

### 4. Controllers
- HTTP endpoint handlers
- Request/response management
- Route registration

## Setup and Running

### Prerequisites
- Go 1.19 or later
- PostgreSQL 12 or later
- Git Bash (for Windows)
- yq (YAML processor)

### Initial Setup
1. Clone the repository:
```bash
git clone <repository-url>
cd TaskManager
```

2. Make scripts executable:
```bash
chmod +x start.sh
chmod +x deployment/scripts/setup_db.sh
```

3. Configure the application:
Edit `config.yaml`:
```yaml
server:
  port: 8080
database:
  host: localhost
  port: 5432
  username: postgres
  password: your_password
  dbname: taskmanager
  sslmode: disable
```

### Running the Application
1. Start the application:
```bash
./start.sh
```
This will:
- Run database setup scripts
- Create necessary schema and tables
- Start the Go application

2. Import tasks using the API:
```bash
POST http://localhost:8080/api/commands/import
```

## Development Guidelines

### Adding New Features
1. Command Operations:
   - Add new command service in Services/CommandServices
   - Implement corresponding repository methods
   - Create new command controller if needed

2. Query Operations:
   - Create new service in Services/QueryServices
   - Implement repository methods in QueryRepository
   - Add query controller in QueryRequest

## Testing
```bash
# Run unit tests
go test ./...

# Test database setup
./deployment/scripts/setup_db.sh
```

## Contributing
1. Follow the established CQRS pattern
2. Keep database initialization separate from application code
3. Add appropriate tests for new features
4. Update documentation as needed

## Troubleshooting
- Check PostgreSQL connection settings in config.yaml
- Ensure all required SQL scripts exist
- Verify file permissions for shell scripts
- Check logs for detailed error messages