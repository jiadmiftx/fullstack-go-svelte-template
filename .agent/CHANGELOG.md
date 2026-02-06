# CHANGELOG - Sekolah Madrasah

Catatan perubahan sistem. Update file ini setiap ada perubahan signifikan.

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
