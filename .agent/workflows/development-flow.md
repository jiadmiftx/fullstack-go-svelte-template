---
description: Workflow pengembangan fitur - Backend & DB dulu, baru Frontend setelah approval
---

# Workflow: Development Flow (Backend-First)

> ⚠️ **ATURAN UTAMA**: Backend & Database harus FINAL dulu sebelum mengerjakan Frontend.

---

## 📋 Overview

Setiap fitur baru dikerjakan dalam **3 FASE** terpisah:

```
┌─────────────────────────────────────────────────────────────┐
│  FASE 1: DESIGN (Arsitektur & Database)                     │
│  ───────────────────────────────────────                    │
│  • Desain struktur tabel                                    │
│  • Definisi relasi antar tabel                              │
│  • Review dengan user → TUNGGU APPROVAL                     │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼ User Approval ✅
┌─────────────────────────────────────────────────────────────┐
│  FASE 2: BACKEND (Implementasi API)                         │
│  ──────────────────────────────────                         │
│  • Schema → Repository → Use Case → Controller → Routes     │
│  • Test API via Swagger                                     │
│  • Unit tests & Review → TUNGGU APPROVAL                    │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼ User Approval ✅
┌─────────────────────────────────────────────────────────────┐
│  FASE 3: FRONTEND (UI & CSS)                                │
│  ──────────────────────────                                 │
│  • Routes & Pages                                           │
│  • Components & Styling                                     │
│  • Integration dengan API                                   │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔶 FASE 1: DESIGN (Arsitektur & Database)

### Langkah 1.1: Analisis Kebutuhan
```
1. Pahami fitur yang diminta user
2. Identifikasi entitas/tabel yang diperlukan
3. Tentukan relasi antar tabel
4. Tentukan field apa saja yang dibutuhkan
```

### Langkah 1.2: Desain Schema
Buat proposal struktur tabel dalam format markdown:

```markdown
### Tabel: `nama_tabel`
| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| id | UUID | PK | Primary Key |
| unit_id | UUID | FK, NOT NULL | Relasi ke units |
| name | VARCHAR(255) | NOT NULL | Nama item |
| created_at | TIMESTAMP | | Auto-fill |
| updated_at | TIMESTAMP | | Auto-fill |
| deleted_at | TIMESTAMP | INDEX | Soft delete |

### Relasi:
- `nama_tabel` belongs to `units` (1:N)
- `nama_tabel` has many `detail_tabel` (1:N)
```

### Langkah 1.3: Desain API Endpoints
```markdown
### Endpoints Baru:
- GET    /api/v1/units/:id/nama-fitur      → List all
- POST   /api/v1/units/:id/nama-fitur      → Create
- GET    /api/v1/units/:id/nama-fitur/:id  → Get by ID
- PUT    /api/v1/units/:id/nama-fitur/:id  → Update
- DELETE /api/v1/units/:id/nama-fitur/:id  → Delete
```

### ⏸️ CHECKPOINT: Minta Approval User
```
Tampilkan proposal ke user:
1. Struktur tabel
2. Relasi antar tabel
3. API endpoints

