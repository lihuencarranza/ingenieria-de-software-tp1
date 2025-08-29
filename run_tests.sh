#!/bin/bash

# Script para ejecutar tests de endpoints de la API Melodía
# Uso: ./run_tests.sh [--build] [--clean]

BUILD=false
CLEAN=false

# Parsear argumentos
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
            echo "Uso: $0 [--build] [--clean]"
            exit 1
            ;;
    esac
done

echo "SCRIPT DE TESTING DE ENDPOINTS - API MELODIA"
echo "============================================================"

# Función para verificar si Docker está corriendo
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        echo "Docker no está corriendo. Por favor, inicia Docker."
        exit 1
    fi
    echo "Docker está corriendo"
}

# Función para verificar si los servicios están corriendo
check_services() {
    local max_attempts=60
    local attempt=0
    
    echo "Esperando a que los servicios estén listos..."
    
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            echo "Servicios Docker iniciados y funcionando!"
            return 0
        fi
        
        attempt=$((attempt + 1))
        echo "Intento $attempt/$max_attempts - Esperando servicios..."
        sleep 2
    done
    
    echo "Los servicios no estuvieron listos en el tiempo esperado"
    return 1
}

# Función para limpiar archivos de test anteriores
clean_test_files() {
    echo "Limpiando archivos de test anteriores..."
    rm -f test_results_*.json
    echo "Archivos de test limpiados"
}

# Función para construir y ejecutar Docker
start_docker_services() {
    echo "Iniciando servicios Docker..."
    
    if [ "$BUILD" = true ]; then
        echo "Construyendo imagen..."
        docker-compose up --build -d
    else
        docker-compose up -d
    fi
    
    if ! check_services; then
        return 1
    fi
    
    return 0
}

# Función para ejecutar los tests
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
    
    return $test_exit_code
}

# Función para mostrar logs de Docker
show_docker_logs() {
    echo "Mostrando logs de Docker..."
    docker-compose logs --tail=20
}

# Función para detener servicios
stop_docker_services() {
    echo "Deteniendo servicios Docker..."
    docker-compose down
    echo "Servicios detenidos"
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
    
    # Ejecutar tests
    if ! run_tests; then
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
    
    if [ $test_exit_code -eq 0 ]; then
        echo "¡Todos los tests pasaron exitosamente!"
    else
        echo "Algunos tests fallaron. Revisa los logs para más detalles."
    fi
}

# Ejecutar función principal
trap 'echo "Deteniendo servicios..."; stop_docker_services; exit 1' INT TERM
main
