# Go Clean Architecture Template

🚀 **Production-ready Go backend template** with Multi-Tenant RBAC, Clean Architecture, and **CORE/BUSINESS separation** for maximum reusability.

## Features

### Core (Reusable Infrastructure)
- ✅ **JWT Authentication** - Secure token-based auth with refresh tokens
- ✅ **User Management** - Complete user CRUD
- ✅ **RBAC System** - Roles and permissions
- ✅ **Middleware** - Auth, CORS, rate limiting, logging
- ✅ **Utilities** - Pagination, validation, helpers

### Business (Example: Multi-tenant)
- 📦 **Organizations** - Multi-tenant organization management
- 📦 **Member Management** - Org-scoped user membership

### Technical Stack
- ✅ **Clean Architecture** - Controller → Use Case → Repository
- ✅ **GORM ORM** - PostgreSQL with auto-migration
- ✅ **Gin Framework** - Fast HTTP router
- ✅ **Elastic APM** - Performance monitoring

## Quick Start

### 1. Clone and Rename

```bash
# Clone the template
git clone https://github.com/your-username/vibe-code-go-template.git my-project
cd my-project

# Run rename script
./scripts/rename-project.sh my-project github.com/my-username/my-project
```

### 2. Setup Environment

```bash
# Copy environment file
cp .env.example .env

# Edit .env with your configuration
nano .env
```

### 3. Run Application

```bash
# Install dependencies
go mod tidy

# Run development server
go run main.go
```

## API Endpoints

### Auth (Public)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/auth/register` | Register |
| POST | `/api/v1/auth/refresh` | Refresh token |

### Users (Protected)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users` | List users |
| GET | `/api/v1/users/:id` | Get user |
| POST | `/api/v1/users` | Create user |
| PUT | `/api/v1/users/:id` | Update user |
| DELETE | `/api/v1/users/:id` | Delete user |

### Roles (Protected)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/roles` | List roles |
| GET | `/api/v1/roles/:id` | Get role with permissions |
| POST | `/api/v1/roles` | Create role |
| PUT | `/api/v1/roles/:id` | Update role |
| DELETE | `/api/v1/roles/:id` | Delete role |

### Permissions (Protected)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/permissions` | List permissions |
| GET | `/api/v1/permissions/:id` | Get permission |
| POST | `/api/v1/permissions` | Create permission |
| DELETE | `/api/v1/permissions/:id` | Delete permission |

### Organizations (Protected)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/organizations` | List organizations |
| GET | `/api/v1/organizations/:id` | Get organization |
| POST | `/api/v1/organizations` | Create organization |
| PUT | `/api/v1/organizations/:id` | Update organization |
| DELETE | `/api/v1/organizations/:id` | Delete organization |
| GET | `/api/v1/organizations/:id/members` | List members |
| POST | `/api/v1/organizations/:id/members` | Add member |
| PUT | `/api/v1/organizations/:id/members/:userId` | Update member |
| DELETE | `/api/v1/organizations/:id/members/:userId` | Remove member |

## Project Structure

```
.
├── CORE_MODULES.md        # 🔧 List of reusable modules
├── BUSINESS_MODULES.md    # 📦 List of business modules
├── app/
│   ├── controller/
│   │   ├── auth_controller/        # 🔧 CORE
│   │   ├── user_controller/        # 🔧 CORE
│   │   ├── role_controller/        # 🔧 CORE
│   │   ├── permission_controller/  # 🔧 CORE
│   │   └── organization_controller/ # 📦 BUSINESS
│   ├── use_case/          # Business logic layer
│   └── repository/        # Data access layer
├── pkg/                   # 🔧 CORE utilities
├── config/                # 🔧 CORE configuration
├── database/              # 🔧 CORE DB connection
├── routes/                # Route registration
└── main.go               # Entry point
```

### Core vs Business Modules

**🔧 CORE (Don't Modify)** - Reusable auth & RBAC infrastructure
- See [`CORE_MODULES.md`](CORE_MODULES.md) for complete list
- Total: 15 modules (auth, users, roles, permissions, middleware)

**📦 BUSINESS (Customize)** - Organization-specific logic
- See [`BUSINESS_MODULES.md`](BUSINESS_MODULES.md) for complete list
- Total: 4 modules (organizations, members)
- Replace with your domain (e-commerce, CMS, etc.)

## Development

```bash
# Build
go build -o bin/app main.go

# Run tests
go test ./...

# Generate Swagger docs
make swagger

# Update version
make update-version
```

## Environment Variables

See `.env.example` for all available configuration options.

**Required:**
- `REST_SECRET` - JWT secret key
- `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`, `DB_NAME` - Database

## License

MIT License - feel free to use for any project.