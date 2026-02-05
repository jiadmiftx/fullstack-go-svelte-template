<script lang="ts">
    import { page } from "$app/stores";
    import { api } from "$lib";
    import { onMount } from "svelte";
    import Modal from "$core/components/Modal.svelte";
    import DataTable from "$core/components/DataTable.svelte";
    import { showToast } from "$core/components/Toast.svelte";

    interface Kegiatan {
        id: string;
        unit_id: string;
        code: string;
        name: string;
        description?: string;
        credits: number;
        is_active: boolean;
        pengurusCount?: number;
        pengurusNames?: string[];
    }

    interface Pengurus {
        id: string;
        user_name: string;
        nip: string;
    }

    interface PengurusKegiatan {
        id: string;
        pengurus_id: string;
        kegiatan_id: string;
        pengurus_name: string;
    }

    interface RT {
        id: string;
        name: string;
    }

    interface Schedule {
        id: string;
        rt_id: string;
        kegiatan_id: string;
        pengurus_id: string;
    }

    $: unitId = $page.params.unitId ?? "";

    let kegiatans: Kegiatan[] = [];
    let isLoading = false;
    let error = "";
    let showModal = false;
    let editingKegiatan: Kegiatan | null = null;

    // Pengurus assignment
    let penguruss: Pengurus[] = [];
    let showPengurussModal = false;
    let selectedKegiatan: Kegiatan | null = null;
    let kegiatanPenguruss: PengurusKegiatan[] = [];
    let pengurusSearch = "";
    let showPengurusDropdown = false;

    // RT filter
    let rts: RT[] = [];
    let schedules: Schedule[] = [];
    let selectedRTFilter = "";

    // Form fields
    let formCode = "";
    let formName = "";
    let formDescription = "";
    let formCredits = 2;

    onMount(async () => {
        await Promise.all([
            loadKegiatans(),
            loadPenguruss(),
            loadRTs(),
            loadSchedules(),
        ]);
    });

    async function loadKegiatans() {
        isLoading = true;
        error = "";
        try {
            const response = await api.getKegiatans({
                unit_id: unitId,
                limit: 100,
            });
            kegiatans = response.data || [];

            // Load pengurus names for each kegiatan
            for (let i = 0; i < kegiatans.length; i++) {
                try {
                    const tsRes = await api.getPengurussByKegiatan(
                        kegiatans[i].id,
                    );
                    const pengurusList = tsRes.data || [];
                    kegiatans[i].pengurusCount = pengurusList.length;
                    kegiatans[i].pengurusNames = pengurusList.map(
                        (t: any) => t.pengurus_name || "Unknown",
                    );
                } catch {
                    kegiatans[i].pengurusCount = 0;
                    kegiatans[i].pengurusNames = [];
                }
            }
            kegiatans = kegiatans; // Trigger reactivity
        } catch (err) {
            error =
                err instanceof Error
                    ? err.message
                    : "Gagal memuat kegiatan";
        } finally {
            isLoading = false;
        }
    }

    async function loadPenguruss() {
        try {
            const response = await api.getPenguruss({
                unit_id: unitId,
                limit: 100,
            });
            penguruss = response.data || [];
            console.log("DEBUG loadPenguruss:", penguruss);
        } catch (err) {
            console.error("Failed to load penguruss:", err);
        }
    }

    async function loadRTs() {
        try {
            const response = await api.getRTs({
                unit_id: unitId,
                limit: 100,
            });
            rts = response.data || [];
        } catch (err) {
            console.error("Failed to load rts:", err);
        }
    }

    async function loadSchedules() {
        try {
            const response = await api.getCalendarEntries({
                unit_id: unitId,
                entry_type: "schedule",
                limit: 500,
            });
            // Map to expected schedule format
            schedules = (response.data || []).map((e: any) => ({
                id: e.id,
                rt_id: e.rt_id,
                kegiatan_id: e.kegiatan_id,
                pengurus_id: e.pengurus_id,
            }));
        } catch (err) {
            console.error("Failed to load schedules:", err);
        }
    }

    // Filter kegiatans by rt (based on schedules)
    $: filteredKegiatans = selectedRTFilter
        ? kegiatans.filter((s) =>
              schedules.some(
                  (sch) =>
                      sch.kegiatan_id === s.id &&
                      sch.rt_id === selectedRTFilter,
              ),
          )
        : kegiatans;

    function openCreateModal() {
        editingKegiatan = null;
        formCode = "";
        formName = "";
        formDescription = "";
        formCredits = 2;
        showModal = true;
    }

    function openEditModal(kegiatan: Kegiatan) {
        editingKegiatan = kegiatan;
        formCode = kegiatan.code;
        formName = kegiatan.name;
        formDescription = kegiatan.description || "";
        formCredits = kegiatan.credits;
        showModal = true;
    }

    async function handleSubmit() {
        error = "";
        try {
            if (editingKegiatan) {
                await api.updateKegiatan(editingKegiatan.id, {
                    code: formCode,
                    name: formName,
                    description: formDescription || undefined,
                    credits: formCredits,
                });
                showToast("Mata pelajaran berhasil diupdate", "success");
            } else {
                await api.createKegiatan({
                    unit_id: unitId,
                    code: formCode,
                    name: formName,
                    description: formDescription || undefined,
                    credits: formCredits,
                });
                showToast("Mata pelajaran berhasil ditambahkan", "success");
            }
            showModal = false;
            await loadKegiatans();
        } catch (err) {
            const errorMsg =
                err instanceof Error
                    ? err.message
                    : "Gagal menyimpan kegiatan";
            showToast(errorMsg, "error");
        }
    }

    async function handleDelete(kegiatan: Kegiatan) {
        if (!confirm(`Hapus kegiatan "${kegiatan.name}"?`)) return;

        try {
            await api.deleteKegiatan(kegiatan.id);
            showToast("Mata pelajaran berhasil dihapus", "success");
            await loadKegiatans();
        } catch (err) {
            const errorMsg =
                err instanceof Error
                    ? err.message
                    : "Gagal menghapus kegiatan";
            showToast(errorMsg, "error");
        }
    }

    async function openPengurussModal(kegiatan: Kegiatan) {
        selectedKegiatan = kegiatan;
        try {
            const response = await api.getPengurussByKegiatan(kegiatan.id);
            kegiatanPenguruss = response.data || [];
        } catch {
            kegiatanPenguruss = [];
        }
        showPengurussModal = true;
    }

    function isPengurusAssigned(pengurusId: string): boolean {
        return kegiatanPenguruss.some((ts) => ts.pengurus_id === pengurusId);
    }

    async function togglePengurus(pengurus: Pengurus) {
        if (!selectedKegiatan) return;

        try {
            if (isPengurusAssigned(pengurus.id)) {
                await api.removeKegiatanFromPengurus(
                    pengurus.id,
                    selectedKegiatan.id,
                );
                kegiatanPenguruss = kegiatanPenguruss.filter(
                    (ts) => ts.pengurus_id !== pengurus.id,
                );
                showToast(
                    `${pengurus.user_name} dihapus dari ${selectedKegiatan.name}`,
                    "success",
                );
            } else {
                const response = await api.assignKegiatanToPengurus(pengurus.id, {
                    kegiatan_id: selectedKegiatan.id,
                });
                kegiatanPenguruss = [...kegiatanPenguruss, response.data];
                showToast(
                    `${pengurus.user_name} ditambahkan ke ${selectedKegiatan.name}`,
                    "success",
                );
            }
            // Refresh kegiatan pengurus count in table
            const idx = kegiatans.findIndex((s) => s.id === selectedKegiatan!.id);
            if (idx >= 0) {
                kegiatans[idx].pengurusCount = kegiatanPenguruss.length;
                kegiatans = kegiatans;
            }
        } catch (err) {
            showToast("Gagal mengubah pengurus", "error");
        }
    }

    // Filtered penguruss for dropdown (exclude already assigned)
    $: filteredPenguruss = penguruss.filter(
        (t) =>
            !isPengurusAssigned(t.id) &&
            ((t.user_name || "")
                .toLowerCase()
                .includes(pengurusSearch.toLowerCase()) ||
                (t.nip || "").includes(pengurusSearch)),
    );

    async function addPengurus(pengurus: Pengurus) {
        if (!selectedKegiatan) return;
        // Prevent duplicate assignment
        if (isPengurusAssigned(pengurus.id)) {
            showToast(
                `${pengurus.user_name} sudah ada di ${selectedKegiatan.name}`,
                "info",
            );
            return;
        }
        try {
            const response = await api.assignKegiatanToPengurus(pengurus.id, {
                kegiatan_id: selectedKegiatan.id,
            });
            kegiatanPenguruss = [...kegiatanPenguruss, response.data];
            showToast(
                `${pengurus.user_name} ditambahkan ke ${selectedKegiatan.name}`,
                "success",
            );
            pengurusSearch = "";
            showPengurusDropdown = false;
            // Refresh kegiatan pengurus count and names
            const idx = kegiatans.findIndex((s) => s.id === selectedKegiatan!.id);
            if (idx >= 0) {
                kegiatans[idx].pengurusCount = kegiatanPenguruss.length;
                kegiatans[idx].pengurusNames = [
                    ...(kegiatans[idx].pengurusNames || []),
                    pengurus.user_name || "Unknown",
                ];
                kegiatans = kegiatans;
            }
        } catch (err) {
            showToast("Gagal menambahkan pengurus", "error");
        }
    }

    async function removePengurus(pengurus: Pengurus) {
        if (!selectedKegiatan) return;
        try {
            await api.removeKegiatanFromPengurus(pengurus.id, selectedKegiatan.id);
            kegiatanPenguruss = kegiatanPenguruss.filter(
                (ts) => ts.pengurus_id !== pengurus.id,
            );
            showToast(
                `${pengurus.user_name} dihapus dari ${selectedKegiatan.name}`,
                "success",
            );
            // Refresh kegiatan pengurus count and names
            const idx = kegiatans.findIndex((s) => s.id === selectedKegiatan!.id);
            if (idx >= 0) {
                kegiatans[idx].pengurusCount = kegiatanPenguruss.length;
                kegiatans[idx].pengurusNames = (
                    kegiatans[idx].pengurusNames || []
                ).filter((n) => n !== (pengurus.user_name || "Unknown"));
                kegiatans = kegiatans;
            }
        } catch (err) {
            showToast("Gagal menghapus pengurus", "error");
        }
    }

    const columns = [
        { key: "code" as keyof Kegiatan, label: "Kode" },
        { key: "name" as keyof Kegiatan, label: "Nama Proyek" },
        {
            key: "pengurusNames" as keyof Kegiatan,
            label: "Pengurus",
            render: (names: string[] | undefined, row: Kegiatan) => {
                const pengurusNames = names || [];
                if (pengurusNames.length === 0) {
                    return '<span class="text-gray-400">-</span>';
                }
                return pengurusNames
                    .map(
                        (name) =>
                            `<span style="display: inline-block; background-color: #10b981; color: white; padding: 2px 8px; border-radius: 12px; font-size: 0.75rem; margin: 2px;">${name}</span>`,
                    )
                    .join(" ");
            },
        },
        {
            key: "credits" as keyof Kegiatan,
            label: "Jam/Minggu",
            render: (credits: number) => `${credits} JP`,
        },
        {
            key: "is_active" as keyof Kegiatan,
            label: "Status",
            render: (active: boolean) =>
                active
                    ? '<span class="bg-green-500/20 text-green-600 px-2 py-1 rounded text-xs">Aktif</span>'
                    : '<span class="bg-red-500/20 text-red-600 px-2 py-1 rounded text-xs">Nonaktif</span>',
        },
    ];

    const actions = [
        {
            label: "Pengurus",
            variant: "secondary" as const,
            onClick: openPengurussModal,
        },
        {
            label: "Edit",
            variant: "primary" as const,
            onClick: openEditModal,
        },
        { label: "Hapus", variant: "danger" as const, onClick: handleDelete },
    ];
