# Script para ejecutar tests de endpoints de la API Melodía
# Uso: .\run_tests.ps1

param(
    [switch]$Build,
    [switch]$Clean
)

Write-Host "SCRIPT DE TESTING DE ENDPOINTS - API MELODIA" -ForegroundColor Cyan
Write-Host "=" * 60 -ForegroundColor Cyan

# Función para verificar si Docker está corriendo
function Test-DockerRunning {
    try {
        $null = docker ps
        return $true
    }
    catch {
        return $false
    }
}

# Función para verificar si los servicios están corriendo
function Test-ServicesRunning {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method GET -TimeoutSec 5
        return $response.StatusCode -eq 200
    }
    catch {
        return $false
    }
}

# Función para limpiar archivos de test anteriores
function Clean-TestFiles {
    Write-Host "Limpiando archivos de test anteriores..." -ForegroundColor Yellow
    Get-ChildItem -Path "." -Filter "test_results_*.json" | Remove-Item -Force
    Write-Host "Archivos de test limpiados" -ForegroundColor Green
}

# Función para construir y ejecutar Docker
function Start-DockerServices {
    Write-Host "Iniciando servicios Docker..." -ForegroundColor Yellow
    
    if ($Build) {
        Write-Host "Construyendo imagen..." -ForegroundColor Yellow
        docker-compose up --build -d
    } else {
        docker-compose up -d
    }
    
    Write-Host "Esperando a que los servicios estén listos..." -ForegroundColor Yellow
    
    # Esperar hasta 2 minutos para que los servicios estén listos
    $maxAttempts = 60
    $attempt = 0
    
    while ($attempt -lt $maxAttempts) {
        if (Test-ServicesRunning) {
            Write-Host "Servicios Docker iniciados y funcionando!" -ForegroundColor Green
            return $true
        }
        
        $attempt++
        Write-Host "Intento $attempt/$maxAttempts - Esperando servicios..." -ForegroundColor Yellow
        Start-Sleep -Seconds 2
    }
    
    Write-Host "Los servicios no estuvieron listos en el tiempo esperado" -ForegroundColor Red
    return $false
}

# Función para ejecutar los tests
function Run-Tests {
    Write-Host "Ejecutando tests de endpoints..." -ForegroundColor Yellow
    
    # Compilar el script de testing
    Write-Host "Compilando script de testing..." -ForegroundColor Yellow
    go build -o test_endpoints.exe scripts/test_endpoints.go
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error compilando script de testing" -ForegroundColor Red
        return $false
    }
    
    Write-Host "Script compilado exitosamente" -ForegroundColor Green
    
    # Ejecutar los tests
    Write-Host "Ejecutando tests..." -ForegroundColor Yellow
    .\test_endpoints.exe
    
    $testExitCode = $LASTEXITCODE
    
    # Limpiar archivo ejecutable
    Remove-Item "test_endpoints.exe" -Force -ErrorAction SilentlyContinue
    
    return $testExitCode -eq 0
}

# Función para mostrar logs de Docker
function Show-DockerLogs {
    Write-Host "Mostrando logs de Docker..." -ForegroundColor Yellow
    docker-compose logs --tail=20
}

# Función para detener servicios
function Stop-DockerServices {
    Write-Host "Deteniendo servicios Docker..." -ForegroundColor Yellow
    docker-compose down
    Write-Host "Servicios detenidos" -ForegroundColor Green
}

# Función principal
function Main {
    # Verificar Docker
    if (-not (Test-DockerRunning)) {
        Write-Host "Docker no está corriendo. Por favor, inicia Docker Desktop." -ForegroundColor Red
        exit 1
    }
    
    Write-Host "Docker está corriendo" -ForegroundColor Green
    
    # Limpiar archivos anteriores si se solicita
    if ($Clean) {
        Clean-TestFiles
    }
    
    # Iniciar servicios
    if (-not (Start-DockerServices)) {
        Write-Host "No se pudieron iniciar los servicios Docker" -ForegroundColor Red
        exit 1
    }
    
    # Ejecutar tests
    $testsPassed = Run-Tests
    
    # Mostrar logs si hay problemas
    if (-not $testsPassed) {
        Write-Host "Algunos tests fallaron. Mostrando logs de Docker..." -ForegroundColor Yellow
        Show-DockerLogs
    }
    
    # Preguntar si mantener servicios corriendo
    Write-Host ""
    $keepRunning = Read-Host "¿Mantener los servicios Docker corriendo? (s/n)"
    
    if ($keepRunning -notmatch "^[Ss]") {
        Stop-DockerServices
    } else {
        Write-Host "Servicios Docker siguen corriendo en background" -ForegroundColor Green
        Write-Host "   Para detenerlos manualmente: docker-compose down" -ForegroundColor Cyan
    }
    
    # Mostrar resumen
    Write-Host ""
    Write-Host "RESUMEN DE LA EJECUCION:" -ForegroundColor Cyan
    Write-Host "   - Servicios Docker: Iniciados" -ForegroundColor Green
    Write-Host "   - Tests ejecutados: Completados" -ForegroundColor Green
    Write-Host "   - Logs guardados: test_results_*.json" -ForegroundColor Green
    
    if ($testsPassed) {
        Write-Host "¡Todos los tests pasaron exitosamente!" -ForegroundColor Green
    } else {
        Write-Host "Algunos tests fallaron. Revisa los logs para más detalles." -ForegroundColor Yellow
    }
}

# Ejecutar función principal
try {
    Main
}
catch {
    Write-Host "Error durante la ejecución: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Deteniendo servicios..." -ForegroundColor Yellow
    Stop-DockerServices
    exit 1
}
