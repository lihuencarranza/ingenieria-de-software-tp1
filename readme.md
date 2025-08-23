# Melodía

## Tabla de Contenido
- [Introducción](#introducción)
- [Desafíos del Proyecto](#desafíos-del-proyecto)
- [Pre-requisitos](#pre-requisitos)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Comandos de Desarrollo](#comandos-de-desarrollo)
- [Testing](#testing)
- [Docker](#docker)
- [Base de Datos](#base-de-datos)

## Introducción


## Desafíos del Proyecto


## Pre-requisitos
- **Go**: Versión 1.25
- **Gin**: Framework web para Go
- **godotenv**: Para cargar variables de entorno desde archivos .env


### Versiones de las dependencias


## Estructura del Proyecto


## Comandos de Desarrollo

### Compilar 
```bash
go build -o bin/melodia cmd/main.go
```

### Ejecutar
 
#### En una terminal:
```bash
go run ./cmd
```




## Testing
### Ejecutar todos los tests
```bash
go test ./... 
```

### Ejecutar una carpeta en específico
```bash
go test ./cmd
```

### Ejecutar tests de un archivo en específico
```bash
go test ./cmd/main_test.go
```

### Ejecutar un test específico
```bash
go test ./cmd -run TestHealthEndpoint
```

### Tests con información detallada
```bash
# Con verbose (más información)
go test -v ./internal/models

# Con coverage (cobertura de código)
go test -cover ./internal/models

# Con coverage detallado
go test -coverprofile=coverage.out ./internal/models
go tool cover -html=coverage.out
```

## Docker
*Nota: Los archivos Docker se agregarán en las siguientes iteraciones*


## Variables de Entorno

## Base de Datos

## Contribución

