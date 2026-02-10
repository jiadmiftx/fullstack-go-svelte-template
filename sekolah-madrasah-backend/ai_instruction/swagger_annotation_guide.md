# 📚 Swaggo Annotation Guide

This guide provides comprehensive instructions for AI agents on how to properly add Swagger/OpenAPI annotations using swaggo library in Go Clean Architecture projects.

## 🎯 Learning Objectives

After completing this guide, AI agents will be able to:
- Add proper Swagger annotations to REST controller **handler functions only**
- Structure API documentation following OpenAPI 3.0 standards
- Define request/response schemas with proper Go struct references
- Handle authentication and authorization in API documentation
- Generate and maintain Swagger documentation automatically

## 📋 Prerequisites

Before using this guide, ensure the following packages are available:
- `github.com/swaggo/swag` - For generating Swagger documentation
- `github.com/swaggo/files` - For serving Swagger UI
- `github.com/swaggo/gin-swagger` - For Gin integration
- Go 1.19+ installed
- Project follows Clean Architecture pattern

## 🚀 Quick Reference

### Basic Annotation Structure (Handler Functions Only)

```go
// FunctionName godoc
// @Summary Short description (max 100 characters)
// @Description Detailed description of what this endpoint does
// @Tags Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Description"
// @Param request body RequestStruct true "Request payload"
// @Success 200 {object} ResponseStruct
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/endpoint [method]
func (c *Controller) FunctionName(ctx *gin.Context) {
    // Handler implementation
}
```

## 📖 Detailed Annotation Guide

### 1. **IMPORTANT: Only Add Swagger Comments to Handler Functions**

**✅ CORRECT: Add Swagger annotations to HTTP handler functions only**
```go
// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user account with the provided details
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "User creation data"
// @Success 201 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
    // HTTP handler implementation
}
```

**❌ WRONG: Do NOT add Swagger annotations to private helper functions**
```go
// validateUserRequest godoc
// @Summary Validate user request
// ❌ DO NOT ADD SWAGGER ANNOTATIONS TO PRIVATE/HELPER FUNCTIONS
func (c *UserController) validateUserRequest(req *CreateUserRequest) error {
    // Helper function implementation
}
```

**❌ WRONG: Do NOT add Swagger annotations to use case functions**
```go
// CreateUseCase godoc
// @Summary Create user use case
// ❌ DO NOT ADD SWAGGER ANNOTATIONS TO USE CASE FUNCTIONS
func (u *UserUseCase) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
    // Use case implementation
}
```

**❌ WRONG: Do NOT add Swagger annotations to repository functions**
```go
// Save godoc
// @Summary Save user to database
// ❌ DO NOT ADD SWAGGER ANNOTATIONS TO REPOSITORY FUNCTIONS
func (r *UserRepository) Save(ctx context.Context, user *User) error {
    // Repository implementation
}
```

**RULE: Swagger annotations should ONLY be added to public HTTP handler functions in controllers that are directly mapped to routes.**

### 2. Function Naming Convention

```go
// ✅ CORRECT: Handler function name matches the comment
// GetUserProfile godoc
func (c *UserController) GetUserProfile(ctx *gin.Context) {
    // Handler implementation
}

// ❌ WRONG: Handler function name doesn't match comment
// GetUserProfile godoc
func (c *UserController) GetProfile(ctx *gin.Context) {
    // Handler implementation
}
```

### 3. Summary and Description

```go
// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user account with the provided details.
// The user will receive a confirmation email to activate their account.
// Password must be at least 8 characters long and contain both letters and numbers.
```

**Rules:**
- `@Summary`: Maximum 100 characters, concise and clear
- `@Description`: Detailed explanation, can span multiple lines
- Use present tense ("Creates", "Retrieves", "Updates")
- Include important validation rules and business logic

### 4. Tags for Grouping

```go
// @Tags User
// @Tags User,Authentication
// @Tags CDN,Management
```

