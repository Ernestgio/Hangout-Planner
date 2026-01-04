# Makefile for file service

APP_NAME=file
COVERAGE_FILE=test/coverage

.PHONY: test lint

## Run tests with coverage
test:
	go test ./internal/... -coverprofile=$(COVERAGE_FILE)

## Run golangci-lint
lint:
	golangci-lint run ./...