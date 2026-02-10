# Feature Request: [Nama Fitur]

> **Untuk User**: Cukup isi bagian yang mudah. AI Agent akan mendesain detail teknis-nya.

---

## 1. Deskripsi Fitur

**Nama Fitur**: [contoh: Classes / Kelas]  
**Deskripsi**: [jelaskan singkat fitur ini untuk apa]  
**Parent**: [Organization / Unit / User] - tabel induk yang berelasi

---

## 2. Data yang Dibutuhkan

Cukup tulis **nama field dan keterangan singkat**. AI Agent akan menentukan tipe data.

### Tabel: `[nama_tabel]`
```
- nama kelas (wajib)
- tingkat/level (angka, contoh: 10, 11, 12)
- wali kelas (referensi ke user, opsional)
- kapasitas maksimal siswa
- deskripsi (opsional)
```

### Tabel: `[nama_tabel_2]` (jika ada)
```
- field 1
- field 2
```

---

## 3. Fitur yang Diinginkan (centang yang perlu)

- [ ] List/Daftar dengan pagination
- [ ] Form Create/Edit
- [ ] Detail view
- [ ] Search/Filter
- [ ] Export data
- [ ] Import data
- [ ] Soft delete
- [ ] Lainnya: ___

---

## 4. Catatan Tambahan (Opsional)

[Tulis hal khusus yang perlu diperhatikan, referensi UI, atau batasan bisnis]

---

## 5. Status (Diisi AI Agent)

| Fase | Status | Tanggal |
|------|--------|---------|
| Design | ⏳ | - |
| Backend | ⏳ | - |
| Frontend | ⏳ | - |

---

# ═══════════════════════════════════════════════════════════════
# CONTOH PENGISIAN (Hapus setelah memahami)
# ═══════════════════════════════════════════════════════════════

## 1. Deskripsi Fitur

**Nama Fitur**: Classes (Kelas/Rombel)  
**Deskripsi**: Manajemen kelas/rombongan belajar per sekolah  
**Parent**: Unit

---

## 2. Data yang Dibutuhkan

### Tabel: `classes`
```
- nama kelas (wajib, contoh: "X IPA 1", "VII A")
- tingkat/level (angka: 1-12 untuk SD-SMA)
- wali kelas (guru, opsional)
- tahun ajaran (contoh: "2025/2026")
- kapasitas maksimal (angka)
- status aktif
```

### Tabel: `class_students` (siswa per kelas)
```
- kelas (referensi)
- siswa (referensi ke user)
- tanggal masuk
- tanggal keluar (opsional, untuk pindah kelas)
```

---

## 3. Fitur yang Diinginkan

- [x] List/Daftar dengan pagination
- [x] Form Create/Edit
- [x] Detail view (lihat daftar siswa)
- [x] Search/Filter (by tingkat, tahun ajaran)
- [ ] Export data
- [ ] Import data
- [x] Soft delete

---

## 4. Catatan Tambahan

- Satu siswa hanya boleh di 1 kelas aktif per tahun ajaran
- Wali kelas harus dari unit member dengan role "guru"