**Best Practices:**
- Use PascalCase for tag names
- Group related endpoints together
- Use multiple tags when endpoint belongs to multiple categories
- Keep tag names consistent across the project

### 5. Accept and Produce Types

```go
// For JSON APIs (most common):
// @Accept json
// @Produce json

// For file uploads:
// @Accept json,multipart/form-data
// @Produce json

// For different response formats:
// @Accept json
// @Produce json,xml
```

### 6. Security Authentication

```go
// For JWT Bearer Token:
// @Security BearerAuth

// For API Key:
// @Security ApiKeyAuth

// For Custom Auth:
// @Security XAuthCron

// Multiple security schemes:
// @Security BearerAuth
// @Security ApiKeyAuth
```

### 7. Parameters Documentation

#### Query Parameters

```go
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param search query string false "Search term to filter results"
// @Param status query string false "Filter by status" Enums(active,inactive)
// @Param from_date query string false "Start date in YYYY-MM-DD format"
```

#### Path Parameters

```go
// @Param id path string true "User ID (UUID format)"
// @Param user_id path int true "Numeric user ID"
// @Param action path string true "Action to perform" Enums(create,update,delete)
```

#### Body Parameters

```go
// @Param request body CreateUserRequest true "User creation data"
// @Param config body UpdateConfigRequest false "Optional configuration update"
```

#### Header Parameters

```go
// @Param X-Request-ID header string false "Unique request identifier"
// @Param Authorization header string true "Bearer token"
```

### 8. Response Documentation

#### Success Responses

```go
// @Success 200 {object} UserResponse
// @Success 201 {object} UserResponse "User created successfully"
// @Success 204 "No Content" "Resource deleted successfully"
// @Success 200 {array} UserListResponse
// @Success 200 {object} map[string]interface{} "Generic response with pagination"
```

#### Error Responses

```go
// @Failure 400 {object} map[string]string "Bad request - Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized - Invalid or missing token"
// @Failure 403 {object} map[string]string "Forbidden - Insufficient permissions"
// @Failure 404 {object} map[string]string "Not found - Resource doesn't exist"
// @Failure 422 {object} ValidationErrorResponse "Validation failed"
// @Failure 500 {object} map[string]string "Internal server error"
```

#### Custom Response Schemas

```go
// @Success 200 {object} PaginationResponse{data=[]UserResponse}
// @Success 200 {object} APIResponse{data=UserResponse,status=string,message=string}
```

### 9. Router Documentation

```go
// @Router /api/v1/users [get]
// @Router /api/v1/users [post]
// @Router /api/v1/users/{id} [get]
// @Router /api/v1/users/{id}/profile [put]
// @Router /api/v1/cdn/{id}/deploy [post]
// @Router /api/v1/cron/cdn/{id}/delete [post]
```

**Rules:**
- Use relative paths (without domain)
- Include HTTP method in brackets
- Match exact route pattern from router configuration
- Include path parameters in curly braces

## 🏗️ Complete Examples

### Example 1: CRUD Operations

```go
// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user account with email, name, and password.
// The email must be unique and valid. Password must be at least 8 characters.
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "User creation data"
// @Success 201 {object} UserResponse "User created successfully"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 409 {object} map[string]string "Email already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
    // HTTP handler implementation only
}

// ListUsers godoc
// @Summary List all users
// @Description Retrieves a paginated list of users with optional filtering
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by name or email"
// @Param status query string false "Filter by status" Enums(active,inactive,suspended)
// @Success 200 {object} PaginationResponse{data=[]UserResponse}
// @Failure 400 {object} map[string]string "Invalid query parameters"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/users [get]
func (c *UserController) ListUsers(ctx *gin.Context) {
    // HTTP handler implementation only
}
```

### Example 2: Complex Business Operations

