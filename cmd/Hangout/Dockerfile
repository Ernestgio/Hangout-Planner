# Start from the official Golang image for building
FROM golang:alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o hangout-planner .

# Use a minimal image for running
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/hangout-planner .

# Expose port (change if your app uses a different port)
EXPOSE 8080

# Run the binary
CMD ["./hangout-planner"]