TUNGGU user bilang "OK" / "Approved" / "Lanjut" sebelum ke FASE 2!
```

---

## 🔶 FASE 2: BACKEND (Implementasi API)

> ⚠️ JANGAN mulai fase ini sebelum FASE 1 di-approve!

### 📚 MANDATORY: Baca Dokumentasi Backend
Sebelum implementasi, **WAJIB** baca dokumentasi di `sekolah-madrasah-backend/ai_instruction/` sesuai urutan:

```
📋 Learning Order (WAJIB DIIKUTI - SEMUA WAJIB):
1. project_architecture.md     → Pahami Clean Architecture
2. app_package.md              → Aturan development app/ folder (99% CRITICAL)
3. main_and_routes_guide.md    → Cara connect semua layers
4. MAP_VALIDATOR_GUIDE.md      → Request validation patterns
5. swagger_annotation_guide.md → Dokumentasi API
6. apm_and_log_guide.md        → Monitoring & logging
```

### 📋 Checklist Sebelum Coding
```
□ Sudah baca project_architecture.md?
□ Sudah baca app_package.md?
□ Sudah baca swagger_annotation_guide.md?
□ Sudah baca apm_and_log_guide.md?
□ Paham Clean Architecture layers?
□ Paham Filter Pattern untuk repository?
□ Paham UUID sebagai default ID type?
□ Paham paginate_utils.Paginate scope?
□ Paham aturan Swagger annotation untuk controller?
□ Paham APM tracing dan structured logging?
```

### Langkah 2.1: Schema
```go
// database/schemas/{nama_fitur}.go
type NamaFitur struct {
    Id        uuid.UUID
    UnitId    uuid.UUID
    // ... fields sesuai desain
}
```

### Langkah 2.2: Repository
```go
// app/repository/{nama_fitur}_repository/repository.go
type NamaFiturRepository interface {
    Create, FindById, FindByUnitId, Update, Delete
}
// Gunakan Filter Pattern untuk query dengan banyak kondisi
```

### Langkah 2.3: Use Case
```go
// app/use_case/{nama_fitur}_use_case/use_case.go
type NamaFiturUseCase struct {
    // business logic
}
// Jangan import controller, hanya repository dan models
```

### Langkah 2.4: Controller
```go
// app/controller/{nama_fitur}_controller/controller.go
// HTTP handlers dengan Swagger annotations
// Gunakan map_validator untuk request validation
```

### Langkah 2.5: Routes
```go
// routes/rest.go
// Daftarkan routes baru dengan proper middleware
```

### Langkah 2.6: Migration
```go
// database/database.go
db.AutoMigrate(&schemas.NamaFitur{})
```

### Langkah 2.7: Generate Swagger Docs
```bash
# Generate swagger documentation
// turbo
make swagger-local
```

### Langkah 2.8: Restart Backend & Test API via Swagger
```bash
# Restart backend untuk load routes baru
go run main.go

