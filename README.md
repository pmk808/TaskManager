# Task Manager API

## Project Structure
```
TaskManager/
├── cmd/
│   ├── main.go                 # Application entry point
│   └── swagger_init.go         # Swagger documentation initialization
├── docs/                       # Swagger generated documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── deployment/                 # Deployment configurations
│   ├── scripts/
│   │   └── setup_db.sh        # Database setup script
│   └── sql/
│       ├── 01_create_database.sql
│       ├── 02_create_schema.sql
│       ├── 03_create_task_table.sql
│       ├── 04_create_status_table.sql
│       └── 05_insert_dummy_data.sql
├── import/                     # CSV import directory
│   └── dummy_tasks.csv
├── Repository/                 # Data access layer
│   ├── CommandRepository/
│   │   ├── interfaces/
│   │   │   └── repository.go
│   │   └── TaskCommandRepository.go
│   └── QueryRepository/
│       ├── interfaces/
│       │   └── repository.go
│       └── TaskQueryRepository.go
├── RequestControllers/         # API Controllers
│   ├── AuthRequest/
│   │   ├── dto/
│   │   │   └── auth.go
│   │   ├── interfaces/
│   │   │   └── controller.go
│   │   └── AuthController.go
│   ├── CommandRequest/
│   │   ├── interfaces/
│   │   │   └── controller.go
│   │   └── CommandApiController.go
│   ├── QueryRequest/
│   │   ├── interfaces/
│   │   │   └── controller.go
│   │   ├── dto/
│   │   │   └── request.go
│   │   └── QueryApiController.go
│   └── httpSetup/
│       ├── config/
│       │   └── setup.go
│       ├── jwt/
│       │   └── jwt.go
│       ├── logger/
│       │   └── setup.go
│       ├── middleware/
│       │   └── jwt_middleware.go
│       └── setup.go
├── Services/                   # Business logic layer
│   ├── CommandServices/
│   │   └── ImportTaskService/
│   │       ├── interfaces/
│   │       │   ├── repository.go
│   │       │   ├── service.go
│   │       │   └── validator.go
│   │       ├── schemas/
│   │       │   ├── dtos.go
│   │       │   ├── models.go
│   │       │   ├── task.go
│   │       ├── validation/
│   │       │   ├── interfaces/
│   │       │   │   └── validator.go
│   │       │   ├── dataValidator_test.go
│   │       │   └── dataValidator.go
│   │       └── ImportTaskService.go
│   └── QueryServices/
│       └── TaskQueryService/
│           ├── interfaces/
│           │   └── service.go
│           ├── validation/
│           │   └── QueryValidator.go
│           └── TaskQueryService.go
├── config.yaml                 # Application configuration
├── go.mod                     # Go module file
├── go.sum                     # Go dependencies
├── start.sh                   # Application startup script
└── README.md                  # Project documentation
```

## Features

- CQRS implementation separating read and write operations
- JWT authentication for secure API access
- Swagger documentation
- PostgreSQL database with proper schema management
- CSV data import functionality
- Client-based task filtering
- Task status history tracking

## Endpoints

### Command Endpoints
- `POST /api/commands/import`: Import tasks from CSV file

### Query Endpoints
- `GET /api/queries/tasks/active`: Get active tasks for a client
- `GET /api/queries/tasks/history`: Get task status history for a client

### Auth Endpoints
- `POST /api/auth/token`: Generate JWT token for authentication

## Requirements

- Go 1.19 or higher
- PostgreSQL 12 or higher
- Git

## Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd TaskManager
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure the application:
- Copy `config.yaml.example` to `config.yaml`
- Update database credentials and other settings

4. Start the application:
```bash
./start.sh
```

## Development

### Adding New Features
1. Create a new feature branch:
```bash
git checkout -b feature/your-feature-name
```

2. Make changes and commit:
```bash
git add .
git commit -m "feat: your feature description"
```

3. Push changes and create pull request:
```bash
git push origin feature/your-feature-name
```

### Database Migrations
- Add new SQL scripts in `deployment/sql/`
- Update `setup_db.sh` if needed
- Test migrations locally before committing

## Testing

### Running Tests
```bash
go test ./...
```

### API Testing
Use Swagger UI at: `http://localhost:8080/swagger/index.html`

## Configuration

### Database Configuration
```yaml
database:
  host: localhost
  port: 5432
  username: postgres
  password: your_password
  dbname: taskmanager
  sslmode: disable
```

### JWT Configuration
```yaml
jwt:
  secret_key: your_secret_key
  expiry_hours: 24
```

## Project Structure Details

### Command Pattern
- Commands represent actions that modify state
- Each command has its own service and handler
- Validation occurs before processing

### Query Pattern
- Queries retrieve data without modifications
- Optimized for read operations
- Client-based filtering

### Authentication
- JWT-based authentication
- Token generation and validation
- Secure parameter handling
