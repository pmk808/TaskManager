# TaskManager
 
## Overview

Task Manager Import Service is a Go-based application that provides a REST API endpoint for importing task data from Excel spreadsheets into a PostgreSQL database. The service includes data validation, error handling, and logging features to ensure data integrity and provide visibility into the import process.

## Features

- POST endpoint for triggering data import
- Excel spreadsheet parsing (XLSX format)
- Data validation before database insertion
- Bulk insert of validated data into PostgreSQL
- Structured logging for monitoring and debugging
- Swagger documentation for API endpoints

## Project Structure

```
taskmanager/
├── cmd/
│   └── main.go           # Application entry point
├── config/
│   └── config.go         # Configuration management
├── handlers/
│   └── task_handler.go   # HTTP request handlers
├── interfaces/
│   ├── repository.go     # Repository interface
│   ├── service.go        # Service interface
│   └── validator.go      # Validator interface
├── repository/
│   └── postgres_repository.go  # PostgreSQL implementation
├── schemas/
│   └── task.go           # Task data structure
├── services/
│   └── task_service.go   # Business logic implementation
├── validation/
│   └── validator.go      # Data validation logic
├── docs/                 # Swagger documentation
├── import/               # Directory for Excel files to import
└── config.yaml           # Application configuration file
```

## Prerequisites

- Go 1.16 or later
- PostgreSQL 12 or later
- Excel files (.xlsx) containing task data

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/taskmanager.git
   ```

2. Navigate to the project directory:
   ```
   cd taskmanager
   ```

3. Install dependencies:
   ```
   go mod tidy
   ```

4. Set up your PostgreSQL database and update the `config.yaml` file with your database credentials.

## Configuration

The application uses a `config.yaml` file for configuration. Ensure this file is properly set up with your database connection details and other necessary configurations.

## Running the Application

1. Start the application:
   ```
   go run cmd/main.go
   ```

2. The server will start on the configured port (default is typically 8080).

3. Access the Swagger UI for API documentation:
   ```
   http://localhost:8080/swagger/index.html
   ```

## API Endpoints

### Import Tasks

- **URL**: `/import`
- **Method**: `POST`
- **Description**: Triggers the import process for task data from Excel files in the configured directory.
- **Success Response**: 
  - **Code**: 200
  - **Content**: `{ "message": "Data import completed successfully" }`
- **Error Response**:
  - **Code**: 500
  - **Content**: `{ "error": "Failed to import data" }`

## Data Import Process

1. The service reads the latest Excel file from the configured directory.
2. It parses the data from the spreadsheet into Task structures.
3. Each task is validated to ensure all required fields are present and correctly formatted.
4. If validation passes, the tasks are bulk inserted into the PostgreSQL database.
5. The process is transactional - if any task fails validation or insertion, no data is committed.

## Logging

The application uses structured logging with logrus. Logs include information about:
- Start and completion of the import process
- Validation successes and failures
- Database operations
- Any errors encountered during the process

## Error Handling

The service includes comprehensive error handling:
- File reading errors (e.g., file not found, incorrect format)
- Data validation errors
- Database connection and operation errors

All errors are logged with relevant context for easier debugging.


```
