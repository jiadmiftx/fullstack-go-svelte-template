# Docker Containerization Guide

## Overview

This guide provides comprehensive instructions for containerizing Go applications using best practices that are project-agnostic and environment-aware. The patterns described here work with any Go project using configuration management and environment variables.

## Core Principles

### 1. Multi-Stage Builds
- **Builder Stage:** Compile the application with all build tools
- **Runtime Stage:** Minimal image with only the compiled binary and runtime dependencies

### 2. Configuration Awareness
- Read all configuration from environment variables
- Support `.env` files for local development
- Never hardcode ports, addresses, or credentials

### 3. Security Best Practices
- Use minimal base images (Alpine)
- Run as non-root user when possible
- Bind to localhost in development
- Use specific version tags, not `latest`

---

## Dockerfile Template

```dockerfile
# ===== BUILD STAGE =====
# IMPORTANT: Always match Go version with go.mod file
# Use ARG to pass Go version dynamically
ARG GO_VERSION=1.21
FROM golang:${GO_VERSION}-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev wget

# Set working directory
WORKDIR /app

# Copy dependency files first (for better layer caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate any required files (Swagger docs, etc.)
# IMPORTANT: Always use production swagger config in containers
# Install make if not available and generate production swagger
RUN apk add --no-cache make && \
    if [ -f Makefile ] && grep -q "swagger-prod" Makefile; then \
        make swagger-prod; \
    else \
        go run github.com/swaggo/swag/cmd/swag@latest init; \
    fi

# Build the application
# CGO_ENABLED=1 is required for some packages (database drivers, etc.)
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# ===== RUNTIME STAGE =====
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget

# Create app user (optional but recommended for production)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Create necessary directories
RUN mkdir -p /app/data && chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy any static files if needed
# COPY --from=builder /app/static ./static

# Expose port (for documentation purposes only)
# Actual port is configured via environment variable
EXPOSE 8080

# Health check (optional but recommended)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT:-8080}/ping || exit 1

# Run the application
CMD ["./main"]
```

### Key Dockerfile Features

1. **Layer Optimization:** Dependencies copied before source code for better caching
2. **Minimal Runtime:** Only essential packages in runtime stage
3. **Security:** Non-root user execution (recommended for production)
4. **Health Checks:** Built-in container health monitoring
5. **Flexible Ports:** Uses environment variable for port configuration

---

## docker-compose.yml Template

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ${PROJECT_NAME:-app-container}
    restart: unless-stopped

    # Environment configuration
    env_file:
      - .env
    environment:
      # Override Docker-specific settings
      DB_HOST: ${DB_HOST:-host.docker.internal}
      DB_PORT: ${DB_PORT:-5432}
      # Add other required environment variables

    # Port binding - localhost only for security
    ports:
      - "127.0.0.1:${PORT:-8080}:${PORT:-8080}"

    # Volume mounts for persistence
    volumes:
      - ./data:/app/data
      - ./logs:/app/logs

    # Host networking for external service access
    extra_hosts:
      - "host.docker.internal:host-gateway"

    # Network configuration
    networks:
      - app-network

    # Resource limits (recommended for production)
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M

    # Health check
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:${PORT:-8080}/ping"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

networks:
  app-network:
    driver: bridge
    # External network option for connecting to other stacks
    # external: true
```

### Key docker-compose Features

1. **Dynamic Naming:** Uses `${PROJECT_NAME}` for container naming
2. **Localhost Binding:** Only exposes service to localhost for security
3. **Environment Management:** Supports `.env` file with defaults
4. **Health Monitoring:** Built-in health checks
5. **Resource Management:** CPU and memory limits
6. **Network Flexibility:** Can use external networks for multi-stack communication

---

## Environment Configuration

### .env File Template

```bash
# ===== APPLICATION CONFIGURATION =====
# Project-specific settings
# PORT and other settings - PROJECT_NAME is auto-detected from go.mod
PORT=8080

# ===== DATABASE CONFIGURATION =====
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=myuser
DB_PASSWORD=mypassword
DB_NAME=mydatabase

