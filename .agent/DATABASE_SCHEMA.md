# Database Schema Documentation

> **Last Updated**: 2026-02-09  
> **Database**: PostgreSQL  
> **ORM**: GORM (Go)

---

## 📊 Overview Relasi

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              USERS                                          │
│  (Pengguna Global - Satu akun bisa masuk ke banyak Yayasan/Sekolah)         │
└─────────────────────────────────────────────────────────────────────────────┘
         │                                    │
         │ (owns)                             │ (joins)
         ▼                                    ▼
┌─────────────────────────┐         ┌─────────────────────────┐
│     ORGANIZATIONS       │         │   ORGANIZATION_MEMBERS  │
│     (Yayasan)           │◄────────│   (Pengurus Yayasan)    │
│                         │  1 : N  │                         │
└─────────────────────────┘         └─────────────────────────┘
         │                                    │
         │ 1 : N                              │
         ▼                                    ▼
┌─────────────────────────┐         ┌─────────────────────────┐
│        UNITS            │         │        ROLES            │
│     (Sekolah)           │         │   (Hak Akses Custom)    │
└─────────────────────────┘         └─────────────────────────┘
         │                                    │
         │                                    │ N : M
         │                                    ▼
         │                          ┌─────────────────────────┐
         │                          │   ROLE_PERMISSIONS      │
         │                          │   (Pivot Table)         │
         │                          └─────────────────────────┘
         │                                    │
         │                                    ▼
         │                          ┌─────────────────────────┐
         │                          │     PERMISSIONS         │
         │                          └─────────────────────────┘
         │
         ├──────────────────┬──────────────────┐
         │ 1 : 1            │ 1 : N            │ 1 : N
         ▼                  ▼                  ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│  UNIT_SETTINGS  │  │  UNIT_MEMBERS   │  │     POSTS       │
└─────────────────┘  └─────────────────┘  └─────────────────┘
                                                  │
                                                  │ 1 : N
                                                  ▼
                                         ┌─────────────────┐
                                         │  POST_COMMENTS  │
                                         │  POST_POLL_*    │
                                         └─────────────────┘
