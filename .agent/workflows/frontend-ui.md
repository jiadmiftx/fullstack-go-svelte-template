---
description: How to build frontend UI for units/sekolah pages - MUST use shared styles and components
---

# Frontend UI Development Workflow

This workflow defines the **mandatory rules** for building frontend UI within the `units/[unitId]/` section of the application.

## 🚨 IMPORTANT RULES

1. **ALWAYS** use shared styles from `$lib/styles/`
2. **ALWAYS** use shared components from `$lib/components/`
3. **NEVER** duplicate styles that already exist in shared CSS
4. **NEVER** create new components if a shared one can be extended

---

## Setup

### 1. Import Shared Styles in Layout

In your `+layout.svelte` or at the top of your page, import the shared styles:

```svelte
<svelte:head>
    <link rel="stylesheet" href="$lib/styles/variables.css" />
    <link rel="stylesheet" href="$lib/styles/components.css" />
</svelte:head>
```

Or import in your main CSS:

```css
@import '$lib/styles/variables.css';
@import '$lib/styles/components.css';
```

### 2. Import Components

```svelte
<script>
    import { 
        PageHeader, 
        SearchBar, 
        DataTable, 
        ActionButton, 
        Modal, 
        Pagination 
    } from '$lib';
</script>
```

---

## Available Components

### PageHeader

Page header with icon, title, subtitle, and action button.

```svelte
<PageHeader
    icon="users"
    title="Data Guru"
    subtitle="Kelola data guru dan tenaga pendidik"
    showAction={canManage}
    actionLabel="Tambah Guru"
    on:action={() => showCreateModal = true}
/>
```

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `icon` | string | "file" | Icon name for header |
| `title` | string | "Page Title" | Page title |
| `subtitle` | string | "" | Subtitle text |
| `showAction` | boolean | true | Show action button |
| `actionLabel` | string | "Tambah" | Action button label |
| `actionIcon` | string | "+" | Action button icon |

---

### SearchBar

Search input with button and stats display.

```svelte
<SearchBar
    placeholder="Cari nama, NIP, atau NUPTK..."
    bind:value={searchQuery}
    totalItems={totalItems}
    itemLabel="guru"
    on:search={handleSearch}
/>
```

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `placeholder` | string | "Cari..." | Input placeholder |
| `value` | string | "" | Search value (bindable) |
| `totalItems` | number | 0 | Total items count |
| `itemLabel` | string | "item" | Label for stats |
| `showStats` | boolean | true | Show stats section |

---

### DataTable

Data table with loading/empty states.

```svelte
<DataTable
    columns={[
        { key: 'name', label: 'Nama' },
        { key: 'email', label: 'Email' },
        { key: 'actions', label: 'Aksi', width: '120px' }
    ]}
    data={items}
    isLoading={isLoading}
    emptyTitle="Belum Ada Data"
    emptyMessage="Mulai dengan menambahkan data baru"
>
    <svelte:fragment slot="row" let:item>
        <td>{item.name}</td>
        <td>{item.email}</td>
        <td>
            <div class="actions">
                <ActionButton type="view" href="/detail/{item.id}" />
                <ActionButton type="edit" on:click={() => edit(item)} />
                <ActionButton type="delete" on:click={() => remove(item)} />
            </div>
        </td>
    </svelte:fragment>
    
    <svelte:fragment slot="empty-action">
        <button class="btn-primary" on:click={create}>+ Tambah</button>
    </svelte:fragment>
</DataTable>
```

---

### ActionButton

Action button for tables.

```svelte
<!-- Link -->
<ActionButton type="view" href="/detail/{item.id}" title="Lihat" />

<!-- Button -->
<ActionButton type="edit" on:click={() => edit(item)} title="Edit" />
<ActionButton type="delete" on:click={() => remove(item)} title="Hapus" />
<ActionButton type="reset" on:click={() => reset(item)} title="Reset" />
```

| Type | Description |
|------|-------------|
| `view` | Eye icon - for viewing details |
| `edit` | Pencil icon - for editing |
| `delete` | Trash icon - red on hover |
| `reset` | Refresh icon - for reset actions |

---

### Modal

