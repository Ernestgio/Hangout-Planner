# Makefile for Hangout Planner

APP_NAME=hangout
COVERAGE_FILE=test/coverage
COVERAGE_HTML_FILE=test/coverage.html

.PHONY: all test lint swag build run clean docker-build docker-run up down

all: lint test build

## Run tests with coverage
test:
	go test ./... -coverprofile=$(COVERAGE_FILE)

## Export coverage to html and Open coverage report in browser
coverage: test
	go tool cover -html=$(COVERAGE_FILE) -o {COVERAGE_HTML_FILE}
	go tool cover -html=$(COVERAGE_FILE)

## Run golangci-lint
lint:
	golangci-lint run ./...

## Generate swagger docs (swag init)
swag:
	swag init -g main.go -o api


## Build the binary
build:
	@echo ">> Building $(APP_NAME)"
	go build -o bin/$(APP_NAME) ./cmd

## Run the app (local)
run: build
	@echo ">> Running $(APP_NAME)"
	./bin/$(APP_NAME)

## Clean build artifacts
clean:
	@echo ">> Cleaning up"
	rm -rf bin $(COVERAGE_FILE) internal/docs

## Build docker image
docker-build:
	@echo ">> Building Docker image $(DOCKER_IMAGE)"
	docker build -t $(DOCKER_IMAGE) .

## Run inside docker container (go run .)
docker-run:
	@echo ">> Running app inside Docker"
	docker run --rm -it -p 8080:8080 $(DOCKER_IMAGE)

## Start docker-compose services (DB, redis, etc.)
up:
	@echo ">> Starting docker-compose"
	docker-compose up -d

## Stop docker-compose services
down:
	@echo ">> Stopping docker-compose"
	docker-compose down