```

---

## 📋 Detail Tabel

### 1. USERS (Pengguna)

**Deskripsi**: Akun pengguna global. Satu user bisa bergabung ke banyak Yayasan/Sekolah.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `email` | VARCHAR(255) | UNIQUE, NOT NULL | Email login |
| `password` | VARCHAR(255) | NOT NULL | Password hash |
| `full_name` | VARCHAR(100) | | Nama lengkap |
| `phone` | VARCHAR(20) | | Nomor HP |
| `avatar` | VARCHAR(500) | | URL foto profil |
| `is_super_admin` | BOOLEAN | DEFAULT false | Flag super admin |
| `is_active` | BOOLEAN | DEFAULT true | Status aktif |
| `email_verified_at` | TIMESTAMP | NULLABLE | Waktu verifikasi email |
| `last_login_at` | TIMESTAMP | NULLABLE | Login terakhir |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Relasi**:
- Has many `organizations` (as owner)
- Has many `organization_members`
- Has many `unit_members`
- Has many `posts` (as author)

---

### 2. ORGANIZATIONS (Yayasan/Lembaga)

**Deskripsi**: Yayasan atau Lembaga Pendidikan. Induk dari Sekolah/Unit.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `owner_id` | UUID | FK → users.id, NOT NULL | Pemilik |
| `name` | VARCHAR(255) | NOT NULL | Nama yayasan |
| `code` | VARCHAR(50) | UNIQUE, NOT NULL | Kode unik (YPI-001) |
| `type` | VARCHAR(50) | NOT NULL | Tipe organisasi |
| `description` | TEXT | | Deskripsi |
| `address` | TEXT | | Alamat |
| `logo` | VARCHAR(500) | | URL logo |
| `is_active` | BOOLEAN | DEFAULT true | Status aktif |
| `settings` | JSONB | DEFAULT '{}' | Pengaturan tambahan |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Relasi**:
- Belongs to `users` (owner)
- Has many `units`
- Has many `organization_members`
- Has many `roles`

---

### 3. UNITS (Sekolah/Satuan Pendidikan)

**Deskripsi**: Sekolah atau Satuan Pendidikan (SD, SMP, SMA, TK, Madrasah).

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `organization_id` | UUID | FK → organizations.id, NOT NULL | Yayasan induk |
| `name` | VARCHAR(255) | NOT NULL | Nama sekolah |
| `code` | VARCHAR(50) | UNIQUE, NOT NULL | NPSN / Kode unik |
| `slug` | VARCHAR(100) | UNIQUE, NULLABLE | URL slug (bisa diisi belakangan) |
| `type` | VARCHAR(50) | NOT NULL | Jenjang (SD/SMP/SMA/TK) |
| `address` | TEXT | | Alamat |
| `phone` | VARCHAR(20) | | Telepon |
| `email` | VARCHAR(255) | | Email |
| `logo` | VARCHAR(500) | | URL logo |
| `is_active` | BOOLEAN | DEFAULT true | Status aktif |
| `settings` | JSONB | DEFAULT '{}' | Pengaturan tambahan |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Relasi**:
- Belongs to `organizations`
- Has one `unit_settings`
- Has many `unit_members`
- Has many `posts`

---

### 4. UNIT_SETTINGS (Pengaturan Sekolah)

**Deskripsi**: Pengaturan akademik per sekolah (tahun ajaran, semester, jadwal).

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `unit_id` | UUID | FK → units.id, UNIQUE, NOT NULL | 1:1 dengan unit |
| `period_duration` | INT | DEFAULT 40 | Durasi jam pelajaran (menit) |
| `start_time` | VARCHAR(10) | DEFAULT '07:00' | Jam mulai |
| `total_periods` | INT | DEFAULT 9 | Jumlah jam pelajaran/hari |
| `break_after_period` | INT | DEFAULT 3 | Istirahat setelah jam ke-N |
| `break_duration` | INT | DEFAULT 15 | Durasi istirahat (menit) |
| `academic_year` | VARCHAR(20) | | Tahun ajaran (2025/2026) |
| `current_semester` | INT | DEFAULT 1 | Semester aktif (1/2) |
| `semester1_start` | DATE | NULLABLE | Mulai semester 1 |
| `semester1_end` | DATE | NULLABLE | Akhir semester 1 |
| `semester2_start` | DATE | NULLABLE | Mulai semester 2 |
| `semester2_end` | DATE | NULLABLE | Akhir semester 2 |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |

**Relasi**:
- Belongs to `units` (1:1)

---

### 5. ORGANIZATION_MEMBERS (Anggota Yayasan)

**Deskripsi**: Keanggotaan user di level Yayasan (Pengurus Yayasan).

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `user_id` | UUID | FK → users.id, NOT NULL | User |
| `organization_id` | UUID | FK → organizations.id, NOT NULL | Yayasan |
| `role_id` | UUID | FK → roles.id, NOT NULL | Role/Jabatan |
| `is_active` | BOOLEAN | DEFAULT true | Status aktif |
| `joined_at` | TIMESTAMP | NOT NULL | Tanggal bergabung |
| `invited_by` | UUID | FK → users.id, NULLABLE | Diundang oleh |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Relasi**:
- Belongs to `users`
- Belongs to `organizations`
- Belongs to `roles`

---

### 6. UNIT_MEMBERS (Anggota Sekolah)

**Deskripsi**: Keanggotaan user di level Sekolah (Guru, Siswa, Staff, Orang Tua).

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `user_id` | UUID | FK → users.id, NOT NULL | User |
| `unit_id` | UUID | FK → units.id, NOT NULL | Sekolah |
| `role` | VARCHAR(20) | NOT NULL, DEFAULT 'staff' | Role enum |
| `is_active` | BOOLEAN | DEFAULT true | Status aktif |
| `joined_at` | TIMESTAMP | NOT NULL | Tanggal bergabung |
| `invited_by` | UUID | FK → users.id, NULLABLE | Diundang oleh |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Role Enum Values**:
- `admin` - Administrator sekolah
- `pengurus` - Pengurus
- `warga` - Warga sekolah umum
- `parent` - Orang tua
- `staff` - Staff

**Relasi**:
- Belongs to `users`
- Belongs to `units`

---

### 7. ROLES (Peran/Jabatan)

**Deskripsi**: Custom role per organisasi untuk RBAC.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `organization_id` | UUID | FK → organizations.id, NULLABLE | Milik org (null = global) |
| `name` | VARCHAR(50) | NOT NULL | Nama role |
| `display_name` | VARCHAR(100) | NOT NULL | Nama tampilan |
| `type` | VARCHAR(20) | NOT NULL, DEFAULT 'custom' | system/custom |
| `level` | INT | NOT NULL, DEFAULT 0 | Level hierarki |
| `description` | TEXT | | Deskripsi |
| `is_default` | BOOLEAN | DEFAULT false | Default role |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Relasi**:
- Belongs to `organizations` (optional)
- Has many `role_permissions`
- Has many `organization_members`

---

### 8. PERMISSIONS (Hak Akses)

**Deskripsi**: Master permission untuk RBAC.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `name` | VARCHAR(100) | UNIQUE, NOT NULL | Nama permission |
| `resource` | VARCHAR(50) | NOT NULL, INDEX | Resource target |
| `action` | VARCHAR(20) | NOT NULL | create/read/update/delete |
| `description` | TEXT | | Deskripsi |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |

**Relasi**:
- Has many `role_permissions`

---

### 9. ROLE_PERMISSIONS (Pivot)

**Deskripsi**: Tabel pivot untuk relasi many-to-many Role ↔ Permission.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `role_id` | UUID | FK → roles.id, NOT NULL | Role |
| `permission_id` | UUID | FK → permissions.id, NOT NULL | Permission |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |

**Primary Key**: Composite (`role_id`, `permission_id`)

---

### 10. POSTS (Pengumuman/Berita)

**Deskripsi**: Pengumuman atau berita. Bisa level Unit atau Org-wide.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `unit_id` | UUID | FK → units.id, NOT NULL | Sekolah asal |
| `author_id` | UUID | FK → users.id, NOT NULL | Penulis |
| `is_org_wide` | BOOLEAN | DEFAULT false | true = tampil di semua unit |
| `title` | VARCHAR(500) | | Judul |
| `content` | TEXT | | Isi konten |
| `post_type` | VARCHAR(20) | NOT NULL, DEFAULT 'text' | text/photo/poll/link |
| `image_url` | VARCHAR(500) | | URL gambar (type: photo) |
| `link_url` | VARCHAR(500) | | URL link (type: link) |
| `link_title` | VARCHAR(255) | | Judul link |
| `link_preview` | TEXT | | Preview link |
| `is_pinned` | BOOLEAN | DEFAULT false | Disematkan |
| `is_important` | BOOLEAN | DEFAULT false | Penting |
| `comment_count` | INT | DEFAULT 0 | Counter komentar |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Post Type Values**:
- `text` - Teks biasa
- `photo` - Dengan gambar
- `poll` - Polling
- `link` - Dengan link preview

**Relasi**:
- Belongs to `units`
- Belongs to `users` (author)
- Has many `post_comments`
- Has many `post_poll_options`

---

### 11. POST_COMMENTS (Komentar)

**Deskripsi**: Komentar pada post. Support 1 level nesting (reply).

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `post_id` | UUID | FK → posts.id, NOT NULL | Post induk |
| `parent_id` | UUID | FK → post_comments.id, NULLABLE | null = top-level |
| `author_id` | UUID | FK → users.id, NOT NULL | Penulis |
| `content` | TEXT | NOT NULL | Isi komentar |
| `reply_count` | INT | DEFAULT 0 | Counter reply |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Relasi**:
- Belongs to `posts`
- Belongs to `users` (author)
- Belongs to `post_comments` (parent, self-referencing)
- Has many `post_comments` (replies)

---

### 12. POST_POLL_OPTIONS (Opsi Polling)

**Deskripsi**: Opsi untuk post type poll.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `post_id` | UUID | FK → posts.id, NOT NULL | Post induk |
| `text` | VARCHAR(255) | NOT NULL | Teks opsi |
| `vote_count` | INT | DEFAULT 0 | Counter vote |
| `urutan` | INT | DEFAULT 0 | Urutan tampil |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |

**Relasi**:
- Belongs to `posts`
- Has many `post_poll_votes`

---

### 13. POST_POLL_VOTES (Vote Polling)

**Deskripsi**: Record vote user pada polling.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `post_id` | UUID | FK → posts.id, NOT NULL | Post |
| `option_id` | UUID | FK → post_poll_options.id, NOT NULL | Opsi yang dipilih |
| `user_id` | UUID | FK → users.id, NOT NULL | User yang vote |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |

**Relasi**:
- Belongs to `posts`
- Belongs to `post_poll_options`
- Belongs to `users`

---

### 14. TEACHER_PROFILES (Profil Guru)

**Deskripsi**: Data profil lengkap untuk guru (NIP, NUPTK, pendidikan, status kepegawaian).

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `user_id` | UUID | FK → users.id, UNIQUE, NOT NULL | 1:1 dengan user |
| `unit_id` | UUID | FK → units.id, NOT NULL | Sekolah |
| `nip` | VARCHAR(30) | NULLABLE | Nomor Induk Pegawai |
| `nuptk` | VARCHAR(30) | NULLABLE | Nomor Unik Pendidik |
| `education_level` | VARCHAR(20) | NULLABLE | S1/S2/S3/D3 |
| `education_major` | VARCHAR(100) | NULLABLE | Jurusan |
| `employment_status` | VARCHAR(20) | DEFAULT 'honorer' | PNS/Honorer/GTY/Kontrak |
| `join_date` | DATE | NULLABLE | Tanggal mulai mengajar |
| `subjects` | JSONB | DEFAULT '[]' | ⚠️ DEPRECATED - Gunakan `teacher_subjects` |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

> [!WARNING]
> Field `subjects` JSONB sudah deprecated. Gunakan relasi `teacher_subjects` untuk data mapel yang diampu.

**Relasi**:
- Belongs to `users` (1:1)
- Belongs to `units`
- Has many `teacher_subjects` (mapel yang diampu)
- Has one `classes` (wali kelas via homeroom_teacher_id)

---

### 15. STUDENT_PROFILES (Profil Siswa)

**Deskripsi**: Data profil lengkap untuk siswa (NIS, NISN, data diri, orang tua).

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `user_id` | UUID | FK → users.id, UNIQUE, NOT NULL | 1:1 dengan user |
| `unit_id` | UUID | FK → units.id, NOT NULL | Sekolah |
| `nis` | VARCHAR(30) | NULLABLE | Nomor Induk Siswa (internal) |
| `nisn` | VARCHAR(20) | NULLABLE | Nomor Induk Siswa Nasional |
| `birth_place` | VARCHAR(100) | NULLABLE | Tempat lahir |
| `birth_date` | DATE | NULLABLE | Tanggal lahir |
| `gender` | VARCHAR(10) | NULLABLE | L/P |
| `religion` | VARCHAR(20) | NULLABLE | Agama |
| `address` | TEXT | NULLABLE | Alamat lengkap |
| `father_name` | VARCHAR(100) | NULLABLE | Nama ayah |
| `mother_name` | VARCHAR(100) | NULLABLE | Nama ibu |
| `guardian_name` | VARCHAR(100) | NULLABLE | Nama wali |
| `parent_phone` | VARCHAR(20) | NULLABLE | Telepon orang tua |
| `enrollment_date` | DATE | NULLABLE | Tanggal masuk sekolah |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Relasi**:
- Belongs to `users` (1:1)
- Belongs to `units`
- Has many `class_enrollments`

---

### 16. CLASSES (Kelas/Rombel)

**Deskripsi**: Kelas atau rombongan belajar per tahun ajaran.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `unit_id` | UUID | FK → units.id, NOT NULL | Sekolah |
| `name` | VARCHAR(50) | NOT NULL | Nama kelas (X IPA 1) |
| `level` | INT | NOT NULL | Tingkat (1-12) |
| `academic_year` | VARCHAR(20) | NOT NULL | Tahun ajaran (2025/2026) |
| `homeroom_teacher_id` | UUID | FK → teacher_profiles.id, NULLABLE | Wali kelas |
| `capacity` | INT | DEFAULT 30 | Kapasitas maksimal |
| `is_active` | BOOLEAN | DEFAULT true | Status aktif |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Relasi**:
- Belongs to `units`
- Belongs to `teacher_profiles` (homeroom teacher)
- Has many `class_enrollments`

---

### 17. CLASS_ENROLLMENTS (Pendaftaran Kelas)

**Deskripsi**: Relasi siswa-kelas per tahun ajaran dengan riwayat status.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `student_profile_id` | UUID | FK → student_profiles.id, NOT NULL | Profil siswa |
| `class_id` | UUID | FK → classes.id, NOT NULL | Kelas |
| `academic_year` | VARCHAR(20) | NOT NULL | Tahun ajaran |
| `status` | VARCHAR(20) | DEFAULT 'active' | active/graduated/transferred/dropped |
| `enrolled_at` | DATE | NOT NULL | Tanggal masuk kelas |
| `left_at` | DATE | NULLABLE | Tanggal keluar |
| `notes` | TEXT | NULLABLE | Catatan |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |

**Status Values**:
- `active` - Aktif di kelas
- `graduated` - Lulus
- `transferred` - Pindah kelas
- `dropped` - Keluar

**Relasi**:
- Belongs to `student_profiles`
- Belongs to `classes`

---

### 18. SUBJECTS (Mata Pelajaran)

**Deskripsi**: Master data mata pelajaran per sekolah.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `unit_id` | UUID | FK → units.id, NOT NULL | Sekolah |
| `name` | VARCHAR(100) | NOT NULL | Nama mapel (Matematika) |
| `code` | VARCHAR(20) | NOT NULL | Kode (MTK) |
| `category` | VARCHAR(50) | NULLABLE | Umum/Jurusan/Mulok |
| `description` | TEXT | NULLABLE | Deskripsi |
| `is_active` | BOOLEAN | DEFAULT true | Status aktif |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Relasi**:
- Belongs to `units`
- Has many `teacher_subjects`

---

### 19. TEACHER_SUBJECTS (Guru-Mapel)

**Deskripsi**: Relasi many-to-many antara guru dan mata pelajaran.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `teacher_profile_id` | UUID | FK → teacher_profiles.id, NOT NULL | Guru |
| `subject_id` | UUID | FK → subjects.id, NOT NULL | Mata pelajaran |
| `is_primary` | BOOLEAN | DEFAULT false | Mapel utama guru |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Relasi**:
- Belongs to `teacher_profiles`
- Belongs to `subjects`

---

### 20. ACTIVITIES (Kegiatan)

**Deskripsi**: Master data kegiatan sekolah (ekstrakurikuler, kajian, event).

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `unit_id` | UUID | FK → units.id, NOT NULL | Sekolah |
| `name` | VARCHAR(100) | NOT NULL | Nama (Pramuka, Kajian Fiqih) |
| `type` | VARCHAR(50) | NOT NULL | ekstrakurikuler/kajian/event |
| `category` | VARCHAR(50) | NULLABLE | halaqah/tahsin/daurah/olahraga/seni/akademik |
| `description` | TEXT | NULLABLE | Deskripsi |
| `start_date` | DATE | NULLABLE | Tanggal mulai kegiatan |
| `end_date` | DATE | NULLABLE | Tanggal berakhir (null = ongoing) |
| `recurrence_type` | VARCHAR(20) | DEFAULT 'none' | none/daily/weekly/monthly |
| `recurrence_days` | INTEGER[] | NULLABLE | Array hari [0-6] untuk weekly, [1-31] untuk monthly |
| `start_time` | VARCHAR(10) | NULLABLE | Jam mulai (14:00) |
| `end_time` | VARCHAR(10) | NULLABLE | Jam selesai (16:00) |
| `location` | VARCHAR(200) | NULLABLE | Lokasi/tempat kegiatan |
| `max_participants` | INT | NULLABLE | Batas maksimal peserta |
| `fee` | DECIMAL(12,2) | NULLABLE | Biaya kegiatan (Rp) |
| `is_active` | BOOLEAN | DEFAULT true | Status aktif |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `updated_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Contoh recurrence_days**:
- Weekly: `[1, 3]` = Senin & Rabu
- Monthly: `[1, 15]` = Tanggal 1 & 15

