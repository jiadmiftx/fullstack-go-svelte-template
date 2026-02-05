<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { isAuthenticated, api, type Organization } from "$lib";
    import Modal from "$core/components/Modal.svelte";
    import DataTable from "$core/components/DataTable.svelte";
    import { showToast } from "$core/components/Toast.svelte";

    let organizations: Organization[] = [];
    let isLoading = false;
    let error = "";
    let showModal = false;
    let editingOrg: Organization | null = null;

    // Form fields
    let formName = "";
    let formCode = "";
    let formType = "";
    let formDescription = "";

    $: if (!$isAuthenticated && typeof window !== "undefined") {
        goto("/auth/login");
    }

    onMount(async () => {
        if ($isAuthenticated) {
            await loadOrganizations();
        }
    });

    async function loadOrganizations() {
        isLoading = true;
        error = "";
        try {
            const response = await api.getOrganizations({ limit: 100 });
            organizations = response.data || [];
        } catch (err) {
            error =
                err instanceof Error
                    ? err.message
                    : "Failed to load organizations";
        } finally {
            isLoading = false;
        }
    }

    function openCreateModal() {
        editingOrg = null;
        formName = "";
        formCode = "";
        formType = "";
        formDescription = "";
        showModal = true;
    }

    function openEditModal(org: Organization) {
        editingOrg = org;
        formName = org.name;
        formCode = org.code;
        formType = org.type;
        formDescription = org.description || "";
        showModal = true;
    }

    async function handleSubmit() {
        error = "";
        try {
            if (editingOrg) {
                // Update only allows: name, description, address, logo, settings
                await api.updateOrganization(editingOrg.id, {
                    name: formName,
                    description: formDescription,
                });
                showToast("Organization updated successfully", "success");
            } else {
                await api.createOrganization({
                    name: formName,
                    code: formCode,
                    type: formType,
                    description: formDescription,
                });
                showToast("Organization created successfully", "success");
            }
            showModal = false;
            await loadOrganizations();
        } catch (err) {
            const errorMsg =
                err instanceof Error
                    ? err.message
                    : "Failed to save organization";
            error = errorMsg;
            showToast(errorMsg, "error");
        }
    }

    async function handleDelete(org: Organization) {
        if (!confirm(`Are you sure you want to delete "${org.name}"?`)) return;

        error = "";
        try {
            await api.deleteOrganization(org.id);
            showToast("Organization deleted successfully", "success");
            await loadOrganizations();
        } catch (err) {
            const errorMsg =
                err instanceof Error
                    ? err.message
                    : "Failed to delete organization";
            error = errorMsg;
            showToast(errorMsg, "error");
        }
    }

    const columns = [
        {
            key: "name" as keyof Organization,
            label: "Name",
            sortable: true,
            render: (name: string, org: Organization) =>
                `<a href="/org/${org.id}/dashboard" class="text-primary-400 hover:underline font-medium">${name}</a>`,
        },
        { key: "code" as keyof Organization, label: "Code", sortable: true },
        { key: "type" as keyof Organization, label: "Type", sortable: true },
        {
            key: "description" as keyof Organization,
            label: "Description",
            render: (desc: string) =>
                desc || '<span class="text-slate-500">-</span>',
        },
        {
            key: "created_at" as keyof Organization,
            label: "Created",
            sortable: true,
            render: (value: string) => new Date(value).toLocaleDateString(),
        },
    ];

    function enterOrganization(org: Organization) {
        goto(`/org/${org.id}/dashboard`);
    }

    const actions = [
        {
            label: "Enter",
            variant: "secondary" as const,
            onClick: enterOrganization,
        },
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
    <title>Organizations - Sekolah Madrasah</title>
</svelte:head>

{#if $isAuthenticated}
    <div class="py-8 px-4">
        <div class="max-w-7xl mx-auto">
            <!-- Header -->
            <div class="flex justify-between items-center mb-8">
                <div>
                    <h1 class="text-3xl font-bold mb-2">Organizations</h1>
                    <p class="text-slate-400">
                        Manage your organizations and their settings
                    </p>
                </div>
                <button on:click={openCreateModal} class="btn-primary">
                    + Create Organization
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

            <!-- Organizations Table -->
            <div class="card">
                {#if isLoading}
                    <div class="text-center py-12 text-slate-400">
                        Loading organizations...
                    </div>
                {:else}
                    <DataTable
                        data={organizations}
                        {columns}
                        {actions}
                        emptyMessage="No organizations found. Create your first one!"
                    />
                {/if}
            </div>
        </div>
    </div>

    <!-- Create/Edit Modal -->
    <Modal
        bind:isOpen={showModal}
        title={editingOrg ? "Edit Organization" : "Create Organization"}
    >
        <form on:submit|preventDefault={handleSubmit} class="space-y-4">
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
                    placeholder="e.g. Sekolah Cendekia"
                    required
                />
            </div>

            {#if !editingOrg}
                <div>
                    <label
                        for="code"
                        class="block text-sm font-medium text-slate-300 mb-1"
                    >
                        Code *
                    </label>
                    <input
                        id="code"
                        type="text"
                        bind:value={formCode}
                        class="input-field"
                        placeholder="e.g. SCK"
                        required
                    />
                </div>

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
                        <option value="">Select Type</option>
                        <option value="school">Sekolah</option>
                        <option value="madrasah">Madrasah</option>
                        <option value="perumahan">Perumahan (Legacy)</option>
                        <option value="university">University</option>
                        <option value="training_center">Training Center</option>
                        <option value="other">Other</option>
                    </select>
                </div>
            {:else}
                <div class="grid grid-cols-2 gap-4">
                    <div>
                        <label
                            class="block text-sm font-medium text-slate-400 mb-1"
                            >Code</label
                        >
                        <div
                            class="input-field bg-slate-700/50 cursor-not-allowed"
                        >
                            {formCode}
                        </div>
                    </div>
                    <div>
                        <label
                            class="block text-sm font-medium text-slate-400 mb-1"
                            >Type</label
                        >
                        <div
                            class="input-field bg-slate-700/50 cursor-not-allowed capitalize"
                        >
                            {formType}
                        </div>
                    </div>
                </div>
                <p class="text-xs text-slate-500">
                    Code and Type cannot be changed after creation.
                </p>
            {/if}

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
                    placeholder="Brief description of the organization"
                    rows="3"
                ></textarea>
            </div>
        </form>

        <div slot="footer">
            <button on:click={() => (showModal = false)} class="btn-secondary">
                Cancel
            </button>
            <button on:click={handleSubmit} class="btn-primary">
                {editingOrg ? "Update" : "Create"}
            </button>
        </div>
    </Modal>
{/if}
