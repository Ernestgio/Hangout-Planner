# Makefile for file service

APP_NAME=file
COVERAGE_FILE=test/coverage
COVERAGE_HTML_FILE=test/coverage.html
AWSLOCAL=aws --endpoint-url=http://localhost:4566 --profile localstack

.PHONY: test coverage lint build run
.PHONY: localstack-up awslocal mysql-run
.PHONY: up down
.PHONY: nginx-reload nginx-test nginx-logs
.PHONY: air diff migrate migrate-status

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
	aws --endpoint-url=http://localhost:4566 --profile localstack $(ARGS)

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
	atlas migrate diff $(NAME) \
  		--to file://migrations \
  		--dev-url "docker://mysql/8/dev" \
  		--format '{{ sql . "  " }}'

migrate:
	atlas migrate apply --env local

migrate-status:
	atlas migrate status --env local 


## Telemetry services
## Otel tollector
otel-run:
	docker compose up -d otelcollector

##prometheus
prometheus-run:
	docker compose up -d prometheus

## Tempo
tempo-run:
	docker compose up -d tempo

#grafana
grafana-run:
	docker compose up -d grafana