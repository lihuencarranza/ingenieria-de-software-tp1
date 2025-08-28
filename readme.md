# Melodía

## Tabla de Contenido
- [Introducción](#introducción)
- [Pre-requisitos](#pre-requisitos)
- [Ejecutar con Docker](#docker)
- [Testing](#testing)
- [Base de Datos](#base-de-datos)
- [Desafíos del Proyecto](#desafíos-del-proyecto)
- [Swagger](#swagger)


## Introducción
Melodía es una API REST desarrollada en Go que permite gestionar canciones y playlists. La solución implementa un sistema de gestión de música con endpoints para crear, leer, actualizar y eliminar tanto canciones como playlists, utilizando Gin como framework web y una base de datos PostgreSQL para la persistencia de datos.


## Pre-requisitos
- **Go**: Versión 1.25
- **Docker**: Versión 20.10 o superior
- **Docker Compose**: Versión 2.0 o superior
- **PostgreSQL**: v16.3 (Base de datos relacional)
- **lib/pq**: v1.10.9 (Driver PostgreSQL para Go)



### Versiones de las dependencias
- **Gin**: v1.10.1 (Framework web)
- **godotenv**: v1.5.1 (Variables de entorno)
- **swaggo/swag**: v1.16.6 (Generador de Swagger)
- **swaggo/gin-swagger**: v1.6.0 (Middleware Swagger para Gin)
- **swaggo/files**: v1.0.1 (Archivos estáticos de Swagger)
- **testify**: v1.9.0 (Framework de testing y assertions)
- **zap**: v1.27.0 (Logging estructurado de alto rendimiento)



## Ejecutar con Docker

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


## Testing

### Librerías de Testing Utilizadas

- **Go Testing Package**: Framework de testing estándar de Go
  - [User Guide](https://golang.org/pkg/testing/)
  - [Repository](https://github.com/golang/go/tree/master/src/testing)

- **Testify**: Framework de testing y assertions para Go
  - [User Guide](https://github.com/stretchr/testify#usage)
  - [Repository](https://github.com/stretchr/testify)

- **Gin**: Framework web para Go (utilizado en tests de endpoints)
  - [User Guide](https://gin-gonic.com/docs/)
  - [Repository](https://github.com/gin-gonic/gin)

### Endpoints

Para ejecutar las pruebas de endpoints es necesario correr el servicio completo con docker. Hay dos maneras de correrlo.

#### Ejecución rápida

- Windows (PowerShell)

```powershell
# Ejecutar tests (Docker ya debe estar corriendo)
.\run_tests.ps1

# Ejecutar tests con rebuild de Docker
.\run_tests.ps1 -Build

# Ejecutar tests limpiando archivos anteriores
.\run_tests.ps1 -Clean
```

- Linux/Mac (Bash)

```bash
# Dar permisos de ejecución
chmod +x run_tests.sh

# Ejecutar tests (Docker ya debe estar corriendo)
./run_tests.sh

# Ejecutar tests con rebuild de Docker
./run_tests.sh --build

# Ejecutar tests limpiando archivos anteriores
./run_tests.sh --clean
```

#### Ejecución manual

1. Levantar docker

```bash
docker compose up --build

# O compilar y levantar
docker compose build
docker compose up
```

Esperamos a que los servicios estén listos y luego para poder ejecutar los tests debemos ingresar a la carpeta de scripts, compilar el archivo de tests y ejecutarlo.

```bash
# Compilar
go build -o test_endpoints scripts/test_endpoints.go
# Ejecutar
./test_endpoints.exe

# Limpiar
rm test_endpoints
```

Los resultados se imprimen en la consola donde se ejecutan y se guardan en la carpeta `scripts/test_results`.

#### Troubleshooting

Error de compilación
```bash
# Verificar versión de Go
go version

# Limpiar cache
go clean -cache

# Verificar dependencias
go mod tidy
```

Tests fallan por timeout
```bash
# Verificar que la API esté respondiendo
curl -v http://localhost:8080/health

# Verificar logs de la aplicación
docker-compose logs melodia-api
```

Docker no responde
```bash
# Verificar estado
docker-compose ps

# Ver logs
docker-compose logs

# Reiniciar
docker-compose down
docker-compose up -d
```

### Modelo

Para ejecutar las pruebas del modelo, se pueden usar los siguientes comandos:
```bash
go test -v ./internal/models
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

## Base de Datos
El proyecto utiliza PostgreSQL como base de datos relacional para persistir canciones y playlists.

### Características
- **Motor**: PostgreSQL 16.3
- **Persistencia**: Volumen Docker para datos persistentes
- **Red**: Red dedicada para comunicación entre contenedores
- **Configuración**: Variables de entorno personalizables en `.env`

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

### Persistencia

Utilizando Docker, los datos quedan guardados en un volumen, de esta manera no se pierden. Solo se eliminan si se pide explícitamente con los comandos.
```bash
# Persisten
docker compose down    # Parar
docker compose up      # Levantar
docker compose restart # Reiniciar

# No persisten
docker compose down -v # Elimina datos
```

## Desiciones de diseño

- Se puede agregar una canción varias veces en una misma playlist.
- Se puede crear una playlist con un nombre ya existente

## Desafíos del Proyecto

Al leer la consigna, sabía que mi desafío iba a estar cerca de la base de datos y la persistencia ya que en el proyecto de Ingeniería de Software I no llegamos a implementar la persistencia. Además cambié de equipo, entonces tuve que instalar todo lo relacionado a PostgreSQL de nuevo. Para no complicarme, una vez que hice el backend, decidí dockerizar todo y crear tests para docker. Entonces dockericé el backend, la base de datos y creé un volumen donde persisten los datos. Por otra parte, este es mi primer proyecto con Go.

### Desafíos Opcionales

1. Implementado

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



