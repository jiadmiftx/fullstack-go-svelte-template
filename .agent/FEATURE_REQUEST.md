# Feature Request: [Nama Fitur]

> **Instruksi untuk AI Agent**: Kerjakan fitur ini sesuai requirements di bawah. Ikuti workflow di `.agent/workflows/development-flow.md`.

---

## 📋 Overview

**Nama Fitur**: [Tuliskan nama fitur]
**Prioritas**: [High / Medium / Low]
**Target**: [Organization / Unit / User]

**Deskripsi Singkat**:
[Jelaskan singkat fitur ini untuk apa]

---

## 📊 Database Requirements

### Tabel Baru
| Tabel | Relasi | Keterangan |
|-------|--------|------------|
| `nama_tabel` | belongs to `units` | Deskripsi |

### Detail Field (per tabel)

#### Tabel: `nama_tabel`
| Field | Type | Nullable | Keterangan |
|-------|------|----------|------------|
| id | UUID | No | Primary Key |
| unit_id | UUID | No | FK → units |
| name | VARCHAR(255) | No | Nama |
| ... | ... | ... | ... |
| created_at | TIMESTAMP | No | Auto |
| updated_at | TIMESTAMP | No | Auto |
| deleted_at | TIMESTAMP | Yes | Soft delete |

---

## 🔌 API Requirements

### Endpoints yang Dibutuhkan
| Method | Path | Keterangan |
|--------|------|------------|
| GET | /api/v1/units/:id/nama-fitur | List all |
| POST | /api/v1/units/:id/nama-fitur | Create new |
| GET | /api/v1/units/:id/nama-fitur/:itemId | Get by ID |
| PUT | /api/v1/units/:id/nama-fitur/:itemId | Update |
| DELETE | /api/v1/units/:id/nama-fitur/:itemId | Delete |

### Request/Response Format (Opsional)
```json
// POST Request Body
{
    "name": "string",
    "field_lain": "value"
}

// Response
{
    "success": true,
    "data": { ... }
}
```

---

## 🖥️ Frontend Requirements

### Halaman yang Dibutuhkan
| Route | Fungsi |
|-------|--------|
| `/org/[orgId]/units/[unitId]/nama-fitur` | List view |
| `/org/[orgId]/units/[unitId]/nama-fitur/[id]` | Detail view |
| `/org/[orgId]/units/[unitId]/nama-fitur/create` | Form create |

### UI Components
- [ ] Tabel dengan pagination
- [ ] Form create/edit
- [ ] Modal konfirmasi delete
- [ ] Search/filter (opsional)

### Menu Navigation
- Label: "Nama Fitur"
- Icon: [nama icon]
- Roles: ["admin", "staff"]

---

## ✅ Acceptance Criteria

- [ ] Tabel database terbuat dengan benar
- [ ] API endpoint berfungsi (CRUD)
- [ ] Frontend terintegrasi dengan API
- [ ] Validasi input berfungsi
- [ ] Error handling proper
- [ ] Responsive design

---

## 📝 Notes / Catatan Tambahan

[Tulis catatan tambahan, referensi, atau hal khusus di sini]

---

## 🚦 Status

| Fase | Status | Tanggal |
|------|--------|---------|
| Design | ⏳ Pending | - |
| Backend | ⏳ Pending | - |
| Frontend | ⏳ Pending | - |

**Legend**: ⏳ Pending | 🔄 In Progress | ✅ Done | ❌ Blocked

---

## 📌 Untuk AI Agent

Setelah selesai mengerjakan setiap fase:
1. Update status di tabel atas
2. Update `.agent/CHANGELOG.md`
3. Minta approval user sebelum lanjut fase berikutnya

**JANGAN lanjut ke fase berikutnya tanpa approval user!**
