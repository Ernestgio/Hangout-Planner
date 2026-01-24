# Hangout Service - Core REST API Microservice

A production-grade RESTful API service for social event management, implementing modern backend patterns with clean architecture, gRPC client integration, and comprehensive observability. Serves as the primary business logic layer in a distributed microservices architecture.

## Overview

The Hangout Service is the main backend service handling user authentication, hangout management, activity coordination, and memory storage orchestration. It integrates with the File Service via gRPC for distributed file operations, demonstrating service mesh patterns and secure inter-service communication.

## Technical Architecture

### Technology Stack

- **Language**: Go 1.24.11
- **Web Framework**: Echo v4 with middleware ecosystem
- **ORM**: GORM with MySQL driver, batch operations optimization
- **Authentication**: JWT with RS256 signing, bcrypt password hashing
- **API Documentation**: Swagger/OpenAPI 3.0 via swag
- **Database Migrations**: Atlas CLI with declarative schema
- **Service Communication**: gRPC client with mTLS for File Service integration
- **Storage Integration**: S3-compatible via File Service abstraction
- **Development Tools**: Air for live reload, golangci-lint for code quality

### Architecture Patterns

- **Clean Architecture**: Layered design with dependency inversion
- **Repository Pattern**: Data access abstraction with transactional support
- **Service Layer Pattern**: Business logic encapsulation and orchestration
- **DTO Pattern**: Request/response mapping with validation
- **Middleware Chain**: Authentication, context propagation, error handling
- **Client-Side Upload Pattern**: Presigned URLs for efficient file transfers

## Core Features

### Authentication & Authorization

- **JWT Authentication**: RS256 asymmetric signing with configurable expiry
- **Secure Password Storage**: bcrypt hashing with adjustable cost factor
- **Middleware Protection**: Route-level authentication enforcement
- **User Context Propagation**: Request-scoped user information across layers
- **Token Claims**: Custom JWT payload with user metadata

### Hangout Management

- **CRUD Operations**: Full lifecycle management for hangout events
- **Ownership Validation**: Users can only modify their own hangouts
- **Listing & Pagination**: Efficient bulk retrieval with cursor-based pagination

### Activity Coordination

- **Activity CRUD**: Manage activities within hangout contexts
- **Bulk Retrieval**: Efficient listing with pagination support

### Memory Management (Client-Side Upload Pattern)

- **Three-Phase Upload Flow**:
  1. **Generate Upload URLs**: Batch create memory records, obtain presigned S3 URLs
  2. **Client Direct Upload**: Client uploads files directly to S3 using presigned URLs
  3. **Confirm Upload**: Mark files as confirmed after successful upload
- **Batch Operations**: Single SQL INSERT for multiple memories
- **Ownership Validation**: Batch fetch memories to verify user access
- **File Service Integration**: gRPC client with mTLS for secure communication
- **Cursor-Based Pagination**: Efficient memory listing with hasMore/nextCursor

### gRPC Client Integration

Integrates with File Service via gRPC with mTLS for secure file operations (generate presigned URLs, confirm uploads, retrieve files, delete files).

## Service Integration

### Nginx API Gateway Flow

```
Client (HTTPS:443)
    ↓
Nginx (TLS termination, HTTP/2)
    ↓
Hangout Service (HTTP:9000)
    ↓ (gRPC + mTLS)
File Service (gRPC:9001)
    ↓
S3 / LocalStack (S3:4566)
```

## Development Setup

### Prerequisites

- Go 1.24.11+
- Docker & Docker Compose
- Make
- golangci-lint
- Atlas CLI
- Swag (for API documentation generation)

### Environment Configuration

Copy `.env.example` to `.env` and configure your environment variables.

### Running the Service

**Start all services with Docker Compose**:

```bash
make up
```

**Local development with Air (live reload)**:

```bash
make air
```

**Run database migrations**:

```bash
make migrate
```

**Generate Swagger documentation**:

```bash
make swag
```

**Access API documentation**:

```
https://localhost/rp-api/hangout-service/swagger/index.html
```

## Testing & Quality

### Unit Testing

```bash
make test
```

### Test Coverage Report

```bash
make coverage
```

Generates HTML report and opens in browser.

### Code Quality

```bash
make lint
```

Runs golangci-lint with project configuration.

### Test Structure

- Table-driven tests for comprehensive scenarios
- Mock repositories for service layer testing

## Future Enhancements

### Features

- Location-based hangout discovery
- Hangout sharing and invitations
- Activity voting and RSVP
- Memory reactions and comments
- Redis session management (prevent concurrent sessions)

### Observability

- Distributed tracing with Jaeger
- Centralized logging with ELK/Loki
- Alerting with Prometheus Alertmanager
- Custom Grafana dashboards
