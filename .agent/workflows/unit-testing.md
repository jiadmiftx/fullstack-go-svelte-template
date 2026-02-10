---
description: Aturan dan template untuk unit testing di backend Go
---

# Workflow: Unit Testing

> ⚠️ **WAJIB**: Setiap fitur baru harus memiliki unit test sebelum dianggap selesai.

---

## 📋 Coverage Requirements

Setiap modul baru harus memiliki test coverage untuk:

| Layer | Test Cases |
|-------|------------|
| **Handler/Controller** | health check, create, get, list, error cases |
| **Use Case** | create, get, list, validation, business limits |
| **Models/Schema** | JSON serialization, filter defaults |

---

## 📁 Struktur File Test

```
app/
├── controller/{module}_controller/
│   ├── controller.go
│   └── controller_test.go      ← Handler tests
├── use_case/{module}_use_case/
│   ├── use_case.go
│   └── use_case_test.go        ← Use case tests
└── repository/{module}_repository/
    ├── repository.go
    └── repository_test.go      ← Repository tests (optional)

database/schemas/
├── {module}.go
└── {module}_test.go            ← Model tests
```

---

## 🔧 Test Patterns

### 1. Handler/Controller Tests

```go
// controller_test.go
package {module}_controller

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock use case
type MockUseCase struct {
    mock.Mock
}

func setupRouter(ctrl *Controller) *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.Default()
    // register routes
    return r
}

func TestGetAll_Success(t *testing.T) {
    mockUC := new(MockUseCase)
    // setup mock expectations
    // make request
    // assert response
}

func TestCreate_ValidationError(t *testing.T) {
    // test invalid input
}

func TestGetById_NotFound(t *testing.T) {
    // test 404 case
}
```

### 2. Use Case Tests

```go
// use_case_test.go
package {module}_use_case

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock repository
type MockRepository struct {
    mock.Mock
}

func TestCreate_Success(t *testing.T) {
    mockRepo := new(MockRepository)
    uc := NewUseCase(mockRepo)
    
    // setup mock
    // call use case
    // assert result
}

func TestCreate_ValidationError(t *testing.T) {
    // test validation rules
}

func TestGetByUnitId_Pagination(t *testing.T) {
    // test pagination limits
}
```

### 3. Model/Schema Tests

```go
// schema_test.go
package schemas

import (
    "encoding/json"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestModel_JSONSerialization(t *testing.T) {
    model := Model{...}
    
    // serialize
    data, err := json.Marshal(model)
    assert.NoError(t, err)
    
    // deserialize
    var result Model
    err = json.Unmarshal(data, &result)
    assert.NoError(t, err)
    assert.Equal(t, model.Field, result.Field)
}

func TestModel_Defaults(t *testing.T) {
    model := Model{}
    // test default values
}
```

---

## 🧪 Table-Driven Tests

Gunakan pattern table-driven untuk multiple test cases:

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateRequest
        wantErr bool
        errMsg  string
    }{
        {
            name:    "valid input",
            input:   CreateRequest{Name: "Test", Level: 10},
            wantErr: false,
        },
        {
            name:    "empty name",
            input:   CreateRequest{Name: "", Level: 10},
            wantErr: true,
            errMsg:  "name is required",
        },
        {
            name:    "invalid level",
            input:   CreateRequest{Name: "Test", Level: 0},
            wantErr: true,
            errMsg:  "level must be between 1 and 12",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validate(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

---

## 🏃 Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./app/use_case/class_use_case/...

# Run with verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

---

## ✅ Checklist Sebelum Merge

- [ ] Semua tests pass (`go test ./...`)
- [ ] Coverage minimal 70% untuk use case layer
- [ ] Test cases mencakup:
  - [ ] Happy path (success cases)
  - [ ] Validation errors
  - [ ] Not found errors
  - [ ] Edge cases (empty list, pagination limits)
- [ ] No skipped tests tanpa alasan

---

## 📝 Template Test File

Gunakan command ini untuk generate skeleton test:

```bash
# Di direktori module
touch controller_test.go
touch use_case_test.go
```

---

## 🔗 Dependencies

Pastikan dependencies testing sudah ada di `go.mod`:

```go
require (
    github.com/stretchr/testify v1.8.0
)
```

Install jika belum:
```bash
go get github.com/stretchr/testify
```