**Relasi**:
- Belongs to `units`
- Has many `activity_teachers`
- Has many `activity_students`

---

### 21. ACTIVITY_TEACHERS (Pembina Kegiatan)

**Deskripsi**: Guru yang ditugaskan sebagai pembina/pengisi kegiatan.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `activity_id` | UUID | FK → activities.id, NOT NULL | Kegiatan |
| `teacher_profile_id` | UUID | FK → teacher_profiles.id, NOT NULL | Guru |
| `role` | VARCHAR(50) | DEFAULT 'pembina' | pembina/pengisi/koordinator |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Relasi**:
- Belongs to `activities`
- Belongs to `teacher_profiles`

---

### 22. ACTIVITY_STUDENTS (Peserta Kegiatan)

**Deskripsi**: Siswa yang terdaftar dalam kegiatan.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| `activity_id` | UUID | FK → activities.id, NOT NULL | Kegiatan |
| `student_profile_id` | UUID | FK → student_profiles.id, NOT NULL | Siswa |
| `is_mandatory` | BOOLEAN | DEFAULT false | Wajib/Pilihan |
| `joined_at` | DATE | NULLABLE | Tanggal gabung |
| `created_at` | TIMESTAMP | NOT NULL | Auto-fill |
| `deleted_at` | TIMESTAMP | INDEX, NULLABLE | Soft delete |

