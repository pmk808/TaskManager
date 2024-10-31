#!/bin/bash

# Color codes for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Base URL
BASE_URL="http://localhost:8080/api"

# Test data
CLIENT_NAME="Client One Corp"
CLIENT_ID="123e4567-e89b-12d3-a456-426614174000" 

# Function to make requests
test_endpoint() {
    local endpoint=$1
    local payload=$2
    local description=$3

    echo -e "\n${GREEN}Testing: ${description}${NC}"
    echo "Endpoint: ${endpoint}"
    echo "Payload: ${payload}"
    
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "${payload}" \
        "${BASE_URL}${endpoint}")
    
    echo -e "Response:\n${response}\n"
}

# Test active tasks endpoint
test_endpoint "/queries/tasks/active" \
    "{\"client_name\":\"${CLIENT_NAME}\",\"client_id\":\"${CLIENT_ID}\"}" \
    "Get Active Tasks"

# Test task history endpoint
test_endpoint "/queries/tasks/history" \
    "{\"client_name\":\"${CLIENT_NAME}\",\"client_id\":\"${CLIENT_ID}\"}" \
    "Get Task History"

# Test with invalid client ID
test_endpoint "/queries/tasks/active" \
    "{\"client_name\":\"${CLIENT_NAME}\",\"client_id\":\"invalid-uuid\"}" \
    "Get Active Tasks (Invalid UUID)"

# Test with empty client name
test_endpoint "/queries/tasks/active" \
    "{\"client_name\":\"\",\"client_id\":\"${CLIENT_ID}\"}" \
    "Get Active Tasks (Empty Client Name)"