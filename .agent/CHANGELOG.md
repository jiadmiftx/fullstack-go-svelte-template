# CHANGELOG - Sekolah Madrasah

Catatan perubahan sistem. Update file ini setiap ada perubahan signifikan.

---

## [2026-02-07] - Student Management UI & Role Updates

### Backend
- ✅ `POST /units/:id/students/with-user` - Create user + student profile in one call
- ✅ `GET /users/me/memberships` - Get user's org & unit memberships with roles
- ✅ Updated Unit Member Roles:
  - Added: `owner` (Kepala Sekolah/Pemilik)
  - Added: `anggota` (Siswa/Member)
  - Removed: `warga`

### Frontend
- ✅ Student list page with Glassmorphism style (`/students`)
- ✅ Student detail page with tabs (`/students/[studentId]`)
- ✅ Create student modal with password generation (NIS + DDMMYYYY)
- ✅ Delete & reset password modals
- ✅ Sidebar menu "Data Siswa" for admin/staff/owner

### Role Hierarchy
| Role | Access Level |
|------|-------------|
| `is_super_admin` | Platform-wide access |
| `owner` | Unit owner, full access |
| `admin` | Unit administrator |
| `staff` | TU/Staff, manage data |
| `pengurus` | Committee, limited access |
| `parent` | Parent, view only |
| `anggota` | Student/member, view only |

---

## [2026-02-06] - User Profiles & Classes

### Database
- ✅ `teacher_profiles` - Profil guru (NIP, NUPTK, pendidikan, status kepegawaian)
- ✅ `student_profiles` - Profil siswa (NIS, NISN, data diri, orang tua)
- ✅ `classes` - Kelas/Rombel per tahun ajaran
- ✅ `class_enrollments` - Pendaftaran siswa ke kelas dengan riwayat status

### API Endpoints
- ✅ Teachers: `GET/POST /units/:id/teachers`, `GET/PUT/DELETE /units/:id/teachers/:teacherId`
- ✅ Students: `GET/POST /units/:id/students`, `GET/PUT/DELETE /units/:id/students/:studentId`
- ✅ Classes: `GET/POST /units/:id/classes`, `GET/PUT/DELETE /units/:id/classes/:classId`
- ✅ Enrollments: `GET/POST /units/:id/classes/:classId/students`, `PUT/DELETE /class-enrollments/:enrollmentId`
- ✅ Transfer: `POST /class-enrollments/:enrollmentId/transfer`

---

## [2026-02-06] - Initial Setup

### Database
- ✅ `organizations` - Yayasan/Lembaga
- ✅ `units` - Sekolah (SD/SMP/SMA/TK)
- ✅ `unit_settings` - Pengaturan akademik (tahun ajaran, semester)
- ✅ `users` - Pengguna global
- ✅ `organization_members` - Pengurus yayasan
- ✅ `unit_members` - Warga sekolah (guru, siswa, staff, ortu)
- ✅ `roles` & `permissions` - RBAC system
- ✅ `posts`, `post_comments`, `post_poll_options`, `post_poll_votes` - Pengumuman

### API Endpoints
- ✅ Auth: login, register, refresh
- ✅ Users: CRUD + /me
- ✅ Organizations: CRUD + members
- ✅ Units: CRUD + members + settings
- ✅ Roles & Permissions: CRUD
- ✅ Posts: CRUD + comments + voting

### Frontend Routes
- ✅ `/auth` - Login, Register
- ✅ `/dashboard` - Super Admin
- ✅ `/org/[orgId]/dashboard` - Org dashboard
- ✅ `/org/[orgId]/units/[unitId]/dashboard` - Unit dashboard
- ✅ `/org/[orgId]/units/[unitId]/pengumuman` - Posts

### Cleanup
- ❌ Removed: iurans, kegiatans, reports, rooms, calendar, pending-registrations, struktur, jabatan, wargas
- ❌ Removed: RT/RW references from schema comments

---

## Template Entry (Copy untuk update baru)

```markdown
## [YYYY-MM-DD] - Nama Perubahan

### Database
- [ ] Tabel baru: `nama_tabel` - Deskripsi
- [ ] Field baru: `tabel.field` - Deskripsi
- [ ] Relasi baru: `tabel_a` → `tabel_b`

### API Endpoints
- [ ] `METHOD /path` - Deskripsi

### Frontend
- [ ] Route baru: `/path` - Deskripsi
- [ ] Component baru: `NamaComponent`

### Bug Fixes
- [ ] Fix: Deskripsi bug yang diperbaiki

### Notes
- Catatan tambahan jika ada
```

---

## Cara Update

1. Buka file ini
2. Tambahkan section baru di ATAS (yang terbaru di atas)
3. Gunakan template di atas
4. Simpan

**AI Agent**: Wajib update file ini setelah menambah fitur baru!
