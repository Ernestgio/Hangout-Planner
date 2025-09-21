# Hangout-Planner

A scalable backend service for planning and managing hangouts, built with Go, Echo, GORM, and MySQL.  
Designed with clean architecture, best practices, and future-proofing in mind.

## 🚀 Tech Stack

- Language: Go 1.23+
- Framework: Echo (HTTP)
- ORM: GORM (MySQL)
- Relational Database: MySQL 8.0
- Infra & Tooling: Docker, Docker Compose, Makefile, Air, Golangci-Lint, Swag, Lefthook
- Github Actions for CI/CD

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

### Environment Variables

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

## Existing Feature

### Project

- Orchestration with docker-compose
  - Network
  - regular health checks
  - fault tolerance (`restart : on-failure`)
  - Dockerfile (multi services setup)
- Github Actions CI/CD
- Lefthook for pre-commit and pre-push actions

### hangout service

#### Module

- Documentation (with echoswagger)
- Unit Tests
  - with mocking and table driven test whenever applicable
  - tests folder containing unit test coverage file in HTML
- Code quality analysis, formatting, and linting with golangci-lint
- Makefile scripts
- Air for project auto reload

#### DB Connectivity

- minified graceful shutdown
- Auto migrate (code-based migration)

#### Server

- standard response
- constants
- sentinel errors
- Clean architecture dependency Injection with interface segregation

## Short Term Plan

### Controller, Services, and repository

- Sign up func and auth service

### DB

- Graceful shutdown
- retry connections
- Atlas versioned migration scripts (up and down)

### Code quality

- refactor and clean up services/hangout/internal/server package
- main cleanup
- separate AppConfig, DbConfig, RedisConfig

### Server settings

- cors middleware
- jwt middleware
- redis initializations

## Long Term Plan

- Features for hangouts, budgets, locations, activities, excel export, sharing
- HTTPS with lets encrypt open source certs
- Nginx API Gateway
- Multiple microservices
- Shared Module for microservices
- Multi db for microservices
- Cloud Deployments
- OAuth / multiple login method
- Excel service export
- Scheduled Notification service
- AWS S3 connectivity for excel file storage (localstack for local development)
- RabbitMQ Docker setup for connection between hangout service and report service
- shared module in pkg/shared
- open source static code analysis (sonarsource / codeQL)
