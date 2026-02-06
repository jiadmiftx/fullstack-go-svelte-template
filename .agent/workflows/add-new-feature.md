---
description: Workflow untuk menambahkan fitur baru (CRUD) ke sistem
---

# Workflow: Menambah Fitur Baru

Gunakan workflow ini saat diminta menambahkan fitur/modul baru ke sistem (contoh: `classes`, `students`, `teachers`).

## Pre-Check

1. Baca `.agent/AI_CONTEXT.md` untuk memahami konteks project
2. Cek apakah fitur terkait dengan `Unit` atau `Organization`
3. Pastikan tidak ada fitur serupa yang sudah ada

---

## Step 1: Database Schema

// turbo
Buat file schema baru di `sekolah-madrasah-backend/database/schemas/{nama_fitur}.go`

**Template:**
```go
package schemas

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type NamaFitur struct {
    Id        uuid.UUID      `gorm:"type:uuid;primaryKey"`
    UnitId    uuid.UUID      `gorm:"type:uuid;not null;index"` // FK ke unit
    Name      string         `gorm:"type:varchar(255);not null"`
    // tambah field lain sesuai kebutuhan
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`

    Unit *Unit `gorm:"foreignKey:UnitId"`
}

func (NamaFitur) TableName() string { return "nama_fiturs" }

func (n *NamaFitur) BeforeCreate(tx *gorm.DB) (err error) {
    if n.Id == uuid.Nil {
        n.Id = uuid.New()
    }
    n.CreatedAt = time.Now()
    n.UpdatedAt = time.Now()
    return
}

func (n *NamaFitur) BeforeUpdate(tx *gorm.DB) (err error) {
    n.UpdatedAt = time.Now()
    return
}
```

---

## Step 2: Repository

// turbo
Buat folder repository di `sekolah-madrasah-backend/app/repository/{nama_fitur}_repository/`

**File: `repository.go`**
```go
package nama_fitur_repository

import (
    "sekolah-madrasah/database/schemas"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type NamaFiturRepository interface {
    Create(item *schemas.NamaFitur) error
    FindById(id uuid.UUID) (*schemas.NamaFitur, error)
    FindByUnitId(unitId uuid.UUID) ([]schemas.NamaFitur, error)
    Update(item *schemas.NamaFitur) error
    Delete(id uuid.UUID) error
}

type namaFiturRepository struct {
    db *gorm.DB
}

func NewNamaFiturRepository(db *gorm.DB) NamaFiturRepository {
    return &namaFiturRepository{db: db}
}

// Implement methods...
```

---

## Step 3: Use Case

// turbo
Buat folder use case di `sekolah-madrasah-backend/app/use_case/{nama_fitur}_use_case/`

**File: `use_case.go`**
```go
package nama_fitur_use_case

import (
    "sekolah-madrasah/app/repository/nama_fitur_repository"
    "sekolah-madrasah/database/schemas"
)

type NamaFiturUseCase struct {
    repo nama_fitur_repository.NamaFiturRepository
}

func NewNamaFiturUseCase(repo nama_fitur_repository.NamaFiturRepository) *NamaFiturUseCase {
    return &NamaFiturUseCase{repo: repo}
}

// Implement business logic methods...
```

---

## Step 4: Controller

// turbo
Buat folder controller di `sekolah-madrasah-backend/app/controller/{nama_fitur}_controller/`

**File: `controller.go`**
```go
package nama_fitur_controller

import (
    "net/http"
    "sekolah-madrasah/app/use_case/nama_fitur_use_case"
    "github.com/gin-gonic/gin"
)

type NamaFiturController struct {
    useCase *nama_fitur_use_case.NamaFiturUseCase
}

func NewNamaFiturController(uc *nama_fitur_use_case.NamaFiturUseCase) *NamaFiturController {
    return &NamaFiturController{useCase: uc}
}

func (c *NamaFiturController) GetAll(ctx *gin.Context) {
    // Implementation
}

func (c *NamaFiturController) Create(ctx *gin.Context) {
    // Implementation
}

// Add other CRUD handlers...
```

---

## Step 5: Routes

Edit file `sekolah-madrasah-backend/routes/rest.go`:

1. Import controller baru
2. Tambahkan ke Container struct
3. Inisialisasi di NewContainer()
4. Daftarkan routes baru

**Contoh tambahan routes:**
```go
// Di dalam units group
namaFiturs := units.Group("/:id/nama-fiturs")
namaFiturs.Use(http_middleware.JWTAuthentication)
{
    namaFiturs.GET("", container.NamaFiturController.GetAll)
    namaFiturs.POST("", container.NamaFiturController.Create)
    namaFiturs.GET("/:itemId", container.NamaFiturController.GetById)
    namaFiturs.PUT("/:itemId", container.NamaFiturController.Update)
    namaFiturs.DELETE("/:itemId", container.NamaFiturController.Delete)
}
```

---

## Step 6: Database Migration

// turbo
Tambahkan schema ke auto-migrate di `sekolah-madrasah-backend/database/database.go`:

```go
db.AutoMigrate(
    // existing...
    &schemas.NamaFitur{},
)
```

---

## Step 7: Frontend Routes

Buat folder route baru di frontend:
```
src/routes/org/[orgId]/units/[unitId]/{nama-fitur}/
├── +page.svelte       # List view
└── [itemId]/
    └── +page.svelte   # Detail view
```

---

## Step 8: Frontend Menu (Optional)

Edit `src/routes/org/[orgId]/units/[unitId]/+layout.svelte`:

Tambahkan item menu baru di array `fullMenu`:
```javascript
{
    href: "nama-fitur",
    label: "Nama Fitur",
    iconType: "icon-name",
    roles: ["admin", "staff"],
},
```

---

## Step 9: Test

// turbo
```bash
# Restart backend untuk trigger migration
cd sekolah-madrasah-backend
go run main.go

# Test endpoint
curl http://localhost:8080/api/v1/units/{unitId}/nama-fiturs
```

---

## Post-Check

1. Pastikan tidak ada error di console backend
2. Cek database untuk memastikan tabel baru terbuat
3. Test CRUD via API atau frontend
4. Update dokumentasi jika perlu
