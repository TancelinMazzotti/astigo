# Astigo - Backend API Boilerplate
[![Docker](https://img.shields.io/badge/docker-ghcr.io%2Ftancelinmazzotti%2Fastigo-blue)](https://github.com/users/TancelinMazzotti/packages/container/package/astigo)
![Go Version](https://img.shields.io/badge/go-1.24.2-blue)
![License](https://img.shields.io/github/license/TancelinMazzotti/astigo)
![Status](https://img.shields.io/badge/status-WIP-orange)



⚠️ **Project Status**: This project is currently under active development. Features and APIs may change without notice. Not recommended for production use yet.

---

Astigo is a production-ready **boilerplate** designed to quickly bootstrap modern and maintainable REST APIs, built on the principles of **Hexagonal Architecture** and **Domain-Driven Design (DDD)**.

This project emphasizes **separation of concerns**, **scalability**, and **extensibility**, while integrating robust technologies for transport, persistence, and configuration management.

---

## 🚀 Key Features

### Architecture & Design
- 📦 **Hexagonal Architecture** for clean separation of concerns
- 🧠 **Domain-Driven Design (DDD)** principles
- 🔌 Well-defined interfaces for all external dependencies

### Core Technologies
- 🔥 High-performance HTTP server using **Gin**
- ⚡ gRPC support for efficient service-to-service communication
- ⚙️ Flexible configuration with **Viper**
  - Environment variables
  - YAML configuration files
  - Command-line flags
- 💻 Intuitive CLI powered by **Cobra**
- 📝 Structured logging with **Zap**


### Data Management
- 🗃️ Persistent storage with **PostgreSQL**
- 🧠 Distributed caching using **Redis**
- 📨 Asynchronous event handling via **NATS**
- 🔐 Authentication and authorization via **Keycloak**

### Testing & Quality
- ✅ Comprehensive unit tests with mocking
- 🧪 Isolated integration tests using Testcontainers
- 📊 Code coverage reporting
- 🔍 Linting and code quality checks

### Operations & Deployment
- 🐳 Docker & Docker Compose support
  - Development environment
  - Production-ready configurations
- 🎯 Kubernetes-ready with health endpoints
  - `/health/liveness` for liveness probes
  - `/health/readiness` for readiness probes
- 📊 Prometheus-compatible metrics at `/metrics`
  - Application metrics
  - Runtime metrics
  - Custom business metrics
- 🔍 **Jaeger** for distributed tracing
  - End-to-end request tracking
  - Performance monitoring
  - Distributed system visualization


---

## Environement variable

| Environment Variable             | Default Value                         | Description                                                 |
|----------------------------------|---------------------------------------|-------------------------------------------------------------|
| `ASTIGO_HTTP_MODE`               | `debug`                               | HTTP server mode (debug/release)                            |
| `ASTIGO_HTTP_PORT`               | `8080`                                | HTTP server listening port                                  |
| `ASTIGO_GRPC_PORT`               | `50051`                               | gRPC server listening port                                  |
| `ASTIGO_HTTP_ISSUER`             | `http://localhost:8080/realms/astigo` | Keycloak realm URL used for JWT token validation            |
| `ASTIGO_HTTP_CLIENT_ID`          | `astigo-api`                          | Keycloak client ID used for API authentication              |
| `ASTIGO_LOG_LEVEL`               | `info`                                | Application logging level (info, debug, error, etc.)        |
| `ASTIGO_LOG_ENCODING`            | `json`                                | Log format encoding (json/console)                          |
| `ASTIGO_JAEGER_URL`              | `localhost:4318`                      | Jaeger collector endpoint URL for distributed tracing       |
| `ASTIGO_JAEGER_SERVICE_NAME`     | `astigo`                              | Service name identifier in Jaeger for tracing visualization |
| `ASTIGO_POSTGRES_HOST`           | `localhost`                           | PostgreSQL server hostname                                  |
| `ASTIGO_POSTGRES_PORT`           | `5432`                                | PostgreSQL connection port                                  |
| `ASTIGO_POSTGRES_DB`             | `astigo`                              | PostgreSQL database name                                    |
| `ASTIGO_POSTGRES_USER`           | `astigo`                              | PostgreSQL username                                         |
| `ASTIGO_POSTGRES_PASSWORD`       | `astigo_password`                     | PostgreSQL password                                         |
| `ASTIGO_POSTGRES_SSLMODE`        | `disable`                             | PostgreSQL SSL mode                                         |
| `ASTIGO_POSTGRES_MAX_OPEN_CONNS` | `10`                                  | PostgreSQL maximum open connections                         |
| `ASTIGO_POSTGRES_MAX_IDLE_CONNS` | `5`                                   | PostgreSQL maximum idle connections                         |
| `ASTIGO_POSTGRES_MAX_LIFETIME`   | `300`                                 | PostgreSQL connection maximum lifetime (seconds)            |
| `ASTIGO_NATS_URL`                | `nats://localhost:4222`               | NATS server connection URL                                  |
| `ASTIGO_REDIS_HOST`              | `localhost`                           | Redis server hostname                                       |
| `ASTIGO_REDIS_PORT`              | `6379`                                | Redis connection port                                       |
| `ASTIGO_REDIS_DB`                | `0`                                   | Redis database index                                        |

All variables are prefixed with `ASTIGO_` to prevent conflicts with other applications. Each variable controls a specific aspect of the application configuration:
- Server settings (HTTP/gRPC)
- Logging
- PostgreSQL database connection
- NATS message broker connection
- Redis cache connection

## 🔐 Keycloak Configuration

After launching the application with docker-compose, follow these steps to configure Keycloak:

1. **Login to Keycloak Admin Console**
  - Access the Keycloak admin interface
  - Login with credentials:
    - Username: `admin`
    - Password: `admin`

2. **Configure Astigo Realm**
  - Navigate to the "astigo" realm
  - Go to "Clients" > "astigo-api" > "Credentials" tab
  - Click "Regenerate" to create a new client secret
  - Save this secret as it will be required for HTTP requests authentication

3. **Create User**
  - Go to "Users" section
  - Click "Add User" and fill in the following details:
    - Username: `astigo`
    - Email: `astigo@gmail.com`
    - First Name: `Asti`
    - Last Name: `Go`
  - Enable "Email Verified"
  - Save the user

4. **Set User Password**
  - Go to the "Credentials" tab
  - Set password to: `astigo`
  - Disable "Temporary" password option
  - Click "Set Password"

5. **Assign Role**
  - Navigate to "Role Mapping" tab
  - Click "Assign role"
  - Select "astigo-api user" role
  - Save changes

6. **Testing API Endpoints**
  - Open the file `e2e/http/Private.http`
  - Replace the existing client secret with your newly generated one
  - You can now execute the HTTP requests to test the API endpoints


These settings will enable proper authentication and authorization for the API endpoints.