# ===== SECURITY CONFIGURATION =====
JWT_SECRET=your-super-secret-jwt-key-here
AES_KEY=your-32-character-aes-key-here
X_AUTH_CRON=your-cron-auth-token

# ===== CORS CONFIGURATION =====
REST_CORS_ORIGIN=http://localhost:3000

# ===== DEBUG/DEVELOPMENT =====
REST_DEBUG_MODE=true

# ===== OPTIONAL SERVICES =====
# Elasticsearch
ELASTIC_HOST=http://localhost:9200
ELASTIC_USERNAME=elastic
ELASTIC_PASSWORD=password
ELASTIC_INDEX=my-index

# Redis (if needed)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# External APIs
API_BASE_URL=https://api.example.com
API_KEY=your-api-key-here
```

### Environment Variable Patterns

1. **Required Variables:** Always have sensible defaults in docker-compose
2. **Security:** Never commit actual secrets to version control
3. **Development vs Production:** Use different `.env` files for different environments
4. **Documentation:** Include comments explaining each variable

---

## Go Application Configuration

### config/init.go Template

```go
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// AppConfig holds all application configuration
type AppConfig struct {
	ServiceName string
	Version     string
	Port        int
	Database    DatabaseConfig
	Security    SecurityConfig
	// Add other config sections as needed
}

// Environment variable helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Invalid integer value for %s, using default: %d", key, defaultValue)
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		log.Printf("Warning: Invalid boolean value for %s, using default: %t", key, defaultValue)
	}
	return defaultValue
}

// Load configuration from environment
func LoadConfig() *AppConfig {
	// Load .env file if it exists (for local development)
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}

	return &AppConfig{
		ServiceName: getEnv("SERVICE_NAME", "My Go Application"),
		Version:     getEnv("APP_VERSION", "1.0.0"),
		Port:        getEnvAsInt("PORT", 8080),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			Username: getEnv("DB_USERNAME", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", ""),
		},
		Security: SecurityConfig{
			JWTSecret: getEnv("JWT_SECRET", ""),
			AESKey:    getEnv("AES_KEY", ""),
		},
	}
}
```

### main.go Integration

```go
package main

import (
	"fmt"
	"log"

	"your-module-path/config"
	"your-module-path/database"
	"your-module-path/routes"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialize database
	db, err := database.Init(cfg)
	if err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}

	// Initialize routes
	router := routes.InitRoutes(db)

	// Start server with configured port
	port := fmt.Sprintf("%d", cfg.Port)
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📚 API Documentation: http://localhost:%s/swagger/index.html", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}
```

---

## Docker Commands Cheat Sheet

### Development Commands

```bash
# Build the image
docker-compose build

# Start the service
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop the service
docker-compose down

# Rebuild and restart
docker-compose up -d --build
```

### Production Commands

```bash
# Use production environment file
docker-compose --env-file .env.production up -d

# Scale the service
docker-compose up -d --scale app=3

# Update with zero downtime
docker-compose up -d --no-deps app
```

### Debugging Commands

```bash
# Check container status
docker-compose ps

# Execute commands in container
docker-compose exec app sh

# View port bindings
docker-compose port app 8080

# Check resource usage
docker stats app-container
```

### Maintenance Commands

```bash
# Clean up unused images
docker image prune -f

# Remove all containers
docker-compose down --remove-orphans

# Force rebuild without cache
docker-compose build --no-cache
```

---

## Multi-Environment Configuration

### File Structure

```
project/
├── .env.example          # Template for new environments
├── .env.local           # Local development secrets
├── .env.development     # Development environment
├── .env.staging         # Staging environment
├── .env.production      # Production environment
├── docker-compose.yml   # Base configuration
├── docker-compose.override.yml  # Development overrides
└── Dockerfile
```

### Environment-Specific Compose Files

**docker-compose.override.yml** (Development)
```yaml
version: '3.8'

