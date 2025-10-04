# Hangout Service - Core Service for Hangout Planner Project

A scalable backend service for planning and managing hangouts, built with Go, Echo, GORM, and MySQL.  
Designed with clean architecture, best practices, and future-proofing in mind.

## 🚀 Tech Stack

- Language: Go 1.23+
- Framework: Echo (HTTP)
- ORM: GORM (MySQL)

## 🏃‍♂️ Local Development

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- golangci-lint
- Make (Makefile)
- ☁️ Air - Live reload for Go apps

### Environment Variables

Copy `.env.example` to `.env` and fill in your configuration.

## Existing Feature

### Module

- Documentation (with echoswagger)
- Unit Tests
  - with mocking and table driven test whenever applicable
  - tests folder containing unit test coverage file in HTML
- Code quality analysis, formatting, and linting with golangci-lint
- Makefile scripts
- Air for project auto reload

### DB Connectivity

- Graceful shutdown
- Auto migrate (code-based migration)

### Server

- standard response
- constants
- sentinel errors
- Clean architecture dependency Injection with interface segregation
- Gzip response compression
