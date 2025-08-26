# Astigo - Backend API Boilerplate
[![Docker](https://img.shields.io/badge/docker-ghcr.io%2Ftancelinmazzotti%2Fastigo-blue)](https://github.com/users/TancelinMazzotti/packages/container/package/astigo)
![Go Version](https://img.shields.io/badge/go-1.24.2-blue)
![License](https://img.shields.io/github/license/TancelinMazzotti/astigo)
![Status](https://img.shields.io/badge/status-WIP-orange)

⚠️ **Project Status**: This project is currently under active development. Features and APIs may change without notice. Not recommended for production use yet.

<img src="astigo.png" alt="Astigo Mascot" width="200"/>

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
| `ASTIGO_AUTH_ISSUER`             | `http://localhost:8080/realms/astigo` | Keycloak realm URL used for JWT token validation            |
| `ASTIGO_AUTH_CLIENT_ID`          | `astigo-api`                          | Keycloak client ID used for API authentication              |
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

## 🔐 Keycloak Access

⚠️ **Important Note about Issuer URL**:
The application uses `host.docker.internal:8090` as the issuer URL instead of `localhost:8090`. This ensures that both the API (running inside Docker) and the client (running on the host machine) use the same issuer URL, which is required for proper OpenID Connect validation.

⚠️ **Important Prerequisite**:
Make sure the `host.docker.internal` entry is present in your hosts file (typically located at `/etc/hosts` on Linux/MacOS or `C:\Windows\System32\drivers\etc\hosts` on Windows). This entry is normally added automatically by Docker Desktop, but it's recommended to verify it. Without this entry, name resolution won't work properly.

The application comes pre-configured with Keycloak and includes all necessary settings:

- **Realm**: `astigo`
- **Client ID**: `astigo-api`
- **Issuer**: `http://host.docker.internal:8090/realms/astigo`
- **Client Secret**: `astigo_secret`
- **Default User**:
    - Username: `astigo`
    - Password: `astigo`
    - Email: `astigo@gmail.com`
    - Email verified: `true`
    - First name: `Asti`
    - Last name: `Go`
    - Role: `astigo-api manager`, `account view-profile`, `account manage-account`, `default-roles-astigo`

These settings are already configured and ready to use in the Docker environment. No additional manual configuration is required.

To access the Keycloak admin interface:
- URL: `http://localhost:8090`
- Admin credentials:
    - Username: `admin`
    - Password: `admin`
- Change the current realm with: `astigo` 

To test the API endpoints, you can directly use the requests in `e2e/http/Private.http` without any modifications needed.
These settings will enable proper authentication and authorization for the API endpoints.
