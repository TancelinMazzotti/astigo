
# Astigo - API REST Boilerplate
![CI](https://github.com/TancelinMazzotti/astigo/actions/workflows/ci.yml/badge.svg?branch=main)
[![Docker](https://img.shields.io/badge/docker-ghcr.io%2Ftancelinmazzotti%2Fastigo-blue)](https://github.com/users/TancelinMazzotti/packages/container/package/astigo)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/TancelinMazzotti/astigo)
![Go Version](https://img.shields.io/badge/go-1.24.2-blue)
![License](https://img.shields.io/github/license/TancelinMazzotti/astigo)


Astigo est un projet **boilerplate** prêt à l’emploi pour démarrer rapidement une API REST moderne et maintenable, basée sur les principes de **l’architecture hexagonale** et du **Domain-Driven Design (DDD)**.

Ce projet met l’accent sur la **séparation des préoccupations**, la **scalabilité**, et l'**extensibilité**, tout en intégrant des technologies robustes pour le transport, la persistance et la configuration.

---

## 🚀 Fonctionnalités principales

- 📦 Architecture **Hexagonale**
- 🧠 Modèle **DDD**
- 🔥 Serveur HTTP basé sur **Gin**
- ⚙️ Configuration flexible via **Viper**
- 💻 CLI intégrée avec **Cobra**
- 📝 Logging structuré avec **Zap**
- 🗃️ Persistance avec **PostgreSQL**
- 🧠 Cache distribué avec **Redis**
- 📨 Événements asynchrones via **NATS**
- ✅ Tests unitaires avec mocking et interfaces
- 🧪 Tests d’intégration isolés avec Testcontainers
- 🐳 Déploiement et développement via Docker & Docker Compose
- ❤️ Endpoints **/health/liveness** et **/health/readiness** pour Kubernetes
- 📊 Endpoint **/metrics** compatible **Prometheus**
