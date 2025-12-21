# Hangout Planner — Scalable Go Backend Platform

Hangout Planner is a Go-based backend platform for planning and managing hangouts. The repo is structured as a monorepo with a core service (`hangout`) and shared packages, fronted by an NGINX edge gateway that handles TLS termination, HTTP/2, and path-based routing.

This project is built to demonstrate production-minded backend engineering: layered architecture, automated database migrations, CI, and a local environment that mirrors common deployment topology (gateway → services → database).

Designed with **clean architecture**, **SOLID principles**, and **future-proof modular design** for microservices scalability.

## Tech Stack

**Backend**

- Go (services/hangout)
- Echo (HTTP server, routing, middleware)
- JWT auth (echo-jwt)
- go-playground/validator (request validation)
- GORM + MySQL driver
- MySQL 8.0

**API & Documentation**

- OpenAPI/Swagger via swag + echo-swagger

**Infrastructure**

- Docker + Docker Compose (local orchestration)
- NGINX (reverse proxy, TLS termination, HTTP/2, gzip compression, header forwarding)

**Engineering Practices**

- GitHub Actions CI (lint + tests + coverage artifact)
- golangci-lint
- Lefthook (local git hooks)
- Atlas migrations (schema diff/apply)
- Make (scripting)

## Repository Layout

- `services/hangout/`: core HTTP API service
- `components/nginx/`: edge gateway (reverse proxy, HTTPS, HTTP/2)
- `components/database/`: local database bootstrap (init script, env)
- `pkg/shared/`: shared Go module (types/constants)

## Local Development

### Prerequisites

- Go 1.24.11
- Docker & Docker Compose
- MySQL (local or via Docker)
- Nginx (local or via Docker)
- Swag CLI for API docs
- golangci-lint
- Make (Makefile)
- Air - Live reload for Go apps
- Lefthook - git hooks for pre-commit / pre-push actions
- Atlas for db auto migration
- mkcert (one-time certificate generation)

### Environment variables

1. Database env

- Copy `components/database/.env.example` to `components/database/.env`

2. Service env

- Copy `services/hangout/.env.example` to `services/hangout/.env`

3. TLS certs for localhost (one-time)

- See `components/nginx/README.md`

### Local deployment with mysql from docker compose and go run

```sh
make mysql-run
make run
```

or use air for auto reload

```sh
make mysql-run
make air
```

### Local deployment fully with docker compose

-- Set DB_HOST to mysql -- utlizing docker network

```sh
make mysql-run (run the database first)
make up
```

### DB Auto Migration

Each services will have its own database, please setup your local environment / system variable to have the connection string value of your db with the variable name `{SERVICE}_DB_URL`. We then can generate diff by executing `make diff NAME={Migration_name}` and apply migration by executing `make migrate`

---

## Existing Features

### Project Infrastructure

- Docker Compose orchestration
- Health checks and container restart policies
- GitHub Actions CI/CD
- Lefthook for local Git workflow automation

### Hangout Service

- Auth, Hangout, and Activity modules
- Swagger auto-docs with echoswagger
- Unit tests (mocking, table-driven)
- Test coverage reports (HTML)
- GolangCI-Lint, Air reload
- Makefile automation
- More details on [Hangout Service Documentation](./services/hangout/README.md).

### Database

- Auto migration with atlas
- Graceful shutdown

### API Gateway

- Nginx with HTTPS
- Nginx as an API gateway, reverse-proxy, and rate limiter

## Roadmap

### Short-Term Goals

- Nginx API gateway, Reverse-proxy, Rate Limiter, and Load balancer + HTTPS
- File service
  - File upload feature (photos attachment for hangout memories!)
  - AWS S3 integration (LocalStack support)
  - gRPC communication between fileservice and hangout service
- Multi db for microservices
- shared module in pkg/shared

### Long-Term Vision

- Excel export service
  - RabbitMQ service interconnect
  - background worker service
- Notification Emails + SMTP
- OAuth / federated logins
- Advanced observability: metrics, tracing, logging
- Redis caching for preventing concurrent login session
- Implement file scanning using opengovsg [lambda-virus-scanner](https://github.com/opengovsg/lambda-virus-scanner) + 2 S3 bucket architecture (dirty and clean bucket)
