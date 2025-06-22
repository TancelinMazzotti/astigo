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