```go
// DeployCDN godoc
// @Summary Deploy CDN configuration
// @Description Triggers deployment of CDN configuration to edge servers.
// This operation may take several minutes to complete.
// You can check deployment status using the status endpoint.
// @Tags CDN
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "CDN ID (UUID format)"
// @Param request body DeployRequest false "Optional deployment parameters"
// @Success 202 {object} DeploymentResponse "Deployment started"
// @Success 204 "No Content" "Deployment completed immediately"
// @Failure 400 {object} map[string]string "Invalid CDN ID or deployment config"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "CDN not found"
// @Failure 409 {object} map[string]string "Deployment already in progress"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/cdn/{id}/deploy [post]
func (c *CDNController) DeployCDN(ctx *gin.Context) {
    // HTTP handler implementation only
}
```

### Example 3: File Upload and Download

```go
// UploadFile godoc
// @Summary Upload file
// @Description Upload a file to the server. Supports multiple file formats and
// automatic virus scanning. Maximum file size is 100MB.
// @Tags File
// @Accept json,multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Param description formData string false "File description"
// @Param tags formData string false "Comma-separated tags"
// @Success 201 {object} FileUploadResponse
// @Failure 400 {object} map[string]string "Invalid file or parameters"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 413 {object} map[string]string "File too large"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/files/upload [post]
func (c *FileController) UploadFile(ctx *gin.Context) {
    // HTTP handler implementation only
}
```

### Example 4: Authentication Endpoints

```go
// Login godoc
// @Summary User login
// @Description Authenticates user credentials and returns a JWT token.
// Token expires in 24 hours. Supports both email and username login.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string "Invalid credentials format"
// @Failure 401 {object} map[string]string "Invalid email or password"
// @Failure 423 {object} map[string]string "Account temporarily locked"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
    // HTTP handler implementation only
}
```

## 📝 Struct Definitions for Documentation

### Request/Response Structs

```go
// Define structs for request/response to generate proper schema
// These structs do NOT need Swagger annotations
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email" example:"user@example.com"`
    Name     string `json:"name" binding:"required,min=2,max=100" example:"John Doe"`
    Password string `json:"password" binding:"required,min=8" example:"SecurePass123!"`
    Role     string `json:"role" binding:"required" example:"user" enum:"user,admin"`
}

type UserResponse struct {
    ID          string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
    Email       string    `json:"email" example:"user@example.com"`
    Name        string    `json:"name" example:"John Doe"`
    Role        string    `json:"role" example:"user"`
    CreatedAt   time.Time `json:"created_at" example:"2023-12-14T10:30:00Z"`
    LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

type PaginationResponse struct {
    Data       interface{} `json:"data"`
    Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
    Page       int   `json:"page" example:"1"`
    Limit      int   `json:"limit" example:"10"`
    Total      int64 `json:"total" example:"150"`
    TotalPages int   `json:"total_pages" example:"15"`
}

type ValidationErrorResponse struct {
    Errors map[string]string `json:"errors" example:"name:Name is required,email:Invalid email format"`
}
```

## 🛠️ Generating Swagger Documentation

### 1. Initialize Swagger

```bash
# Initialize swagger in the project
swag init -g cmd/main.go -o docs

# Or with custom options
swag init -g cmd/main.go \
    -o docs \
    --parseDependency \
    --parseInternal \
    --parseDepth 1
```

### 2. Main.go Setup

```go
package main

import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"

    _ "github.com/Lintasarta/ai-cdn-services/docs" // Import generated docs
)

// @title AI CDN Services API
// @version 1.0
// @description API documentation for AI CDN Services
// @termsOfService https://lintasarta.com/terms

// @contact.name API Support
// @contact.url https://lintasarta.com/support
// @contact.email support@lintasarta.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host api.cdn.example.com
// @BasePath /api/v1
// @schemes https http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token for authentication

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API key for service-to-service authentication

