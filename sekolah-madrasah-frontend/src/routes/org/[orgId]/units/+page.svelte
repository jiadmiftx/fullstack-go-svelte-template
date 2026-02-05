<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { page } from "$app/stores";
    import { isAuthenticated, api, selectedOrganization } from "$lib";
    import {
        canCreatePerumahan,
        canUpdatePerumahan,
        canDeletePerumahan,
    } from "$core/stores/permissions";
    import Modal from "$core/components/Modal.svelte";
    import DataTable from "$core/components/DataTable.svelte";
    import { showToast } from "$core/components/Toast.svelte";

    interface Perumahan {
        id: string;
        organization_id: string;
        name: string;
        code: string;
        address?: string;
        phone?: string;
        email?: string;
        logo?: string;
        is_active: boolean;
    }

    // Get organization ID from URL
    $: orgId = $page.params.orgId || $selectedOrganization?.id || "";

    let perumahans: Perumahan[] = [];
    let isLoading = false;
    let error = "";
    let showModal = false;
    let editingPerumahan: Perumahan | null = null;

    // Form fields
    let formName = "";
    let formCode = "";
    let formAddress = "";
    let formPhone = "";
    let formEmail = "";

    $: if (!$isAuthenticated && typeof window !== "undefined") {
        goto("/auth/login");
    }

    // Reload when orgId changes
    $: if (orgId && $isAuthenticated) {
        loadPerumahans();
    }

    onMount(async () => {
        if ($isAuthenticated && orgId) {
            await loadPerumahans();
        }
    });

    async function loadPerumahans() {
        if (!orgId) return;

        isLoading = true;
        error = "";
        try {
            const response = await api.getPerumahans({
                limit: 100,
                organization_id: orgId,
            });
            perumahans = response.data || [];
        } catch (err) {
            error = err instanceof Error ? err.message : "Gagal memuat data RT";
        } finally {
            isLoading = false;
        }
    }

    function openCreateModal() {
        editingPerumahan = null;
        formName = "";
        formCode = "";
        formAddress = "";
        formPhone = "";
        formEmail = "";
        showModal = true;
    }

    function openEditModal(perumahan: Perumahan) {
        editingPerumahan = perumahan;
        formName = perumahan.name;
        formCode = perumahan.code;
        formAddress = perumahan.address || "";
        formPhone = perumahan.phone || "";
        formEmail = perumahan.email || "";
        showModal = true;
    }

    async function handleSubmit() {
        error = "";
        try {
            if (editingPerumahan) {
                await api.updatePerumahan(editingPerumahan.id, {
                    name: formName,
                    address: formAddress,
                    phone: formPhone,
                    email: formEmail,
                });
                showToast("RT berhasil diupdate", "success");
            } else {
                await api.createPerumahan({
                    organization_id: orgId,
                    name: formName,
                    code: formCode,
                    address: formAddress,
                    phone: formPhone,
                    email: formEmail,
                });
                showToast("RT berhasil ditambahkan", "success");
            }
            showModal = false;
            await loadPerumahans();
        } catch (err) {
            const errorMsg =
                err instanceof Error ? err.message : "Gagal menyimpan RT";
            error = errorMsg;
            showToast(errorMsg, "error");
        }
    }

    async function handleDelete(perumahan: Perumahan) {
        if (!confirm(`Yakin ingin menghapus "${perumahan.name}"?`)) return;

        error = "";
        try {
            await api.deletePerumahan(perumahan.id);
            showToast("RT berhasil dihapus", "success");
            await loadPerumahans();
        } catch (err) {
            const errorMsg =
                err instanceof Error ? err.message : "Gagal menghapus RT";
            error = errorMsg;
            showToast(errorMsg, "error");
        }
    }

    const columns = [
        { key: "name" as keyof Perumahan, label: "Nama RT", sortable: true },
        { key: "code" as keyof Perumahan, label: "Kode", sortable: true },
        {
            key: "address" as keyof Perumahan,
            label: "Alamat",
            render: (v: string) => v || '<span class="text-slate-500">-</span>',
        },
        {
            key: "phone" as keyof Perumahan,
            label: "Telepon",
            render: (v: string) => v || '<span class="text-slate-500">-</span>',
        },
        {
            key: "is_active" as keyof Perumahan,
            label: "Status",
            render: (active: boolean) =>
                active
                    ? '<span class="bg-green-500/20 text-green-400 px-2 py-1 rounded text-xs font-medium">Aktif</span>'
                    : '<span class="bg-red-500/20 text-red-400 px-2 py-1 rounded text-xs font-medium">Nonaktif</span>',
        },
    ];

    function enterPerumahan(perumahan: Perumahan) {
        goto(`/org/${orgId}/units/${perumahan.id}`);
    }

    // Dynamic actions based on permissions
    $: actions = [
        {
            label: "Masuk",
            variant: "secondary" as const,
            onClick: enterPerumahan,
        },
        ...($canUpdatePerumahan
            ? [
                  {
                      label: "Edit",
                      variant: "primary" as const,
                      onClick: openEditModal,
                  },
              ]
            : []),
        ...($canDeletePerumahan
            ? [
                  {
                      label: "Hapus",
                      variant: "danger" as const,
                      onClick: handleDelete,
                  },
              ]
            : []),
    ];
