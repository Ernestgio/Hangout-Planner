FROM golang:1.24.11-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/file-service .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/file-service ./bin/file-service

EXPOSE 9001

RUN adduser -D appuser
USER appuser

CMD ["./bin/file-service"]
