# Stage 1: Build the Go application
FROM golang:1.22.5 AS build-stage

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api/main.go

# Stage 2: Development environment with air
FROM golang:1.22.5 AS dev-stage

WORKDIR /app

# Install dependencies
RUN apt-get update && apt-get install -y curl git && apt-get clean

# Install air
RUN go install github.com/air-verse/air@latest

# Copy the application files
COPY . .

# Install Go dependencies
RUN go mod download

# Expose the application port
EXPOSE 8080

# Command to run air for development
CMD ["air", "-c", ".air.toml"]

# Stage 3: Production environment
FROM alpine:3.20 AS prod-stage

WORKDIR /app

# Copy the built application binary from build-stage
COPY --from=build-stage /api /app/api

# Copy the frontend static files
COPY --from=build-stage /app/frontend /app/frontend

# Expose the application port
EXPOSE 8080

# Command to run the application binary
ENTRYPOINT ["/app/api"]

