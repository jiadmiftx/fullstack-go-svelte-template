# Frontend SvelteKit Template

A production-ready SvelteKit template with **RBAC (Role-Based Access Control)**, authentication, and clean architecture separating **core** infrastructure from **business** logic.

## 🎯 Features

### Core (Reusable Infrastructure)
- ✅ **Authentication** - Login/logout with JWT
- ✅ **RBAC** - Role-based access control
- ✅ **User Management** - CRUD operations
- ✅ **Reusable Components** - DataTable, Toast, Modal, StatsCard
- ✅ **API Client** - Axios-like wrapper with auth
- ✅ **Type Safety** - Full TypeScript support

### Business (Example Implementation)
- 📦 **Multi-tenant Organizations**
- 📊 **Charts** - Bar & Line charts
- 🎨 **Light Mode Theme** - Neutral gray/black palette
- 🎨 **Hugeicons** - SVG icons (Stroke Rounded)

---

## 📁 Project Structure

```
frontend-svelte-template/
├── src/
│   ├── core/                    # 🔧 Reusable (Don't modify for new projects)
│   │   ├── components/          # DataTable, Toast, Modal, etc.
│   │   ├── api/                 # API client
│   │   ├── stores/              # Auth & organization stores
│   │   ├── types/               # User, Role, Permission types
│   │   └── index.ts             # Core exports
│   │
│   ├── business/                # 📦 Customize for your use case
│   │   ├── components/          # Your domain components
│   │   ├── types/               # Your domain types
│   │   └── index.ts             # Business exports
│   │
│   ├── lib/                     # Re-exports core + business
│   │   └── index.ts
│   │
│   └── routes/                  # Your pages
│       ├── +layout.svelte       # Main layout with sidebar
│       ├── auth/                # Login/register pages
│       └── ...                  # Your app pages
```

---

## 🚀 Quick Start

### 1. Clone & Install
```bash
git clone <your-repo-url>
cd frontend-svelte-template
npm install
```

### 2. Configure Backend API
Update the API URL in `src/core/api/client.ts`:
```typescript
export const API_BASE_URL = 'http://localhost:8080/api/v1';  // Your backend URL
```

### 3. Run Development Server
```bash
npm run dev
```
Open [http://localhost:5173](http://localhost:5173)

---

## 🛠️ How to Customize

### Option 1: Keep Business Logic (Multi-tenant)
If you're building a multi-tenant app with organizations:
1. Keep `src/business/` as-is
2. Customize pages in `src/routes/`
3. Add your domain logic to `src/business/`

### Option 2: Replace Business Logic
If you're building something else (e.g., e-commerce, CMS):

1. **Clear business directory:**
   ```bash
   rm -rf src/business/*
   mkdir -p src/business/components src/business/types
   ```

2. **Create your domain types** in `src/business/types/index.ts`:
   ```typescript
   // Example: E-commerce
   export interface Product {
     id: string;
     name: string;
     price: number;
     // ...
   }
   ```

3. **Create your components** in `src/business/components/`:
   ```svelte
   <!-- ProductCard.svelte -->
   <script lang="ts">
     import type { Product } from '$business/types';
     export let product: Product;
   </script>
   ```

4. **Export from** `src/business/index.ts`:
   ```typescript
   export { default as ProductCard } from './components/ProductCard.svelte';
   export type { Product } from './types';
   ```

5. **Update routes** to use your new business logic

---

## 📦 Importing Components

### Core Components (Always Available)
```svelte
<script>
  import DataTable from "$core/components/DataTable.svelte";
  import { showToast } from "$core/components/Toast.svelte";
  import { auth, isAuthenticated } from "$core";
</script>
```

### Business Components (Your Domain)
```svelte
<script>
  import ProductCard from "$business/components/ProductCard.svelte";
  import type { Product } from "$business";
</script>
```

---

## 🎨 UI Components

### DataTable
Powerful table with search, sort, pagination:
```svelte
<DataTable 
  data={users} 
  columns={[
    { key: 'id', label: 'ID' },
    { key: 'email', label: 'Email' }
  ]}
  searchable
  itemsPerPage={10}
/>
```

### Toast Notifications
```svelte
<script>
  import { showToast } from "$core/components/Toast.svelte";
  
  function handleSuccess() {
    showToast('Success!', 'success');
  }
</script>
```

### Modal
```svelte
<Modal bind:isOpen={showModal} title="Add User">
  <!-- Your form here -->
</Modal>
```

---

## 🔐 Authentication

The template includes a complete auth system:

```svelte
<script>
  import { auth, isAuthenticated, currentUser } from "$core";
  
  async function login() {
    await auth.login('user@example.com', 'password');
  }
  
  function logout() {
    auth.logout();
  }
</script>

{#if $isAuthenticated}
  <p>Welcome, {$currentUser?.full_name}</p>
  <button on:click={logout}>Logout</button>
{/if}
```

---

## 🏗️ Backend Requirements

This frontend expects a REST API with the following endpoints:

### Auth
- `POST /auth/login` - Login with email/password
- `POST /auth/register` - Register new user

### Users
- `GET /users` - List users
- `POST /users` - Create user
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

### Roles & Permissions
- `GET /roles` - List roles
- `GET /permissions` - List permissions

See the [API Client](src/core/api/client.ts) for full endpoint list.

---

## 🎨 Theming

The template uses **Tailwind CSS** with a neutral light theme. To customize:

1. **Colors** - Edit `src/app.css`
2. **Layout** - Edit `src/routes/+layout.svelte`
3. **Icons** - Replace Hugeicons SVG in components

---

## 📝 Scripts

```bash
npm run dev          # Start dev server
npm run build        # Build for production
npm run preview      # Preview production build
npm run check        # Type check
```

---

## 🤝 Contributing

This is a template. Fork it and make it your own!

---

## 📄 License

MIT License - feel free to use for commercial projects.

---

## 🙏 Credits

Built with:
- [SvelteKit](https://kit.svelte.dev/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Chart.js](https://www.chartjs.org/)
- [Hugeicons](https://hugeicons.com/)
