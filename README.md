# Hangout Planner — Scalable Go Backend Platform

Hangout Planner is a microservices-based backend platform for planning and managing hangouts with photo memory uploads. The architecture features a REST API service for business logic and a dedicated gRPC File Service for storage operations, fronted by an NGINX gateway handling TLS termination, HTTP/2, and path-based routing.

This project is built to demonstrate production-minded backend engineering: layered architecture, automated database migrations, CI, grpc microservices integration (with mTLS), file upload with S3-compatible storage, and a local environment that mirrors common deployment topology (gateway → services → database → object storage).

Designed with **clean architecture**, **SOLID principles**, and **future-proof modular design** for microservices scalability.

## Tech Stack

**Backend**

- Go (services/hangout, services/file)
- Echo (HTTP server, routing, middleware)
- gRPC + Protocol Buffers (service-to-service communication)
- mTLS (mutual TLS authentication)
- JWT auth (echo-jwt)
- go-playground/validator (request validation)
- GORM + MySQL driver
- MySQL 8.0 (separate databases per service)
- Localstack

**API & Documentation**

- OpenAPI/Swagger via swag + echo-swagger

**Infrastructure**

- Docker + Docker Compose (local orchestration)
- NGINX (reverse proxy, TLS termination, HTTP/2, gzip compression, header forwarding)
- LocalStack (local AWS S3 emulation for development)
- AWS SDK for Go v2 (S3 client, presigned URLs)

**Engineering Practices**

- GitHub Actions CI (lint + tests + coverage artifact)
- golangci-lint
- Lefthook (local git hooks)
- Atlas migrations (schema diff/apply)
- Make (scripting)
- mkcert for one time cert generation

## Repository Layout

- `services/hangout/`: REST API service (auth, hangouts, activities, memories)
- `services/file/`: gRPC File Service (file lifecycle, S3 operations)
- `pkg/shared/`: shared contracts (Protocol Buffers, enums, types)
- `components/nginx/`: edge gateway (reverse proxy, HTTPS, HTTP/2)
- `components/database/`: local database bootstrap

## Local Development

### Prerequisites

- Go 1.24.11+
- Docker & Docker Compose
- Make
- Protocol Buffer compiler (protoc) with Go plugins
- Atlas CLI (migrations)
- Swag CLI (API docs)
- mkcert (one-time TLS certificate generation)

### Environment Setup

1. Copy `.env.example` to `.env` in each service directory
2. Generate TLS certificates (one-time): see `components/nginx/README.md`
3. Start services: `make up`
4. Run migrations: `cd services/hangout && make migrate && cd ../file && make migrate`

Each service has its own database. Set `{SERVICE}_DB_URL` environment variables for Atlas migrations.

---

## Existing Features

### Microservices Architecture

- **Hangout Service**: REST API with JWT auth, batch DB operations, Swagger docs
- **File Service**: gRPC file lifecycle management with mTLS authentication
- **Protocol Buffers**: Shared contracts in `pkg/shared` for service communication
- **Separate databases**: Each service owns its schema with Atlas migrations
- **Service discovery**: Docker DNS for internal routing

### File Upload Architecture

- Client-side upload via presigned S3 URLs (no file bytes through API)
- Batch operations for multi-file uploads with rollback support
- LocalStack S3 for development (AWS-compatible)
- AES-256 encryption, MD5 checksums, 15-minute URL expiry
- Max 10 files per upload, 10MB per file
- Supported formats: `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp`

### Security

- mTLS for service-to-service communication (Hangout ↔ File)
- JWT authentication for client requests
- TLS termination at Nginx gateway
- Certificate-based authentication with CA validation

### Infrastructure

- **Docker Compose**: Multi-service orchestration with health checks and restart policies
- **Nginx**: API gateway with reverse proxy, HTTP/2, TLS termination, rate limiting, and load balancing
- **LocalStack**: AWS-compatible S3 emulation for local development with persistence
- **GitHub Actions**: CI/CD pipeline for linting, testing, and coverage reporting

### Tooling

- **Make**: Automation scripts for common tasks (build, test, migrate, up/down)
- **Air**: Live reload for Go services during development
- **Swag**: OpenAPI/Swagger documentation generation from code annotations
- **golangci-lint**: Comprehensive Go linting with multiple checkers
- **Lefthook**: Git hooks for pre-commit and pre-push automation
- **mkcert**: Local TLS certificate generation for HTTPS development
- **Atlas**: Database schema migrations with diff and apply capabilities

## Roadmap

### Short-Term Goals

- **Observability stack**: Prometheus (metrics), Grafana (visualization), Grafana Tempo (distributed tracing), OpenTelemetry Collector (instrumentation)

### Long-Term Vision

- Security, between services, port protection, cors, permissions and roles for services and users
- Excel export service
  - RabbitMQ service interconnect
  - background worker service
- Notification Emails + SMTP
- OAuth / federated logins
- Role-based access control (RBAC) for multi-user scenarios
- Redis caching layer for session management and rate limiting
- Implement file scanning using opengovsg [lambda-virus-scanner](https://github.com/opengovsg/lambda-virus-scanner) + 3 S3 buckets architecture (dirty, clean, and thumbnail / resized image bucket)
