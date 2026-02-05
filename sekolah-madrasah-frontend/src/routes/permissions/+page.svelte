<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { isAuthenticated, api, type Permission } from "$lib";
    import Modal from "$core/components/Modal.svelte";
    import DataTable from "$core/components/DataTable.svelte";
    import { showToast } from "$core/components/Toast.svelte";

    let permissions: Permission[] = [];
    let isLoading = false;
    let error = "";
    let showModal = false;

    // Form fields
    let formResource = "";
    let formAction = "";
    let formDescription = "";

    $: if (!$isAuthenticated && typeof window !== "undefined") {
        goto("/auth/login");
    }

    onMount(async () => {
        if ($isAuthenticated) {
            await loadPermissions();
        }
    });

    async function loadPermissions() {
        isLoading = true;
        error = "";
        try {
            const response = await api.getPermissions({ limit: 100 });
            permissions = response.data || [];
        } catch (err) {
            error =
                err instanceof Error
                    ? err.message
                    : "Failed to load permissions";
        } finally {
            isLoading = false;
        }
    }

    function openCreateModal() {
        formResource = "";
        formAction = "";
        formDescription = "";
        showModal = true;
    }

    async function handleSubmit() {
        error = "";
        try {
            await api.createPermission({
                resource: formResource,
                action: formAction,
                description: formDescription || undefined,
            });
            showToast("Permission created successfully", "success");
            showModal = false;
            await loadPermissions();
        } catch (err) {
            const errorMsg =
                err instanceof Error
                    ? err.message
                    : "Failed to create permission";
            error = errorMsg;
            showToast(errorMsg, "error");
        }
    }

    async function handleDelete(perm: Permission) {
        if (!confirm(`Are you sure you want to delete "${perm.name}"?`)) return;

        error = "";
        try {
            await api.deletePermission(perm.id);
            showToast("Permission deleted successfully", "success");
            await loadPermissions();
        } catch (err) {
            const errorMsg =
                err instanceof Error
                    ? err.message
                    : "Failed to delete permission";
            error = errorMsg;
            showToast(errorMsg, "error");
        }
    }

    const columns = [
        {
            key: "resource" as keyof Permission,
            label: "Resource",
            sortable: true,
        },
        { key: "action" as keyof Permission, label: "Action", sortable: true },
        {
            key: "name" as keyof Permission,
            label: "Name",
            render: (name: string) => name || "-",
        },
        {
            key: "description" as keyof Permission,
            label: "Description",
            render: (desc: string) =>
                desc || '<span class="text-slate-500">-</span>',
        },
        {
            key: "created_at" as keyof Permission,
            label: "Created",
            sortable: true,
            render: (value: string) => new Date(value).toLocaleDateString(),
        },
    ];

    const actions = [
        {
            label: "Delete",
            variant: "danger" as const,
            onClick: handleDelete,
        },
    ];
</script>

<svelte:head>
    <title>Permissions - Multi Tenant RBAC</title>
</svelte:head>

{#if $isAuthenticated}
    <div class="py-8 px-4">
        <div class="max-w-7xl mx-auto">
            <!-- Header -->
            <div class="flex justify-between items-center mb-8">
                <div>
                    <h1 class="text-3xl font-bold mb-2">Permissions</h1>
                    <p class="text-slate-400">Manage system permissions</p>
                </div>
                <button on:click={openCreateModal} class="btn-primary">
                    + Create Permission
                </button>
            </div>

            <!-- Error Alert -->
            {#if error}
                <div
                    class="bg-red-500/20 border border-red-500 text-red-400 px-4 py-3 rounded-lg mb-6"
                >
                    {error}
                </div>
            {/if}

            <!-- Permissions Table -->
            <div class="card">
                {#if isLoading}
                    <div class="text-center py-12 text-slate-400">
                        Loading permissions...
                    </div>
                {:else}
                    <DataTable
                        data={permissions}
                        {columns}
                        {actions}
                        emptyMessage="No permissions found. Create your first one!"
                    />
                {/if}
            </div>
        </div>
    </div>

    <!-- Create Modal -->
    <Modal bind:isOpen={showModal} title="Create Permission">
        <form on:submit|preventDefault={handleSubmit} class="space-y-4">
            <div>
                <label
                    for="resource"
                    class="block text-sm font-medium text-slate-300 mb-1"
                >
                    Resource *
                </label>
                <input
                    id="resource"
                    type="text"
                    bind:value={formResource}
                    class="input-field"
                    placeholder="e.g. users"
                    required
                />
                <p class="text-xs text-slate-500 mt-1">
                    Resource name (lowercase, singular)
                </p>
            </div>

            <div>
                <label
                    for="action"
                    class="block text-sm font-medium text-slate-300 mb-1"
                >
                    Action *
                </label>
                <input
                    id="action"
                    type="text"
                    bind:value={formAction}
                    class="input-field"
                    placeholder="e.g. create"
                    required
                />
                <p class="text-xs text-slate-500 mt-1">
                    Action name (create, read, update, delete, manage, etc.)
                </p>
            </div>

            <div>
                <label
                    for="description"
                    class="block text-sm font-medium text-slate-300 mb-1"
                >
                    Description
                </label>
                <textarea
                    id="description"
                    bind:value={formDescription}
                    class="input-field"
                    placeholder="Describe what this permission allows..."
                    rows="3"
                />
            </div>
        </form>

        <div slot="footer">
            <button on:click={() => (showModal = false)} class="btn-secondary">
                Cancel
            </button>
            <button on:click={handleSubmit} class="btn-primary">
                Create
            </button>
        </div>
    </Modal>
{/if}