services:
  app:
    volumes:
      - .:/app                    # Live code mounting
      - /app/vendor               # Exclude vendor
    environment:
      REST_DEBUG_MODE: "true"
      GIN_MODE: "debug"
```

**docker-compose.prod.yml** (Production)
```yaml
version: '3.8'

services:
  app:
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '2.0'
          memory: 1G
    environment:
      REST_DEBUG_MODE: "false"
      GIN_MODE: "release"
```

---

## Best Practices Checklist

### Dockerfile
- [ ] Use multi-stage builds
- [ ] **Use dynamic Go version matching with go.mod** (via ARG)
- [ ] Use specific version tags (not `latest`)
- [ ] Copy `.dockerignore` file
- [ ] Optimize layer caching
- [ ] Use non-root user in production
- [ ] Include health checks

### docker-compose.yml
- [ ] Use environment variables for all configuration
- [ ] **Use dynamic container naming based on project name**
- [ ] Bind to localhost in development
- [ ] Set resource limits
- [ ] Use external networks for multi-stack communication
- [ ] Include restart policies

### Swagger Documentation
- [ ] **Always use production swagger configuration in containers**
- [ ] Use `make swagger-prod` for container builds
- [ ] Use `make swagger-local` only for local development
- [ ] Set correct host, BasePath, and schemes for each environment
- [ ] Automate swagger generation in build process

### Security
- [ ] Never commit secrets to version control
- [ ] Use `.env.example` as template
- [ ] Run as non-root user
- [ ] Use minimal base images
- [ ] Regular security updates

### Performance
- [ ] Implement proper health checks
- [ ] Set appropriate resource limits
- [ ] Use volume mounts for persistence
- [ ] Optimize Docker layer caching

---

## Common Issues and Solutions

### Port Already in Use
```bash
# Find and kill process using the port
lsof -ti :8080 | xargs kill -9

# Or use different port
PORT=8081 docker-compose up
```

### Permission Issues
```bash
# Fix file permissions for volumes
sudo chown -R $USER:$USER ./data

# Or use Docker volumes instead of bind mounts
```

### Environment Variables Not Loading
```bash
# Check if .env file exists
ls -la .env*

# Verify environment variables
docker-compose config
```

### Build Issues
```bash
# Clear build cache
docker builder prune -a

# Rebuild without cache
docker-compose build --no-cache
```

---

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: Build and Deploy

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Login to Registry
        uses: docker/login-action@v2
        with:
          registry: ${{ secrets.REGISTRY_URL }}
          username: ${{ secrets.REGISTRY_USERNAME }}
          password: ${{ secrets.REGISTRY_PASSWORD }}

      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: ${{ secrets.REGISTRY_URL }}/myapp:${{ github.sha }}

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to production
        run: |
          docker-compose -f docker-compose.prod.yml pull
          docker-compose -f docker-compose.prod.yml up -d
```

---

## Monitoring and Logging

### Logging Best Practices

```go
// Use structured logging
log.Printf("Server starting on port %d", cfg.Port)

// Include request ID in logs for tracing
log.Printf("[%s] Processing request", requestID)

// Log configuration on startup (without secrets)
log.Printf("Database host: %s, Port: %d", cfg.DB.Host, cfg.DB.Port)
```

### Health Check Endpoints

```go
// Add to your routes
router.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "healthy",
        "timestamp": time.Now(),
        "version": cfg.Version,
    })
})

router.GET("/ping", func(c *gin.Context) {
    c.String(200, "pong")
})
```

---

## Advanced Docker Patterns

### 1. Multi-Architecture Builds

```dockerfile
# For ARM64 and AMD64 support
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder
ARG TARGETPLATFORM
ARG BUILDPLATFORM

# Build for specific target platform
RUN CGO_ENABLED=1 GOOS=linux GOARCH=$TARGETARCH go build -o main .
```

### 2. Dynamic Go Version Matching

```dockerfile
# Extract Go version from go.mod automatically
ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine AS builder
```

**Build Commands with Auto-Detection:**

