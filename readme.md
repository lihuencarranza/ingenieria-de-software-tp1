# Melodía

## Tabla de Contenido
- [Introducción](#introducción)
- [Desafíos del Proyecto](#desafíos-del-proyecto)
- [Pre-requisitos](#pre-requisitos)
- [Comandos de Desarrollo](#comandos-de-desarrollo)
- [Docker](#docker)
- [Testing](#testing)
- [Swagger](#swagger)

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


##


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
El proyecto incluye Docker Compose para ejecutar tanto la aplicación como la base de datos PostgreSQL y los tests.

### Configuración de Variables de Entorno

Las siguientes variables pueden modificarse en el archivo `.env`
- `HOST_PORT`: Puerto del host para la aplicación (default: 8080)
- `DATABASE_PORT`: Puerto del host para PostgreSQL (default: 5432)
- `DATABASE_NAME`: Nombre de la base de datos (default: melodiadb)
- `DATABASE_USER`: Usuario de la base de datos (default: melodia_admin)
- `DATABASE_PASSWORD`: Contraseña de la base de datos (default: melodia_password)

### Servicios Incluidos
- **melodia**: Servicio de la aplicación API
- **postgres**: Base de datos PostgreSQL 16.3

### Comandos Principales

#### Levantar todo el stack (aplicación + base de datos)
```bash
docker compose up --build
```

#### Parar todos los servicios
```bash
docker compose down
```

#### Levantar solo la base de datos
```bash
docker compose up postgres
```

#### Levantar solo la aplicación (requiere base de datos corriendo)
```bash
docker compose up melodia
```

### Otros Comandos Útiles
```bash
# Construir solo la imagen de la aplicación
docker compose build melodia

# Ver logs de todos los servicios
docker compose logs -f

# Ver logs de un servicio específico
docker compose logs -f melodia
docker compose logs -f postgres

# Ver estado de los servicios
docker compose ps

# Reiniciar un servicio específico
docker compose restart melodia
```


### Probar la Aplicación
- Salud: http://localhost:8080/health
- Swagger UI: http://localhost:8080/swagger/index.html
- Base de datos: localhost:5432 (usando pgAdmin o psql)

### Acceso a la Base de Datos
```bash
# Conectar desde el host usando psql
psql -h localhost -p 5432 -U melodia_admin -d melodiadb

# O conectarse al contenedor
docker compose exec postgres psql -U melodia_admin -d melodiadb
```
```


## Base de Datos
El proyecto utiliza PostgreSQL como base de datos relacional para persistir canciones y playlists.

### Características
- **Motor**: PostgreSQL 16.3
- **Persistencia**: Volumen Docker para datos persistentes
- **Red**: Red dedicada para comunicación entre contenedores
- **Configuración**: Variables de entorno personalizables

### Estructura de la Base de Datos
- **Tabla songs**: Almacena información de canciones (id, title, artist)
- **Tabla playlists**: Almacena playlists (id, name, description, isPublished, publishedAt)
- **Tabla playlist_songs**: Relación many-to-many entre playlists y canciones con timestamp de agregado

### Conexión desde la Aplicación
La aplicación se conecta automáticamente a la base de datos usando las variables de entorno:
- `DATABASE_HOST`: postgres (nombre del servicio en Docker Compose)
- `DATABASE_PORT`: 5432
- `DATABASE_NAME`: melodiadb
- `DATABASE_USER`: melodia_admin
- `DATABASE_PASSWORD`: melodia_password

### Migraciones y Esquema
El esquema de la base de datos se crea automáticamente al iniciar la aplicación por primera vez.

## Contribución

