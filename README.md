# Sekolah Madrasah - SvelteKit & Go (Clean Architecture)

Modern school management system platform built with SvelteKit (Frontend) and Go (Backend).

## Project Structure
- `sekolah-madrasah-backend/` - Go backend with Clean Architecture
- `sekolah-madrasah-frontend/` - SvelteKit + TailwindCSS frontend

## Prerequisites
- Docker & Docker Compose
- Go 1.22+
- Node.js 18+

## Quick Start

### 1. Start Database (Docker)
```bash
cd sekolah-madrasah-backend
docker-compose up -d postgres
```

### 2. Start Backend
```bash
cd sekolah-madrasah-backend
cp .env.example .env
# Edit .env and ensure DB_PORT=5434 (if using docker-compose provided)
go run main.go
```

### 3. Start Frontend
```bash
cd sekolah-madrasah-frontend
cp .env.example .env
npm install
npm run dev
```

## Default Credentials
- **Super Admin**: superadmin@mail.com / admin123
