# Hangout Service - Core Service for Hangout Planner Project

The **core backend service** responsible for creating, managing, and listing hangouts.  
Implements clean architecture principles and production-ready practices using Go, Echo, and GORM.

## Tech Stack

- Go (module: `services/hangout`)
- Echo + echo-jwt
- GORM + MySQL
- go-playground/validator (request validation)
- Swagger/OpenAPI via swag + echo-swagger
- Atlas migrations
- golangci-lint
- Air (Live reload)
- AWS SDK for Go v2 (S3 client)
- LocalStack (local S3 emulation)

## Architecture

The service follows a layered, standard-convention dependency-inverted structure:

- `handlers/`: HTTP handlers (request binding/validation, response mapping)
- `services/`: application use-cases
- `repository/`: persistence boundaries
- `domain/`: entities and core types
- `dto/` + `mapper/`: transport models and mapping
- `middlewares/`: auth + request context wiring
- `internal/http/`: shared request/response utilities (validator, sanitizer, response envelope)

The edge gateway terminates TLS/HTTP/2 and forwards requests to this service over the Docker network (HTTP).

## Running Locally

### Prerequisites

- Docker Desktop
- Make (optional but recommended)
- Go 1.24.11
- Docker & Docker Compose
- golangci-lint
- Air - Live reload for Go apps
- Swag (API documentation)
- Atlas

### Environment

Copy `.env.example` to `.env` and fill in your configuration.

## Features

### Modules

- Auth Modules
- Hangout Modules
- Activity modules
- **Memory Modules** — Photo upload system with concurrent processing and S3 storage
  - Concurrent file processing with goroutines (6-10x faster than sequential)
  - Partial success handling for better user experience
  - S3-compatible storage with AES-256 encryption and MD5 checksums
  - Presigned URLs for secure file access (15-minute expiry)
  - Cursor-based pagination for memory listing

### Core

- RESTful API built on Echo
- Swagger API documentation
- Graceful server shutdown
- Dependency injection via interfaces
- Auto DB migration
- Context propagation across all layers (for timeouts, cancellation, and future observability/tracing)
- Full Hangout CRUD
- Signup and Signin
- **File Upload System**
  - Multipart/form-data handling
  - Transport-agnostic validation layer
  - Thread-safe concurrent uploads
  - Atomic per-file transactions
  - Designed for future microservice extraction

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
- Share hangout features
