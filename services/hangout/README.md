# Hangout Service - Core REST API Microservice

A production-grade RESTful API service for social event management, implementing modern backend patterns with clean architecture, secure gRPC client integration, and full-stack observability.

The Hangout Service acts as the primary business logic layer in a distributed microservices architecture, coordinating authentication, hangout lifecycle management, and memory orchestration.

## Overview

The Hangout Service is responsible for:

- User authentication and authorization

- Hangout and activity management

- Hangout Memories lifecycle orchestration

- Coordinating file operations via the File Service (gRPC + mTLS)

It demonstrates real-world microservice design principles including service boundaries, contract-driven communication, observability instrumentation, and layered architecture.

## Technical Architecture

### Technology Stack

- **Language**: Go 1.24.11
- **Web Framework**: Echo v4 with middleware ecosystem
- **ORM**: GORM with MySQL driver, batch operations optimization
- **Authentication**: JWT with RS256 signing, bcrypt password hashing
- **API Documentation**: Swagger/OpenAPI 3.0 via swag
- **Database Migrations**: Atlas CLI with declarative schema
- **Service Communication**: gRPC client with mTLS for File Service integration
- **Development Tools**: Air for live reload, golangci-lint for code quality

---

### Architecture Patterns

- **Clean Architecture**: Layered design with dependency inversion
- **Repository Pattern**: Data access abstraction with transactional support
- **Service Layer Pattern**: Business logic encapsulation and orchestration
- **DTO Pattern**: Request/response mapping with validation
- **Middleware Chain**: Authentication, context propagation, error handling
- **Client-Side Upload Pattern**: Presigned URLs for efficient file transfers

## Core Features

### Authentication & Authorization

- JWT authentication (RS256 asymmetric signing)
- Configurable token expiration
- Secure password hashing via bcrypt
- Route-level middleware enforcement
- User context propagation across request lifecycle
- Custom JWT claims payload

---

### Hangout Management

- **CRUD Operations**: Full lifecycle management for hangout events
- **Ownership Validation**: Users can only modify their own hangouts
- **Listing & Pagination**: Efficient bulk retrieval with cursor-based pagination
- Optimized DB queries for bulk retrieval

---

### Activity Coordination

- **Activity CRUD**: Manage activities within hangout contexts
- **Bulk Retrieval**: Efficient listing with pagination support

---

### Memory Management (Client-Side Upload Pattern)

- **Three-Phase Upload Flow**:
  1. **Generate Upload URLs**: Batch create memory records, obtain presigned S3 URLs
  2. **Client Direct Upload**: Client uploads files directly to S3 using presigned URLs
  3. **Confirm Upload**: Mark files as confirmed after successful upload
- **Batch Operations**: Single SQL INSERT for multiple memories
- **Ownership Validation**: Batch fetch memories to verify user access
- **File Service Integration**: gRPC client with mTLS for secure communication
- **Cursor-Based Pagination**: Efficient memory listing with hasMore/nextCursor

---

## API Documentation

The Hangout Service exposes a fully documented OpenAPI 3.0 specification generated directly from code annotations.

Interactive Swagger UI:

![Hangout Service Swagger UI](hangout-service-swagger.jpg)

Access locally:

https://localhost/rp-api/hangout-service/swagger/index.html

Swagger documentation includes:

- Request and response schemas
- Authentication requirements
- Pagination contracts
- Error response formats
- Example payloads

Generate documentation:

```bash
make swag
```

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

gRPC Client Integration

The Hangout Service communicates with the File Service via:

- mTLS-secured gRPC connection
- Shared Protocol Buffer contracts
- Context propagation for tracing
- Timeout and error handling boundaries

Operations include:

- Generate presigned upload URLs
- Confirm upload completion
- Retrieve file metadata
- Delete files

## Observability

The Hangout Service is instrumented using OpenTelemetry and exports telemetry to a centralized observability stack.

### Metrics

- HTTP request duration (latency histograms)
- Request throughput
- Error rate tracking
- Database query timing
- gRPC client latency

### Distributed Tracing

- OpenTelemetry HTTP middleware instrumentation

- Context propagation across service boundaries

- End-to-end trace visibility:
  - Incoming HTTP request
  - Business logic execution
  - gRPC call to File Service
  - Downstream S3 interaction (via File Service)

Traces are exported to Grafana Tempo and visualized in Grafana dashboards.

This enables:

- Latency bottleneck detection
- Cross-service request correlation
- Root cause analysis for production failures

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

Configure:

- Database connection
- JWT keys and configs
- gRPC TLS certificates
- File Service endpoint
- Observability exporter configuration

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

- Table-driven tests
- Mock repositories for service isolation
- Layer-specific validation tests
- Business logic validation without HTTP layer coupling

## Future Enhancements

### Features

- Location-based hangout discovery
- Hangout sharing and invitations
- Activity voting and RSVP
- Memory reactions and comments
- Redis session management (prevent concurrent sessions)
