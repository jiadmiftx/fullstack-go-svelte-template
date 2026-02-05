<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import { api } from "$lib";

    $: orgId = $page.params.orgId;
    $: unitId = $page.params.unitId;

    interface AnggotaJabatan {
        jabatan_warga_id: string;
        warga_profile_id: string;
        nama_warga: string;
        is_ketua: boolean;
    }

    interface JabatanStruktur {
        id: string;
        nama: string;
        tipe: string;
        is_unique: boolean;
        urutan: number;
        anggota: AnggotaJabatan[];
    }

    interface Warga {
        id: string;
        user_name?: string;
        blok_rumah?: string;
        nomor_rumah?: string;
    }

    let struktur: JabatanStruktur[] = [];
    let wargas: Warga[] = [];
    let isLoading = true;
    let error = "";
    let successMessage = "";

    // Modal state
    let showAssignModal = false;
    let selectedJabatan: JabatanStruktur | null = null;
    let selectedWargaId = "";
    let isKetua = false;
    let isAssigning = false;

    onMount(async () => {
        await Promise.all([loadStruktur(), loadWargas()]);
    });

    async function loadStruktur() {
        isLoading = true;
        error = "";
        try {
            const res = await api.get(`/units/${unitId}/struktur`);
            struktur = res.data || [];
            if (struktur.length === 0) {
                // Initialize default jabatans if empty
                await initializeJabatans();
            }
        } catch (err) {
            error =
                err instanceof Error ? err.message : "Gagal memuat struktur RT";
        } finally {
            isLoading = false;
        }
    }

    async function loadWargas() {
        try {
            const res = await api.get(
                `/warga-profiles?unit_id=${unitId}&limit=200`,
            );
            wargas = res.data || [];
        } catch (err) {
            console.error("Failed to load wargas:", err);
        }
    }

    async function initializeJabatans() {
        try {
            await api.post(`/units/${unitId}/jabatans/initialize`);
            await loadStruktur();
        } catch (err) {
            console.error("Failed to initialize jabatans:", err);
        }
    }

    function openAssignModal(jabatan: JabatanStruktur) {
        selectedJabatan = jabatan;
        selectedWargaId = "";
        isKetua = jabatan.tipe === "seksi"; // Default to ketua for seksi
        showAssignModal = true;
    }

    async function handleAssign() {
        if (!selectedJabatan || !selectedWargaId) return;
        isAssigning = true;
        error = "";
        try {
            await api.post(`/jabatans/${selectedJabatan.id}/assign`, {
                warga_profile_id: selectedWargaId,
                is_ketua: isKetua,
            });
            successMessage = "Warga berhasil ditambahkan ke jabatan";
            showAssignModal = false;
            await loadStruktur();
            setTimeout(() => (successMessage = ""), 3000);
        } catch (err) {
            error =
                err instanceof Error ? err.message : "Gagal menambahkan warga";
        } finally {
            isAssigning = false;
        }
    }

    async function handleRemove(jabatanWargaId: string) {
        if (!confirm("Yakin ingin menghapus warga dari jabatan ini?")) return;
        error = "";
        try {
            await api.delete(`/jabatan-warga/${jabatanWargaId}`);
            successMessage = "Warga berhasil dihapus dari jabatan";
            await loadStruktur();
            setTimeout(() => (successMessage = ""), 3000);
        } catch (err) {
            error =
                err instanceof Error ? err.message : "Gagal menghapus warga";
        }
    }

    function getAvailableWargas(jabatan: JabatanStruktur) {
        const assignedIds = jabatan.anggota.map((a) => a.warga_profile_id);
        return wargas.filter((w) => !assignedIds.includes(w.id));
    }

    function getJabatanIcon(tipe: string, nama: string): string {
        if (nama.includes("Ketua")) return "👤";
        if (nama.includes("Wakil")) return "👥";
        if (nama.includes("Bendahara")) return "💰";
        if (nama.includes("Sekretaris")) return "📝";
        if (nama.includes("Keamanan")) return "🔒";
        if (nama.includes("Kebersihan")) return "🧹";
        return tipe === "seksi" ? "📋" : "👔";
    }
</script>

<svelte:head>
    <title>Struktur RT</title>
</svelte:head>

