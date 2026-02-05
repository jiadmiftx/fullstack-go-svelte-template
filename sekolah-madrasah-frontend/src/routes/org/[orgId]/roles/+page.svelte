<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import {
        api,
        selectedOrganization,
        type Role,
        type Permission,
    } from "$lib";
    import Modal from "$core/components/Modal.svelte";
    import DataTable from "$core/components/DataTable.svelte";
    import { showToast } from "$core/components/Toast.svelte";

    let orgId: string = "";
    let roles: Role[] = [];
    let permissions: Permission[] = [];
    let isLoading = false;
    let error = "";
    let showModal = false;
    let editingRole: Role | null = null;

    // Form fields
    let formName = "";
    let formDisplayName = "";
    let formType = "custom";
    let formLevel = 1;
    let formDescription = "";
    let formPermissionIds: string[] = [];

    $: orgId = $page.params.orgId ?? "";
    $: organization = $selectedOrganization;

    onMount(async () => {
        if (orgId) {
            await Promise.all([loadRoles(), loadPermissions()]);
        }
    });

    async function loadRoles() {
        isLoading = true;
        error = "";
        try {
            console.log("Loading roles for orgId:", orgId);
            const response = await api.getRoles({
                limit: 100,
                organization_id: orgId,
            });
            console.log("Roles response:", response);
            // Force reactivity by creating new array reference
            roles = [...(response.data || [])];
            console.log("Roles loaded:", roles.length, roles);
        } catch (err) {
            console.error("Error loading roles:", err);
            error = err instanceof Error ? err.message : "Failed to load roles";
        } finally {
            isLoading = false;
        }
    }

    async function loadPermissions() {
        try {
            const response = await api.getPermissions({ limit: 100 });
            permissions = response.data || [];
        } catch (err) {
            console.error("Failed to load permissions:", err);
        }
    }

    function openCreateModal() {
        editingRole = null;
        formName = "";
        formDisplayName = "";
        formType = "custom";
        formLevel = 1;
        formDescription = "";
        formPermissionIds = [];
        showModal = true;
    }

    async function openEditModal(role: Role) {
        editingRole = role;
        formName = role.name;
        formDisplayName = role.display_name || "";
        formType = role.type || "custom";
        formLevel = role.level;
        formDescription = role.description || "";

        // Fetch role with full permissions data
        try {
            console.log("Fetching role with id:", role.id);
            const fullRole = await api.getRole(role.id);
            console.log("Full role response:", fullRole);
            console.log(
                "Permissions from response:",
                fullRole.data?.permissions,
            );
            formPermissionIds =
                fullRole.data?.permissions?.map((p: Permission) => p.id) || [];
            console.log("Parsed permission IDs:", formPermissionIds);
        } catch (err) {
            console.error("Failed to load role permissions:", err);
            formPermissionIds = role.permissions?.map((p) => p.id) || [];
        }
        showModal = true;
    }

    async function handleSubmit() {
        error = "";
        try {
            if (editingRole) {
                console.log("Updating role:", editingRole.id);
                await api.updateRole(editingRole.id, {
                    name: formName,
                    display_name: formDisplayName,
                    description: formDescription,
                    level: formLevel,
                    permission_ids: formPermissionIds,
                });
                showToast("Role updated successfully", "success");
            } else {
                const payload = {
                    organization_id: orgId,
                    name: formName,
                    display_name: formDisplayName,
                    type: formType,
                    level: formLevel,
                    description: formDescription,
                    permission_ids: formPermissionIds,
                };
                console.log("Creating role with payload:", payload);
                const response = await api.createRole(payload);
                console.log("Create role response:", response);
                showToast("Role created successfully", "success");
            }
            showModal = false;
            await loadRoles();
        } catch (err) {
            const errorMsg =
                err instanceof Error ? err.message : "Failed to save role";
            error = errorMsg;
            showToast(errorMsg, "error");
        }
    }

    async function handleDelete(role: Role) {
        if (!confirm(`Are you sure you want to delete "${role.name}"?`)) return;
        try {
            await api.deleteRole(role.id);
            showToast("Role deleted successfully", "success");
            await loadRoles();
        } catch (err) {
            const errorMsg =
                err instanceof Error ? err.message : "Failed to delete role";
            showToast(errorMsg, "error");
        }
    }

    function togglePermission(permId: string) {
        if (formPermissionIds.includes(permId)) {
            formPermissionIds = formPermissionIds.filter((id) => id !== permId);
        } else {
            formPermissionIds = [...formPermissionIds, permId];
        }
    }

    function selectAllPermissions() {
        formPermissionIds = permissions.map((p) => p.id);
    }

    function deselectAllPermissions() {
        formPermissionIds = [];
    }

    const columns = [
        { key: "name" as keyof Role, label: "Name", sortable: true },
        {
            key: "display_name" as keyof Role,
            label: "Display Name",
            render: (name: string) => name || "-",
        },
        {
            key: "type" as keyof Role,
            label: "Type",
            render: (type: string) =>
                `<span class="bg-gray-500/20 text-gray-600 px-2 py-1 rounded text-xs capitalize">${type}</span>`,
        },
        { key: "level" as keyof Role, label: "Level", sortable: true },
        {
            key: "permissions" as keyof Role,
            label: "Permissions",
            render: (perms: Permission[]) => perms?.length || 0,
        },
    ];

    const actions = [
        {
            label: "Edit",
            variant: "primary" as const,
            onClick: openEditModal,
        },
        {
            label: "Delete",
            variant: "danger" as const,
            onClick: handleDelete,
        },
    ];
