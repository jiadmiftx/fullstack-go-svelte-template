# 📋 API Response Standards

This guide documents the standard response formats for all API endpoints.

## 🎯 Purpose

Ensure consistent API response structure across all endpoints for easier frontend integration.

---

## 🏷️ JSON Tag Rules (MANDATORY)

**ALL Go struct fields that are serialized to JSON MUST have explicit `json` tags with snake_case naming.**

### Why?
- Go uses PascalCase (e.g., `FullName`)
- Frontend expects snake_case (e.g., `full_name`)
- Without json tags, Go serializes to PascalCase, breaking frontend

### Rules

```go
// ✅ CORRECT - Every field has json tag with snake_case
type User struct {
    Id       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    FullName string    `gorm:"type:varchar(100)" json:"full_name"`
    Email    string    `json:"email"`
    Password string    `json:"-"`  // Hidden from JSON response
}

// ❌ WRONG - Missing json tags (will serialize as PascalCase)
type User struct {
    Id       uuid.UUID `gorm:"type:uuid;primaryKey"`
    FullName string    `gorm:"type:varchar(100)"`  // Becomes "FullName" in JSON!
}
```

### Special Tags

| Tag | Purpose |
|-----|---------|
| `json:"field_name"` | Standard snake_case field |
| `json:"-"` | Exclude from JSON (e.g., Password, DeletedAt) |
| `json:"field,omitempty"` | Omit if empty/nil (use for relations) |

### Naming Conventions

| Go Field | JSON Tag |
|----------|----------|
| `FullName` | `json:"full_name"` |
| `CreatedAt` | `json:"created_at"` |
| `OrganizationId` | `json:"organization_id"` |
| `IsActive` | `json:"is_active"` |
| `NIS` | `json:"nis"` (acronyms lowercase) |
| `NUPTK` | `json:"nuptk"` |

---

### 2. List/Paginated Response

**⚠️ MANDATORY for all list endpoints** - Use `DataWithPaginateResponse` for ALL endpoints returning arrays/lists.

```go
import (
    "sekolah-madrasah/pkg/gin_utils"
    "sekolah-madrasah/pkg/paginate_utils"
)

// Calculate total pages
totalPages := int(total) / limit
if int(total) % limit > 0 {
    totalPages++
}

ctx.JSON(http.StatusOK, gin_utils.DataWithPaginateResponse{
    DataResponse: gin_utils.DataResponse{
        Message: "Items retrieved successfully",
        Data:    items, // Direct array, NOT nested
    },
    Paginate: &paginate_utils.PaginateData{
        Page:       page,
        Limit:      limit,
        TotalData:  total,
        TotalPages: totalPages,
    },
})
```

**JSON Response:**
```json
{
  "message": "Students retrieved successfully",
  "data": [
    { "id": "uuid1", "name": "Student 1" },
    { "id": "uuid2", "name": "Student 2" }
  ],
  "paginate": {
    "page": 1,
    "limit": 10,
    "total_data": 100,
    "total_pages": 10
  }
}
```

### 3. Message Only Response

Use for DELETE or actions without data return.

```go
gin_utils.MessageResponse{
    Message: "Item deleted successfully",
}
```

**JSON Response:**
```json
{
  "message": "Item deleted successfully"
}
```

## ❌ Anti-Patterns (DON'T DO THIS!)

### Nested Data in List Response

```go
// ❌ WRONG - Creates nested structure
ctx.JSON(http.StatusOK, gin_utils.DataResponse{
    Message: "Success",
    Data: gin.H{
        "data":  items,    // NO! Don't nest data inside data
        "total": total,
        "page":  page,
    },
})
```

This creates confusing structure:
```json
{
  "message": "Success",
  "data": {
    "data": [...],  // ❌ Nested data!
    "total": 100
  }
}
```

### Inconsistent Field Names

```go
// ❌ WRONG
Paginate: &paginate_utils.PaginateData{
    CurrentPage: 1,  // Should be "page"
    PerPage: 10,     // Should be "limit"
}
```

## ✅ Required Fields for Pagination

| Field | Type | Description |
|-------|------|-------------|
| `page` | int | Current page number (1-indexed) |
| `limit` | int | Items per page |
| `total_data` | int64 | Total count of all items |
| `total_pages` | int | Total number of pages |

## 🔧 Frontend Handling

Standard way to handle paginated responses in frontend:

```typescript
const response = await api.get('/items?page=1&limit=10');

// Access data directly
const items = response.data;

// Access pagination info
if (response.paginate) {
    const { page, limit, total_data, total_pages } = response.paginate;
}
```

## 📋 Checklist Before Creating List Endpoint

- [ ] Import `paginate_utils` package
- [ ] Use `DataWithPaginateResponse` (not `DataResponse`)
- [ ] Put array directly in `Data` field (no nesting)
- [ ] Include `Paginate` with all 4 fields
- [ ] Calculate `TotalPages` correctly
- [ ] Test with Swagger UI

## 🔄 Migration Notes

If you find old endpoints with nested data structure, update them to use `DataWithPaginateResponse` following this guide.