<div class="page-container">
    <div class="page-header">
        <h1>🏛️ Struktur Pengurus RT</h1>
        <p class="subtitle">Kelola jabatan dan kepengurusan RT</p>
    </div>

    {#if error}
        <div class="alert alert-error">{error}</div>
    {/if}
    {#if successMessage}
        <div class="alert alert-success">{successMessage}</div>
    {/if}

    {#if isLoading}
        <div class="loading">
            <div class="spinner"></div>
            <p>Memuat struktur...</p>
        </div>
    {:else}
        <div class="struktur-grid">
            {#each struktur as jabatan}
                <div class="jabatan-card {jabatan.tipe}">
                    <div class="jabatan-header">
                        <span class="jabatan-icon"
                            >{getJabatanIcon(jabatan.tipe, jabatan.nama)}</span
                        >
                        <h3>{jabatan.nama}</h3>
                        <span class="jabatan-badge {jabatan.tipe}"
                            >{jabatan.tipe}</span
                        >
                    </div>

                    <div class="anggota-list">
                        {#if jabatan.anggota.length === 0}
                            <p class="empty-state">Belum ada yang menjabat</p>
                        {:else}
                            {#each jabatan.anggota as anggota}
                                <div class="anggota-item">
                                    <div class="anggota-info">
                                        {#if jabatan.tipe === "seksi" && anggota.is_ketua}
                                            <span class="ketua-badge"
                                                >Ketua</span
                                            >
                                        {/if}
                                        <span class="anggota-name"
                                            >{anggota.nama_warga ||
                                                "Nama tidak tersedia"}</span
                                        >
                                    </div>
                                    <button
                                        class="btn-remove"
                                        on:click={() =>
                                            handleRemove(
                                                anggota.jabatan_warga_id,
                                            )}
                                        title="Hapus dari jabatan"
                                    >
                                        ✕
                                    </button>
                                </div>
                            {/each}
                        {/if}
                    </div>

                    {#if !jabatan.is_unique || jabatan.anggota.length === 0}
                        <button
                            class="btn-add"
                            on:click={() => openAssignModal(jabatan)}
                        >
                            + Tambah {jabatan.tipe === "seksi"
                                ? "Anggota"
                                : "Pejabat"}
                        </button>
                    {/if}
                </div>
            {/each}
        </div>
    {/if}
</div>

<!-- Assign Modal -->
{#if showAssignModal && selectedJabatan}
    <div class="modal-overlay" on:click={() => (showAssignModal = false)}>
        <div class="modal" on:click|stopPropagation>
            <div class="modal-header">
                <h3>Tambah ke {selectedJabatan.nama}</h3>
                <button
                    class="btn-close"
                    on:click={() => (showAssignModal = false)}>✕</button
                >
            </div>
            <div class="modal-body">
                <div class="form-group">
                    <label for="warga">Pilih Warga</label>
                    <select
                        id="warga"
                        bind:value={selectedWargaId}
                        class="form-select"
                    >
                        <option value="">-- Pilih Warga --</option>
                        {#each getAvailableWargas(selectedJabatan) as warga}
                            <option value={warga.id}>
                                {warga.user_name || "Tanpa Nama"}
                                {#if warga.blok_rumah}
                                    (Blok {warga.blok_rumah}/{warga.nomor_rumah})
                                {/if}
                            </option>
                        {/each}
                    </select>
                </div>
                {#if selectedJabatan.tipe === "seksi"}
                    <div class="form-group">
                        <label class="checkbox-label">
                            <input type="checkbox" bind:checked={isKetua} />
                            <span>Jadikan Ketua Seksi</span>
                        </label>
                    </div>
                {/if}
            </div>
            <div class="modal-footer">
                <button
                    class="btn-secondary"
                    on:click={() => (showAssignModal = false)}>Batal</button
                >
                <button
                    class="btn-primary"
                    on:click={handleAssign}
                    disabled={!selectedWargaId || isAssigning}
                >
                    {isAssigning ? "Menyimpan..." : "Simpan"}
                </button>
            </div>
        </div>
    </div>
{/if}

<style>
    .page-container {
        padding: 1.5rem;
        max-width: 1200px;
        margin: 0 auto;
    }
    .page-header {
        margin-bottom: 2rem;
    }
    .page-header h1 {
        font-size: 1.75rem;
        font-weight: 700;
        color: #1e293b;
        margin: 0;
    }
    .subtitle {
        color: #64748b;
        margin: 0.25rem 0 0;
    }

    .alert {
        padding: 0.75rem 1rem;
        border-radius: 0.5rem;
        margin-bottom: 1rem;
    }
    .alert-error {
        background: #fef2f2;
        border: 1px solid #fecaca;
        color: #dc2626;
    }
    .alert-success {
        background: #f0fdf4;
        border: 1px solid #bbf7d0;
        color: #16a34a;
    }

    .loading {
        text-align: center;
        padding: 3rem;
        color: #64748b;
    }
    .spinner {
        width: 40px;
        height: 40px;
        border: 3px solid #e2e8f0;
        border-top-color: #7c3aed;
        border-radius: 50%;
        animation: spin 1s linear infinite;
        margin: 0 auto 1rem;
    }
    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }

    .struktur-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 1.25rem;
    }

    .jabatan-card {
        background: #fff;
        border-radius: 0.75rem;
        padding: 1.25rem;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        border: 1px solid #e2e8f0;
    }
    .jabatan-card.seksi {
        border-left: 4px solid #8b5cf6;
    }
    .jabatan-card.struktural {
        border-left: 4px solid #0ea5e9;
    }

    .jabatan-header {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin-bottom: 1rem;
        flex-wrap: wrap;
    }
    .jabatan-icon {
        font-size: 1.5rem;
    }
    .jabatan-header h3 {
        font-size: 1rem;
        font-weight: 600;
        color: #1e293b;
        margin: 0;
        flex: 1;
    }
    .jabatan-badge {
        font-size: 0.65rem;
        padding: 0.25rem 0.5rem;
        border-radius: 1rem;
        text-transform: uppercase;
        font-weight: 600;
    }
    .jabatan-badge.struktural {
        background: #e0f2fe;
        color: #0369a1;
    }
    .jabatan-badge.seksi {
        background: #ede9fe;
        color: #6d28d9;
    }

    .anggota-list {
        min-height: 60px;
    }
    .empty-state {
        color: #94a3b8;
        font-size: 0.875rem;
        font-style: italic;
        text-align: center;
        padding: 1rem 0;
    }

    .anggota-item {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0.5rem 0.75rem;
        background: #f8fafc;
        border-radius: 0.5rem;
        margin-bottom: 0.5rem;
    }
    .anggota-info {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }
    .anggota-name {
        font-size: 0.875rem;
        font-weight: 500;
        color: #334155;
    }
    .ketua-badge {
        font-size: 0.65rem;
        background: #fef3c7;
        color: #92400e;
        padding: 0.125rem 0.375rem;
        border-radius: 0.25rem;
        font-weight: 600;
    }

    .btn-remove {
        background: none;
        border: none;
        color: #94a3b8;
        cursor: pointer;
        padding: 0.25rem;
        font-size: 0.875rem;
    }
    .btn-remove:hover {
        color: #ef4444;
    }

    .btn-add {
        width: 100%;
        padding: 0.625rem;
        border: 2px dashed #cbd5e1;
        border-radius: 0.5rem;
        background: transparent;
        color: #64748b;
        font-size: 0.875rem;
        cursor: pointer;
        margin-top: 0.75rem;
        transition: all 0.2s;
    }
    .btn-add:hover {
        border-color: #7c3aed;
        color: #7c3aed;
        background: #f5f3ff;
    }

    /* Modal */
    .modal-overlay {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
        padding: 1rem;
    }
    .modal {
        background: #fff;
        border-radius: 0.75rem;
        width: 100%;
        max-width: 400px;
        box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
    }
    .modal-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 1rem 1.25rem;
        border-bottom: 1px solid #e2e8f0;
    }
    .modal-header h3 {
        margin: 0;
        font-size: 1rem;
        font-weight: 600;
    }
    .btn-close {
        background: none;
        border: none;
        font-size: 1.25rem;
        cursor: pointer;
        color: #64748b;
    }
    .modal-body {
        padding: 1.25rem;
    }
    .modal-footer {
        display: flex;
        gap: 0.75rem;
        justify-content: flex-end;
        padding: 1rem 1.25rem;
        border-top: 1px solid #e2e8f0;
    }

    .form-group {
        margin-bottom: 1rem;
    }
    .form-group label {
        display: block;
        font-size: 0.875rem;
        font-weight: 500;
        color: #374151;
        margin-bottom: 0.375rem;
    }
    .form-select {
        width: 100%;
        padding: 0.625rem;
        border: 1px solid #d1d5db;
        border-radius: 0.5rem;
        font-size: 0.875rem;
    }
    .form-select:focus {
        outline: none;
        border-color: #7c3aed;
        box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.15);
    }

    .checkbox-label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        cursor: pointer;
    }
    .checkbox-label input {
        width: 1rem;
        height: 1rem;
        accent-color: #7c3aed;
    }

    .btn-secondary {
        padding: 0.5rem 1rem;
        background: #f1f5f9;
        color: #475569;
        border: none;
        border-radius: 0.5rem;
        font-weight: 500;
        cursor: pointer;
    }
    .btn-primary {
        padding: 0.5rem 1rem;
        background: linear-gradient(135deg, #7c3aed, #a855f7);
        color: #fff;
        border: none;
        border-radius: 0.5rem;
        font-weight: 500;
        cursor: pointer;
    }
    .btn-primary:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }
</style>
