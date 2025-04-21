# Swagger Troubleshooting Guide

This guide provides solutions to common issues with the Swagger/OpenAPI integration in the NutriMatch API.

## Common Issues

### 1. "Internal Server Error" when accessing `/swagger/doc.json`

**Symptoms:**
- Swagger UI shows "Fetch error" or "Internal Server Error" when trying to load
- The `/swagger/doc.json` endpoint returns a 500 error

**Solutions:**

1. **Check if the doc.json file exists:**
   ```bash
   ls -la docs/doc.json
   ```
   If it doesn't exist, generate it:
   ```bash
   make swagger
   ```

2. **Check file permissions:**
   ```bash
   chmod 644 docs/doc.json
   ```

3. **Manually serve the file:**
   If the automatic generation isn't working, you can manually create a doc.json file in the docs directory.

4. **Check server logs:**
   Look for any errors in the server logs that might indicate why the file isn't being served.

### 2. Swagger UI not loading

**Symptoms:**
- The Swagger UI page is blank or shows errors
- JavaScript console shows errors

**Solutions:**

1. **Check browser console for errors:**
   Open your browser's developer tools and check the console for any JavaScript errors.

2. **Verify the Swagger UI route:**
   Make sure the route for Swagger UI is correctly configured in the server.go file.

3. **Check if the index.html file exists:**
   ```bash
   ls -la docs/index.html
   ```
   If it doesn't exist, create it or copy it from the Swagger UI distribution.

### 3. API endpoints not showing in Swagger UI

**Symptoms:**
- Swagger UI loads but doesn't show any API endpoints
- The doc.json file is empty or missing endpoint definitions

**Solutions:**

1. **Check your annotations:**
   Make sure your handler functions have the correct Swagger annotations.

2. **Regenerate the Swagger documentation:**
   ```bash
   make swagger
   ```

3. **Check for syntax errors in annotations:**
   Look for any syntax errors in your Swagger annotations that might prevent proper generation.

4. **Manually update the doc.json file:**
   If automatic generation isn't working, you can manually edit the doc.json file to include your endpoints.

### 4. "Could not render this component" error

**Symptoms:**
- Swagger UI shows "Could not render this component" for some or all endpoints

**Solutions:**

1. **Check your model definitions:**
   Make sure all models referenced in your Swagger annotations are properly defined.

2. **Check for circular references:**
   Circular references in your model definitions can cause rendering issues.

3. **Simplify complex models:**
   Try simplifying complex model structures that might be causing rendering issues.

## Advanced Troubleshooting

### Direct File Access

You can directly access the Swagger files to check if they're being served correctly:

- `/swagger/doc.json` - Should return the OpenAPI specification
- `/swagger/ui/` - Should display the Swagger UI
- `/swagger-static/swagger.json` - Alternative location for the specification

### Manual Swagger Generation

If the automatic generation isn't working, you can manually install and run the Swagger tools:

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

### Debugging the Swagger Handler

You can add debug logging to the Swagger handler routes in server.go:

```go
r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
    logger.Info().Msg("Serving Swagger doc.json")
    w.Header().Set("Content-Type", "application/json")
    http.ServeFile(w, r, filepath.Join(swaggerDir, "doc.json"))
})
```

## Additional Resources

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
