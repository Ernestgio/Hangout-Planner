# File Service - Distributed File Management Microservice

A production-grade gRPC microservice dedicated to secure file management and S3 operations. Extracted from the Hangout Service to enforce separation of concerns, enable horizontal scalability, and establish clear service boundaries.

The File Service operates as a storage abstraction layer in a distributed microservices architecture and communicates with upstream services via gRPC secured with mutual TLS (mTLS).

## Overview

The File Service provides:

- Presigned S3 URL generation
- Upload confirmation workflow
- File metadata persistence
- Secure service-to-service communication
- Distributed tracing instrumentation

It ensures that file bytes never flow through the core API service, enforcing an efficient client-side upload architecture.

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
- **Observability**: OpenTelemetry instrumentation for distributed tracing and Prometheus metrics
- **Database Migrations**: Atlas CLI
- **Live Reload**: Air for development hot-reloading

---

### Design Patterns

- **Clean Architecture**: Layered separation (handlers, services, repository, domain)
- **Dependency Injection**: Interface-based abstractions for testability
- **Repository Pattern**: Data access abstraction with transactional support
- **Service Layer Pattern**: Business logic orchestration

- **Storage Abstraction Layer**: Encapsulates S3 client interactions
- **DTO Mapping Layer**: Explicit mapping between domain models and protobuf messages
- **Context Propagation**: Distributed trace continuity across service boundaries

## Core Responsibilites

### 1. File Lifecycle Management

- Create file metadata records
- Track upload state (pending, uploaded, etc.)
- Validate ownership and associations
- Support deletion and retrieval operations

### 2. Presigned URL Generation

The service generates time-bound presigned URLs using AWS SDK v2.

Characteristics:

- 15-minute expiration
- Strict content-type validation
- MD5 checksum support
- Enforced maximum file size (10MB)
- Limited batch size (max 10 files per request)

### 3. Upload Confirmation Workflow

After client upload:

1. Client notifies Hangout Service
2. Hangout Service calls ConfirmUpload via gRPC
3. File Service:
   - Verifies file exists in storage
   - Updates file state

## Service Architecture

### Layer Responsibilities

- **gRPC Handlers**: Request validation, response marshaling, error translation, trace propagation
- **Services**: Business logic, orchestration between storage and database
- **Repository**: Database operations with transactional support
- **Storage**: S3 client operations (presigned URLs, object operations)
- **Mapper**: Transform between domain entities and Protocol Buffer messages

## Integration with Hangout Service

### Ednd-to-End Request Flow

```
Client (HTTPS:443)
    ↓
Nginx (TLS termination)
    ↓
Hangout Service (HTTP:9000)
    ↓ (gRPC + mTLS)
File Service (gRPC:9001)
    ↓
S3 / LocalStack (4566)
```

Memory Upload Sequence

1. Client requests upload via Hangout Service
2. Hangout Service creates memory records
3. Hangout Service calls GenerateUploadURLs
4. File Service:
   - Creates file records
   - Generates presigned S3 URLs
5. Client uploads directly to S3
6. Client confirms upload
7. Hangout Service calls ConfirmUpload
8. File Service updates file status to uploaded

## Observability

The File Service is instrumented using OpenTelemetry and exports telemetry to a centralized observability stack.

### Metrics (Prometheus)

- gRPC request latency histograms
- Request throughput
- Error rate tracking
- Database query timing
- S3 interaction latency

Metrics are scraped by Prometheus and visualized in Grafana dashboards.

### Distributed Tracing

- gRPC server instrumentation via OpenTelemetry
- Trace context propagation from Hangout Service
- Visibility into:
  - gRPC method execution
  - Database operations
  - S3 client calls

Traces are exported to Grafana Tempo and allow full cross-service trace inspection.

## Development Setup

### Prerequisites

- Go 1.24.11+
- Docker & Docker Compose
- Make
- Protocol Buffer Compiler (`protoc`)
- Atlas CLI for migrations
- mkcert (for certificate generation)

### Environment Configuration

Copy `.env.example` to `.env` and configure the environment varaibles

Configure:

- Database connection
- S3 endpoint (LocalStack)
- TLS certificate paths
- gRPC server configuration
- Observability exporter endpoints

---

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
- Multi-bucket architecture (dirty → clean → thumbnails)
- Virus scanning integration
- Object lifecycle policies
- CDN integration for edge delivery
