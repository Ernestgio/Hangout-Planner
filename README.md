# 🌍 Hangout Planner — Scalable Go Backend Platform

A **production-grade backend platform** for planning and managing hangouts — built in **Go** with **Echo**, **GORM**, and **MySQL**.  
Designed with **clean architecture**, **SOLID principles**, and **future-proof modular design** for microservices scalability.

## 🚀 Tech Stack

**Core:**

- 🟦 Language: Go 1.23+
- ⚙️ Framework: Echo (HTTP)
- 🗄️ ORM: GORM
- 💾 Database: MySQL 8.0

**Infra & Dev Tooling:**

- 🐳 Docker & Docker Compose
- 🧰 Makefile (automated scripts)
- 🌀 Air (live reload)
- 🧹 GolangCI-Lint (code linting)
- 🧾 Swag (OpenAPI documentation)
- 🪝 Lefthook (pre-commit & pre-push hooks)
- 🧪 CodeQL & GitHub Actions (CI/CD)
- Atlas for DB auto migration

## 🏃‍♂️ Local Development

### Prerequisites

- Go 1.24.11
- Docker & Docker Compose
- MySQL (local or via Docker)
- Swag CLI for API docs
- golangci-lint
- Make (Makefile)
- ☁️ Air - Live reload for Go apps
- Lefthook - git hooks for pre-commit / pre-push actions
- Atlas for db auto migration

### Mysql Environment Variables

Copy `components/database/.env.example` to `components/database/.env.example` and fill in your configuration

### Application Environment Variables

Copy `services/hangout/.env.example` to `services/hangout/.env` and fill in your configuration.

### Local deployment with mysql from docker compose and go run

```sh
make mysql-run
make run
```

### Local deployment fully with docker compose

-- Set DB_HOST to mysql -- utlizing docker network

```sh
make up
```

### DB Auto Migration

Each services will have its own database, please setup your local environment / system variable to have the connection string value of your db with the variable name `{SERVICE}_DB_URL`. We then can generate diff by executing `make diff NAME={Migration_name}` and apply migration by executing `make migrate`

---

## ⚡ Existing Features

### 🔧 Project Infrastructure

- Docker Compose orchestration
- Health checks and container restart policies
- GitHub Actions CI/CD
- Lefthook for local Git workflow automation
- Atlas for db auto migration

### 💬 Hangout Service

- Swagger auto-docs with echoswagger
- Unit tests (mocking, table-driven)
- Test coverage reports (HTML)
- GolangCI-Lint, Air reload
- Makefile automation
- Auth, Hangout, and Activity modules

### 💾 Database

- Auto migration with atlas and gorm
- Graceful shutdown

### 🌐 Server

- Standardized JSON response builder
- Centralized constants & sentinel errors
- Dependency injection (interfaces for all layers)
- Context propagation across all layers (for timeouts, cancellation, and future observability/tracing)

## 🧭 Roadmap

### 🧩 Short-Term Goals

- Retryable DB connections
- CORS middleware

### 🌐 Long-Term Vision

- Excel export service
  - RabbitMQ service interconnect
- Notification Emails
- File service
  - File upload feature (photos attachment for hangout memories!)
  - Memcached cluster caching presigned URL
  - AWS S3 integration (LocalStack support)
  - gRPC communication between fileservice and hangout service
- Multi db for microservices
- shared module in pkg/shared
- OAuth / federated logins
- Nginx API gateway + HTTPS (Let’s Encrypt)
- Advanced observability: metrics, tracing, logging
- Redis caching for File PreSignedURL and preventing concurrent login session