</script>

<svelte:head>
    <title>Manajemen RT - {$selectedOrganization?.name || "Perumahan"}</title>
</svelte:head>

{#if $isAuthenticated}
    <div class="py-8 px-4">
        <div class="max-w-7xl mx-auto">
            <!-- Header -->
            <div class="flex justify-between items-center mb-8">
                <div>
                    <h1 class="text-3xl font-bold mb-2">🏫 Manajemen RT</h1>
                    <p class="text-slate-400">
                        Kelola unit RT di bawah {$selectedOrganization?.name ||
                            "organisasi ini"}
                    </p>
                </div>
                {#if $canCreatePerumahan}
                    <button on:click={openCreateModal} class="btn-primary">
                        + Tambah RT
                    </button>
                {/if}
            </div>

            <!-- Error Alert -->
            {#if error}
                <div
                    class="bg-red-500/20 border border-red-500 text-red-400 px-4 py-3 rounded-lg mb-6"
                >
                    {error}
                </div>
            {/if}

            <!-- Stats Cards -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
                <div class="card p-4">
                    <div class="text-2xl font-bold text-blue-400">
                        {perumahans.length}
                    </div>
                    <div class="text-slate-400 text-sm">Total RT</div>
                </div>
                <div class="card p-4">
                    <div class="text-2xl font-bold text-green-400">
                        {perumahans.filter((s) => s.is_active).length}
                    </div>
                    <div class="text-slate-400 text-sm">Aktif</div>
                </div>
            </div>

            <!-- Perumahans Table -->
            <div class="card">
                {#if isLoading}
                    <div class="text-center py-12 text-slate-400">
                        <div class="animate-pulse">Memuat data RT...</div>
                    </div>
                {:else}
                    <DataTable
                        data={perumahans}
                        {columns}
                        {actions}
                        emptyMessage="Belum ada RT terdaftar. Klik 'Tambah RT' untuk menambahkan yang pertama!"
                    />
                {/if}
            </div>
        </div>
    </div>

    <!-- Create/Edit Modal -->
    <Modal
        bind:isOpen={showModal}
        title={editingPerumahan ? "Edit RT" : "Tambah RT Baru"}
        size="lg"
    >
        <form on:submit|preventDefault={handleSubmit} class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label
                        for="name"
                        class="block text-sm font-medium text-slate-300 mb-1"
                    >
                        Nama RT *
                    </label>
                    <input
                        id="name"
                        type="text"
                        bind:value={formName}
                        class="input-field"
                        placeholder="e.g. SD Negeri 1 Jakarta"
                        required
                    />
                </div>

                <div>
                    <label
                        for="code"
                        class="block text-sm font-medium text-slate-300 mb-1"
                    >
                        Kode/NPSN *
                    </label>
                    <input
                        id="code"
                        type="text"
                        bind:value={formCode}
                        class="input-field"
                        placeholder="e.g. 20100001"
                        disabled={!!editingPerumahan}
                        required={!editingPerumahan}
                    />
                    {#if editingPerumahan}
                        <p class="text-xs text-slate-500 mt-1">
                            Kode tidak dapat diubah
                        </p>
                    {/if}
                </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label
                        for="phone"
                        class="block text-sm font-medium text-slate-300 mb-1"
                    >
                        Telepon
                    </label>
                    <input
                        id="phone"
                        type="text"
                        bind:value={formPhone}
                        class="input-field"
                        placeholder="e.g. 021-1234567"
                    />
                </div>
            </div>

            <div>
                <label
                    for="email"
                    class="block text-sm font-medium text-slate-300 mb-1"
                >
                    Email
                </label>
                <input
                    id="email"
                    type="email"
                    bind:value={formEmail}
                    class="input-field"
                    placeholder="e.g. admin@RT.sch.id"
                />
            </div>

            <div>
                <label
                    for="address"
                    class="block text-sm font-medium text-slate-300 mb-1"
                >
                    Alamat
                </label>
                <textarea
                    id="address"
                    bind:value={formAddress}
                    class="input-field"
                    placeholder="Alamat lengkap RT"
                    rows="2"
                ></textarea>
            </div>
        </form>

        <div slot="footer">
            <button on:click={() => (showModal = false)} class="btn-secondary">
                Batal
            </button>
            <button on:click={handleSubmit} class="btn-primary">
                {editingPerumahan ? "Simpan Perubahan" : "Tambah RT"}
            </button>
        </div>
    </Modal>
{/if}
