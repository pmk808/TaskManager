#!/bin/bash

# Get the script's directory path
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SQL_DIR="$SCRIPT_DIR/../sql"

# Database connection parameters
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="postgres"
DB_PASSWORD="peemak"
DB_NAME="taskmanager"

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored messages
log_message() {
    local type=$1
    local message=$2
    case $type in
        "info")
            echo -e "${GREEN}[INFO]${NC} $message"
            ;;
        "error")
            echo -e "${RED}[ERROR]${NC} $message"
            ;;
        "warning")
            echo -e "${YELLOW}[WARNING]${NC} $message"
            ;;
    esac
}

# Function to execute SQL files
execute_sql_file() {
    local file=$1
    local dbname=$2
    local message=$3
    
    log_message "info" "$message"
    
    if [ ! -f "$file" ]; then
        log_message "error" "SQL file not found: $file"
        exit 1
    fi
    
    if [ -n "$dbname" ]; then
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $dbname -f "$file"
    else
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -f "$file"
    fi
    
    if [ $? -eq 0 ]; then
        log_message "info" "Successfully executed $file"
    else
        log_message "error" "Failed to execute $file"
        exit 1
    fi
}

# Main execution
log_message "info" "Starting database setup..."

# Verify SQL files exist
log_message "info" "Checking SQL files..."
for sql_file in "$SQL_DIR"/*.sql; do
    if [ ! -f "$sql_file" ]; then
        log_message "error" "SQL file not found: $sql_file"
        exit 1
    else
        log_message "info" "Found SQL file: $sql_file"
    fi
done

# Create database
execute_sql_file "$SQL_DIR/01_create_database.sql" "" "Creating database..."

# Verify database exists
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1" >/dev/null 2>&1
if [ $? -eq 0 ]; then
    log_message "info" "Database created successfully"
    
    # Create schema
    execute_sql_file "$SQL_DIR/02_create_schema.sql" $DB_NAME "Creating schema..."
    
    # Create tables
    execute_sql_file "$SQL_DIR/03_create_table.sql" $DB_NAME "Creating tables..."
    
    log_message "info" "Database setup completed successfully!"
else
    log_message "error" "Failed to create database"
    exit 1
fi