# Buka Swagger UI
# http://localhost:8080/swagger/index.html
```

### Langkah 2.9: Unit Tests
```bash
# Buat test files sesuai workflow unit-testing.md
// turbo
go test ./app/use_case/{nama_fitur}_use_case/... -v
// turbo
go test ./database/schemas/... -v
```

> Lihat `.agent/workflows/unit-testing.md` untuk detail coverage requirements.

### ⏸️ CHECKPOINT: Minta Approval User
```
Konfirmasi ke user:
1. Backend sudah selesai
2. API bisa ditest via Swagger (http://localhost:8080/swagger/index.html)
3. Tabel sudah terbuat di database
4. Unit tests sudah dibuat dan passing

TUNGGU user bilang "OK lanjut ke frontend" sebelum ke FASE 3!
```

---

## 🔶 FASE 3: FRONTEND (UI & CSS)

> ⚠️ JANGAN mulai fase ini sebelum FASE 2 di-approve!

### Langkah 3.1: Routes
```
src/routes/org/[orgId]/units/[unitId]/{nama-fitur}/
├── +page.svelte       # List view
└── [itemId]/
    └── +page.svelte   # Detail view
```

### Langkah 3.2: Menu Navigation
```javascript
// +layout.svelte
{
    href: "nama-fitur",
    label: "Label Menu",
    iconType: "icon-type",
    roles: ["admin", "staff"],
}
```

### Langkah 3.3: Components & Styling
```css
/* Gunakan Vanilla CSS, BUKAN TailwindCSS */
.component-name {
    /* styles */
}
```

### Langkah 3.4: Integration
```javascript
// Panggil API backend
const response = await api.get(`/units/${unitId}/nama-fitur`);
```

---

## 📝 Contoh Eksekusi

### User Request: "Buatkan fitur Jadwal Pelajaran"

**Agent Response (FASE 1):**
```
Saya akan membuat fitur Jadwal Pelajaran. Berikut proposal desainnya:

### Tabel: `schedules`
| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| id | UUID | PK | |
| unit_id | UUID | FK | Sekolah |
| class_id | UUID | FK | Kelas |
| subject | VARCHAR(100) | NOT NULL | Mata pelajaran |
| day_of_week | INT | NOT NULL | 1-7 (Senin-Minggu) |
| start_time | TIME | NOT NULL | Jam mulai |
| end_time | TIME | NOT NULL | Jam selesai |

### API Endpoints:
- GET/POST /api/v1/units/:id/schedules
- GET/PUT/DELETE /api/v1/units/:id/schedules/:scheduleId

Apakah struktur ini sudah sesuai?
```

**User:** "OK lanjut"

**Agent (FASE 2):**
```
Sebelum implementasi, saya akan baca:
1. ✅ ai_instruction/project_architecture.md
2. ✅ ai_instruction/app_package.md
3. ✅ ai_instruction/main_and_routes_guide.md
4. ✅ ai_instruction/MAP_VALIDATOR_GUIDE.md
5. ✅ ai_instruction/swagger_annotation_guide.md
6. ✅ ai_instruction/apm_and_log_guide.md

Memulai implementasi backend...
```

**Agent (setelah FASE 2 selesai):**
```
Backend sudah selesai:
- ✅ Schema `schedule.go` dibuat
- ✅ Repository dengan Filter Pattern
- ✅ UseCase dengan business logic
- ✅ Controller dengan Swagger annotations
- ✅ Routes didaftarkan
- ✅ Swagger docs: http://localhost:8080/swagger/index.html
- ✅ Unit tests passing

Silakan test API via Swagger. Lanjut ke Frontend?
```

**User:** "Ya, lanjut frontend"

**Agent (FASE 3):** Implementasi frontend...

---

## ⚠️ Aturan Penting

1. **WAJIB** tunggu approval user sebelum pindah fase
2. **WAJIB** baca `ai_instruction/` sesuai urutan sebelum coding backend
3. **JANGAN** langsung kerja frontend jika backend belum final
4. **JANGAN** modifikasi schema setelah masuk FASE 3 (kecuali user minta)
5. **SELALU** update `CHANGELOG.md` setelah setiap fase selesai
6. **JANGAN** push ke git tanpa konfirmasi user
7. **WAJIB** buat unit tests untuk use case dan model (lihat `unit-testing.md`)
8. **WAJIB** generate Swagger docs dan test API sebelum minta approval FASE 2

---

## 🔗 Referensi

### Backend Development (SEMUA WAJIB)
- **Learning Order**: `sekolah-madrasah-backend/ai_instruction/instruction_order.md`
- **Clean Architecture**: `sekolah-madrasah-backend/ai_instruction/project_architecture.md`
- **App Package Rules**: `sekolah-madrasah-backend/ai_instruction/app_package.md`
- **Routes & DI**: `sekolah-madrasah-backend/ai_instruction/main_and_routes_guide.md`
- **Validation**: `sekolah-madrasah-backend/ai_instruction/MAP_VALIDATOR_GUIDE.md`
- **Swagger Docs**: `sekolah-madrasah-backend/ai_instruction/swagger_annotation_guide.md`
- **APM & Logging**: `sekolah-madrasah-backend/ai_instruction/apm_and_log_guide.md`

### Project Documentation
- Database Schema: `.agent/DATABASE_SCHEMA.md`
- AI Context: `.agent/AI_CONTEXT.md`
- Unit Testing: `.agent/workflows/unit-testing.md`
- Update Docs: `.agent/workflows/update-documentation.md`

### APIs
- Swagger UI: `http://localhost:8080/swagger/index.html`
- API Base URL: `http://localhost:8080/api/v1`