Modal dialog with slots.

```svelte
<Modal
    show={showModal}
    title="Tambah Data"
    size="default"
    on:close={() => showModal = false}
>
    <!-- Modal body content -->
    <div class="form-grid">
        <div class="form-group">
            <label>Nama</label>
            <input type="text" bind:value={name} />
        </div>
    </div>
    
    <!-- Actions slot -->
    <svelte:fragment slot="actions">
        <button class="btn-secondary" on:click={() => showModal = false}>
            Batal
        </button>
        <button class="btn-primary" on:click={save}>
            Simpan
        </button>
    </svelte:fragment>
</Modal>
```

| Size | Max Width |
|------|-----------|
| `small` | 400px |
| `default` | 600px |
| `large` | 800px |

---

### Pagination

Pagination controls.

```svelte
<Pagination
    currentPage={currentPage}
    totalPages={totalPages}
    on:change={(e) => {
        currentPage = e.detail.page;
        loadData();
    }}
/>
```

---

## CSS Variables

All components use CSS variables defined in `$lib/styles/variables.css`. Key variables:

### Colors
```css
--color-primary: #00ced1;
--color-text-primary: #111827;
--color-text-secondary: #6b7280;
--color-bg-page: #f5f5f7;
--color-bg-card: #ffffff;
--color-border: #e5e7eb;
--color-danger: #dc2626;
```

### Typography
```css
--font-size-sm: 0.875rem;
--font-size-base: 1rem;
--font-size-3xl: 1.75rem;  /* Page titles */
```

### Spacing
```css
--spacing-sm: 0.5rem;
--spacing-lg: 1rem;
--spacing-xl: 1.5rem;
--page-padding: 2rem;
```

### Border Radius
```css
--radius-md: 8px;
--radius-lg: 10px;
--radius-xl: 12px;
--radius-2xl: 16px;
```

---

## Shared CSS Classes

Use these classes from `$lib/styles/components.css`:

| Class | Description |
|-------|-------------|
| `.glass-card` | White card with shadow |
| `.page-container` | Page wrapper with padding |
| `.btn-primary` | Primary teal button |
| `.btn-secondary` | Secondary gray button |
| `.btn-danger` | Red danger button |
| `.btn-action` | Table action button |
| `.form-grid` | 2-column form layout |
| `.form-group` | Form field container |
| `.badge` | Status badge |
| `.avatar` | User avatar |
| `.error-message` | Error alert |

---

## Example Page Structure

```svelte
<script lang="ts">
    import { PageHeader, SearchBar, DataTable, ActionButton, Modal, Pagination } from '$lib';
    
    let items = [];
    let isLoading = true;
    let searchQuery = '';
    let currentPage = 1;
    let totalPages = 1;
    let showCreateModal = false;
    
    // ... load data, handlers, etc.
</script>

<div class="page-container">
    <PageHeader
        icon="users"
        title="Data Master"
        subtitle="Kelola data master"
        on:action={() => showCreateModal = true}
    />
    
    <SearchBar
        placeholder="Cari..."
        bind:value={searchQuery}
        totalItems={items.length}
        on:search={loadData}
    />
    
    <DataTable
        columns={columns}
        data={items}
        {isLoading}
    >
        <svelte:fragment slot="row" let:item>
            <td>{item.name}</td>
            <td>
                <div class="actions">
                    <ActionButton type="edit" on:click={() => edit(item)} />
                    <ActionButton type="delete" on:click={() => remove(item)} />
                </div>
            </td>
        </svelte:fragment>
    </DataTable>
    
    <Pagination
        {currentPage}
        {totalPages}
        on:change={(e) => { currentPage = e.detail.page; loadData(); }}
    />
</div>

<Modal show={showCreateModal} title="Tambah Data" on:close={() => showCreateModal = false}>
    <!-- form -->
</Modal>

<style>
    /* Only page-specific styles here! */
    /* Use shared classes from components.css */
</style>
```

---

## Extending Components

If you need to customize a component:

1. **First** check if you can use CSS variables
2. **Second** check if adding a prop makes sense
3. **Last resort**: Create page-specific styles (document why!)

---

// turbo-all
