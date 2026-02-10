# Elastic APM and Elastic Log Guide

> **🤖 AI Agent Reference**: Complete guide for implementing Elastic APM tracking and structured logging in Go Clean Architecture applications.
>
> **⚠️ IMPORTANT**: Read `ai_instruction/instruction_order.md` for learning sequence. This guide should only be implemented AFTER `app_package.md` is 99% complete and core features are working.

## 🎯 Purpose

This guide provides:
- Elastic APM implementation patterns for function tracking
- Elastic Search logging with correlation
- General examples that work across ALL layers (Controller, Use Case, Repository)
- Best practices for observability and debugging

## 📋 Table of Contents

1. [Prerequisites](#prerequisites)
2. [APM Function Tracking](#apm-function-tracking)
3. [Elastic Logging](#elastic-logging)
4. [Correlation ID Management](#correlation-id-management)
5. [Layer-Specific Examples](#layer-specific-examples)
6. [APM Spam Prevention](#apm-spam-prevention)
7. [Best Practices](#best-practices)

## 🔥 Prerequisites

**🛑 MANDATORY Dependencies Required Before APM Implementation:**

### 1. Configuration Package (CRITICAL)
```go
// config/app.go - MUST exist and be initialized
type AppConfig struct {
    APM APMConfig  // APM configuration
    // ... other config fields
}

// config/init.go - MUST exist for config initialization
func InitializeConfig() error

// APM Configuration Structure
type APMConfig struct {
    Enabled         bool   `json:"enabled"`         // From ELASTIC_APM_ENABLED
    ServiceName     string `json:"service_name"`     // From ELASTIC_APM_SERVICE_NAME
    ServerURL       string `json:"server_url"`       // From ELASTIC_APM_SERVER_URL
    SecretToken     string `json:"secret_token"`     // From ELASTIC_APM_SECRET_TOKEN
    VerifyServerCert bool   `json:"verify_server_cert"` // From ELASTIC_APM_VERIFY_SERVER_CERT
}
```

### 2. Elasticsearch Service
```go
// services/elasticsearch_service/interfaces.go
type ElasticsearchService interface {
    Connect(ctx context.Context) error
    IndexDocument(ctx context.Context, indexName, docID string, document interface{}) error
    // ... other methods
}
```

### 3. APM Helper Package

#### Required Files and Implementation

**File: pkg/apm_helper/function_tracker.go**
```go
package apm_helper

import (
    "context"
    "go.elastic.co/apm/v2"
)

type FunctionTracker struct {
    Span *apm.Span
}

// TrackFunction starts tracking a function for APM monitoring
func TrackFunction(ctx context.Context, functionName string) *FunctionTracker {
    span, _ := apm.StartSpan(ctx, functionName)
    return &FunctionTracker{Span: span}
}

// End completes the function tracking successfully
func (ft *FunctionTracker) End() {
    if ft.Span != nil {
        ft.Span.End()
    }
}

// EndWithError completes the function tracking with an error
func (ft *FunctionTracker) EndWithError(err error) {
    if ft.Span != nil {
        ft.Span.SetOutcome("failure")
        if err != nil {
            ft.Span.Context.SetCustom("error", err.Error())
        }
        ft.Span.End()
    }
}

// SetLabel adds a custom label to the current span
func (ft *FunctionTracker) SetLabel(key string, value interface{}) {
    if ft.Span != nil {
        ft.Span.SetLabel(key, value)
    }
}
```

**File: pkg/apm_helper/logger.go**
```go
package apm_helper

import (
    "context"
    "go.elastic.co/apm/v2"
)

// LogToTransaction logs a message to the current APM transaction
func LogToTransaction(ctx context.Context, level, message string, fields map[string]interface{})

// LogInfo logs an info level message
func LogInfo(ctx context.Context, message string, fields ...map[string]interface{})

// LogWarning logs a warning level message
func LogWarning(ctx context.Context, message string, fields ...map[string]interface{})

// LogError logs an error level message
func LogError(ctx context.Context, message string, err error, fields ...map[string]interface{})

// LogDebug logs a debug level message
func LogDebug(ctx context.Context, message string, fields ...map[string]interface{})
```

**⚠️ IMPORTANT Function Signatures:**
- `LogError` requires `err` parameter BEFORE fields: `LogError(ctx, message, err, fields...)`
- Warning function is `LogWarning` (NOT `LogWarn`)
- Other logs use variadic fields: `LogInfo(ctx, message, fields...)`
- Correlation ID is managed by ElasticAPM middleware automatically

### 3. HTTP Middleware
```go
// pkg/http_middleware/http.go
func ElasticAPM(c *gin.Context)
func SetCorrelationID(c *gin.Context)
```

### 4. Required Go Modules
```go
// go.mod dependencies
require (
    go.elastic.co/apm/v2 v2.7.2
    github.com/elastic/go-elasticsearch/v8 v8.11.1
    github.com/gin-gonic/gin v1.11.0
    github.com/labstack/gommon/log v0.4.2
)
```

**⚠️ AI Agent Rule**: DO NOT proceed with APM implementation unless ALL prerequisites exist!

## 🔧 Configuration Management (MANDATORY)

### Config Package Structure
- `config/` package provides centralized configuration
- All config access through `config.APP` global variable
- APM configuration in `config.APP.APM` with fields:
  - `Enabled` - APM tracking enabled/disabled
  - `ServiceName` - Service name for APM
  - `ServerURL` - APM server URL
  - `SecretToken` - APM secret token
  - `VerifyServerCert` - SSL verification

### Additional Config Fields for Environment Detection
- `config.APP.DebugMode` - Debug mode enabled (boolean)
- `config.APP.TestMode` - Test mode enabled (boolean)
- `config.APP.ServiceName` - Service identifier
- `config.APP.Environment` - Environment name (development/production)

### Proper Config Usage:
```go
// ✅ CORRECT: Use centralized config
if config.APP.APM.Enabled {
    tracker := apm_helper.TrackFunction(ctx, "Function")
}

// ✅ CORRECT: Use debug config
if config.APP.DebugMode {
    log.Debugf("Debug information")
}

// ✅ CORRECT: Check environment
if config.APP.DebugMode || config.APP.TestMode {
    // Development/test specific logic
}

// ❌ WRONG: Direct environment access
if os.Getenv("APM_ENABLED") == "true" {
    // Bad pattern - use config.APP.APM.Enabled instead
}

if os.Getenv("ENV") == "development" {
    // Bad pattern - use config.APP.DebugMode instead
}
```

### Configuration Best Practices
- **Type Safety**: Config provides compile-time type checking
- **Validation**: Config validates values at startup
- **Defaults**: Config provides sensible default values
- **Centralization**: Single source of truth for all configuration
- **Testability**: Easy to mock and test with different configs

## 📍 APM Function Tracking

### Basic APM Tracking Pattern (ALL LAYERS)

```go
package any_layer

import (
    "context"
    "<PROJECT_NAME>/pkg/apm_helper"
)

func AnyFunction(ctx context.Context, param string) (result string, err error) {
    // Start APM tracking
    tracker := apm_helper.TrackFunction(ctx, "AnyFunction")
    defer tracker.End()

    // Your function logic here
    result = "processed: " + param

    return result, nil
}
```

### APM Tracking with Error Handling

```go
func AnyFunctionWithError(ctx context.Context, param string) (string, error) {
    tracker := apm_helper.TrackFunction(ctx, "AnyFunctionWithError")

    result, err := processSomething(param)
    if err != nil {
        tracker.EndWithError(err)
        return "", err
    }

    tracker.End()
    return result, nil
}
```

### APM Tracking with Custom Span Name

```go
func BusinessProcess(ctx context.Context, data interface{}) error {
    tracker := apm_helper.TrackFunction(ctx, "business.process.data")

    // Add custom context to span
    tracker.Span.SetLabel("data_type", "user_profile")
    tracker.Span.SetLabel("processing_mode", "async")

    err := doProcessing(data)
    tracker.End()

    return err
}
```

## 📊 Elastic Logging

### Basic Logging with APM Correlation

```go
func AnyFunction(ctx context.Context) error {
    // Info log with correlation
    apm_helper.LogInfo(ctx, "Starting any function processing")

    // Warning log
    apm_helper.LogWarning(ctx, "Deprecated feature used", map[string]interface{}{
        "feature": "legacy_api",
        "alternative": "new_api_v2",
    })

    // Error log with details (NOTE: LogError needs err parameter first)
    err := processSomething()
    if err != nil {
        apm_helper.LogError(ctx, "Processing failed", err, map[string]interface{}{
            "error_code": "VALIDATION_ERROR",
            "input_data": "invalid_format",
        })
        return err
    }

    return nil
}
```

### Structured Logging Patterns

```go
// Success case
apm_helper.LogInfo(ctx, "Operation completed successfully", map[string]interface{}{
    "operation": "user_creation",
    "user_id": "12345",
    "duration_ms": 150,
    "records_processed": 100,
})

// Error case with stack trace (NOTE: LogError needs err parameter first)
err := database.Query(...)
if err != nil {
    apm_helper.LogError(ctx, "Database operation failed", err, map[string]interface{}{
        "operation": "user_query",
        "query": "SELECT * FROM users WHERE id = ?",
        "retry_count": 3,
    })
}
```

### Custom Log Levels

```go
func processWithLogging(ctx context.Context) {
    // Correlation ID is automatically set by middleware
    // Access via context if needed, but logs will be automatically correlated

    // Debug level (development only)
    if config.APP.DebugMode {
        log.Debugf("Debug: Processing started")
    }

    // Info level
    log.Infof("Processing data item: %s", itemID)

    // Warning level
    log.Warnf("Slow operation detected: %dms", duration)

    // Error level
    log.Errorf("Failed to process: %v", err)
}
```

## 🔗 Correlation ID Management

### Automatic Correlation (from Middleware)

```go
func HandleRequest(c *gin.Context) {
    // Correlation ID is automatically set by middleware
    ctx := c.Request.Context()

    log.InfoContext(ctx, "Request started")

    // All downstream calls will have the same correlation ID
    result := callUseCase(ctx, someData)

    c.JSON(200, result)
}
```

### Manual Correlation ID Setting

```go
func BackgroundTask() {
    // Create new context with correlation ID
    correlationID := generateUUID()
    ctx := context.WithValue(context.Background(), "correlation_id", correlationID)

    // All logs and APM traces will use this correlation ID
    tracker := apm_helper.TrackFunction(ctx, "BackgroundTask")

    processSomething(ctx)

    tracker.End()
}
```

### Cross-Service Correlation

```go
func CallExternalAPI(ctx context.Context) (*APIResponse, error) {
    tracker := apm_helper.TrackFunction(ctx, "ExternalAPI.Call")

    req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)
    // Headers are automatically set by the logger/transport layer

    resp, err := httpClient.Do(req)
    if err != nil {
        tracker.EndWithError(err)
        return nil, err
    }

    tracker.End()
    return parseResponse(resp), nil
}
```

## 🎯 Layer-Specific Examples

### Controller Layer Example

```go
package user_controller

func (uc *UserController) CreateUser(c *gin.Context) {
    ctx := c.Request.Context()
    tracker := apm_helper.TrackFunction(ctx, "UserController.CreateUser")

    // Log request received
    apm_helper.LogInfo(ctx, "Create user request received", map[string]interface{}{
        "user_agent": c.GetHeader("User-Agent"),
        "ip_address": c.ClientIP(),
    })

    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        tracker.EndWithError(err)
        apm_helper.LogError(ctx, "Invalid request body", err, map[string]interface{}{
            "validation_errors": err.Error(),
        })
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    result, err := uc.userUseCase.CreateUser(ctx, req)
    if err != nil {
        tracker.EndWithError(err)
        apm_helper.LogError(ctx, "Failed to create user", err, map[string]interface{}{
            "email": req.Email,
            "error": err.Error(),
        })
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    tracker.End()
    apm_helper.LogInfo(ctx, "User created successfully", map[string]interface{}{
        "user_id": result.ID,
        "email": result.Email,
    })

    c.JSON(201, result)
}
```

### Use Case Layer Example

```go
package user_use_case

func (uc *UserUseCase) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
    tracker := apm_helper.TrackFunction(ctx, "UserUseCase.CreateUser")

    apm_helper.LogInfo(ctx, "Starting user creation process", map[string]interface{}{
        "email": req.Email,
        "name": req.Name,
    })

    // Validate business rules
    if err := uc.validateUserRequest(req); err != nil {
        tracker.EndWithError(err)
        return nil, err
    }

    // Check if user already exists
    existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
    if err != nil && !errors.Is(err, ErrNotFound) {
        tracker.EndWithError(err)
        return nil, err
    }
    if existingUser != nil {
        tracker.End()
        apm_helper.LogWarning(ctx, "User already exists", map[string]interface{}{
            "email": req.Email,
        })
        return nil, ErrUserAlreadyExists
    }

    // Create user
    user, err := uc.userRepo.Create(ctx, req)
    if err != nil {
        tracker.EndWithError(err)
        return nil, err
    }

    tracker.End()
    apm_helper.LogInfo(ctx, "User creation completed", map[string]interface{}{
        "user_id": user.ID,
    })

    return toUserResponse(user), nil
}
```

### Repository Layer Example

```go
package user_repository

func (r *UserRepository) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
    tracker := apm_helper.TrackFunction(ctx, "UserRepository.Create")

    log.DebugContext(ctx, "Creating user in database")

    user := &User{
        ID:        generateUUID(),
        Email:     req.Email,
        Name:      req.Name,
        CreatedAt: time.Now(),
    }

    if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
        tracker.EndWithError(err)
        log.LogError(ctx, "Failed to create user in database", err,
            log.WithField("email", req.Email),
        )
        return nil, err
    }

    tracker.End()
    log.InfoContext(ctx, "User created successfully",
        log.WithField("user_id", user.ID),
    )

    return user, nil
}
```

## 🚫 APM Spam Prevention

### Selective Function Tracking

```go
// ONLY track important functions
func ImportantBusinessProcess(ctx context.Context) error {
    tracker := apm_helper.TrackFunction(ctx, "BusinessProcess.Important")
    defer tracker.End()

    // Track expensive operations
    tracker.Span.SetLabel("cost_estimate", "high")
    tracker.Span.SetLabel("impact", "critical")

    return doSomething()
}

// DON'T track simple getters
func (u *User) GetName() string {
    return u.name  // No APM tracking needed
}

// DON'T track internal helpers
func formatErrorMessage(err error) string {
    return fmt.Sprintf("Error: %v", err)  // No APM tracking needed
}
```

### Conditional APM Tracking

```go
func ConditionalTracking(ctx context.Context, data interface{}) error {
    // Only track when APM is enabled
    if config.APP.APM.Enabled {
        tracker := apm_helper.TrackFunction(ctx, "ConditionalOperation")
        defer tracker.End()
    }

    // Always log important events
    apm_helper.LogInfo(ctx, "Conditional operation executed")

    return processData(data)
}
```

### High-Frequency Operations

```go
func HighFrequencyLoop(ctx context.Context, items []Item) error {
    // Track once for the entire batch, not per item
    tracker := apm_helper.TrackFunction(ctx, "BatchProcessor.ProcessBatch")
    defer tracker.End()

    tracker.Span.SetLabel("batch_size", len(items))
    tracker.Span.SetLabel("operation", "bulk_process")

    for i, item := range items {
        // No individual APM tracking for each item
        if i%1000 == 0 {
            apm_helper.LogInfo(ctx, "Batch processing progress", map[string]interface{}{
                "processed": i,
                "total": len(items),
            })
        }
        processItem(item)
    }

    return nil
}
```

## 💡 Best Practices

### 1. APM Naming Conventions
```go
// GOOD: Descriptive and hierarchical
tracker := apm_helper.TrackFunction(ctx, "UserService.CreateUser")
tracker := apm_helper.TrackFunction(ctx, "PaymentGateway.ProcessTransaction")
tracker := apm_helper.TrackFunction(ctx, "CacheService.InvalidateUserCache")

// BAD: Generic or unclear
tracker := apm_helper.TrackFunction(ctx, "function")
tracker := apm_helper.TrackFunction(ctx, "doStuff")
tracker := apm_helper.TrackFunction(ctx, "process")
```

### 2. Span Labels Usage
```go
func enrichSpan(tracker *apm_helper.FunctionTracker) {
    // Add relevant labels
    tracker.Span.SetLabel("service", "user-management")
    tracker.Span.SetLabel("version", "v1.2.0")
    tracker.Span.SetLabel("environment", config.APP.Environment)

    // Add custom labels for business context
    tracker.Span.SetLabel("user_tier", "premium")
    tracker.Span.SetLabel("operation_type", "write")
    tracker.Span.SetLabel("resource_type", "user_profile")
}
```

### 3. Error Handling Patterns
```go
func handleErrorsGracefully(ctx context.Context) error {
    tracker := apm_helper.TrackFunction(ctx, "ErrorHandlingExample")

    result, err := someOperation()
    if err != nil {
        // Categorize errors for better analytics
        switch {
        case errors.Is(err, ErrValidation):
            tracker.Span.SetLabel("error_type", "validation")
            tracker.Span.SetLabel("error_category", "client_error")
        case errors.Is(err, ErrDatabase):
            tracker.Span.SetLabel("error_type", "database")
            tracker.Span.SetLabel("error_category", "infrastructure_error")
        default:
            tracker.Span.SetLabel("error_type", "unknown")
            tracker.Span.SetLabel("error_category", "system_error")
        }

        // Add context to error
        apm_helper.LogError(ctx, "Operation failed", err, map[string]interface{}{
            "error_category": tracker.Span.Labels()["error_category"],
            "error_type": tracker.Span.Labels()["error_type"],
            "recoverable": true,
        })

        tracker.EndWithError(err)
        return err
    }

    tracker.End()
    return nil
}
```

### 4. Performance Considerations
```go
// Avoid tracking in hot paths
func hotPathFunction(data []byte) {
    // NO APM tracking in performance-critical functions
    process(data)
}

// Use sampling for frequent operations
func frequentOperation(ctx context.Context) {
    if shouldSample(0.1) { // 10% sampling
        tracker := apm_helper.TrackFunction(ctx, "FrequentOperation")
        defer tracker.End()
    }

    doWork()
}

func shouldSample(rate float64) bool {
    return rand.Float64() < rate
}
```

### 5. Debugging Integration
```go
func debuggableFunction(ctx context.Context, config Config) error {
    tracker := apm_helper.TrackFunction(ctx, "DebuggableFunction")
    defer tracker.End()

    // Add debug information in development
    if config.APP.DebugMode {
        tracker.Span.SetLabel("debug_config", config)
        tracker.Span.SetLabel("debug_headers", getHeaders(ctx))

        apm_helper.LogDebug(ctx, "Debug information", map[string]interface{}{
            "config": config,
        })
    }

    return processWithConfig(config)
}
```

## 🚀 Quick Setup Checklist

### ✅ Before Implementation:
- [ ] config/ package exists and is initialized
- [ ] config.APP.APM is properly configured with enabled flag
- [ ] ElasticsearchService is initialized and connected
- [ ] APM helper package exists with all required files
- [ ] ElasticAPM middleware is configured in routes
- [ ] Structured logging is configured
- [ ] Go modules are properly configured

### ✅ Implementation:
- [ ] Add APM tracking to critical business functions
- [ ] Add structured logging with proper context
- [ ] Implement error categorization
- [ ] Set up appropriate span labels
- [ ] Configure sampling for high-frequency operations

### ✅ Verification:
- [ ] Check APM UI for incoming traces
- [ ] Verify correlation ID propagation
- [ ] Test error handling with proper categorization
- [ ] Monitor performance impact
- [ ] Review log aggregation in Elasticsearch

## 🔍 AI Agent Instructions

**When adding APM to existing functions:**

1. **ALWAYS** check if prerequisites exist first
2. **ONLY** track meaningful business functions
3. **NEVER** track simple getters/setters
4. **USE** descriptive function names in APM
5. **ADD** relevant context labels
6. **HANDLE** errors appropriately with categorization
7. **AVOID** APM spam in hot paths
8. **LOG** with proper structured logging consistently

**Warning Pattern for Missing APM:**
```
❌ APM prerequisites missing. Cannot add APM tracking.
Required: config/ package, pkg/apm_helper/, services/elasticsearch_service/, ElasticAPM middleware

❌ Configuration access error: Use config.APP.APM.Enabled instead of os.Getenv()
Required: Centralized configuration through config package
```