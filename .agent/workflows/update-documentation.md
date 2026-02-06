---
description: Workflow wajib setelah menambah fitur baru - update dokumentasi
---

# Workflow: Update Dokumentasi

**WAJIB dijalankan setelah menambah fitur/perubahan signifikan ke sistem.**

---

## Step 1: Update CHANGELOG.md

// turbo
Edit file `.agent/CHANGELOG.md`:

1. Tambahkan section baru di **ATAS** (yang terbaru di atas)
2. Format tanggal: `[YYYY-MM-DD]`
3. Isi sesuai template yang sudah disediakan

**Contoh entry:**
```markdown
## [2026-02-07] - Fitur Classes

### Database
- [x] Tabel baru: `classes` - Kelas/Rombel

### API Endpoints
- [x] `GET/POST /api/v1/units/:id/classes`
- [x] `GET/PUT/DELETE /api/v1/units/:id/classes/:classId`

### Frontend
- [x] Route baru: `/org/[orgId]/units/[unitId]/classes`
```

---

## Step 2: Update DATABASE_SCHEMA.md (Jika Ada Tabel/Field Baru)

// turbo
Edit file `.agent/DATABASE_SCHEMA.md` jika ada:

1. **Tabel baru** - Tambahkan section baru dengan format yang sama
2. **Field baru di tabel existing** - Update section tabel tersebut
3. **Relasi baru** - Update diagram relasi di atas dan tabel ringkasan

**Template tabel baru:**
```markdown
### N. NAMA_TABEL (Deskripsi)

**Deskripsi**: Penjelasan singkat tabel.

| Field | Type | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | Primary Key |
| ... | ... | ... | ... |

**Relasi**:
- Belongs to `parent_table`
- Has many `child_table`
```

---

## Step 3: Update AI_CONTEXT.md (Jika Perlu)

Edit file `.agent/AI_CONTEXT.md` jika ada perubahan pada:

1. **Section 2 (Database)** - Jika ada tabel baru
2. **Section 3 (Struktur Kode)** - Jika ada folder baru
3. **Section 6 (API Endpoints)** - Jika ada endpoint baru
4. **Section 9 (Roadmap)** - Update checklist ✅

---

## Step 4: Commit Perubahan

// turbo
```bash
git add .agent/
git commit -m "Docs: Update changelog, schema, and context for [nama fitur]"
```

**JANGAN push tanpa konfirmasi user!**

---

## Reminder

Dokumentasi yang baik memastikan:
- ✅ Agent baru langsung paham konteks
- ✅ User bisa tracking progress
- ✅ Tidak ada informasi yang hilang antar sesi
