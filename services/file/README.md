# File Service - Distributed File Management Microservice

A production-grade gRPC microservice dedicated to secure file management and S3 operations, extracted from the monolithic Hangout Service to enable horizontal scalability and separation of concerns.

## Overview

The File Service provides centralized file lifecycle management with presigned URL generation, S3 orchestration, and distributed tracing capabilities. It communicates with upstream services via gRPC with mutual TLS authentication, ensuring secure service-to-service communication in a microservices architecture.

## Technical Architecture

### Communication Protocol

- **gRPC with HTTP/2**: Binary protocol for efficient inter-service communication
- **Protocol Buffers**: Strongly-typed contracts defined in `pkg/shared/proto`
- **mTLS (Mutual TLS)**: Bidirectional certificate authentication for zero-trust networking
- **Certificate Infrastructure**: CA-signed certificates for both server and client authentication

### Technology Stack

- **Language**: Go 1.24.11
- **RPC Framework**: gRPC with Protocol Buffers
- **ORM**: GORM with MySQL driver
- **Storage**: AWS S3 / LocalStack (S3-compatible)
- **Security**: mTLS with x509 certificates, AES-256 encryption at rest
- **Observability**: OpenTelemetry instrumentation for distributed tracing
- **Database Migrations**: Atlas CLI
- **Live Reload**: Air for development hot-reloading

### Design Patterns

- **Clean Architecture**: Layered separation (handlers, services, repository, domain)
- **Dependency Injection**: Interface-based abstractions for testability
- **Repository Pattern**: Data access abstraction with transactional support
- **Service Layer Pattern**: Business logic encapsulation
- **DTO Pattern**: Data transfer objects for gRPC message mapping

## Core Features

- File Lifecycle Management
- S3 Integration
- grpc servier
- mTLS communication

## Service Architecture

### Layer Responsibilities

- **gRPC Handlers**: Request validation, response marshaling, error translation
- **Services**: Business logic, orchestration between storage and database
- **Repository**: Database operations with transactional support
- **Storage**: S3 client operations (presigned URLs, object operations)
- **Mapper**: Transform between domain entities and Protocol Buffer messages

## Development Setup

### Prerequisites

- Go 1.24.11+
- Docker & Docker Compose
- Make
- Protocol Buffer Compiler (`protoc`)
- Atlas CLI for migrations
- OpenSSL (for certificate generation)

### Environment Configuration

Copy `.env.example` to `.env` and configure the environment varaibles

### Running the Service

**Start with Docker Compose** (recommended):

```bash
make up
```

**Local Development with Air**:

```bash
make air
```

**Run Database Migrations**:

```bash
make migrate
```

**Check Migration Status**:

```bash
make migrate-status
```

## Integration with Hangout Service

### Request Flow

1. **Client** → HTTPS request to Nginx (port 443)
2. **Nginx** → HTTP request to Hangout Service (port 9000)
3. **Hangout Service** → gRPC + mTLS to File Service (port 9001)
4. **File Service** → S3 operations via AWS SDK

### Memory Upload Sequence

1. Hangout service receives upload request from client
2. Hangout service creates memory records in its database
3. Hangout service calls `GenerateUploadURLs` gRPC method
4. File service creates file records and returns presigned URLs
5. Hangout service updates memories with `file_id` references
6. Client uploads files directly to S3 using presigned URLs
7. Client calls confirm endpoint on Hangout service
8. Hangout service calls `ConfirmUpload` gRPC method
9. File service marks files as `uploaded`

## Testing

### Unit Tests

```bash
make test
```

### Test Coverage

```bash
make coverage
```

Opens HTML coverage report in browser.

## Future Enhancements

### Features

- Image resizing and thumbnail generation
- File virus scanning integration
- CDN integration for faster delivery
- Object storage lifecycle policies

### Observability and Monitoring

- Prometheus metrics exposition
- Grafana dashboards
- Distributed tracing with Jaeger/Tempo
- Log aggregation with ELK/Loki
