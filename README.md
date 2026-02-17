# Codebase Golang

This is a Go boilerplate/codebase project designed for learning and rapid API development.

This repository follows a clean architecture pattern and Domain-Driven Design (DDD) principles, inspired by established Go codebase templates: [1](https://github.com/fahminlb33/devoria1-wtc-backend).

## Features

- ✅ Clean Architecture & Domain-Driven Design (DDD)
- ✅ RESTful API with Fiber framework
- ✅ PostgreSQL database with GORM ORM
- ✅ Redis caching support
- ✅ JWT & Basic Authentication
- ✅ Structured logging with Zap
- ✅ Configuration management with Viper
- ✅ Auto-generated Swagger documentation
- ✅ Unit testing with Testify and Mockery
- ✅ Health check endpoints
- ✅ Graceful shutdown
- ✅ Docker support

## Libraries Used

This project uses several libraries to facilitate development and ensure code quality:

- **Fiber**: Fast and lightweight web framework for building APIs in Go, used for routing and handling HTTP requests
- **GORM**: Feature-rich ORM library for database interactions with PostgreSQL
- **Viper**: Configuration management library for handling environment variables and config files
- **Zap**: High-performance structured logger for application events
- **Testify**: Testing toolkit providing assertions and mocking capabilities
- **Mockery**: Mock code auto-generator for creating test doubles
- **Swagger**: API documentation generator with interactive UI
- **Redis**: In-memory data structure store for caching

## Project Structure

The project follows a clean architecture structure:

```
go-codebase/
│
├── cmd/main/              # Application entry point
│   ├── bootstrap.go       # Application initialization and dependency injection
│   └── main.go            # Main function and server setup
│
├── pkg/                   # Reusable library code (can be imported by external projects)
│   ├── authentication/    # Authentication utilities (JWT, Basic Auth, Password)
│   ├── binding/           # Request binding helpers
│   ├── database/          # Database connection and utilities
│   ├── deps/              # Dependency injection container
│   ├── health/            # Health check handlers
│   ├── logger/            # Logging utilities
│   ├── middleware/        # HTTP middlewares
│   ├── redis/             # Redis client setup
│   ├── utils/             # Common utilities
│   ├── validator/         # Request validation
│   └── wrapper/           # Response wrapper utilities
│
├── internal/              # Private application code (domain logic)
│   └── example/           # Example domain (book management)
│       ├── handler/       # HTTP handlers
│       ├── repository/    # Data access layer
│       ├── schema/        # Request/response schemas
│       └── usecase/       # Business logic
│
├── model/                 # Database models
│   ├── book.go            # Book model
│   └── user.go            # User model
│
├── config/                # Configuration management
│   ├── config.go          # Config loader
│   ├── default.go         # Default values
│   └── type.go            # Config type definitions
│
├── api/                   # API documentation (auto-generated)
│   ├── docs.go            # Swagger documentation code
│   ├── swagger.json       # Swagger JSON spec
│   └── swagger.yaml       # Swagger YAML spec
│
├── db/                    # Database migrations
│   └── migration/         # SQL migration files
│
├── Dockerfile             # Production Docker image
├── Makefile              # Build and development commands
├── go.mod                # Go module dependencies
└── .env.example          # Environment variables template
```

## Getting Started

### Prerequisites

- Go 1.23 or higher
- PostgreSQL 14 or higher
- Redis 7 or higher (optional)
- Make (optional, for using Makefile commands)

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/Alwanly/go-codebase.git
   cd go-codebase
   ```

2. **Install dependencies**

   ```bash
   make install
   # or manually:
   go mod download
   ```

3. **Set up environment variables**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Configure database**

   Update the following in your `.env` file:
   ```env
   POSTGRES_URI=postgres://user:password@localhost:5432/yourdb
   REDIS_URI=redis://localhost:6379
   ```

5. **Generate RSA keys for JWT (optional)**

   ```bash
   # Generate private key
   openssl genrsa -out private.pem 2048
   
   # Generate public key
   openssl rsa -in private.pem -pubout -out public.pem
   
   # Copy keys to .env as single-line base64 or load them in code
   ```

### Running the Application

**Development mode with hot reload:**

```bash
make dev
```

**Standard mode:**

```bash
make run
```

**Build binary:**

```bash
make build
./app
```

**Using Docker:**

```bash
docker build -t go-codebase .
docker run -p 9000:9000 --env-file .env go-codebase
```

### API Documentation

Once the application is running in development mode, access the Swagger documentation at:

```
http://localhost:9000/swagger/index.html
```

### Health Checks

The application provides three health check endpoints:

- `GET /health` - General health check
- `GET /ready` - Readiness probe (for Kubernetes)
- `GET /live` - Liveness probe (for Kubernetes)

## Development

### Generating API Documentation

```bash
make docs
```

This uses `swag` to scan your code and generate Swagger documentation in the `./api` directory.

### Running Tests

```bash
make test
```

### Generating Test Coverage

```bash
make coverage
```

Coverage report will be generated in `coverage/coverage.html`.

### Generating Mocks

```bash
make mock
```

### Code Formatting

```bash
make format
```

### Linting

```bash
make lint
```

## API Endpoints

### Authentication

All `/books/v1/*` endpoints require JWT authentication.

**Example book endpoints:**

- `POST /books/v1/` - Create a new book
- `GET /books/v1/` - List all books
- `GET /books/v1/:id` - Get book by ID
- `PUT /books/v1/:id` - Update book
- `DELETE /books/v1/:id` - Delete book

## Environment Variables

Key environment variables (see `.env.example` for complete list):

| Variable | Description | Default |
|----------|-------------|---------|
| `ENV` | Environment (development/production) | development |
| `PORT` | HTTP server port | 9000 |
| `SERVICE_NAME` | Service name for logging | go-codebase |
| `POSTGRES_URI` | PostgreSQL connection string | Required |
| `REDIS_URI` | Redis connection string | Optional |
| `JWT_ISSUER` | JWT token issuer | codebase |
| `JWT_AUDIENCE` | JWT token audience | codebase |
| `PRIVATE_KEY` | RSA private key for JWT signing | Required |
| `PUBLIC_KEY` | RSA public key for JWT verification | Required |

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Acknowledgments

Inspired by clean architecture patterns and go codebase from the Go community.