```bash
# Extract Go version and project name from go.mod
GO_VERSION=$(grep "^go " go.mod | awk '{print $2}' | cut -d'.' -f1,2)
PROJECT_NAME=$(grep "^module " go.mod | awk '{print $2}' | awk -F'/' '{print $NF}')
docker build \
  --build-arg GO_VERSION=$GO_VERSION \
  -t $PROJECT_NAME:latest .

# One-liner for complete setup
PROJECT_NAME=$(grep "^module " go.mod | awk '{print $2}' | awk -F'/' '{print $NF}') && \
echo "PROJECT_NAME=$PROJECT_NAME" > .env && \
docker-compose up -d --build
```

# Or using a script
#!/bin/bash
GO_VERSION=$(grep "^go " go.mod | awk '{print $2}' | cut -d'.' -f1,2)
docker build \
  --build-arg GO_VERSION=$GO_VERSION \
  --build-arg VERSION=$(git describe --tags --always) \
  --build-arg COMMIT_SHA=$(git rev-parse HEAD) \
  -t myapp:$(git describe --tags --always) .
```

**docker-compose.yml with Dynamic Configuration:**

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      args:
        GO_VERSION: ${GO_VERSION:-$(grep "^go " go.mod | awk '{print $2}' | cut -d'.' -f1,2)}
        VERSION: ${VERSION:-dev}
        COMMIT_SHA: ${COMMIT_SHA:-local}
    container_name: ${PROJECT_NAME:-$(grep "^module " go.mod | awk '{print $2}' | awk -F'/' '{print $NF}')}
```

**Quick Commands for Project Setup:**

```bash
# One-liner to set up and run with auto-detected project name
export PROJECT_NAME=$(grep "^module " go.mod | awk '{print $2}' | awk -F'/' '{print $NF}') && \
docker-compose up -d --build

# Or run with specific project name
PROJECT_NAME=my-custom-name docker-compose up -d
```

### 3. Build Args for Versioning

```dockerfile
ARG VERSION=dev
ARG COMMIT_SHA
LABEL version=${VERSION}
LABEL commit=${COMMIT_SHA}
```

### 4. Automated Project Detection Script

**build.sh:**
```bash
#!/bin/sh

# Auto-detect project name from go.mod
if [ -f "go.mod" ]; then
    # Extract module name and get last part (project name)
    MODULE_NAME=$(grep "^module " go.mod | awk '{print $2}')
    PROJECT_NAME=$(echo $MODULE_NAME | awk -F'/' '{print $NF}')
    echo "Detected project name: $PROJECT_NAME"

    # Auto-detect Go version from go.mod
    GO_VERSION=$(grep "^go " go.mod | awk '{print $2}' | cut -d'.' -f1,2)
    echo "Detected Go version: $GO_VERSION"
else
    PROJECT_NAME="my-app"
    GO_VERSION="1.21"
    echo "Using defaults - Project: $PROJECT_NAME, Go: $GO_VERSION"
fi

# Auto-detect git info
if command -v git >/dev/null 2>&1 && [ -d ".git" ]; then
    VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
    COMMIT_SHA=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
    echo "Git version: $VERSION"
    echo "Git commit: $COMMIT_SHA"
else
    VERSION="dev"
    COMMIT_SHA="unknown"
    echo "No git repository found, using defaults"
fi

# Create/update .env with project name
if [ ! -f ".env" ]; then
    echo "Creating .env file..."
    cat > .env << EOF
PROJECT_NAME=$PROJECT_NAME
PORT=8080
EOF
else
    # Update PROJECT_NAME in .env if exists
    if grep -q "PROJECT_NAME=" .env; then
        sed -i.bak "s/PROJECT_NAME=.*/PROJECT_NAME=$PROJECT_NAME/" .env
    else
        echo "PROJECT_NAME=$PROJECT_NAME" >> .env
    fi
fi

# Build the image
docker build \
    --build-arg GO_VERSION=$GO_VERSION \
    --build-arg VERSION=$VERSION \
    --build-arg COMMIT_SHA=$COMMIT_SHA \
    -t $PROJECT_NAME:$VERSION \
    -t $PROJECT_NAME:latest .

echo "Built image: $PROJECT_NAME:$VERSION"
echo "To run: docker-compose up -d"
echo "Container will be named: $PROJECT_NAME"
```

