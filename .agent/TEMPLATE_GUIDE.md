# Base Template: Core System (RBAC + Multi-tenant)

This project (`sekolah-madrasah`) is a **Base Template** designed for rapid development of SaaS-like applications requiring:
1.  **Authentication** (JWT, Refresh Tokens)
2.  **RBAC** (Role Based Access Control: Users, Roles, Permissions)
3.  **Multi-tenancy** (Organizations, Units/Branches)
4.  **Clean Architecture** (Backend) & **SvelteKit** (Frontend)

## Architecture Overview

### Backend (`/backend`)
Built with **Go (Golang)** using **Clean Architecture** principles.

*   **`app/controller`**: HTTP Handlers (Gin). Input validation, response formatting.
*   **`app/usecase`**: Business Logic. Co-ordinates repositories.
*   **`app/repository`**: Data Access Layer (GORM). database queries.
*   **`database/schemas`**: Database Models/Entities.
*   **`routes`**: URL routing and Middleware wiring.

**Key Core Modules:**
*   `auth`: Login, Register, Token Management.
*   `user`: User profile management.
*   `organization`: Multi-tenant containers (e.g., Schools, Companies).
*   `unit`: Sub-units within organizations (e.g., Classes, Departments).
*   `role` & `permission`: Dynamic RBAC system.
*   `post`: Basic announcement/content system (Core feature).

### Frontend (`/frontend`)
Built with **SvelteKit** + **TailwindCSS**.

*   **`src/routes`**: File-based routing.
    *   `/auth`: Login/Register pages.
    *   `/dashboard`: Global Super Admin dashboard.
    *   `/org/[orgId]`: Organization-specific dashboard context.
*   **`src/lib`**: Shared utilities, stores, and API client.
*   **`src/core`**: Core UI components (Buttons, Modals, Tables).

---

## How to Clone & Rename (New Project Guide)

To use this as a template for a new project (e.g., `klinik-kesehatan`), follow these steps:

### 1. File Copy
Copy the entire folder structure, excluding `node_modules`, `.git`, `tmp`, and build artifacts.
```bash
cp -r sekolah-madrasah klinik-kesehatan
cd klinik-kesehatan
```

### 2. Backend Renaming
Find and replace the module name in all `.go`, `go.mod` files.
*   **Search**: `sekolah-madrasah`
*   **Replace**: `klinik-kesehatan`

**Command (Mac/Linux):**
```bash
cd klinik-kesehatan-backend
find . -type f -name "*.go" -print0 | xargs -0 sed -i '' 's/sekolah-madrasah/klinik-kesehatan/g'
sed -i '' 's/sekolah-madrasah/klinik-kesehatan/g' go.mod
```

### 3. Frontend Branding
Updates the visible names and branding.
*   **Search**: `Sekolah Madrasah` -> `Klinik Kesehatan`
*   **Search**: `sekolah-madrasah` -> `klinik-kesehatan` (for keys/constants)
*   **Key Files**:
    *   `src/routes/+layout.svelte` (Sidebar titles)
    *   `src/routes/+page.svelte` (Landing page)
    *   `src/routes/auth/...` (Login/Register titles)
    *   `.env` (API URLs if changed)

### 4. Database Setup
1.  Update `docker-compose.yml`:
    *   Change container names: `sekolah_madrasah_postgres` -> `klinik_postgres`
    *   Change DB name/user in `environment`.
    *   **IMPORTANT**: Change exposed port (e.g., `5434` -> `5435`) to avoid conflicts.
2.  Update `.env` in backend/frontend to match new DB credentials and ports.

### 5. Domain Modeling (Adding New Features)
The system separates **Core** (Stability) from **Business** (Innovation).
*   **Do NOT modify**: `auth`, `role`, `permission` (unless necessary).
*   **Add New**: Create new folders in `app/controller`, `app/usecase`, etc. for your specific domain (e.g., `patients`, `appointments`).
