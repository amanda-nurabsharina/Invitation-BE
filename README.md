# Invitation API

Backend API for Invitation Online built with Go and Fiber framework.

## Project Structure

```
BE Invitation/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── database/
│   │   └── database.go         # Database connection and migrations
│   ├── domain/
│   │   ├── auth/
│   │   │   └── auth.go         # Auth domain models
│   │   └── user/
│   │       ├── user.go         # User domain model
│   │       └── profile.go      # User profile model
│   ├── handler/
│   │   └── auth/
│   │       └── auth_handler.go # Auth HTTP handlers
│   ├── middleware/
│   │   └── auth_middleware.go  # Authentication middleware
│   ├── repository/
│   │   ├── auth/
│   │   │   └── auth_repository.go    # Auth repository
│   │   └── user/
│   │       └── user_repository.go    # User repository
│   └── service/
│       ├── auth/
│       │   └── auth_service.go       # Auth business logic
│       ├── seed_initial_user/
│       │   └── seed_initial_user_service.go  # Initial user seeder
│       └── user/
│           └── user_service.go       # User business logic
├── pkg/
│   ├── config/
│   │   └── config.go           # Configuration management
│   └── utils/
│       ├── jwt.go              # JWT utilities
│       └── password.go         # Password hashing utilities
├── .env.example
├── .gitignore
├── go.mod
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher

### Setup

1. Clone the repository

2. Copy environment file:

```bash
cp .env.example .env
```

3. Update `.env` with your database credentials

4. Create the database:

```sql
CREATE DATABASE invitation_db;
```

5. Install dependencies:

```bash
go mod tidy
```

6. Run the application:

```bash
go run cmd/main.go
```

## API Endpoints

### Public Endpoints

| Method | Endpoint                     | Description          |
| ------ | ---------------------------- | -------------------- |
| GET    | `/health`                    | Health check         |
| POST   | `/api/v1/auth/register`      | Register new user    |
| POST   | `/api/v1/auth/login`         | Login user           |
| POST   | `/api/v1/auth/refresh-token` | Refresh access token |

### Protected Endpoints (Requires Authentication)

| Method | Endpoint               | Description      |
| ------ | ---------------------- | ---------------- |
| POST   | `/api/v1/auth/logout`  | Logout user      |
| GET    | `/api/v1/auth/profile` | Get user profile |

### Admin Endpoints (Requires Admin Role)

| Method | Endpoint                  | Description     |
| ------ | ------------------------- | --------------- |
| GET    | `/api/v1/admin/dashboard` | Admin dashboard |

## Authentication

The API uses JWT Bearer token authentication. Include the token in the Authorization header:

```
Authorization: Bearer <access_token>
```

## Default Admin User

When running in development mode, an initial admin user is created:

- **Email:** admin@invitation.com
- **Username:** admin
- **Password:** admin1234
- **Role:** admin

## Request/Response Examples

### Register

**Request:**

```json
POST /api/v1/auth/register
{
  "email": "user@example.com",
  "username": "newuser",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "role": "user"
}
```

**Response:**

```json
{
  "error": false,
  "message": "User registered successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "abc123...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

### Login

**Request:**

```json
POST /api/v1/auth/login
{
  "username": "admin",
  "password": "admin1234"
}
```

**Response:**

```json
{
  "error": false,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "abc123...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

## Docker

### Production

Build and run with Docker Compose:

```bash
# Start all services (app + PostgreSQL)
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop all services
docker-compose down
```

### Development

Development mode with hot-reload and Adminer database UI:

```bash
# Start development environment
docker-compose -f docker-compose.dev.yml up -d

# View logs
docker-compose -f docker-compose.dev.yml logs -f app

# Rebuild and restart
docker-compose -f docker-compose.dev.yml up -d --build

# Stop development environment
docker-compose -f docker-compose.dev.yml down
```

**Development Services:**

- **API:** http://localhost:8080
- **Adminer (Database UI):** http://localhost:8081

### Using Makefile

```bash
# Production
make docker-up        # Start production
make docker-down      # Stop production
make docker-logs      # View logs

# Development
make docker-dev-up    # Start development
make docker-dev-down  # Stop development
make docker-dev-logs  # View logs

# Build only
make docker-build     # Build Docker image
```

### Environment Variables

Set environment variables in production:

```bash
# Create .env file for production
DB_USER=postgres
DB_PASSWORD=your-secure-password
DB_NAME=invitation_db
JWT_SECRET_KEY=your-secure-jwt-secret
```

## Local Development (without Docker)

```bash
# Install dependencies
go mod tidy

# Run with hot-reload (requires Air)
go install github.com/air-verse/air@latest
air -c .air.toml

# Or run directly
go run cmd/main.go
```

## License

MIT License