### 5. Custom Entrypoint Script

```dockerfile
# Copy entrypoint script
COPY --from=builder /app/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["./main"]
```

**entrypoint.sh:**
```bash
#!/bin/sh
# Wait for database (if using PostgreSQL, MySQL, etc.)
# while ! nc -z $DB_HOST $DB_PORT; do sleep 1; done

echo "Starting application..."
exec "$@"
```

---

## Production Deployment Patterns

### 1. Docker Swarm Example

```yaml
version: '3.8'

services:
  app:
    image: your-registry/app:${VERSION:-latest}
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
        order: start-first
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
    networks:
      - app-network
    secrets:
      - db_password
      - jwt_secret

secrets:
  db_password:
    external: true
  jwt_secret:
    external: true
```

### 2. Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-app
  template:
    metadata:
      labels:
        app: go-app
    spec:
      containers:
      - name: app
        image: your-registry/app:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: host
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /ping
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

---

## Performance Optimization

### 1. Image Size Reduction

```dockerfile
# Use distroless for even smaller images
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/main .
USER nonroot:nonroot

EXPOSE 8080
CMD ["./main"]
```

### 2. Build Cache Optimization

```dockerfile
# Separate dependency build and app build
FROM golang:1.21-alpine AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY --from=deps /app/go.mod /app/go.sum ./
COPY --from=deps /app/go/pkg/mod /app/go/pkg/mod
COPY . .
RUN CGO_ENABLED=1 go build -o main .
```

### 3. Parallel Builds

```dockerfile
# Build with parallel compilation
RUN go build -ldflags="-s -w" -o main .
```

---

## Security Hardening

### 1. Security Scanning

```bash
# Scan for vulnerabilities
docker scan your-app:latest

# Use Trivy for comprehensive scanning
trivy image your-app:latest
```

### 2. Minimal Permissions

```dockerfile
# Drop all capabilities and add only what's needed
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S appgroup && \
    adduser -S -G appgroup appuser

# Copy as non-root and set proper permissions
COPY --chown=appuser:appgroup --from=builder /app/main .
USER appuser
```

### 3. Runtime Security

```yaml
# docker-compose.yml with security options
services:
  app:
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp:noexec,nosuid,size=100m
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE
```

---

## Troubleshooting Guide

### 1. Debug Commands

```bash
# Check container details
docker inspect <container-id>

# View all environment variables
docker-compose exec app printenv

# Check network connectivity
docker-compose exec app ping google.com

# View filesystem
docker-compose exec app ls -la /app
```

### 2. Common Debugging Scenarios

```bash
# Port binding issues
netstat -tulpn | grep :8080

# Permission issues
docker-compose exec app id
docker-compose exec app ls -la /app

# Memory issues
docker stats <container-id>

# Log issues
docker-compose logs --tail=100 app

# Go version mismatch issues
go version  # Check local Go version
grep "^go " go.mod  # Check required Go version in go.mod
docker run --rm golang:1.21-alpine go version  # Check Docker Go version
```

---

## Conclusion

This guide provides a comprehensive framework for containerizing Go applications that can be adapted to any project with configuration management and environment variables.

### Key Takeaways:
1. **Always use multi-stage builds** to minimize image size
2. **Always match Go version with go.mod** - use dynamic ARG for flexibility
3. **Auto-detect project name from go.mod** - no more hardcoded container names
4. **Never hardcode configuration** - use environment variables
5. **Implement health checks** for production readiness
6. **Follow security best practices** from development to production
7. **Monitor and log** appropriately for debugging
8. **Test thoroughly** in different environments

### Next Steps:
1. Adapt these templates to your specific project needs
2. Set up CI/CD pipelines for automated builds
3. Implement monitoring and alerting
4. Regular security scans and updates
5. Documentation updates as your architecture evolves