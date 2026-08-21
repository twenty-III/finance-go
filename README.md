# Finance Tracker Backend

A high-performance, robust REST and GraphQL API for personal financial tracking, written in Go. The project utilizes idiomatic Go directory structures (`cmd/` and `internal/`) and is organized using **feature-based slicing** for simplicity and maintainability.

## Features

- **Dual APIs**: Exposes both a RESTful JSON API (via Gorilla Mux) and a GraphQL API (via gqlgen) offering unified business logic.
- **Idiomatic Go Structure**: Uses feature-based slicing under `internal/` (e.g., `internal/user`, `internal/expense`) to group models, repositories, and handlers logically.
- **Robust Security**: Built-in stateless JWT-based authentication via HTTP-only cookies, robust password hashing with `bcrypt`, and password strength validation.
- **Distributed Rate Limiting**: Intelligent token-bucket rate-limiting middleware to guard against abuse.
- **PostgreSQL Persistence**: Fully-featured Postgres repository layer with comprehensive indexing and structured database migrations (using `golang-migrate`).
- **Docker-Ready**: Packaged with a multi-stage Dockerfile and a complete `docker-compose.yml` for instantaneous local deployment.

## Architecture

```mermaid
flowchart TD
    Client[Client / REST & GraphQL] --> Router[API Layer: Mux & Gqlgen]
    Router --> Middleware[Middleware: JWT Auth & IP Rate Limiting]
    Middleware --> Feature[Feature Slices: User, Expense]
    Feature --> DB[(PostgreSQL)]
```

### File Structure Overview

- `cmd/api/`: Application entry point (`main.go`).
- `internal/`: Private application code.
  - `user/`: All logic, handlers, services, and repositories related to Users.
  - `expense/`: All logic, handlers, services, and repositories related to Expenses.
  - `auth/`: Authentication services, JWT handling, and password hashing.
  - `server/`: HTTP routing, base handlers, and rate-limiting middleware.
  - `db/`: Database connection logic and migrations.
  - `graph/`: GraphQL schemas and resolvers.
  - `config/`: Environment configuration loading.

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL
- Docker & Docker Compose (optional, but recommended)

### Running Locally (Without Docker)

1. Clone this repository to your workspace.
2. Ensure Go and PostgreSQL are installed.
3. Configure your environment variables in a `.env` file (you can use the provided `.env` as a base).
4. Install dependencies:
   ```bash
   go mod tidy
   ```
5. Run the application (utilizes `air` for hot-reloading if you use `make run`):
   ```bash
   make run
   ```
   The API will be available at `http://localhost:8080`.

### Running with Docker

1. Ensure Docker and Docker Compose are installed.
2. Build and start the services:
   ```bash
   docker-compose up -d --build
   ```
3. The API will be available at `http://localhost:8080` and the Postgres instance will be bound to port `5432`.

## Documentation

Comprehensive API documentation is available in the `docs/` directory:

- [REST API Specifications](./docs/API_DEFINITION.md)
- [GraphQL API Schema](./docs/GRAPH_API_SCHEMA.md)
- [Database Design](./docs/DB_DESIGN.md)

## Testing

Run the full test suite using:

```bash
make test
```

## Contributing

Pull requests are always welcome! For major changes, please open an issue first to discuss your proposed modifications. Ensure to update tests as appropriate.


