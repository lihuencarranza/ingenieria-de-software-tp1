#!/bin/bash

# Script to run endpoint tests for the Melodia API
# Usage: ./run_tests.sh [--build] [--clean]

BUILD=false
CLEAN=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --build)
            BUILD=true
            shift
            ;;
        --clean)
            CLEAN=true
            shift
            ;;
        *)
            echo "Usage: $0 [--build] [--clean]"
            exit 1
            ;;
    esac
done

echo "ENDPOINT TESTING SCRIPT - MELODIA API"
echo "============================================================"

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        echo "Docker is not running. Please start Docker."
        exit 1
    fi
    echo "Docker is running"
}

# Function to check if services are running
check_services() {
    local max_attempts=60
    local attempt=0
    
    echo "Waiting for services to be ready..."
    
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            echo "Docker services started and running!"
            return 0
        fi
        
        attempt=$((attempt + 1))
        echo "Attempt $attempt/$max_attempts - Waiting for services..."
        sleep 2
    done
    
    echo "Services were not ready in the expected time"
    return 1
}

# Function to clean previous test files
clean_test_files() {
    echo "Cleaning previous test files..."
    rm -f test_results_*.json
    echo "Test files cleaned"
}

# Function to build and run Docker
start_docker_services() {
    echo "Starting Docker services..."
    
    if [ "$BUILD" = true ]; then
        echo "Building image..."
        docker-compose up --build -d
    else
        docker-compose up -d
    fi
    
    if ! check_services; then
        return 1
    fi
    
    return 0
}

# Function to run tests
run_tests() {
    echo "Ejecutando tests de endpoints..."
    
    # Compilar el script de testing
    echo "Compilando script de testing..."
    if ! go build -o test_endpoints scripts/test_endpoints.go; then
        echo "Error compilando script de testing"
        return 1
    fi
    
    echo "Script compilado exitosamente"
    
    # Ejecutar los tests
    echo "Ejecutando tests..."
    ./test_endpoints
    local test_exit_code=$?
    
    # Limpiar archivo ejecutable
    rm -f test_endpoints
    
    # Retornar el código de salida de los tests
    return $test_exit_code
}

# Function to show Docker logs
show_docker_logs() {
    echo "Showing Docker logs..."
    docker-compose logs --tail=20
}

# Function to stop services
stop_docker_services() {
    echo "Stopping Docker services..."
    docker-compose down
    echo "Services stopped"
}

main() {
    # Check Docker
    check_docker
    
    # Clean previous files if requested
    if [ "$CLEAN" = true ]; then
        clean_test_files
    fi
    
    # Start services
    if ! start_docker_services; then
        echo "Could not start Docker services"
        exit 1
    fi
    
    # Run tests and show logs if there are problems
    local test_result=0
    if ! run_tests; then
        test_result=$?
        echo "Some tests failed or there was a problem. Showing Docker logs..."
        show_docker_logs
    fi
    
    # Automatically close Docker services when tests finish
    echo "Cerrando servicios Docker automáticamente..."
    stop_docker_services
    
    # Show summary
    echo ""
    echo "EXECUTION SUMMARY:"
    echo "   - Docker Services: Started and closed automatically"
    echo "   - Tests executed: Completed"
    echo "   - Logs saved: test_results_*.json"
    
    if [ $test_result -eq 0 ]; then
        echo "All tests passed successfully!"
    else
        echo "Some tests failed or there was a problem. Check logs for more details."
    fi
}

# Run main function
trap 'echo "Stopping services..."; stop_docker_services; exit 1' INT TERM
main
