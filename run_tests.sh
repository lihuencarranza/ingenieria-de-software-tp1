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

# Función principal
main() {
    # Verificar Docker
    check_docker
    
    # Limpiar archivos anteriores si se solicita
    if [ "$CLEAN" = true ]; then
        clean_test_files
    fi
    
    # Iniciar servicios
    if ! start_docker_services; then
        echo "No se pudieron iniciar los servicios Docker"
        exit 1
    fi
    
    # Ejecutar tests y capturar el resultado
    local test_result=0
    if ! run_tests; then
        test_result=$?
        echo "Algunos tests fallaron. Mostrando logs de Docker..."
        show_docker_logs
    fi
    
    # Cerrar automáticamente los servicios Docker al terminar los tests
    echo "Cerrando servicios Docker automáticamente..."
    stop_docker_services
    
    # Mostrar resumen
    echo ""
    echo "RESUMEN DE LA EJECUCION:"
    echo "   - Servicios Docker: Iniciados y cerrados automáticamente"
    echo "   - Tests ejecutados: Completados"
    echo "   - Logs guardados: test_results_*.json"
    
    if [ $test_result -eq 0 ]; then
        echo "¡Todos los tests pasaron exitosamente!"
    else
        echo "Algunos tests fallaron. Revisa los logs para más detalles."
    fi
}

# Run main function
trap 'echo "Stopping services..."; stop_docker_services; exit 1' INT TERM
main
