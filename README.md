# Hangout-Planner

A scalable backend service for planning and managing hangouts, built with Go, Echo, GORM, and MySQL.  
Designed with clean architecture, best practices, and future-proofing in mind.

## 🚀 Tech Stack

- Go
- Echo
- GORM
- MySQL
- Docker

## 🏃‍♂️ Local Development

### Prerequisites

- Go 1.23+
- Docker
- MySQL (local or Dockerized)

### Run Hangout Service

```sh
docker build -t hangout . \
  && docker run --rm --env-file .env -p 9000:9000 hangout \
  && docker rmi hangout
```

### Environment Variables

Copy `.env.example` to `.env` and fill in your configuration.

---

## Long Term Plan

- Multiple microservices
- Github Actions CI/CD
- Cloud Deployments
