#!/bin/bash

# Get the script's directory path
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$SCRIPT_DIR/../.."
SQL_DIR="$SCRIPT_DIR/../sql"
CONFIG_FILE="$PROJECT_ROOT/config.yaml"

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

# Function to load configuration
load_config() {
    if [ ! -f "$CONFIG_FILE" ]; then
        log_message "error" "Config file not found: $CONFIG_FILE"
        exit 1
    fi

    # Load database configuration
    DB_HOST=$(yq eval '.database.host' "$CONFIG_FILE")
    DB_PORT=$(yq eval '.database.port' "$CONFIG_FILE")
    DB_USER=$(yq eval '.database.username' "$CONFIG_FILE")
    DB_PASSWORD=$(yq eval '.database.password' "$CONFIG_FILE")
    DB_NAME=$(yq eval '.database.dbname' "$CONFIG_FILE")
    DB_SSLMODE=$(yq eval '.database.sslmode' "$CONFIG_FILE")

    # Load logging configuration
    LOG_LEVEL=$(yq eval '.logging.level' "$CONFIG_FILE")
    LOG_FORMAT=$(yq eval '.logging.format' "$CONFIG_FILE")

    # Verify configuration
    if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_NAME" ]; then
        log_message "error" "Missing required database configuration"
        exit 1
    fi

    log_message "info" "Configuration loaded successfully"
    log_message "info" "Using database: $DB_NAME on $DB_HOST:$DB_PORT"
    log_message "info" "Logging level: $LOG_LEVEL, format: $LOG_FORMAT"
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
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d "$dbname" -f "$file"
    else
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -f "$file"
    fi
    
    local result=$?
    if [ $result -eq 0 ]; then
        log_message "info" "Successfully executed $file"
        return 0
    else
        log_message "error" "Failed to execute $file"
        return 1
    fi
}

# Main execution
main() {
    log_message "info" "Starting database setup..."

    # Load configuration
    load_config

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
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d "$DB_NAME" -c "SELECT 1" >/dev/null 2>&1
    if [ $? -eq 0 ]; then
        log_message "info" "Database created successfully"
        
        # Create schema
        execute_sql_file "$SQL_DIR/02_create_schema.sql" "$DB_NAME" "Creating schema..."
        
        # Create tables
        execute_sql_file "$SQL_DIR/03_create_table.sql" "$DB_NAME" "Creating tables..."
        
        log_message "info" "Database setup completed successfully!"
    else
        log_message "error" "Failed to create database"
        exit 1
    fi
}

# Run main function
main