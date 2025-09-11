# Hangout-Planner

A scalable backend service for planning and managing hangouts, built with Go, Echo, GORM, and MySQL.  
Designed with clean architecture, best practices, and future-proofing in mind.

## 🚀 Tech Stack

- Language: Go 1.23+
- Framework: Echo (HTTP)
- ORM: GORM (MySQL)
- Relational Database: MySQL 8.0
- Infra & Tooling: Docker, Docker Compose, Makefile, Air, Golangci-Lint, Swag

## 🏃‍♂️ Local Development

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- MySQL (local or via Docker)
- Swag CLI for API docs
- golangci-lint
- Make (Makefile)

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

### hangout service

#### Module

- Documentation (with echoswagger)
- Unit Tests
  - with mocking and table driven test whenever applicable
  - tests folder containing unit test coverage file in HTML
- Code quality analysis, formatting, and linting with golangci-lint
- Makefile scripts

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

### Dev Dependencies

- air for dev hot reload
- pre-commit actions

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
- HTTPS
- Nginx API Gateway
- Multiple microservices
- Shared Module for microservices
- Multi db for microservices
- Github Actions CI/CD
- Cloud Deployments
- OAuth / multiple login method
- Excel service export
- Scheduled Notification service
- AWS S3 connectivity for excel file storage
