# Hangout Service - Core Service for Hangout Planner Project

The **core backend service** responsible for creating, managing, and listing hangouts.  
Implements clean architecture principles and production-ready practices using Go, Echo, and GORM.

## ⚙️ Tech Stack

- 🟦 Go 1.23+
- ⚙️ Echo (HTTP Web Framework)
- 🗄️ GORM (ORM)
- 💾 MySQL (8.0)
- 🧪 GolangCI-Lint
- 🧰 Air (Live reload)
- 🧾 Swag (API documentation)

## 🏃‍♂️ Local Development

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- golangci-lint
- Make (Makefile)
- ☁️ Air - Live reload for Go apps

### Environment Variables

Copy `.env.example` to `.env` and fill in your configuration.

## ✅ Features

### 💡 Core

- RESTful API built on Echo
- Swagger API documentation
- Graceful server shutdown
- Dependency injection via interfaces
- Auto DB migration

### 🧪 Testing & Quality

- Unit tests (table-driven + mocks)
- HTML test coverage reports
- GolangCI-Lint configuration
- Makefile automation
- Live reload with Air

🧰 Server Layer

- Standard JSON response format
- Sentinel error design
- Request validator integration

🧭 Future Enhancements

- JWT authentication middleware
- Pagination, filtering, sorting
- Centralized error handling middleware