**Relasi**:
- Belongs to `activities`
- Belongs to `student_profiles`

---

## 📋 Ringkasan Relasi

| Parent | Child | Tipe | Keterangan |
|--------|-------|------|------------|
| `users` | `organizations` | 1:N | User owns Organizations |
| `organizations` | `units` | 1:N | Yayasan has Sekolah |
| `organizations` | `organization_members` | 1:N | Yayasan has Pengurus |
| `organizations` | `roles` | 1:N | Yayasan defines Roles |
| `units` | `unit_settings` | 1:1 | Sekolah has Settings |
| `units` | `unit_members` | 1:N | Sekolah has Warga |
| `units` | `posts` | 1:N | Sekolah has Posts |
| `units` | `teacher_profiles` | 1:N | Sekolah has Guru |
| `units` | `student_profiles` | 1:N | Sekolah has Siswa |
| `units` | `classes` | 1:N | Sekolah has Kelas |
| `users` | `organization_members` | 1:N | User joins Yayasan |
| `users` | `unit_members` | 1:N | User joins Sekolah |
| `users` | `teacher_profiles` | 1:1 | User has Teacher Profile |
| `users` | `student_profiles` | 1:1 | User has Student Profile |
| `roles` | `role_permissions` | N:M | Role has Permissions |
| `posts` | `post_comments` | 1:N | Post has Comments |
| `posts` | `post_poll_options` | 1:N | Post has Poll Options |
| `post_poll_options` | `post_poll_votes` | 1:N | Option has Votes |
| `student_profiles` | `class_enrollments` | 1:N | Student has Enrollments |
| `classes` | `class_enrollments` | 1:N | Class has Enrollments |
| `units` | `subjects` | 1:N | Sekolah has Mapel |
| `subjects` | `teacher_subjects` | 1:N | Mapel has Guru |
| `teacher_profiles` | `teacher_subjects` | 1:N | Guru has Mapel |
| `teacher_profiles` | `classes` | 1:N | Guru as Wali Kelas |
| `units` | `activities` | 1:N | Sekolah has Kegiatan |
| `activities` | `activity_teachers` | 1:N | Kegiatan has Pembina |
| `activities` | `activity_students` | 1:N | Kegiatan has Peserta |
| `teacher_profiles` | `activity_teachers` | 1:N | Guru as Pembina |
| `student_profiles` | `activity_students` | 1:N | Siswa in Kegiatan |

---

## 🔧 Catatan Teknis

### Soft Delete
Semua tabel utama menggunakan soft delete (`deleted_at`). Data tidak benar-benar dihapus, hanya ditandai dengan timestamp.

### UUID
Semua primary key menggunakan UUID v4 untuk keamanan dan skalabilitas.

### JSONB
Field `settings` di `organizations` dan `units` menggunakan JSONB untuk fleksibilitas penyimpanan konfigurasi.

### Timestamps
- `created_at`: Auto-fill saat create (via GORM hook)
- `updated_at`: Auto-update saat update (via GORM hook)

---

**File Location**: `sekolah-madrasah-backend/database/schemas/`
