#!/bin/bash

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

# Check if PostgreSQL is running
check_postgres() {
    log_message "info" "Checking PostgreSQL status..."
    
    # Check if psql is installed
    if ! command -v psql &> /dev/null; then
        log_message "error" "PostgreSQL is not installed or not in PATH"
        return 1
    fi

    # Check if PostgreSQL server is running
    if ! psql -U postgres -c "SELECT 1" &> /dev/null; then
        log_message "error" "PostgreSQL is not running or connection failed"
        return 1
    fi

    log_message "info" "PostgreSQL is running and accessible"
    return 0
}

# Main execution
main() {
    log_message "info" "Starting application setup..."

    # Check PostgreSQL
    if ! check_postgres; then
        log_message "error" "PostgreSQL check failed"
        exit 1
    fi

    # Run database setup script
    log_message "info" "Setting up database..."
    if ! bash deployment/scripts/setup_db.sh; then
        log_message "error" "Database setup failed"
        exit 1
    fi
    log_message "info" "Database setup completed"

    # Start the Go application
    log_message "info" "Starting Go application..."
    go run cmd/main.go
}

# Handle script interruption (Ctrl+C)
trap 'echo -e "\n${YELLOW}[WARNING]${NC} Script interrupted by user"; exit 1' INT

# Run main function
main