# AI CDN Services - Project Architecture Documentation

## 📋 Table of Contents

1. [Overview](#overview)
2. [Clean Architecture Implementation](#clean-architecture-implementation)
3. [Project Structure](#project-structure)
4. [Layer Details](#layer-details)
5. [Data Flow](#data-flow)
6. [External Integrations](#external-integrations)
7. [Security Architecture](#security-architecture)
8. [Database Architecture](#database-architecture)
9. [API Architecture](#api-architecture)
10. [Deployment Architecture](#deployment-architecture)
11. [Monitoring & Observability](#monitoring--observability)
12. [Design Patterns Used](#design-patterns-used)
13. [Architecture Decisions](#architecture-decisions)
14. [Usage Guide for LLM](#usage-guide-for-llm)

## Overview

This is a **Clean Architecture-based Go application** that provides Content Delivery Network (CDN) management services with enterprise-grade features. The project demonstrates best practices for building scalable, maintainable, and testable Go applications.

### Key Characteristics

- **Clean Architecture Pattern** with strict layer separation
- **Interface-driven design** for loose coupling and testability
- **Domain models in each layer** to minimize inter-layer dependencies
- **Multi-database support** (PostgreSQL primary, SQLite for testing)
- **External service integrations** (Varnish, Ansible, AWS S3, SSL services)
- **Enterprise-grade security** with JWT authentication
- **Production-ready deployment** with Docker and Kubernetes
- **Comprehensive monitoring** with Elastic APM integration

## Clean Architecture Implementation

### Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│                    PRESENTATION LAYER                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │ Controllers  │  │ HTTP Routes  │  │ Middleware   │ │
│  │   (REST)     │  │   Config     │  │   (CORS,     │ │
│  │              │  │              │  │   JWT, APM)  │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────┐
│                    USE CASE LAYER                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   Business   │  │ Validation   │  │  Error       │ │
│  │   Logic      │  │   Rules      │  │  Handling    │ │
│  │              │  │              │  │              │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────┐
│                    REPOSITORY LAYER                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   Database   │  │   Data       │  │  Query       │ │
│  │   Models     │  │   Mapping    │  │  Builders    │ │
│  │              │  │              │  │              │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────┐
│                     SERVICE LAYER                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   External   │  │   Third      │  │  Message     │ │
│  │   APIs       │  │   Party      │  │  Broker      │ │
│  │  (External)  │  │   Services  │  │ (Async/Queue) │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### Dependency Rules

1. **Dependencies point inward**: Outer layers depend on inner layers
2. **Business logic independence**: Use cases don't depend on external details
3. **Interface segregation**: Services are defined as interfaces
4. **Domain model isolation**: Each layer has its own models to reduce coupling
5. **Single responsibility**: Each layer has one clear purpose

## Project Structure

### Directory Layout

```
ai-cdn-services/
├── main.go                     # Application entry point
├── config/                     # Configuration management
│   ├── app.go                 # Configuration structures
│   ├── init.go                # Environment loading
│   └── init_helpers.go        # Validation helpers
├── app/                        # Application layers (Clean Architecture)
│   ├── controller/            # Presentation layer (HTTP handlers)
│   ├── use_case/              # Business logic layer
│   └── repository/            # Data access layer
├── services/                   # External service integrations
│   ├── cdn_service/           # Varnish CDN integration
│   ├── ansbile_service/       # Ansible automation
│   ├── sp_backend/            # SP Backend integration
│   ├── bucket_service/        # AWS S3 storage
│   └── ssl_service/           # SSL certificate management
├── pkg/                       # Shared utilities and helpers
│   ├── paginate_utils/        # Pagination utilities
│   ├── common/                # Common helpers
│   ├── database/              # Database connection
│   ├── apm_helper/            # Elastic APM
│   ├── http_utils/            # HTTP client utilities
│   ├── auth_utils/            # Authentication helpers
│   └── validation_utils/      # Validation helpers
├── routes/                    # Route configuration
│   ├── init.go               # Route initialization
│   └── docs.go               # Swagger documentation setup
├── database/                  # Database setup and migrations
│   ├── init.go               # Database initialization
│   └── migrations/           # Database migration files
├── docs/                      # Documentation
│   ├── project_architecture.md # Architecture documentation
│   ├── app_package.md          # App package detailed documentation
│   └── swagger/              # API documentation
└── Deployment configuration files
```

### App Package Structure

For detailed implementation of the Clean Architecture layers (Controller, Use Case, Repository), please refer to:
**[app_package.md](app_package.md)**

This document focuses on the overall system architecture, while app_package.md provides:
- Detailed layer implementations
- Code examples and patterns
- Model transformations
- Error handling patterns
- Pagination and filtering implementation

## Layer Details

### Clean Architecture Layers Overview

#### 1. Controller Layer (Presentation)
- **Purpose**: Handle HTTP requests, responses, and routing
- **Responsibilities**:
  - Parse HTTP requests and validate parameters
  - Transform HTTP requests to use case requests
  - Transform use case responses to HTTP responses
  - Handle HTTP-specific errors and status codes
- **Implementation Details**: See [docs/app_package.md](./app_package.md#controller-layer)

#### 2. Use Case Layer (Business Logic)
- **Purpose**: Implement application business rules and orchestrate workflows
- **Responsibilities**:
  - Apply business rules and validation
  - Orchestrate repository operations
  - Coordinate external services
  - Manage business transactions
  - Perform authorization checks
- **Implementation Details**: See [docs/app_package.md](./app_package.md#use-case-layer)

#### 3. Repository Layer (Data Access)
- **Purpose**: Manage data persistence and retrieval operations
- **Responsibilities**:
  - Database CRUD operations
  - Data mapping between layers
  - Query optimization and filtering
  - Transaction management
- **Implementation Details**: See [docs/app_package.md](./app_package.md#repository-layer)

#### 4. Service Layer (External Integrations)
- **Purpose**: Abstract external service integrations
- **Components**:
  - **Varnish CDN**: Content delivery network management
  - **Ansible**: Infrastructure automation
  - **SP Backend**: Backend service integration
  - **AWS S3**: Cloud storage operations
  - **SSL Service**: Certificate management

## Data Flow

### Request Flow Architecture

```
HTTP Request → Controller Models → Use Case Models → Repository Models → Database
                    ↓                   ↓                 ↓
                HTTP Validation   Business Logic   Data Mapping
                    ↓                   ↓                 ↓
                Transformation   Orchestration  Persistence
                    ↓                   ↓                 ↓
                Use Case Call    Repository Call  Transaction
```

### Response Flow Architecture

```
Database → Repository Models → Use Case Models → Controller Models → HTTP Response
    ↓              ↓                   ↓                 ↓
Data Mapping   Business Logic     Transformation     JSON Serialization
    ↓              ↓                   ↓                 ↓
Repository     Use Case Return     Controller Return    HTTP Status Code
Response       Response            Response            Response
```

## Why Models in Each Layer?

### 1. Loose Coupling
- Each layer has its own models, preventing direct dependency between layers
- Changes in one layer's models don't affect other layers
- Easier to evolve each layer independently

### 2. Layer-Specific Concerns
- **Controller Models**: HTTP-specific concerns (validation tags, JSON names)
- **Use Case Models**: Business logic concerns (business rules, domain concepts)
- **Repository Models**: Database-specific concerns (GORM tags, relationships)
- **Service Models**: External API-specific concerns (API request/response formats)

### 3. Testability
- Each layer can be tested with its own test data
- Mocks can use layer-specific interfaces and models
- Easier to create test data for each layer

## External Integrations

### 1. Varnish CDN Service

**Purpose**: Manage Varnish caching server operations

**Key Features**:
- CDN domain lifecycle management
- VCL (Varnish Configuration Language) management
- Cache purging and invalidation
- Statistics and analytics
- Session management with Varnish

**Architecture**:
```go
type VarnishService interface {
    CreateDomain(ctx context.Context, session *Session, domain string) (int, error)
    DeleteDomain(ctx context.Context, session *Session, domain string) (int, error)
    DeployVCLGroup(ctx context.Context, session *Session, domain string) (int, error)
    PurgeCache(ctx context.Context, url string, password string) (int, error)
    GetStats(ctx context.Context, session *Session, filter StatFilter) ([]StatResponse, int, error)
}
```

### 2. Ansible Service

**Purpose**: Automate infrastructure deployment and configuration

**Key Features**:
- Automated CDN deployment
- Configuration management
- Rollback capabilities
- Integration with Varnish
- Deployment tracking

**Architecture**:
```go
type AnsibleService interface {
    AnsibleCreateCDN(req CreateCDNRequest, cdnId string, sslData *SSLData) (int, error)
    AnsibleDeleteCDN(domainId, vclId, deploymentId, vclgroupId int64, cdnData Cdn, cdnId string) (int, error)
    AnsibleDeployCDN(vclGroupId int64, cdn Cdn, cdnId string) (int, error)
    AnsibleUpdateDeployment(params UpdateDeploymentParams) (int, error)
}
```

### 3. AWS S3 Service

**Purpose**: Cloud storage operations for CDN assets

**Key Features**:
- File upload and download
- Bucket management
- Access control
- Metadata management
- Multi-part upload support

**Architecture**:
```go
type BucketService interface {
    GetFile(ctx context.Context, location string) (File, int, error)
    UploadFile(ctx context.Context, location string, data []byte) (int, error)
    DeleteFile(ctx context.Context, location string) (int, error)
    ListFiles(ctx context.Context, prefix string) ([]FileInfo, int, error)
}
```

### 4. SSL Service

**Purpose**: SSL certificate management

**Key Features**:
- Certificate generation
- Certificate renewal
- Certificate validation
- Integration with CDN
- Automated provisioning

**Architecture**:
```go
type SSLRepository interface {
    GetSSL(ctx context.Context, filter SSLFilter) (Ssl, int, error)
    CreateSSL(ctx context.Context, ssl Ssl) (int, error)
    UpdateSSL(ctx context.Context, filter SSLFilter, ssl Ssl) (int, error)
    DeleteSSL(ctx context.Context, filter SSLFilter) (int, error)
}
```

### 5. SP Backend Service

**Purpose**: Integration with SP backend services

**Key Features**:
- Project suspend availability check
- Audit log management
- WebSocket progress updates
- Backend API communication

**Architecture**:
```go
type SPBackend interface {
    CheckProjectSuspendAvailability(ctx context.Context, projectID uuid.UUID) (ResponseProjectSuspendAvailability, int, error)
    ErrorToAuditLog(payload AuditLogPayload)
    SendWebsocketProgress(payload Progress)
}
```

## Security Architecture

### Authentication & Authorization

#### JWT-Based Authentication
```go
type JWTMiddleware struct {
    secretKey       []byte
    tokenService    TokenService
    authRepo         AuthRepository
}

func (j *JWTMiddleware) ValidateToken(tokenString string) (*TokenClaims, error) {
    // Parse and validate JWT token
    // Check token expiration
    // Validate user permissions
    // Return user context
}
```

#### Role-Based Access Control (RBAC)
- **User**: Basic CDN user access
- **Admin**: Full CDN management access
- **Super Admin**: System-wide access

#### Security Features
- Password hashing with bcrypt
- Session management with refresh tokens
- Rate limiting on sensitive endpoints
- CORS configuration
- SQL injection prevention
- XSS protection

### Data Security

#### Encryption
- Sensitive data encryption at rest
- HTTPS enforcement
- API key management
- Secure configuration storage

#### Audit Trail
- Comprehensive logging with Elastic APM
- Audit log for all critical operations
- Error tracking and monitoring
- Security event logging

## Database Architecture

### Multi-Database Setup

#### Primary Databases (PostgreSQL)
1. **main**: Primary application database
2. **core**: Core application entities
3. **ssl**: SSL certificate management
4. **billing**: Billing and usage tracking

#### Database Schema Patterns
```go
type CDN struct {
    ID             uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    ProjectId     uuid.UUID  `gorm:"not null"`
    Domain         string     `gorm:"uniqueIndex;not null"`
    Status         CDNStatus  `gorm:"not null;default:'pending'"`
    OriginUrl      string     `gorm:"not null"`
    CreatedAt      time.Time
    UpdatedAt      time.Time
    DeletedAt      gorm.DeletedAt `gorm:"index"`

    // Relationships
    SSL            *SSL        `gorm:"foreignKey:CdnId"`
    Statistics     []CDNStat  `gorm:"foreignKey:CdnId"`
    Deployments    []Deployment `gorm:"foreignKey:CdnId"`
}
```

**🆔 UUID-First Rule**: All ID fields use UUID by default. Never use int/auto-increment unless explicitly justified. UUID provides better security, scalability, and distributed system support.

### Connection Management
- Connection pooling with pgx
- Health checks
- Automatic reconnection
- Transaction management
- Migration support

### Data Access Patterns
- Repository pattern with interfaces
- GORM ORM for database operations
- Optimized queries with proper indexing
- Pagination support for large datasets
- Soft deletes for data retention

## API Architecture

### RESTful API Design

#### API Versioning
```
/api/v1/cdn          # CDN management
/api/v1/settings       # Configuration
/api/v1/system         # System administration
/api/v1/analytics      # Usage analytics
```

#### Standard HTTP Methods
- **GET**: Resource retrieval (with pagination)
- **POST**: Resource creation
- **PUT**: Resource updates (full)
- **PATCH**: Partial resource updates
- **DELETE**: Resource deletion

#### Response Format
```json
{
  "success": true,
  "data": {
    // Response data
  },
  "message": "Operation successful",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### API Documentation
- OpenAPI/Swagger specification
- Interactive API documentation
- Code examples for multiple languages
- Postman collection for testing

### Rate Limiting
- Per-endpoint rate limits
- User-based quotas
- Anonymous user limits
- Burst handling

## Deployment Architecture

### Containerization

#### Dockerfile Structure
```dockerfile
# Multi-stage build
FROM golang:1.21-alpine AS builder
# Build stage...

FROM alpine:latest
# Runtime stage...
```

#### Container Configuration
- Minimal runtime image
- Security scanning
- Health checks
- Graceful shutdown
- Environment-specific builds

### Kubernetes Deployment

#### Pod Configuration
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ai-cdn-services
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ai-cdn-services
  template:
    spec:
      containers:
      - name: ai-cdn-services
        image: ai-cdn-services:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
```

#### Service Configuration
- Load balancing
- Service discovery
- Config maps
- Secrets management
- Horizontal pod autoscaling

### CI/CD Pipeline

#### GitHub Actions
- Automated testing
- Code quality checks
- Security scanning
- Container image building
- Kubernetes deployment
- Rollback strategies

## Monitoring & Observability

### Elastic APM Integration

#### Application Performance Monitoring
- Request tracing
- Error tracking
- Performance metrics
- Custom instrumentation
- Distributed tracing

#### Centralized Logging
- Structured logging
- Log aggregation
- Search and filter capabilities
- Alert configuration
- Log retention policies

### Health Checks

#### Application Health
- Database connectivity
- External service availability
- System resource monitoring
- Custom health indicators
- Health endpoint

#### Infrastructure Monitoring
- CPU and memory usage
- Network latency
- Disk I/O performance
- Request throughput
- Error rates

## Design Patterns Used

### 1. Repository Pattern
Implementation: All data access through repository interfaces
Purpose: Decouple application logic from data storage

### 2. Service Layer Pattern
Implementation: External service integrations as interfaces
Purpose: Decouple business logic from external dependencies

### 3. Dependency Injection
Implementation: Constructor injection in initialization
Purpose: Enable testing, loose coupling, easier configuration

### 4. Factory Pattern
Implementation: Service and controller factories
Purpose: Centralized object creation, consistent initialization

### 5. Strategy Pattern
Implementation: Different authentication strategies
Purpose: Encapsulate algorithms, make them interchangeable

### 6. Observer Pattern
Implementation: Event logging and notifications
Purpose: Decouple event generation from event handling

## Architecture Decisions

### Why Models in Each Layer?

1. **Maintainability**: Changes in one layer don't cascade to others
2. **Testability**: Each layer can be tested independently
3. **Flexibility**: Easy to swap implementations
4. **Clarity**: Each layer's models express its specific concerns

### Why Clean Architecture?

1. **Testability**: Business logic isolated from external dependencies
2. **Flexibility**: Easy to swap implementations (databases, services)
3. **Maintainability**: Clear separation of concerns
4. **Scalability**: Components can be scaled independently

### Technology Choices

1. **Go**: Performance, concurrency, ecosystem
2. **Gin**: Lightweight HTTP framework
3. **GORM**: Feature-rich ORM with good PostgreSQL support
4. **PostgreSQL**: Robust, scalable relational database
5. **Docker**: Containerization and deployment
6. **Kubernetes**: Orchestration and scaling

## Usage Guide for LLM

This architecture template provides a comprehensive foundation for building enterprise-grade Go applications. When using this template:

1. **Layer Models**: Always define models, interfaces, and implementations in separate files
2. **Transformations**: Always transform between layer models at boundaries
3. **Interfaces**: Define interfaces for all major components
4. **Dependencies**: Use dependency injection for loose coupling
5. **Business Logic**: Keep business logic in use cases, not in controllers or repositories
6. **Pagination**: Use the paginate_utils pattern for all list endpoints
7. **Error Handling**: All layers below controller must return `(result, int, error)` where the int is HTTP status code
8. **🚨 NO COMMENTS**: NEVER add ANY comments to the code - EXCEPT Swagger annotations ONLY for API documentation

The architecture is designed to be:
- **Scalable**: Can handle growth in users, data, and complexity
- **Maintainable**: Easy to understand, modify, and extend
- **Testable**: Components can be tested in isolation
- **Deployable**: Production-ready with modern DevOps practices

*This documentation provides a comprehensive overview of the Go Clean Architecture template with complete interface-driven design, layer-specific models, and robust pagination/filtering patterns.*