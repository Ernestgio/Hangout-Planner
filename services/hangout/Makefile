# Makefile for Hangout Planner

APP_NAME=hangout
COVERAGE_FILE=test/coverage
COVERAGE_HTML_FILE=test/coverage.html

.PHONY: all test lint swag build run clean docker-build docker-run up down air

all: lint test build

## Run tests with coverage
test:
	go test ./internal/... -coverprofile=$(COVERAGE_FILE)

## Export coverage to html and Open coverage report in browser
coverage: test
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML_FILE)
	go tool cover -html=$(COVERAGE_FILE)

## Run golangci-lint
lint:
	golangci-lint run ./...

## Generate swagger docs (swag init)
swag:
	swag init -g main.go -o api --parseDependency --parseInternal


## Build the binary
build:
	go build -o bin/$(APP_NAME) .

## Run the app (local)
run: 
	go run .


## Run inside docker container (go run .)
mysql-run:
	docker compose up -d mysql


## Start docker-compose services (DB, redis, etc.)
up:
	docker-compose up --build

## Stop docker-compose services
down:
	docker-compose down -v

## Air live reload
air:
	air -c .air.toml

## Atlas diff - Usage: make diff NAME=migration_name
diff:
ifndef NAME
	$(error NAME is required. Usage: make diff NAME=migration_name)
endif
	atlas migrate diff $(NAME) --env local

migrate:
	atlas migrate apply --env local

migrate-status:
	atlas migrate status --env local 