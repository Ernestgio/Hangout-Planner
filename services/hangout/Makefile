# Makefile for Hangout Planner

APP_NAME=hangout
COVERAGE_FILE=test/coverage
COVERAGE_HTML_FILE=test/coverage.html
AWSLOCAL=aws --endpoint-url=http://localhost:4566 --profile localstack

.PHONY: all test coverage lint swag build run localstack-up awslocal mysql-run up down nginx-reload nginx-test nginx-logs air diff migrate migrate-status

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

## Localstack
localstack-up:
	docker compose up -d localstack

## Run AWS CLI commands against LocalStack - Usage: make awslocal ARGS="s3 ls" or ARGS="s3 cp file.txt s3://bucket/"
awslocal:
	aws --endpoint-url=http://localhost:4566 --profile localstack

## Run inside docker container (go run .)
mysql-run:
	docker compose up -d mysql


## Start docker-compose services (DB, redis, etc.)
up:
	docker-compose up --build

## Stop docker-compose services
down:
	docker-compose down

## Restart NGINX only (after config changes)
nginx-reload:
	docker exec hangout-nginx nginx -s reload

## Test NGINX configuration
nginx-test:
	docker exec hangout-nginx nginx -t

## View NGINX logs
nginx-logs:
	docker logs hangout-nginx -f

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