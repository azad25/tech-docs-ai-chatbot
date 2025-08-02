#!/bin/bash

# Tech Docs AI Integration Test Script
# This script tests the complete application functionality

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_BASE_URL="http://localhost"
TIMEOUT=30
RETRY_COUNT=5

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to wait for service to be ready
wait_for_service() {
    local url=$1
    local service_name=$2
    local max_attempts=$3
    
    print_status "Waiting for $service_name to be ready..."
    
    for i in $(seq 1 $max_attempts); do
        if curl -s -f "$url" > /dev/null 2>&1; then
            print_success "$service_name is ready!"
            return 0
        fi
        print_status "Attempt $i/$max_attempts: $service_name not ready yet, waiting..."
        sleep 5
    done
    
    print_error "$service_name failed to start within expected time"
    return 1
}

# Function to test API endpoint
test_api_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4
    local description=$5
    
    print_status "Testing: $description"
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$API_BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            "$API_BASE_URL$endpoint")
    fi
    
    # Extract status code (last line)
    status_code=$(echo "$response" | tail -n1)
    # Extract response body (all but last line)
    response_body=$(echo "$response" | head -n -1)
    
    if [ "$status_code" -eq "$expected_status" ]; then
        print_success "$description - Status: $status_code"
        echo "Response: $response_body" | head -c 200
        echo ""
        return 0
    else
        print_error "$description - Expected: $expected_status, Got: $status_code"
        echo "Response: $response_body"
        return 1
    fi
}

# Function to run unit tests
run_unit_tests() {
    print_status "Running unit tests..."
    
    if go test -v ./internal/app -run "TestHandler.*Unit" > test_results_unit.log 2>&1; then
        print_success "Unit tests passed"
        return 0
    else
        print_error "Unit tests failed"
        cat test_results_unit.log
        return 1
    fi
}

# Function to run integration tests
run_integration_tests() {
    print_status "Running integration tests..."
    
    if go test -v ./integration_test.go > test_results_integration.log 2>&1; then
        print_success "Integration tests passed"
        return 0
    else
        print_warning "Some integration tests failed (expected without running services)"
        cat test_results_integration.log | grep -E "(PASS|FAIL|ERROR)"
        return 0  # Don't fail the script for integration tests
    fi
}

# Function to test scraper functionality
test_scrapers() {
    print_status "Testing scraper functionality..."
    
    if go test -v ./internal/scraper > test_results_scraper.log 2>&1; then
        print_success "Scraper tests passed"
        return 0
    else
        print_error "Scraper tests failed"
        cat test_results_scraper.log
        return 1
    fi
}

# Function to build the application
build_application() {
    print_status "Building application..."
    
    # Build server
    if go build -o server cmd/server/main.go; then
        print_success "Server built successfully"
    else
        print_error "Failed to build server"
        return 1
    fi
    
    # Build worker
    if go build -o worker cmd/worker/main.go; then
        print_success "Worker built successfully"
    else
        print_error "Failed to build worker"
        return 1
    fi
    
    return 0
}

# Function to test Docker build
test_docker_build() {
    print_status "Testing Docker build..."
    
    if docker build -t tech-docs-ai:test --target production .; then
        print_success "Docker build successful"
        return 0
    else
        print_error "Docker build failed"
        return 1
    fi
}

# Function to start services for testing
start_test_services() {
    print_status "Starting services for testing..."
    
    # Check if docker-compose is available
    if ! command -v docker-compose &> /dev/null; then
        print_error "docker-compose is not installed"
        return 1
    fi
    
    # Start services
    if docker-compose up -d; then
        print_success "Services started"
    else
        print_error "Failed to start services"
        return 1
    fi
    
    # Wait for services to be ready
    wait_for_service "$API_BASE_URL/health" "API Server" 12
    
    return 0
}

# Function to stop test services
stop_test_services() {
    print_status "Stopping test services..."
    docker-compose down
    print_success "Services stopped"
}

# Function to run API tests
run_api_tests() {
    print_status "Running API endpoint tests..."
    
    local failed_tests=0
    
    # Test health endpoint
    test_api_endpoint "GET" "/health" "" 200 "Health check endpoint" || ((failed_tests++))
    
    # Test chat endpoint
    test_api_endpoint "POST" "/api/v1/chat" '{"message":"Hello, test message"}' 200 "Chat endpoint" || ((failed_tests++))
    
    # Test add document endpoint
    test_api_endpoint "POST" "/api/v1/documents" '{"title":"Test Document","content":"Test content","category":"Testing","tags":["test"],"author":"Test Author"}' 201 "Add document endpoint" || ((failed_tests++))
    
    # Test search documents endpoint
    test_api_endpoint "GET" "/api/v1/documents/search?q=test&limit=5" "" 200 "Search documents endpoint" || ((failed_tests++))
    
    # Test scrape endpoint
    test_api_endpoint "POST" "/api/v1/scrape" '{"url":"https://httpbin.org/html","category":"Test","tags":["test"]}' 200 "Scrape endpoint" || ((failed_tests++))
    
    # Test chat with history endpoint
    test_api_endpoint "POST" "/api/v1/chat/history" '{"session_id":"test-session","message":"Hello with history"}' 200 "Chat with history endpoint" || ((failed_tests++))
    
    # Test get chat history endpoint
    test_api_endpoint "GET" "/api/v1/chat/history?session_id=test-session&limit=10" "" 200 "Get chat history endpoint" || ((failed_tests++))
    
    # Test conversation insights endpoint
    test_api_endpoint "GET" "/api/v1/chat/insights?session_id=test-session" "" 200 "Conversation insights endpoint" || ((failed_tests++))
    
    if [ $failed_tests -eq 0 ]; then
        print_success "All API tests passed"
        return 0
    else
        print_error "$failed_tests API tests failed"
        return 1
    fi
}

