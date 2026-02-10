---
description: Panduan untuk memulai sesi pengembangan dengan AI Agent
---

# Workflow: Onboarding AI Agent

Gunakan workflow ini di awal sesi dengan AI Agent baru untuk memastikan konteks yang benar.

---

## Step 1: Baca Konteks Project

// turbo
```bash
cat .agent/AI_CONTEXT.md
```

Pastikan agent memahami:
- Struktur hierarki (Organization → Unit)
- Database schema yang aktif
- Coding conventions

---

## Step 2: Cek Status Pengembangan

// turbo
```bash
cat .agent/TEMPLATE_GUIDE.md
```

---

## Step 3: Cek Struktur Backend

// turbo
```bash
ls -la sekolah-madrasah-backend/app/controller/
ls -la sekolah-madrasah-backend/database/schemas/
```

---

## Step 4: Cek Struktur Frontend

// turbo
```bash
ls -la sekolah-madrasah-frontend/src/routes/org/\[orgId\]/units/\[unitId\]/
```

---

## Step 5: Cek Server Status

// turbo
```bash
curl -s http://localhost:8080/api/v1/ping
```

Jika tidak ada response, jalankan server:
```bash
# Terminal 1
cd sekolah-madrasah-backend && go run main.go

# Terminal 2
cd sekolah-madrasah-frontend && npm run dev
```

---

## Reminders untuk Agent

Setelah onboarding, agent harus mengikuti aturan:

1. **JANGAN** push ke git tanpa konfirmasi user
2. **JANGAN** modifikasi `auth_controller` tanpa diminta
3. **CSS Styling**: Boleh kombinasi TailwindCSS utility classes + custom Vanilla CSS
4. **SELALU** gunakan UUID untuk primary key
5. **SELALU** tambahkan soft delete (`DeletedAt`)
6. **SELALU** ikuti pola kode yang sudah ada