</script>

<svelte:head>
    <title>Mata Pelajaran - SeekOlah</title>
</svelte:head>

<div class="py-8 px-4">
    <div class="max-w-7xl mx-auto">
        <!-- Header -->
        <div class="flex justify-between items-center mb-8">
            <div>
                <h1 class="text-3xl font-bold mb-2 text-gray-900">
                    📚 Mata Pelajaran
                </h1>
                <p class="text-gray-600">
                    Kelola daftar kegiatan RT
                </p>
            </div>
            <button on:click={openCreateModal} class="btn-primary">
                + Tambah Proyek
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

        <!-- RT Filter -->
        <div class="filter-bar mb-4">
            <label class="filter-label">Filter berdasarkan RT:</label>
            <select bind:value={selectedRTFilter} class="filter-select">
                <option value="">Semua RT</option>
                {#each rts as cls}
                    <option value={cls.id}>{cls.name}</option>
                {/each}
            </select>
            {#if selectedRTFilter && filteredKegiatans.length !== kegiatans.length}
                <span class="filter-count">
                    Menampilkan {filteredKegiatans.length} dari {kegiatans.length}
                    proyek
                </span>
            {/if}
        </div>

        <!-- Kegiatans Table -->
        <div
            class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden"
        >
            {#if isLoading}
                <div class="text-center py-12 text-gray-600">
                    Memuat data...
                </div>
            {:else}
                <DataTable
                    data={filteredKegiatans}
                    {columns}
                    {actions}
                    emptyMessage={selectedRTFilter
                        ? "Tidak ada proyek yang diajarkan di tim ini"
                        : "Belum ada kegiatan. Tambahkan proyek pertama!"}
                />
            {/if}
        </div>
    </div>
</div>

<!-- Create/Edit Modal -->
<Modal
    bind:isOpen={showModal}
    title={editingKegiatan ? "Edit Mata Pelajaran" : "Tambah Mata Pelajaran"}
    size="md"
>
    <form on:submit|preventDefault={handleSubmit} class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
            <div>
                <label
                    for="code"
                    class="block text-sm font-medium text-gray-700 mb-1"
                    >Kode *</label
                >
                <input
                    id="code"
                    type="text"
                    bind:value={formCode}
                    class="input-field"
                    placeholder="MTK"
                    required
                    maxlength="20"
                />
            </div>
            <div>
                <label
                    for="credits"
                    class="block text-sm font-medium text-gray-700 mb-1"
                    >Jam/Minggu *</label
                >
                <input
                    id="credits"
                    type="number"
                    bind:value={formCredits}
                    class="input-field"
                    required
                    min="1"
                    max="20"
                />
            </div>
        </div>

        <div>
            <label
                for="name"
                class="block text-sm font-medium text-gray-700 mb-1"
                >Nama Proyek *</label
            >
            <input
                id="name"
                type="text"
                bind:value={formName}
                class="input-field"
                placeholder="Matematika"
                required
            />
        </div>

        <div>
            <label
                for="description"
                class="block text-sm font-medium text-gray-700 mb-1"
                >Deskripsi</label
            >
            <textarea
                id="description"
                bind:value={formDescription}
                class="input-field"
                rows="3"
                placeholder="Deskripsi kegiatan..."
            ></textarea>
        </div>
    </form>

    <div slot="footer">
        <button on:click={() => (showModal = false)} class="btn-secondary"
            >Batal</button
        >
        <button on:click={handleSubmit} class="btn-primary">
            {editingKegiatan ? "Update" : "Simpan"}
        </button>
    </div>
</Modal>

<!-- Pengurus Assignment Modal -->
<Modal
    bind:isOpen={showPengurussModal}
    title={selectedKegiatan ? `Pengurus - ${selectedKegiatan.name}` : "Assign Pengurus"}
    size="lg"
>
    <div class="penguruss-modal">
        <!-- Assigned Penguruss (Chips) -->
        <div class="assigned-section">
            <span class="section-label"
                >Pengurus yang Mengajar ({kegiatanPenguruss.length}):</span
            >
            {#if kegiatanPenguruss.length === 0}
                <p class="empty-text">Belum ada pengurus yang ditugaskan</p>
            {:else}
                <div class="assigned-chips">
                    {#each kegiatanPenguruss as ts}
                        {@const pengurus = penguruss.find(
                            (t) => t.id === ts.pengurus_id,
                        )}
                        {#if pengurus}
                            <span class="assigned-chip">
                                <span class="chip-nip"
                                    >{pengurus.nip || "-"}</span
                                >
                                <span class="chip-name"
                                    >{pengurus.user_name || "Tanpa Nama"}</span
                                >
                                <button
                                    type="button"
                                    class="chip-remove"
                                    on:click={() => removePengurus(pengurus)}
                                    title="Hapus pengurus">×</button
                                >
                            </span>
                        {/if}
                    {/each}
                </div>
            {/if}
        </div>

        <!-- Add Pengurus Section -->
        <div class="add-section">
            <span class="section-label">Tambah Pengurus:</span>
            <div class="search-container">
                <input
                    type="text"
                    bind:value={pengurusSearch}
                    class="search-input"
                    placeholder="Cari berdasarkan nama atau NIP..."
                    on:focus={() => (showPengurusDropdown = true)}
                />
                {#if showPengurusDropdown && filteredPenguruss.length > 0}
                    <div class="pengurus-dropdown">
                        {#each filteredPenguruss as pengurus}
                            <button
                                type="button"
                                class="dropdown-item"
                                on:click={() => addPengurus(pengurus)}
                            >
                                <span class="item-nip"
                                    >{pengurus.nip || "-"}</span
                                >
                                <span class="item-name"
                                    >{pengurus.user_name || "Tanpa Nama"}</span
                                >
                            </button>
                        {/each}
                    </div>
                {/if}
            </div>
        </div>
    </div>
    <div slot="footer">
        <button
            on:click={() => {
                showPengurussModal = false;
                showPengurusDropdown = false;
                pengurusSearch = "";
            }}
            class="btn-primary"
        >
            Selesai
        </button>
    </div>
</Modal>

<style>
    h1 {
        color: #1e293b;
    }
    .card {
        background: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 1rem;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    }
    .input-field {
        width: 100%;
        padding: 0.625rem 0.875rem;
        background: #ffffff;
        border: 1px solid #d1d5db;
        border-radius: 0.5rem;
        color: #1f2937;
        font-size: 0.875rem;
        transition: all 0.2s;
    }
    .input-field:focus {
        outline: none;
        border-color: #7c3aed;
        box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.15);
    }
    .btn-primary {
        padding: 0.625rem 1.25rem;
        background: linear-gradient(135deg, #7c3aed 0%, #6d28d9 100%);
        color: white;
        font-weight: 600;
        border: none;
        border-radius: 0.5rem;
        cursor: pointer;
        transition: all 0.2s;
    }
    .btn-primary:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(124, 58, 237, 0.4);
    }
    .btn-secondary {
        padding: 0.625rem 1.25rem;
        background: #f1f5f9;
        color: #475569;
        font-weight: 500;
        border: none;
        border-radius: 0.5rem;
        cursor: pointer;
        transition: all 0.2s;
    }
    .btn-secondary:hover {
        background: #e2e8f0;
    }
    /* Pengurus Assignment Modal Styles */
    .penguruss-modal {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
    }
    .assigned-section,
    .add-section {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }
    .section-label {
        font-weight: 600;
        color: #374151;
        font-size: 0.875rem;
    }
    .empty-text {
        color: #9ca3af;
        font-size: 0.875rem;
        font-style: italic;
    }
    .assigned-chips {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
    }
    .assigned-chip {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.5rem 0.75rem;
        background: #dcfce7;
        border: 1px solid #86efac;
        border-radius: 0.5rem;
        color: #166534;
    }
    .chip-nip {
        font-weight: 600;
        font-size: 0.7rem;
        background: #bbf7d0;
        padding: 0.125rem 0.375rem;
        border-radius: 0.25rem;
    }
    .chip-name {
        font-size: 0.875rem;
    }
    .chip-remove {
        background: none;
        border: none;
        color: #dc2626;
        font-size: 1.25rem;
        font-weight: bold;
        cursor: pointer;
        padding: 0;
        line-height: 1;
        margin-left: 0.25rem;
    }
    .chip-remove:hover {
        color: #b91c1c;
    }
    .search-container {
        position: relative;
    }
    .search-input {
        width: 100%;
        padding: 0.625rem 0.875rem;
        background: #ffffff;
        border: 1px solid #d1d5db;
        border-radius: 0.5rem;
        color: #1f2937;
        font-size: 0.875rem;
    }
    .search-input:focus {
        outline: none;
        border-color: #7c3aed;
        box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.15);
    }
    .pengurus-dropdown {
        position: absolute;
        top: 100%;
        left: 0;
        right: 0;
        background: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 0.5rem;
        box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1);
        max-height: 200px;
        overflow-y: auto;
        z-index: 50;
        margin-top: 0.25rem;
    }
    .dropdown-item {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        width: 100%;
        padding: 0.625rem 0.875rem;
        background: none;
        border: none;
        text-align: left;
        cursor: pointer;
        transition: background 0.15s;
    }
    .dropdown-item:hover {
        background: #f3e8ff;
    }
    .item-nip {
        font-weight: 600;
        font-size: 0.75rem;
        background: #f1f5f9;
        padding: 0.125rem 0.5rem;
        border-radius: 0.25rem;
        color: #475569;
        min-width: 70px;
    }
    .item-name {
        color: #1e293b;
        font-size: 0.875rem;
    }
    .filter-bar {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        background: #ffffff;
        padding: 0.75rem 1rem;
        border-radius: 0.75rem;
        border: 1px solid #e2e8f0;
    }
    .filter-label {
        font-size: 0.875rem;
        color: #475569;
        white-space: nowrap;
    }
    .filter-select {
        padding: 0.5rem 1rem;
        background: #f8fafc;
        border: 1px solid #d1d5db;
        border-radius: 0.5rem;
        color: #1f2937;
        font-size: 0.875rem;
        min-width: 180px;
    }
    .filter-select:focus {
        outline: none;
        border-color: #7c3aed;
        box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.15);
    }
    .filter-count {
        font-size: 0.75rem;
        color: #7c3aed;
        background: #ede9fe;
        padding: 0.25rem 0.75rem;
        border-radius: 1rem;
    }
</style>