</script>

<svelte:head>
    <title
        >Roles - {organization?.name || "Organization"} - Multi Tenant RBAC</title
    >
</svelte:head>

<div class="py-8 px-4">
    <div class="max-w-7xl mx-auto">
        <!-- Header -->
        <div class="flex justify-between items-center mb-8">
            <div>
                <h1 class="text-3xl font-bold mb-2">Organization Roles</h1>
                <p class="text-slate-400">
                    Manage roles for {organization?.name || "this organization"}
                </p>
            </div>
            <button on:click={openCreateModal} class="btn-primary">
                + Create Role
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

        <!-- Roles Table -->
        <div class="card">
            {#if isLoading}
                <div class="text-center py-12 text-slate-400">
                    Loading roles...
                </div>
            {:else}
                <DataTable
                    data={roles}
                    {columns}
                    {actions}
                    emptyMessage="No roles in this organization yet. Create one to get started."
                />
            {/if}
        </div>
    </div>
</div>

<!-- Create/Edit Modal -->
<Modal
    bind:isOpen={showModal}
    title={editingRole ? "Edit Role" : "Create Role"}
    size="lg"
>
    <form on:submit|preventDefault={handleSubmit} class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
            <div>
                <label
                    for="name"
                    class="block text-sm font-medium text-slate-300 mb-1"
                >
                    Name *
                </label>
                <input
                    id="name"
                    type="text"
                    bind:value={formName}
                    class="input-field"
                    placeholder="e.g. pengurus"
                    required
                />
            </div>

            <div>
                <label
                    for="displayName"
                    class="block text-sm font-medium text-slate-300 mb-1"
                >
                    Display Name *
                </label>
                <input
                    id="displayName"
                    type="text"
                    bind:value={formDisplayName}
                    class="input-field"
                    placeholder="e.g. Pengurus"
                    required
                />
            </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div>
                <label
                    for="type"
                    class="block text-sm font-medium text-slate-300 mb-1"
                >
                    Type *
                </label>
                <select
                    id="type"
                    bind:value={formType}
                    class="input-field"
                    disabled={!!editingRole}
                >
                    <option value="owner">Owner</option>
                    <option value="admin">Admin</option>
                    <option value="member">Member</option>
                    <option value="custom">Custom</option>
                </select>
            </div>

            <div>
                <label
                    for="level"
                    class="block text-sm font-medium text-slate-300 mb-1"
                >
                    Level *
                </label>
                <input
                    id="level"
                    type="number"
                    bind:value={formLevel}
                    class="input-field"
                    min="1"
                    max="100"
                    required
                />
            </div>
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
                rows="2"
                placeholder="Role description..."
            ></textarea>
        </div>

        <!-- Permissions -->
        <div>
            <div class="flex justify-between items-center mb-2">
                <label class="block text-sm font-medium text-slate-300">
                    Permissions ({formPermissionIds.length} selected)
                </label>
                <div class="flex gap-2">
                    <button
                        type="button"
                        on:click={selectAllPermissions}
                        class="text-xs px-2 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded"
                    >
                        Select All
                    </button>
                    <button
                        type="button"
                        on:click={deselectAllPermissions}
                        class="text-xs px-2 py-1 bg-gray-600 hover:bg-gray-700 text-white rounded"
                    >
                        Deselect All
                    </button>
                </div>
            </div>
            <div
                class="border border-slate-600 rounded-lg p-3 max-h-48 overflow-y-auto bg-slate-800/50"
            >
                {#if permissions.length === 0}
                    <p class="text-slate-400 text-sm">
                        No permissions available. Create some first.
                    </p>
                {:else}
                    <div class="grid grid-cols-2 gap-2">
                        {#each permissions as perm}
                            <label
                                class="flex items-center gap-2 text-sm cursor-pointer hover:bg-slate-700/50 p-1 rounded"
                            >
                                <input
                                    type="checkbox"
                                    checked={formPermissionIds.includes(
                                        perm.id,
                                    )}
                                    on:change={() => togglePermission(perm.id)}
                                    class="rounded"
                                />
                                <span>{perm.resource}.{perm.action}</span>
                            </label>
                        {/each}
                    </div>
                {/if}
            </div>
        </div>

        <div class="flex justify-end gap-3 pt-4">
            <button
                type="button"
                on:click={() => (showModal = false)}
                class="btn-secondary"
            >
                Cancel
            </button>
            <button type="submit" class="btn-primary">
                {editingRole ? "Update" : "Create"} Role
            </button>
        </div>
    </form>
</Modal>
