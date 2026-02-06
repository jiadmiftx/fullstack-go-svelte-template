# AI Agent Context: Sekolah-Madrasah

> **PENTING**: File ini adalah dokumen onboarding untuk AI Agent. Baca keseluruhan sebelum memulai pekerjaan apapun.

---

## 1. PROJECT OVERVIEW

### Deskripsi Singkat
Sistem Manajemen Sekolah/Madrasah berbasis web dengan arsitektur:
- **Backend**: Go (Golang) + Gin + GORM + PostgreSQL
- **Frontend**: SvelteKit + CSS (without TailwindCSS)

### Konsep Hierarki Utama
```
USERS (Global)
    │
    ├── owns → ORGANIZATIONS (Yayasan/Lembaga)
    │               │
    │               ├── has → UNITS (Sekolah/Satuan Pendidikan)
    │               │           │
    │               │           ├── → UNIT_SETTINGS (Pengaturan Akademik)
    │               │           ├── → UNIT_MEMBERS (Guru, Siswa, Staff, Ortu)
    │               │           └── → POSTS (Pengumuman)
    │               │
    │               ├── has → ORGANIZATION_MEMBERS (Pengurus Yayasan)
    │               └── defines → ROLES & PERMISSIONS (Hak Akses Custom)
    │
    └── joins → (via ORGANIZATION_MEMBERS atau UNIT_MEMBERS)
```

### Terminologi Penting
| Istilah di Kode | Artinya di Dunia Nyata |
|-----------------|------------------------|
| `Organization` | Yayasan / Lembaga Pendidikan |
| `Unit` | Sekolah / Satuan Pendidikan (SD, SMP, SMA, TK) |
| `UnitMember` | Warga Sekolah (Guru, Siswa, Staff, Orang Tua) |
| `OrganizationMember` | Pengurus Yayasan |
| `Post` | Pengumuman / Berita |

---

## 2. STRUKTUR DATABASE

### Tabel Aktif
| Tabel | Deskripsi | Relasi Ke |
|-------|-----------|-----------|
| `users` | Akun pengguna global | - |
| `organizations` | Yayasan | `users.id` (owner) |
| `units` | Sekolah | `organizations.id` |
| `unit_settings` | Pengaturan sekolah (1:1 dengan unit) | `units.id` |
| `organization_members` | Keanggotaan di Yayasan | `users`, `organizations`, `roles` |
| `unit_members` | Keanggotaan di Sekolah | `users`, `units` |
| `roles` | Custom role per organization | `organizations.id` |
| `permissions` | Master permission | - |
| `role_permissions` | Pivot role-permission | `roles`, `permissions` |
| `posts` | Pengumuman | `units`, `users` |
| `post_comments` | Komentar | `posts`, `users` |
| `post_poll_options` | Opsi polling | `posts` |
| `post_poll_votes` | Vote user | `post_poll_options`, `users` |

### Tabel yang TIDAK ADA (Sudah Dihapus)
- ❌ `rt`, `rt_warga` (konsep RT/RW)
- ❌ `iuran`, `kegiatan` (fitur perumahan)
- ❌ `warga_profile`, `pengurus_profile`
- ❌ `calendar_entry`, `schedule`, `event`, `room`
- ❌ `attendance`, `jabatan`

---

## 3. STRUKTUR KODE

### Backend (`sekolah-madrasah-backend/`)
```
app/
├── controller/          # HTTP Handlers (Gin)
│   ├── auth_controller/
│   ├── user_controller/
│   ├── organization_controller/
│   ├── unit_controller/
│   ├── unit_member_controller/
│   ├── unit_settings_controller/
│   ├── role_controller/
│   ├── permission_controller/
│   └── post_controller/
├── use_case/            # Business Logic
├── repository/          # Data Access (GORM)
└── service/             # Shared services

database/
└── schemas/             # GORM Models (10 files)

routes/
└── rest.go              # All API routes defined here

config/                  # App configuration
pkg/                     # Shared utilities (JWT, CORS, etc)
```