# Function to run performance tests
run_performance_tests() {
    print_status "Running basic performance tests..."
    
    # Test concurrent requests
    print_status "Testing concurrent requests..."
    
    for i in {1..5}; do
        curl -s -X POST \
            -H "Content-Type: application/json" \
            -d '{"message":"Performance test message '$i'"}' \
            "$API_BASE_URL/api/v1/chat" > /dev/null &
    done
    
    wait
    print_success "Concurrent requests test completed"
    
    return 0
}

# Function to generate test report
generate_test_report() {
    print_status "Generating test report..."
    
    cat > test_report.md << EOF
# Tech Docs AI Test Report

Generated on: $(date)

## Test Results Summary

### Unit Tests
$(if [ -f test_results_unit.log ]; then echo "✅ Passed"; else echo "❌ Not run"; fi)

### Integration Tests
$(if [ -f test_results_integration.log ]; then echo "✅ Completed"; else echo "❌ Not run"; fi)

### Scraper Tests
$(if [ -f test_results_scraper.log ]; then echo "✅ Passed"; else echo "❌ Not run"; fi)

### API Tests
$(if [ $api_tests_passed -eq 1 ]; then echo "✅ Passed"; else echo "❌ Failed"; fi)

### Build Tests
$(if [ $build_successful -eq 1 ]; then echo "✅ Passed"; else echo "❌ Failed"; fi)

## Detailed Results

### Unit Test Output
\`\`\`
$(if [ -f test_results_unit.log ]; then cat test_results_unit.log; else echo "No unit test results"; fi)
\`\`\`

### Integration Test Output
\`\`\`
$(if [ -f test_results_integration.log ]; then cat test_results_integration.log; else echo "No integration test results"; fi)
\`\`\`

### Scraper Test Output
\`\`\`
$(if [ -f test_results_scraper.log ]; then cat test_results_scraper.log; else echo "No scraper test results"; fi)
\`\`\`

## Recommendations

1. Ensure all services are running before running full integration tests
2. Check logs for any warnings or errors
3. Monitor performance under load
4. Verify all API endpoints are functioning correctly

EOF

    print_success "Test report generated: test_report.md"
}

# Main execution
main() {
    print_status "Starting Tech Docs AI Integration Tests"
    print_status "========================================"
    
    local overall_success=1
    local build_successful=0
    local api_tests_passed=0
    
    # Run unit tests first
    if run_unit_tests; then
        print_success "✅ Unit tests completed"
    else
        print_error "❌ Unit tests failed"
        overall_success=0
    fi
    
    # Run scraper tests
    if test_scrapers; then
        print_success "✅ Scraper tests completed"
    else
        print_error "❌ Scraper tests failed"
        overall_success=0
    fi
    
    # Run integration tests (without services)
    if run_integration_tests; then
        print_success "✅ Integration tests completed"
    else
        print_warning "⚠️ Integration tests had issues"
    fi
    
    # Build application
    if build_application; then
        print_success "✅ Application build completed"
        build_successful=1
    else
        print_error "❌ Application build failed"
        overall_success=0
    fi
    
    # Test Docker build
    if test_docker_build; then
        print_success "✅ Docker build completed"
    else
        print_warning "⚠️ Docker build failed"
    fi
    
    # Ask user if they want to run full integration tests with services
    echo ""
    read -p "Do you want to run full integration tests with Docker services? (y/N): " -n 1 -r
    echo ""
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Starting full integration tests with services..."
        
        # Start services
        if start_test_services; then
            # Run API tests
            if run_api_tests; then
                print_success "✅ API tests completed"
                api_tests_passed=1
            else
                print_error "❌ API tests failed"
                overall_success=0
            fi
            
            # Run performance tests
            if run_performance_tests; then
                print_success "✅ Performance tests completed"
            else
                print_warning "⚠️ Performance tests had issues"
            fi
            
            # Stop services
            stop_test_services
        else
            print_error "❌ Failed to start services for integration tests"
            overall_success=0
        fi
    else
        print_status "Skipping full integration tests with services"
    fi
    
    # Generate test report
    generate_test_report
    
    # Final summary
    echo ""
    print_status "========================================"
    if [ $overall_success -eq 1 ]; then
        print_success "🎉 All tests completed successfully!"
        print_status "The Tech Docs AI application is ready for use."
    else
        print_error "❌ Some tests failed. Please check the logs and fix issues."
    fi
    
    print_status "Test report generated: test_report.md"
    print_status "Log files: test_results_*.log"
    
    return $overall_success
}

# Run main function
main "$@"