func main() {
    r := gin.Default()

    // Setup routes here...
    setupRoutes(r)

    // Add Swagger UI route
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    r.Run(":8080")
}
```

### 3. Regenerating Documentation

```bash
# Always run this command after adding/modifying Swagger annotations
swag init -g cmd/main.go -o docs

# Or create a Makefile target
make swagger:
    swag init -g cmd/main.go -o docs
```

## ✅ Best Practices

### Do's
- ✅ Add Swagger comments **ONLY** to HTTP handler functions in controllers
- ✅ Use meaningful and consistent tag names
- ✅ Include all possible HTTP status codes
- ✅ Provide clear descriptions for parameters
- ✅ Define proper request/response structs
- ✅ Use examples in struct fields
- ✅ Group related endpoints with same tags
- ✅ Document authentication requirements
- ✅ Keep comments up-to-date with code changes

### Don'ts
- ❌ **NEVER** add Swagger comments to private/helper functions
- ❌ **NEVER** add Swagger comments to use case functions
- ❌ **NEVER** add Swagger comments to repository functions
- ❌ **NEVER** add Swagger comments to utility functions
- ❌ Don't use vague or generic descriptions
- ❌ Don't forget to document error responses
- ❌ Don't use inconsistent parameter names
- ❌ Don't mix up `@Param` and `@Success` types
- ❌ Don't forget to regenerate docs after changes
- ❌ Don't use hardcoded URLs in @Router
- ❌ Don't include sensitive information in descriptions

## 🔍 Common Issues and Solutions

### Issue: Generated schema is empty
```go
// ❌ WRONG: Struct fields not exported (lowercase)
type userResponse struct {
    id string `json:"id"`
    email string `json:"email"`
}

// ✅ CORRECT: Exported fields (uppercase)
type UserResponse struct {
    ID    string `json:"id"`
    Email string `json:"email"`
}
```

### Issue: Parameter not showing in Swagger UI
```go
// ❌ WRONG: Missing required field
// @Param id path string "User ID"

// ✅ CORRECT: Include all required fields
// @Param id path string true "User ID"
```

### Issue: Multiple security schemes not working
```go
// ❌ WRONG: Only one security line
// @Security BearerAuth

// ✅ CORRECT: Multiple security lines for multiple schemes
// @Security BearerAuth
// @Security ApiKeyAuth
```

### Issue: Array responses not documented correctly
```go
// ❌ WRONG: Doesn't specify array properly
// @Success 200 {object} UserResponse

// ✅ CORRECT: Specify array type
// @Success 200 {array} UserResponse

// For paginated arrays:
// @Success 200 {object} PaginationResponse{data=[]UserResponse}
```

## 🔧 AI Agent Instructions

**CRITICAL RULE: When working with Swagger annotations:**

1. **ONLY add Swagger annotations to HTTP handler functions** in controller files
2. **NEVER add Swagger annotations to:**
   - Private/helper functions (functions starting with lowercase)
   - Use case functions (in use_case directories)
   - Repository functions (in repository directories)
   - Utility functions (any non-HTTP handler functions)
   - Service layer functions
   - Model or struct definitions

3. **Check for existing patterns** in the codebase first
4. **Follow the project's naming conventions** for tags and routes
5. **Add complete documentation** for all new HTTP endpoints
6. **Generate new docs** using `swag init` after changes
7. **Test the Swagger UI** at `/swagger/index.html`
8. **Keep comments synchronized** with code changes

**Warning Pattern for Incorrect Swagger Usage:**
```
❌ Swagger annotations found on non-handler function
❌ NEVER add Swagger annotations to private/use case/repository functions
❌ Swagger annotations should ONLY be added to HTTP handler functions in controllers

❌ Missing Swagger documentation for public HTTP handler
Required: Add complete swaggo annotations with @Summary, @Description, @Tags, @Param, @Success, @Failure, @Router
```

## 📚 Additional Resources

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
- [Gin-Swagger Examples](https://github.com/swaggo/gin-swagger/tree/master/examples)