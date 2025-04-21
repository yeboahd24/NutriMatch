# NutriMatch API Documentation Guide

This guide explains how to use and maintain the OpenAPI/Swagger documentation for the NutriMatch API.

## Overview

NutriMatch uses [Swagger](https://swagger.io/) (OpenAPI) for API documentation. The documentation is generated from annotations in the Go code using the [swaggo/swag](https://github.com/swaggo/swag) package.

## Accessing the API Documentation

When the server is running, you can access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

This provides an interactive interface to explore and test the API endpoints.

## Setting Up Swagger

### Installing Dependencies

To install the required Swagger dependencies:

```bash
# Using make
make swagger-deps

# Or directly
./scripts/install_swagger_deps.sh
```

### Generating Documentation

To generate or update the Swagger documentation:

```bash
# Using make
make swagger

# Or directly
./scripts/generate_swagger.sh
```

This will parse the annotations in the code and generate the Swagger specification files in the `docs` directory.

## Adding Documentation to Endpoints

### General API Information

General API information is defined in the `cmd/api/main.go` file:

```go
// Package main is the entry point for the NutriMatch API server.
//
// @title NutriMatch API
// @version 1.0
// @description A personalized nutrition recommendation system built on the OpenNutrition dataset
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.nutrimatch.com/support
// @contact.email support@nutrimatch.com
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.
package main
```

### Endpoint Documentation

Each endpoint should be documented using Swagger annotations. Here's an example:

```go
// @Summary Get food by ID
// @Description Get detailed information about a specific food
// @Tags foods
// @Accept json
// @Produce json
// @Param id path string true "Food ID"
// @Success 200 {object} docs.Response{data=docs.FoodDetailResponse}
// @Failure 404 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /foods/{id} [get]
func (h *FoodHandler) GetFood(w http.ResponseWriter, r *http.Request) {
    // Implementation...
}
```

### Model Definitions

Models used in the API documentation are defined in `internal/api/docs/swagger_models.go`. These are used only for documentation purposes and don't affect the actual code.

## Best Practices

1. **Keep Documentation Updated**: Update the Swagger annotations whenever you change an endpoint.
2. **Document All Parameters**: Include all parameters, whether they're in the path, query, or request body.
3. **Include Examples**: Where possible, include examples to make the API easier to understand.
4. **Document Error Responses**: Document all possible error responses for each endpoint.
5. **Use Tags**: Group related endpoints using tags for better organization.

## Troubleshooting

If you encounter issues with Swagger generation:

1. **Check Annotations**: Ensure your annotations are correctly formatted.
2. **Check Model Definitions**: Make sure all models referenced in annotations are defined.
3. **Regenerate Documentation**: Run `make swagger` to regenerate the documentation.
4. **Check Console Output**: Look for error messages during generation.

## References

- [Swagger Specification](https://swagger.io/specification/)
- [swaggo/swag Documentation](https://github.com/swaggo/swag)
- [Annotation Format Guide](https://github.com/swaggo/swag#declarative-comments-format)
