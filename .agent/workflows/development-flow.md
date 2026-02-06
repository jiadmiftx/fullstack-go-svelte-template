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
│  • Test API endpoints                                       │
│  • Review dengan user → TUNGGU APPROVAL                     │
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
```

### Langkah 2.3: Use Case
```go
// app/use_case/{nama_fitur}_use_case/use_case.go
type NamaFiturUseCase struct {
    // business logic
}
```

### Langkah 2.4: Controller
```go
// app/controller/{nama_fitur}_controller/controller.go
// HTTP handlers
```

### Langkah 2.5: Routes
```go
// routes/rest.go
// Daftarkan routes baru
```

### Langkah 2.6: Migration
```go
// database/database.go
db.AutoMigrate(&schemas.NamaFitur{})
```

### Langkah 2.7: Test API
```bash
# Test endpoints
curl http://localhost:8080/api/v1/units/{unitId}/nama-fitur
```

### ⏸️ CHECKPOINT: Minta Approval User
```
Konfirmasi ke user:
1. Backend sudah selesai
2. API bisa ditest
3. Tabel sudah terbuat di database

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

### User Request: "Buatkan fitur Kelas/Rombel"

**Agent Response (FASE 1):**
```
Baik, saya akan membuat fitur Classes. Berikut proposal desainnya:

### Tabel: `classes`
| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| id | UUID | PK | |
| unit_id | UUID | FK | Sekolah |
| name | VARCHAR(100) | NOT NULL | "X IPA 1" |
| level | INT | NOT NULL | Tingkat (10,11,12) |
| homeroom_teacher_id | UUID | FK, NULLABLE | Wali Kelas |

### API Endpoints:
- GET/POST /api/v1/units/:id/classes
- GET/PUT/DELETE /api/v1/units/:id/classes/:classId

Apakah struktur ini sudah sesuai? Jika ya, saya akan lanjut implementasi Backend.
```

**User:** "OK lanjut"

**Agent (FASE 2):** Implementasi backend...

**Agent (setelah FASE 2 selesai):**
```
Backend sudah selesai:
- ✅ Schema `class.go` dibuat
- ✅ Repository, UseCase, Controller dibuat
- ✅ Routes didaftarkan
- ✅ Tabel `classes` sudah terbuat di database

Apakah mau saya lanjut ke Frontend?
```

**User:** "Ya, lanjut frontend"

**Agent (FASE 3):** Implementasi frontend...

---

## ⚠️ Aturan Penting

1. **WAJIB** tunggu approval user sebelum pindah fase
2. **JANGAN** langsung kerja frontend jika backend belum final
3. **JANGAN** modifikasi schema setelah masuk FASE 3 (kecuali user minta)
4. **SELALU** update `CHANGELOG.md` setelah setiap fase selesai
5. **JANGAN** push ke git tanpa konfirmasi user

---

## 🔗 Referensi
- Desain Database: `.agent/AI_CONTEXT.md` (Section 2)
- Coding Patterns: `.agent/AI_CONTEXT.md` (Section 4-5)
- Update Docs: `.agent/workflows/update-documentation.md`
