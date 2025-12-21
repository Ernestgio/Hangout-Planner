# Hangout Service - Core Service for Hangout Planner Project

The **core backend service** responsible for creating, managing, and listing hangouts.  
Implements clean architecture principles and production-ready practices using Go, Echo, and GORM.

## ⚙️ Tech Stack

- Go 1.24.11
- Echo (HTTP Web Framework)
- GORM (ORM)
- MySQL (8.0)
- GolangCI-Lint
- Air (Live reload)
- Swag (API documentation)
- Atlas for DB auto migration

## 🏃‍♂️ Local Development

### Prerequisites

- Go 1.24.11
- Docker & Docker Compose
- golangci-lint
- Make (Makefile)
- Air - Live reload for Go apps
- Swag (API documentation)
- Atlas

### Environment Variables

Copy `.env.example` to `.env` and fill in your configuration.

## Features

### Modules

- Auth Modules
- Hangout Modules
- Activity modules

### Core

- RESTful API built on Echo
- Swagger API documentation
- Graceful server shutdown
- Dependency injection via interfaces
- Auto DB migration
- Context propagation across all layers (for timeouts, cancellation, and future observability/tracing)
- Full Hangout CRUD
- Signup and Signin

### Testing & Quality

- Unit tests (table-driven + mocks)
- HTML test coverage reports
- GolangCI-Lint configuration
- Makefile automation
- Live reload with Air

### Server Layer

- Standard JSON response format
- Sentinel error design
- Request validator integration
- JWT authentication middleware
- Centralized error handling middleware

### Future Enhancements

- Redis to prevent concurrent session
- Location tagging
- Hangout memories
- Share hangout features
