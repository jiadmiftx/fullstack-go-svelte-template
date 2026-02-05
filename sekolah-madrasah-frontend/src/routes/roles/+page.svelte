<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import {
        isAuthenticated,
        isSuperAdmin,
        api,
        selectedOrgId,
        selectedOrganization,
        type Role,
        type Permission,
        type Organization,
    } from "$lib";
    import Modal from "$core/components/Modal.svelte";
    import DataTable from "$core/components/DataTable.svelte";
    import { showToast } from "$core/components/Toast.svelte";

    let roles: Role[] = [];
    let permissions: Permission[] = [];
    let organizations: Organization[] = [];
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
    let formOrgId = ""; // For Super Admin to select org

    $: if (!$isAuthenticated && typeof window !== "undefined") {
        goto("/auth/login");
    }

    onMount(async () => {
        if ($isAuthenticated) {
            await Promise.all([
                loadRoles(),
                loadPermissions(),
                loadOrganizations(),
            ]);
        }
    });

    async function loadRoles() {
        isLoading = true;
        error = "";
        try {
            const response = await api.getRoles({ limit: 100 });
            console.log("Roles API Response:", response);
            console.log("Roles data:", response.data);
            roles = response.data || [];
            console.log("Parsed roles:", roles);
        } catch (err) {
            error = err instanceof Error ? err.message : "Failed to load roles";
            console.error("Load roles error:", err);
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

    async function loadOrganizations() {
        try {
            const response = await api.getOrganizations({ limit: 100 });
            organizations = response.data || [];
        } catch (err) {
            console.error("Failed to load organizations:", err);
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
        formOrgId =
            $selectedOrgId ||
            (organizations.length > 0 ? organizations[0].id : "");
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
            const fullRole = await api.getRole(role.id);
            formPermissionIds =
                fullRole.data?.permissions?.map((p: Permission) => p.id) || [];
        } catch (err) {
            console.error("Failed to load role permissions:", err);
            formPermissionIds = role.permissions?.map((p) => p.id) || [];
        }
        showModal = true;
    }

    async function handleSubmit() {
        error = "";

        // For new roles, require an organization
        const orgId = editingRole ? editingRole.organization_id : formOrgId;
        if (!editingRole && !orgId) {
            error = "Please select an organization first";
            showToast("Please select an organization first", "error");
            return;
        }

        try {
            if (editingRole) {
                // Update only allows: name, display_name, description, level, permission_ids
                await api.updateRole(editingRole.id, {
                    name: formName,
                    display_name: formDisplayName,
                    description: formDescription,
                    level: formLevel,
                    permission_ids: formPermissionIds,
                });
                showToast("Role updated successfully", "success");
            } else {
                await api.createRole({
                    organization_id: formOrgId,
                    name: formName,
                    display_name: formDisplayName,
                    type: formType,
                    level: formLevel,
                    description: formDescription,
                    permission_ids: formPermissionIds,
                });
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

        error = "";
        try {
            await api.deleteRole(role.id);
            showToast("Role deleted successfully", "success");
            await loadRoles();
        } catch (err) {
            const errorMsg =
                err instanceof Error ? err.message : "Failed to delete role";
            error = errorMsg;
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
        { key: "level" as keyof Role, label: "Level", sortable: true },
        {
            key: "permissions" as keyof Role,
            label: "Permissions",
            render: (perms: Permission[]) => {
                if (!perms || perms.length === 0)
                    return '<span class="text-slate-500">None</span>';
                return `<span class="text-slate-400">${perms.length} permissions</span>`;
            },
        },
        {
            key: "is_active" as keyof Role,
            label: "Status",
            render: (active: boolean) =>
                active
                    ? '<span class="bg-green-500/20 text-green-400 px-2 py-1 rounded text-xs">Active</span>'
                    : '<span class="bg-red-500/20 text-red-400 px-2 py-1 rounded text-xs">Inactive</span>',
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
    <title>Roles - Multi Tenant RBAC</title>
</svelte:head>

{#if $isAuthenticated}
    <div class="py-8 px-4">
        <div class="max-w-7xl mx-auto">
            <!-- Header -->
            <div class="flex justify-between items-center mb-8">
                <div>
                    <h1 class="text-3xl font-bold mb-2">Roles</h1>
                    <p class="text-slate-400">
                        Manage roles and assign permissions
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
                        emptyMessage="No roles found. Create your first one!"
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
            <!-- Organization selector for new roles -->
            {#if !editingRole}
                <div>
                    <label
                        for="organization"
                        class="block text-sm font-medium text-slate-300 mb-1"
                    >
                        Organization *
                    </label>
                    <select
                        id="organization"
                        bind:value={formOrgId}
                        class="input-field"
                        required
                    >
                        <option value="">Select Organization</option>
                        {#each organizations as org}
                            <option value={org.id}>{org.name}</option>
                        {/each}
                    </select>
                    {#if organizations.length === 0}
                        <p class="text-yellow-400 text-sm mt-1">
                            No organizations found. <a
                                href="/organizations"
                                class="underline">Create one first</a
                            >.
                        </p>
                    {/if}
                </div>
            {/if}

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
                        placeholder="e.g. admin"
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
                        placeholder="e.g. Administrator"
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
                        required
                    >
                        <option value="system">System</option>
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
                    <p class="text-xs text-slate-500 mt-1">
                        Higher level = more authority
                    </p>
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
                    placeholder="Role description..."
                    rows="2"
                />
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
                    class="bg-slate-900 rounded-lg p-4 max-h-64 overflow-y-auto"
                >
                    {#if permissions.length === 0}
                        <p class="text-slate-500 text-sm">
                            No permissions available
                        </p>
                    {:else}
                        <div class="grid grid-cols-2 gap-2">
                            {#each permissions as perm}
                                <label
                                    class="flex items-center gap-2 p-2 hover:bg-slate-800 rounded cursor-pointer"
                                >
                                    <input
                                        type="checkbox"
                                        checked={formPermissionIds.includes(
                                            perm.id,
                                        )}
                                        on:change={() =>
                                            togglePermission(perm.id)}
                                        class="w-4 h-4 rounded border-slate-600 bg-slate-700 text-primary-600 focus:ring-2 focus:ring-primary-500"
                                    />
                                    <div>
                                        <div class="text-sm">{perm.name}</div>
                                        <div class="text-xs text-slate-500">
                                            {perm.resource}.{perm.action}
                                        </div>
                                    </div>
                                </label>
                            {/each}
                        </div>
                    {/if}
                </div>
                <p class="text-xs text-slate-500 mt-1">
                    Selected: {formPermissionIds.length} permission(s)
                </p>
            </div>
        </form>

        <div slot="footer">
            <button on:click={() => (showModal = false)} class="btn-secondary">
                Cancel
            </button>
            <button on:click={handleSubmit} class="btn-primary">
                {editingRole ? "Update" : "Create"}
            </button>
        </div>
    </Modal>
{/if}
