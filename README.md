# Hangout-Planner

A scalable backend service for planning and managing hangouts, built with Go, Echo, GORM, and MySQL.  
Designed with clean architecture, best practices, and future-proofing in mind.

## 🚀 Tech Stack

- Go
- Echo
- GORM
- MySQL
- Docker

## 🏃‍♂️ Local Development

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- MySQL (local or Dockerized)
- Go Swag CLI tool
- Golangci-Lint tool

### Environment Variables

Copy `.env.example` to `.env` and fill in your configuration.

Local deployment with mysql from docker compose and go run

```sh
make mysql-run
make run
```

Local deployment fully with docker compose
-- Set DB_HOST to mysql -- utlizing docker network

```sh
make up
```

---

## Existing Feature

### Project

- Documentation (with echoswagger)
- Orchestration with docker compose
  - Network
  - regular health checks
  - fault tolerance (`restart : on-failure`)
  - Dockerfile (multi services setup)
- Unit Tests
  - with mocking and table driven test whenever applicable
  - tests folder containing unit test coverage file in HTML

### DB Connectivity

- minified graceful shutdown
- Auto migrate (code-embedded)

### Server

- standard response
- constants
- sentinel errors
- Clean architecture dependency Injection with interface segregation

## Short Term Plan

### Controller, Services, and repository

- Transaction wrapper
- Sign up func

### Dev Dependencies

- Linter
- go fmt
- Unit tests (coverage, mocking, out files)
- Code quality analysis
- air for pre-commit actions

### DB

- Graceful shutdown
- retry connections
- Atlas versioned migration scripts (up and down)

### Code quality

- refactor and clean up internal/server package
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
- Github Actions CI/CD
- Cloud Deployments
- OAth / multiple login method
