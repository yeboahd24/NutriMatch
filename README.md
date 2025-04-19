# NutriMatch

NutriMatch is a personalized nutrition recommendation system built on the OpenNutrition dataset. It provides rule-based food recommendations based on user profiles and preferences.

## Features

- User profile management with health parameters, dietary restrictions, and preferences
- Food database management using the OpenNutrition dataset
- Rule-based food recommendation engine
- RESTful API for client applications
- Authentication and authorization
- User data management and privacy controls

## Technology Stack

- **Backend**: Go (Golang) with Chi router
- **Database**: PostgreSQL
- **Database Access**: sqlc for type-safe SQL
- **Configuration**: Viper
- **Authentication**: JWT with secure token management
- **API Documentation**: Swagger/OpenAPI
- **Logging**: Zerolog
- **Security**: OWASP security practices, rate limiting, input validation
- **Containerization**: Docker
- **Deployment**: Kubernetes-ready configuration

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Docker and Docker Compose (optional)
- OpenNutrition dataset

### Installation

1. Clone the repository:

```bash
git clone https://github.com/yeboahd24/nutrimatch.git
cd nutrimatch
```

2. Install dependencies:

```bash
go mod download
```

3. Create a `.env` file based on `.env.example`:

```bash
cp .env.example .env
```

4. Set up the database:

```bash
# Run PostgreSQL using Docker
docker-compose up -d db

# Run migrations
go run scripts/migrate.go up
```

5. Import the OpenNutrition dataset:

```bash
go run scripts/import.go /path/to/opennutrition_foods.tsv
```

6. Run the application:

```bash
go run cmd/api/main.go
```

### Using Docker

You can also run the entire application using Docker Compose:

```bash
docker-compose up -d
```

## API Documentation

API documentation is available at `/swagger/index.html` when the server is running.

## Development

### Running Tests

```bash
go test ./...
```

### Generating SQL Code

```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate code
sqlc generate
```

### Running Migrations

```bash
# Apply migrations
go run scripts/migrate.go up

# Rollback last migration
go run scripts/migrate.go down
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- OpenNutrition dataset: https://www.opennutrition.app
