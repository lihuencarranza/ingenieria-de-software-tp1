# Melodía

## Tabla de Contenido
- [Introducción](#introducción)
- [Desafíos del Proyecto](#desafíos-del-proyecto)
- [Pre-requisitos](#pre-requisitos)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Comandos de Desarrollo](#comandos-de-desarrollo)
- [Testing](#testing)
- [Swagger](#swagger)
- [Docker](#docker)
- [Base de Datos](#base-de-datos)

## Introducción


## Desafíos del Proyecto


## Pre-requisitos
- **Go**: Versión 1.25

### Versiones de las dependencias
- **Gin**: v1.10.1 (Framework web)
- **godotenv**: v1.5.1 (Variables de entorno)
- **swaggo/swag**: v1.16.6 (Generador de Swagger)
- **swaggo/gin-swagger**: v1.6.0 (Middleware Swagger para Gin)
- **swaggo/files**: v1.0.1 (Archivos estáticos de Swagger)

## Estructura del Proyecto


## Comandos de Desarrollo

### Compilar 
```bash
go build -o bin/melodia cmd/main.go
```

### Ejecutar
```bash
go run cmd/main.go
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

## Swagger
El proyecto incluye documentación interactiva de la API usando Swagger/OpenAPI.

### Acceso a Swagger UI
Una vez que el servidor esté ejecutándose, puedes acceder a la documentación interactiva en:
```
http://localhost:8080/swagger/index.html
```

### Características de Swagger
- ✅ **Documentación interactiva** de todos los endpoints
- ✅ **Modelos de datos** completos (Song, Playlist, PlaylistSong)
- ✅ **Ejemplos de requests/responses** para cada endpoint
- ✅ **Testing directo** de la API desde la interfaz web
- ✅ **Esquemas OpenAPI 2.0** compatibles con herramientas estándar

### Archivos de Swagger
- `swagger.json` - Configuración principal en formato JSON
- `docs/swagger.yaml` - Configuración en formato YAML
- `docs/docs.go` - Código Go generado para Swagger

### Generar Documentación
Para regenerar la documentación de Swagger:
```bash
# Instalar swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Generar documentación
swag init -g cmd/main.go
```

## Docker
*Nota: Los archivos Docker se agregarán en las siguientes iteraciones*


## Variables de Entorno

## Base de Datos

## Contribución

