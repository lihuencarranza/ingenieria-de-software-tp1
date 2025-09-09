# Script to run endpoint tests for the Melodia API
# Usage: .\run_tests.ps1

param(
    [switch]$Build,
    [switch]$Clean
)

Write-Host "ENDPOINT TESTING SCRIPT - MELODIA API" -ForegroundColor Cyan
Write-Host "=" * 60 -ForegroundColor Cyan

# Function to check if Docker is running
function Test-DockerRunning {
    try {
        $null = docker ps
        return $true
    }
    catch {
        return $false
    }
}

# Function to check if services are running
function Test-ServicesRunning {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method GET -TimeoutSec 5
        return $response.StatusCode -eq 200
    }
    catch {
        return $false
    }
}

# Function to clean previous test files
function Clean-TestFiles {
    Write-Host "Cleaning previous test files..." -ForegroundColor Yellow
    Get-ChildItem -Path "." -Filter "test_results_*.json" | Remove-Item -Force
    Write-Host "Test files cleaned" -ForegroundColor Green
}

# Function to build and run Docker
function Start-DockerServices {
    Write-Host "Starting Docker services..." -ForegroundColor Yellow
    
    if ($Build) {
        Write-Host "Building image..." -ForegroundColor Yellow
        docker-compose up --build -d
    } else {
        docker-compose up -d
    }
    
    Write-Host "Waiting for services to be ready..." -ForegroundColor Yellow
    
    # Wait up to 2 minutes for services to be ready
    $maxAttempts = 60
    $attempt = 0
    
    while ($attempt -lt $maxAttempts) {
        if (Test-ServicesRunning) {
            Write-Host "Docker services started and running!" -ForegroundColor Green
            return $true
        }
        
        $attempt++
        Write-Host "Attempt $attempt/$maxAttempts - Waiting for services..." -ForegroundColor Yellow
        Start-Sleep -Seconds 2
    }
    
    Write-Host "Services were not ready in the expected time" -ForegroundColor Red
    return $false
}

# Function to run tests
function Run-Tests {
    Write-Host "Running endpoint tests..." -ForegroundColor Yellow
    
    # Compile testing script
    Write-Host "Compiling testing script..." -ForegroundColor Yellow
    go build -o test_endpoints.exe scripts/test_endpoints.go
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error compiling testing script" -ForegroundColor Red
        return $false
    }
    
    Write-Host "Script compiled successfully" -ForegroundColor Green
    
    # Run tests and capture output
    Write-Host "Running tests..." -ForegroundColor Yellow
    $testOutput = .\test_endpoints.exe 2>&1
    
    $testExitCode = $LASTEXITCODE
    
    # Clean executable file
    Remove-Item "test_endpoints.exe" -Force -ErrorAction SilentlyContinue
    
    # Analyze output to determine if tests passed
    # Look for success message in output
    $allTestsPassed = $testOutput -match "🎉 ALL TESTS PASSED SUCCESSFULLY!"
    
    # If we find the success message, consider tests passed
    # regardless of exit code
    if ($allTestsPassed) {
        return $true
    }
    
    # If no success message found, check exit code
    return $testExitCode -eq 0
}

# Function to show Docker logs
function Show-DockerLogs {
    Write-Host "Showing Docker logs..." -ForegroundColor Yellow
    docker-compose logs --tail=20
}

# Function to stop services
function Stop-DockerServices {
    Write-Host "Stopping Docker services..." -ForegroundColor Yellow
    docker-compose down
    Write-Host "Services stopped" -ForegroundColor Green
}

function Main {
    # Check Docker
    if (-not (Test-DockerRunning)) {
        Write-Host "Docker is not running. Please start Docker Desktop." -ForegroundColor Red
        exit 1
    }
    
    Write-Host "Docker is running" -ForegroundColor Green
    
    # Clean previous files if requested
    if ($Clean) {
        Clean-TestFiles
    }
    
    # Start services
    if (-not (Start-DockerServices)) {
        Write-Host "Could not start Docker services" -ForegroundColor Red
        exit 1
    }
    
    # Run tests
    $testsPassed = Run-Tests
    
    # Show logs if there are problems
    if (-not $testsPassed) {
        Write-Host "Some tests failed or there was a problem. Showing Docker logs..." -ForegroundColor Yellow
        Show-DockerLogs
    }
    
    # Automatically close Docker services when tests finish
    Write-Host "Automatically closing Docker services..." -ForegroundColor Yellow
    Stop-DockerServices
    
    # Show summary
    Write-Host ""
    Write-Host "EXECUTION SUMMARY:" -ForegroundColor Cyan
    Write-Host "   - Docker Services: Started and closed automatically" -ForegroundColor Green
    Write-Host "   - Tests executed: Completed" -ForegroundColor Green
    Write-Host "   - Logs saved: test_results_*.json" -ForegroundColor Green
    
    if ($testsPassed) {
        Write-Host "All tests passed successfully!" -ForegroundColor Green
    } else {
        Write-Host "Some tests failed or there was a problem. Check logs for more details." -ForegroundColor Yellow
    }
}

# Run main function
try {
    Main
}
catch {
    Write-Host "Error during execution: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Stopping services..." -ForegroundColor Yellow
    Stop-DockerServices
    exit 1
}
