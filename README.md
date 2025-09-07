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

```sh
docker-compose up --build
```

### Environment Variables

Copy `.env.example` to `.env` and fill in your configuration.

---

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
- migration scripts (up and down)

### Code quality

- refactor and clean up internal/server package
- main cleanup

### Server settings

- cors middleware
- jwt middleware
- redis initializations

## Long Term Plan

- Multiple microservices
- Github Actions CI/CD
- Cloud Deployments
