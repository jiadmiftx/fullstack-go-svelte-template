<script lang="ts">
    import { page } from "$app/stores";
    import { api } from "$lib";
    import { onMount } from "svelte";
    import Modal from "$core/components/Modal.svelte";
    import DataTable from "$core/components/DataTable.svelte";
    import { showToast } from "$core/components/Toast.svelte";

    interface Iuran {
        id: string;
        warga_id: string;
        kegiatan_id: string;
        rt_id: string;
        pengurus_id: string;
        academic_year: string;
        semester: number;
        type: string;
        score: number;
        max_score: number;
        notes?: string;
    }

    interface RT {
        id: string;
        name: string;
    }
    interface Kegiatan {
        id: string;
        name: string;
        code: string;
    }
    interface Warga {
        id: string;
        name: string;
        nis: string;
    }

    $: unitId = $page.params.unitId ?? "";

    let iurans: Iuran[] = [];
    let rts: RT[] = [];
    let kegiatans: Kegiatan[] = [];
    let wargas: Warga[] = [];
    let isLoading = false;
    let error = "";
    let showModal = false;
    let editingIuran: Iuran | null = null;

    // Filters
    let filterRTId = "";
    let filterKegiatanId = "";
    let filterType = "";
    let filterAcademicYear = new Date().getFullYear().toString();
    let filterSemester = 1;

    // Form fields
    let formWargaId = "";
    let formKegiatanId = "";
    let formRTId = "";
    let formType = "daily";
    let formScore = 0;
    let formMaxScore = 100;
    let formNotes = "";

    const iuranTypes = [
        { value: "daily", label: "Nilai Harian" },
        { value: "midterm", label: "UTS" },
        { value: "final", label: "UAS" },
        { value: "assignment", label: "Tugas" },
    ];

    onMount(async () => {
        await Promise.all([
            loadIurans(),
            loadRTs(),
            loadKegiatans(),
            loadWargas(),
        ]);
    });

    async function loadIurans() {
        isLoading = true;
        error = "";
        try {
            const params: any = {
                unit_id: unitId,
                academic_year: filterAcademicYear,
                semester: filterSemester,
                limit: 200,
            };
            if (filterRTId) params.rt_id = filterRTId;
            if (filterKegiatanId) params.kegiatan_id = filterKegiatanId;
            if (filterType) params.type = filterType;

            const response = await api.getIurans(params);
            iurans = response.data || [];
        } catch (err) {
            error = err instanceof Error ? err.message : "Gagal memuat nilai";
        } finally {
            isLoading = false;
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

    async function loadKegiatans() {
        try {
            const response = await api.getKegiatans({
                unit_id: unitId,
                limit: 100,
            });
            kegiatans = response.data || [];
        } catch (err) {
            console.error("Failed to load kegiatans:", err);
        }
    }

    async function loadWargas() {
        try {
            const response = await api.getWargas({
                unit_id: unitId,
                limit: 200,
            });
            wargas = response.data || [];
        } catch (err) {
            console.error("Failed to load wargas:", err);
        }
    }

    function openCreateModal() {
        editingIuran = null;
        formWargaId = wargas[0]?.id || "";
        formKegiatanId = kegiatans[0]?.id || "";
        formRTId = rts[0]?.id || "";
        formType = "daily";
        formScore = 0;
        formMaxScore = 100;
        formNotes = "";
        showModal = true;
    }

    function openEditModal(iuran: Iuran) {
        editingIuran = iuran;
        formWargaId = iuran.warga_id;
        formKegiatanId = iuran.kegiatan_id;
        formRTId = iuran.rt_id;
        formType = iuran.type;
        formScore = iuran.score;
        formMaxScore = iuran.max_score;
        formNotes = iuran.notes || "";
        showModal = true;
    }

    async function handleSubmit() {
        error = "";
        try {
            if (editingIuran) {
                await api.updateIuran(editingIuran.id, {
                    score: formScore,
                    max_score: formMaxScore,
                    notes: formNotes || undefined,
                });
                showToast("Nilai berhasil diupdate", "success");
            } else {
                await api.createIuran({
                    unit_id: unitId,
                    warga_id: formWargaId,
                    kegiatan_id: formKegiatanId,
                    rt_id: formRTId,
                    pengurus_id: "", // TODO: get from auth
                    academic_year: filterAcademicYear,
                    semester: filterSemester,
                    type: formType,
                    score: formScore,
                    max_score: formMaxScore,
                    notes: formNotes || undefined,
                });
                showToast("Nilai berhasil ditambahkan", "success");
            }
            showModal = false;
            await loadIurans();
        } catch (err) {
            const errorMsg =
                err instanceof Error ? err.message : "Gagal menyimpan nilai";
            showToast(errorMsg, "error");
        }
    }

    async function handleDelete(iuran: Iuran) {
        if (!confirm("Apakah Anda yakin ingin menghapus nilai ini?")) return;

        try {
            await api.deleteIuran(iuran.id);
            showToast("Nilai berhasil dihapus", "success");
            await loadIurans();
        } catch (err) {
            const errorMsg =
                err instanceof Error ? err.message : "Gagal menghapus nilai";
            showToast(errorMsg, "error");
        }
    }

    function getWargaName(wargaId: string): string {
        return wargas.find((s) => s.id === wargaId)?.name || "-";
    }

    function getKegiatanName(kegiatanId: string): string {
        return kegiatans.find((s) => s.id === kegiatanId)?.name || "-";
    }

    function getRTName(rtId: string): string {
        return rts.find((c) => c.id === rtId)?.name || "-";
    }

    function getTypeName(type: string): string {
        return iuranTypes.find((t) => t.value === type)?.label || type;
    }

    function getScoreColor(score: number, max: number): string {
        const pct = (score / max) * 100;
        if (pct >= 80) return "text-green-600";
        if (pct >= 60) return "text-yellow-600";
        return "text-red-600";
    }

    const columns = [
        {
            key: "warga_id" as keyof Iuran,
            label: "Karyawan",
            render: (wargaId: string) => getWargaName(wargaId),
        },
        {
            key: "kegiatan_id" as keyof Iuran,
            label: "Mata Pelajaran",
            render: (kegiatanId: string) => getKegiatanName(kegiatanId),
        },
        {
            key: "rt_id" as keyof Iuran,
            label: "RT",
            render: (rtId: string) => getRTName(rtId),
        },
        {
            key: "type" as keyof Iuran,
            label: "Tipe",
            render: (type: string) =>
                `<span class="bg-gray-500/20 px-2 py-1 rounded text-xs">${getTypeName(type)}</span>`,
        },
        {
            key: "score" as keyof Iuran,
            label: "Nilai",
            render: (score: number, iuran: Iuran) =>
                `<span class="font-bold ${getScoreColor(score, iuran.max_score)}">${score}/${iuran.max_score}</span>`,
        },
    ];

    const actions = [
        {
            label: "Edit",
            variant: "secondary" as const,
            onClick: openEditModal,
        },
        { label: "Hapus", variant: "danger" as const, onClick: handleDelete },
    ];
</script>

<svelte:head>
    <title>Nilai Karyawan - SeekOlah</title>
</svelte:head>

<div class="py-8 px-4">
    <div class="max-w-7xl mx-auto">
        <!-- Header -->
        <div class="flex justify-between items-center mb-8">
            <div>
                <h1 class="text-3xl font-bold mb-2">📊 Nilai Karyawan</h1>
                <p class="text-gray-600">
                    Kelola nilai warga per kegiatan
                </p>
            </div>
            <button on:click={openCreateModal} class="btn-primary">
                + Input Nilai
            </button>
        </div>

        <!-- Filters -->
        <div class="card mb-6">
            <div class="p-4 flex gap-4 flex-wrap items-end">
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1"
                        >Tahun Ajaran</label
                    >
                    <input
                        type="text"
                        bind:value={filterAcademicYear}
                        on:change={loadIurans}
                        class="input-field w-[120px]"
                        placeholder="2024"
                    />
                </div>
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1"
                        >Semester</label
                    >
                    <select
                        bind:value={filterSemester}
                        on:change={loadIurans}
                        class="input-field"
                    >
                        <option value={1}>Ganjil</option>
                        <option value={2}>Genap</option>
                    </select>
                </div>
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1"
                        >RT</label
                    >
                    <select
                        bind:value={filterRTId}
                        on:change={loadIurans}
                        class="input-field min-w-[150px]"
                    >
                        <option value="">Semua RT</option>
                        {#each rts as cls}
                            <option value={cls.id}>{cls.name}</option>
                        {/each}
                    </select>
                </div>
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1"
                        >Proyek</label
                    >
                    <select
                        bind:value={filterKegiatanId}
                        on:change={loadIurans}
                        class="input-field min-w-[150px]"
                    >
                        <option value="">Semua Proyek</option>
                        {#each kegiatans as kegiatan}
                            <option value={kegiatan.id}>{kegiatan.name}</option>
                        {/each}
                    </select>
                </div>
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1"
                        >Tipe</label
                    >
                    <select
                        bind:value={filterType}
                        on:change={loadIurans}
                        class="input-field min-w-[150px]"
                    >
                        <option value="">Semua Tipe</option>
                        {#each iuranTypes as type}
                            <option value={type.value}>{type.label}</option>
                        {/each}
                    </select>
                </div>
            </div>
        </div>

        <!-- Error Alert -->
        {#if error}
            <div
                class="bg-red-500/20 border border-red-500 text-red-400 px-4 py-3 rounded-lg mb-6"
            >
                {error}
            </div>
        {/if}

        <!-- Iurans Table -->
        <div class="card">
            {#if isLoading}
                <div class="text-center py-12 text-gray-600">
                    Memuat nilai...
                </div>
            {:else}
                <DataTable
                    data={iurans}
                    {columns}
                    {actions}
                    emptyMessage="Belum ada data nilai. Input nilai pertama!"
                />
            {/if}
        </div>
    </div>
</div>

<!-- Create/Edit Modal -->
<Modal
    bind:isOpen={showModal}
    title={editingIuran ? "Edit Nilai" : "Input Nilai"}
    size="md"
>
    <form on:submit|preventDefault={handleSubmit} class="space-y-4">
        {#if !editingIuran}
            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label
                        for="wargaId"
                        class="block text-sm font-medium text-gray-700 mb-1"
                        >Karyawan *</label
                    >
                    <select
                        id="wargaId"
                        bind:value={formWargaId}
                        class="input-field"
                        required
                    >
                        {#each wargas as warga}
                            <option value={warga.id}
                                >{warga.name} ({warga.nis})</option
                            >
                        {/each}
                    </select>
                </div>
                <div>
                    <label
                        for="rtId"
                        class="block text-sm font-medium text-gray-700 mb-1"
                        >Tim *</label
                    >
                    <select
                        id="rtId"
                        bind:value={formRTId}
                        class="input-field"
                        required
                    >
                        {#each rts as cls}
                            <option value={cls.id}>{cls.name}</option>
                        {/each}
                    </select>
                </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label
                        for="kegiatanId"
                        class="block text-sm font-medium text-gray-700 mb-1"
                        >Mata Pelajaran *</label
                    >
                    <select
                        id="kegiatanId"
                        bind:value={formKegiatanId}
                        class="input-field"
                        required
                    >
                        {#each kegiatans as kegiatan}
                            <option value={kegiatan.id}>{kegiatan.name}</option>
                        {/each}
                    </select>
                </div>
                <div>
                    <label
                        for="type"
                        class="block text-sm font-medium text-gray-700 mb-1"
                        >Tipe Nilai *</label
                    >
                    <select
                        id="type"
                        bind:value={formType}
                        class="input-field"
                        required
                    >
                        {#each iuranTypes as type}
                            <option value={type.value}>{type.label}</option>
                        {/each}
                    </select>
                </div>
            </div>
        {/if}

        <div class="grid grid-cols-2 gap-4">
            <div>
                <label
                    for="score"
                    class="block text-sm font-medium text-gray-700 mb-1"
                    >Nilai *</label
                >
                <input
                    id="score"
                    type="number"
                    bind:value={formScore}
                    class="input-field"
                    required
                    min="0"
                    max={formMaxScore}
                    step="0.01"
                />
            </div>
            <div>
                <label
                    for="maxScore"
                    class="block text-sm font-medium text-gray-700 mb-1"
                    >Nilai Max</label
                >
                <input
                    id="maxScore"
                    type="number"
                    bind:value={formMaxScore}
                    class="input-field"
                    min="1"
                    step="0.01"
                />
            </div>
        </div>

        <div>
            <label
                for="notes"
                class="block text-sm font-medium text-gray-700 mb-1"
                >Catatan</label
            >
            <textarea
                id="notes"
                bind:value={formNotes}
                class="input-field"
                rows="2"
                placeholder="Catatan tambahan..."
            ></textarea>
        </div>
    </form>

    <div slot="footer">
        <button on:click={() => (showModal = false)} class="btn-secondary"
            >Batal</button
        >
        <button on:click={handleSubmit} class="btn-primary">
            {editingIuran ? "Update" : "Simpan"}
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
        padding: 0.5rem 0.75rem;
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
</style>