### Frontend (`sekolah-madrasah-frontend/`)
```
src/
├── routes/              # SvelteKit file-based routing
│   ├── +layout.svelte   # Root layout (sidebar, auth check)
│   ├── +page.svelte     # Landing page
│   ├── auth/            # Login, Register
│   ├── dashboard/       # Super Admin dashboard
│   └── org/[orgId]/     # Organization context
│       ├── dashboard/
│       └── units/[unitId]/   # Unit context
│           ├── dashboard/
│           ├── pengumuman/
│           └── profile/
├── core/                # Reusable UI components
│   └── components/      # Toast, Modal, Table, etc
└── lib/
    └── index.ts         # Exports from core
```

---

## 4. CODING CONVENTIONS

### Backend (Go)
1. **Naming**: PascalCase untuk exported, camelCase untuk internal
2. **Error Handling**: Return error, jangan panic
3. **UUID**: Semua primary key menggunakan UUID v4
4. **Timestamps**: `CreatedAt`, `UpdatedAt`, `DeletedAt` (soft delete)
5. **GORM Hooks**: `BeforeCreate`, `BeforeUpdate` untuk auto-set timestamps

### Frontend (Svelte)
1. **State**: Gunakan Svelte stores (`$lib`)
2. **CSS**: Vanilla CSS (TIDAK menggunakan TailwindCSS)
3. **API Calls**: Gunakan `api` dari `$lib`
4. **Components**: Simpan di `src/core/components/`

### Pola Penamaan File
| Layer | Format Nama File | Contoh |
|-------|------------------|--------|
| Controller | `{domain}_controller/` | `user_controller/` |
| Use Case | `{domain}_use_case/` | `user_use_case/` |
| Repository | `{domain}_repository/` | `user_repository/` |
| Schema | `{domain}.go` | `user.go` |

---

## 5. DEVELOPMENT WORKFLOW

### Menambah Fitur Baru (Contoh: `classes` untuk Kelas/Rombel)

**Step 1: Schema (Database)**
```go
// database/schemas/class.go
type Class struct {
    Id        uuid.UUID
    UnitId    uuid.UUID  // FK
    Name      string     // "X IPA 1"
    Level     int        // 10, 11, 12
    // ...
}
```

**Step 2: Repository**
```go
// app/repository/class_repository/repository.go
type ClassRepository interface {
    Create(class *schemas.Class) error
    FindByUnitId(unitId uuid.UUID) ([]schemas.Class, error)
    // ...
}
```

**Step 3: Use Case**
```go
// app/use_case/class_use_case/use_case.go
type ClassUseCase struct {
    repo class_repository.ClassRepository
}

func (uc *ClassUseCase) CreateClass(req dto.CreateClassRequest) error {
    // business logic
}
```

**Step 4: Controller**
```go
// app/controller/class_controller/controller.go
func (c *ClassController) CreateClass(ctx *gin.Context) {
    // parse request, call use case, return response
}
```

**Step 5: Routes**
```go
// routes/rest.go
classes := units.Group("/:id/classes")
{
    classes.GET("", container.ClassController.GetClasses)
    classes.POST("", container.ClassController.CreateClass)
}
```

**Step 6: Frontend**
```
src/routes/org/[orgId]/units/[unitId]/classes/
├── +page.svelte       # List & CRUD UI
└── [classId]/
    └── +page.svelte   # Detail page
```

---

## 6. API ENDPOINTS (Current)

