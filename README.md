🌍 Hangout Planner — Scalable Go Backend Platform

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

## 🏃‍♂️ Local Development

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- MySQL (local or via Docker)
- Swag CLI for API docs
- golangci-lint
- Make (Makefile)
- ☁️ Air - Live reload for Go apps
- Lefthook - git hooks for pre-commit / pre-push actions

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

---

## ⚡ Existing Features

### 🔧 Project Infrastructure

- Docker Compose orchestration
- Health checks and container restart policies
- GitHub Actions CI/CD
- Lefthook for local Git workflow automation

### 💬 Hangout Service

- Swagger auto-docs with echoswagger
- Unit tests (mocking, table-driven)
- Test coverage reports (HTML)
- GolangCI-Lint, Air reload
- Makefile automation

### 💾 Database

- Auto migration
- Graceful shutdown
- Future migration tooling ready (Atlas)

### 🌐 Server

- Standardized JSON response builder
- Centralized constants & sentinel errors
- Dependency injection (interfaces for all layers)

## 🧭 Roadmap

### 🧩 Short-Term Goals

- Graceful shutdown improvements
- Retryable DB connections
- Atlas migration (up/down)
- CORS & JWT middleware
- Full Hangout CRUD
- Budget, activity, and location modules

### 🌐 Long-Term Vision

- Full Hangout CRUD & collaboration
- Excel export microservice
- Notification scheduling
- AWS S3 integration (LocalStack support)
- File upload feature (photos attachment for hangout memories!)
- Multi db for microservices
- shared module in pkg/shared
- RabbitMQ service interconnect
- Memcached cluster caching
- OAuth / federated logins
- File uploads for hangout memories
- Nginx API gateway + HTTPS (Let’s Encrypt)
- Context propagation refactor + Advanced observability: metrics, tracing, logging