### Auth
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/refresh`

### Users
- `GET/POST /api/v1/users`
- `GET/PUT/DELETE /api/v1/users/:id`
- `GET /api/v1/users/me`

### Organizations
- `GET/POST /api/v1/organizations`
- `GET/PUT/DELETE /api/v1/organizations/:id`
- `GET/POST /api/v1/organizations/:id/members`
- `PUT/DELETE /api/v1/organizations/:id/members/:userId`

### Units
- `GET/POST /api/v1/units`
- `GET/PUT/DELETE /api/v1/units/:id`
- `GET/POST /api/v1/units/:id/members`
- `GET/PUT/DELETE /api/v1/units/:id/members/:memberId`
- `GET/PUT /api/v1/units/:id/settings`

### Roles & Permissions
- `GET/POST /api/v1/roles`
- `GET/PUT/DELETE /api/v1/roles/:id`
- `GET/POST/DELETE /api/v1/permissions`

### Posts
- `GET/POST /api/v1/posts`
- `GET/PUT/DELETE /api/v1/posts/:id`
- `GET/POST /api/v1/posts/:id/comments`
- `DELETE /api/v1/posts/:id/comments/:commentId`
- `POST /api/v1/posts/:id/vote`

---

## 7. ATURAN UNTUK AI AGENT

### DO (Lakukan)
✅ Baca file ini sebelum memulai task apapun
✅ Ikuti pola kode yang sudah ada (lihat contoh di controller/use_case lain)
✅ Gunakan UUID untuk semua primary key
✅ Tambahkan `CreatedAt`, `UpdatedAt`, `DeletedAt` di setiap schema baru
✅ Update `routes/rest.go` saat menambah endpoint baru
✅ Commit dengan pesan yang jelas (format: `[Scope]: Deskripsi`)
✅ **UPDATE `.agent/CHANGELOG.md` setelah setiap perubahan signifikan** (wajib!)
✅ **UPDATE `.agent/DATABASE_SCHEMA.md` jika ada tabel/field baru** (wajib!)
✅ **UPDATE file ini (`AI_CONTEXT.md`) jika ada perubahan struktur database/API**

### DON'T (Jangan)
❌ Modifikasi `auth_controller` kecuali diminta secara eksplisit
❌ Menggunakan TailwindCSS di frontend (gunakan vanilla CSS)
❌ Auto-push ke repository tanpa konfirmasi user
❌ Membuat tabel dengan field `rt`, `warga`, `perumahan` (sudah deprecated)
❌ Menghapus field `DeletedAt` (soft delete wajib)

### Saat Tidak Yakin
1. Tanyakan ke user sebelum membuat perubahan besar
2. Lihat implementasi serupa di modul lain
3. Cek file ini untuk referensi

### Dokumentasi Terkait
- **FEATURE_REQUEST.md** - **Template untuk request fitur baru** ← User mengisi ini
- **DATABASE_SCHEMA.md** - **Dokumentasi lengkap semua tabel & relasi database**
- **CHANGELOG.md** - Catatan semua perubahan sistem (`.agent/CHANGELOG.md`)
- **TEMPLATE_GUIDE.md** - Panduan clone & rename project
- **workflows/development-flow.md** - **WAJIB BACA** - Alur pengembangan (Backend-first)
- **workflows/** - Workflow untuk task tertentu

---

## 8. ENVIRONMENT

### Backend (.env)
```
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5434
DB_USER=sekolah_user
DB_PASSWORD=sekolah_password
DB_NAME=sekolah_db
JWT_SECRET=your_jwt_secret
```

### Frontend (.env)
```
PUBLIC_API_URL=http://localhost:8080/api/v1
```

### Docker
```bash
docker-compose up -d postgres  # Start DB
docker-compose down            # Stop DB
```

---

## 9. RENCANA PENGEMBANGAN (Roadmap)

### Fase 1: Core (✅ SELESAI)
- [x] Auth, User, Role, Permission
- [x] Organization, Unit, Unit Members
- [x] Posts & Comments

### Fase 2: Akademik (🔄 BELUM)
- [ ] Classes (Kelas/Rombel)
- [ ] Teachers (Guru) dengan subject
- [ ] Students (Siswa)
- [ ] Parents (Orang Tua)
- [ ] Academic Calendar

### Fase 3: Operasional (🔄 BELUM)
- [ ] Attendance (Absensi)
- [ ] Schedule (Jadwal Pelajaran)
- [ ] Grades (Nilai)

---

## 10. QUICK REFERENCE

### Start Development
```bash
# Terminal 1: Backend
cd sekolah-madrasah-backend
go run main.go

# Terminal 2: Frontend
cd sekolah-madrasah-frontend
npm run dev
```

### Check API
```bash
curl http://localhost:8080/api/v1/ping
# Expected: {"message":"pong"}
```

### Database Migration
Migration otomatis saat backend start (`database/database.go` → `AutoMigrate`)

---

**Last Updated**: 2026-02-06
**Maintained By**: AI Agent + User
