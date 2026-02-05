<script lang="ts">
    import { page } from "$app/stores";
    import { api } from "$lib";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { showToast } from "$core/components/Toast.svelte";

    interface Warga {
        id: string;
        unit_member_id: string;
        unit_id: string;
        user_id: string;
        rt_id?: string;
        nis: string;
        nisn?: string;
        gender?: string;
        birth_place?: string;
        address?: string;
        parent_name?: string;
        parent_phone?: string;
        status: string;
        user_name: string;
        email?: string;
        phone?: string;
    }

    interface Iuran {
        id: string;
        kegiatan_id: string;
        type: string;
        score: number;
        max_score: number;
        kegiatan?: { name: string };
    }

    interface Attendance {
        id: string;
        date: string;
        status: string;
    }

    interface WargaRT {
        id: string;
        rt_id: string;
        rt_type: string;
        rt_name: string;
        is_active: boolean;
    }

    interface RTHistory {
        id: string;
        rt_name: string;
        academic_year: string;
        joined_at: string;
        left_at?: string;
        status: string; // active, promoted, transferred, graduated, dropped
        notes?: string;
    }

    $: unitId = $page.params.unitId ?? "";
    $: orgId = $page.params.orgId ?? "";
    $: wargaId = $page.params.wargaId ?? "";

    let warga: Warga | null = null;
    let iurans: Iuran[] = [];
    let attendances: Attendance[] = [];
    let rts: { id: string; name: string }[] = [];
    let wargaRTs: WargaRT[] = [];
    let rtHistory: RTHistory[] = [];
    let isLoading = true;
    let error = "";
    let isEditing = false;
    let activeTab = "iurans";

    // Form fields
    let formNis = "";
    let formNisn = "";
    let formParentName = "";
    let formParentPhone = "";
    let formAddress = "";

    onMount(async () => {
        await Promise.all([
            loadWarga(),
            loadIurans(),
            loadAttendances(),
            loadRTs(),
            loadWargaRTs(),
            loadRTHistory(),
        ]);
    });

    async function loadWarga() {
        isLoading = true;
        error = "";
        try {
            const response = await api.getWarga(wargaId);
            warga = response.data;
            formNis = warga?.nis || "";
            formNisn = warga?.nisn || "";
            formParentName = warga?.parent_name || "";
            formParentPhone = warga?.parent_phone || "";
            formAddress = warga?.address || "";
        } catch (err) {
            error =
                err instanceof Error ? err.message : "Gagal memuat data warga";
        } finally {
            isLoading = false;
        }
    }

    async function loadIurans() {
        try {
            const response = await api.getIurans({
                unit_id: unitId,
                warga_id: wargaId,
                limit: 50,
            });
            iurans = response.data || [];
        } catch (err) {
            console.error("Failed to load iurans:", err);
        }
    }

    async function loadAttendances() {
        try {
            const response = await api.getAttendances({
                unit_id: unitId,
                warga_id: wargaId,
                limit: 50,
            });
            attendances = response.data || [];
        } catch (err) {
            console.error("Failed to load attendances:", err);
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

    async function loadWargaRTs() {
        try {
            const response = await api.getRTsByWarga(wargaId);
            wargaRTs = response.data || [];
        } catch (err) {
            console.error("Failed to load warga rts:", err);
        }
    }

    async function loadRTHistory() {
        try {
            const response = await api.getWargaRTHistory(wargaId);
            rtHistory = response.data || [];
        } catch (err) {
            console.error("Failed to load rt history:", err);
        }
    }

    async function handleSave() {
        if (!warga) return;

        try {
            await api.updateWarga(warga.id, {
                nis: formNis,
                nisn: formNisn || undefined,
                parent_name: formParentName || undefined,
                parent_phone: formParentPhone || undefined,
                address: formAddress || undefined,
            });
            showToast("Profil warga berhasil diupdate", "success");
            isEditing = false;
            await loadWarga();
        } catch (err) {
            const errorMsg =
                err instanceof Error ? err.message : "Gagal menyimpan profil";
            showToast(errorMsg, "error");
        }
    }

    function goBack() {
        goto(`/org/${orgId}/units/${unitId}/wargas`);
    }

    function getRTName(rtId?: string): string {
        if (!rtId) return "-";
        return rts.find((c) => c.id === rtId)?.name || "-";
    }

    function getStatusColor(status: string): string {
        switch (status) {
            case "present":
                return "text-green-600";
            case "sick":
                return "text-yellow-600";
            case "permission":
                return "text-blue-600";
            case "absent":
                return "text-red-600";
            default:
                return "text-gray-600";
        }
    }

    function getStatusLabel(status: string): string {
        switch (status) {
            case "present":
                return "Hadir";
            case "sick":
                return "Sakit";
            case "permission":
                return "Izin";
            case "absent":
                return "Alpha";
            default:
                return status;
        }
    }

    // Calculate attendance stats
    $: attendanceStats = {
        present: attendances.filter((a) => a.status === "present").length,
        sick: attendances.filter((a) => a.status === "sick").length,
        permission: attendances.filter((a) => a.status === "permission").length,
        absent: attendances.filter((a) => a.status === "absent").length,
    };

    // Calculate iuran average
    $: iuranAverage =
        iurans.length > 0
            ? (
                  iurans.reduce(
                      (sum, g) => sum + (g.score / g.max_score) * 100,
                      0,
                  ) / iurans.length
              ).toFixed(1)
            : "-";
</script>

<svelte:head>
    <title>{warga?.user_name || "Profile Karyawan"} - SeekOlah</title>
</svelte:head>

<div class="py-8 px-4">
    <div class="max-w-4xl mx-auto">
        <!-- Back Button -->
        <button
            on:click={goBack}
            class="flex items-center gap-2 text-gray-600 hover:text-slate-600 mb-6"
        >
            ← Kembali ke Daftar Karyawan
        </button>

        {#if isLoading}
            <div class="card">
                <div class="text-center py-12 text-gray-600">
                    Memuat data...
                </div>
            </div>
        {:else if error}
            <div class="card">
                <div class="text-center py-12 text-red-400">{error}</div>
            </div>
        {:else if warga}
            <!-- Profile Header -->
            <div class="card mb-6">
                <div class="p-6">
                    <div class="flex items-start gap-6">
                        <div
                            class="w-24 h-24 rounded-full bg-gradient-to-br from-green-500 to-teal-600 flex items-center justify-center text-white text-3xl font-bold"
                        >
                            {(warga.user_name || "S").charAt(0).toUpperCase()}
                        </div>
                        <div class="flex-1">
                            <div class="flex justify-between items-start">
                                <div>
                                    <h1
                                        class="text-2xl font-bold text-gray-800"
                                    >
                                        {warga.user_name || "Tanpa Nama"}
                                    </h1>
                                    <p class="text-gray-500 mt-1">
                                        🎓 Karyawan • RT {getRTName(
                                            warga.rt_id,
                                        )}
                                    </p>
                                </div>
                                <button
                                    on:click={() => (isEditing = !isEditing)}
                                    class="btn-secondary text-sm"
                                >
                                    {isEditing ? "Batal" : "Edit Profil"}
                                </button>
                            </div>

                            <div class="grid grid-cols-2 gap-4 mt-4">
                                <div>
                                    <div
                                        class="text-xs text-gray-400 uppercase"
                                    >
                                        NIS
                                    </div>
                                    {#if isEditing}
                                        <input
                                            type="text"
                                            bind:value={formNis}
                                            class="input-field mt-1"
                                        />
                                    {:else}
                                        <div class="font-medium text-gray-700">
                                            {warga.nis || "-"}
                                        </div>
                                    {/if}
                                </div>
                                <div>
                                    <div
                                        class="text-xs text-gray-400 uppercase"
                                    >
                                        NISN
                                    </div>
                                    {#if isEditing}
                                        <input
                                            type="text"
                                            bind:value={formNisn}
                                            class="input-field mt-1"
                                        />
                                    {:else}
                                        <div class="font-medium text-gray-700">
                                            {warga.nisn || "-"}
                                        </div>
                                    {/if}
                                </div>
                                <div>
                                    <div
                                        class="text-xs text-gray-400 uppercase"
                                    >
                                        Nama Orang Tua
                                    </div>
                                    {#if isEditing}
                                        <input
                                            type="text"
                                            bind:value={formParentName}
                                            class="input-field mt-1"
                                        />
                                    {:else}
                                        <div class="font-medium text-gray-700">
                                            {warga.parent_name || "-"}
                                        </div>
                                    {/if}
                                </div>
                                <div>
                                    <div
                                        class="text-xs text-gray-400 uppercase"
                                    >
                                        Telepon Orang Tua
                                    </div>
                                    {#if isEditing}
                                        <input
                                            type="text"
                                            bind:value={formParentPhone}
                                            class="input-field mt-1"
                                        />
                                    {:else}
                                        <div class="font-medium text-gray-700">
                                            {warga.parent_phone || "-"}
                                        </div>
                                    {/if}
                                </div>
                            </div>

                            {#if isEditing}
                                <div class="mt-4">
                                    <button
                                        on:click={handleSave}
                                        class="btn-primary"
                                    >
                                        Simpan Perubahan
                                    </button>
                                </div>
                            {/if}
                        </div>
                    </div>
                </div>
            </div>

            <!-- RT yang Diikuti -->
            <div class="card mb-6">
                <div class="px-6 py-4 border-b border-gray-200">
                    <h2 class="text-lg font-semibold text-gray-800">
                        🏫 RT yang Diikuti
                    </h2>
                </div>
                <div class="p-6">
                    {#if wargaRTs.length === 0}
                        <div class="text-center py-4 text-gray-400">
                            Belum ada tim yang diikuti
                        </div>
                    {:else}
                        <div class="flex flex-wrap gap-2">
                            {#each wargaRTs as sc}
                                <span
                                    class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-sm font-medium {sc.rt_type ===
                                    'regular'
                                        ? 'bg-blue-100 text-blue-700'
                                        : 'bg-purple-100 text-purple-700'}"
                                >
                                    {sc.rt_name}
                                    <span class="text-xs opacity-75">
                                        ({sc.rt_type === "regular"
                                            ? "Reguler"
                                            : "Ekstrakurikuler"})
                                    </span>
                                </span>
                            {/each}
                        </div>
                    {/if}
                </div>
            </div>

            <!-- Quick Stats -->
            <div class="grid grid-cols-2 gap-4 mb-6">
                <div class="card p-4">
                    <div class="text-sm text-gray-400">Rata-rata Nilai</div>
                    <div class="text-3xl font-bold text-blue-600">
                        {iuranAverage}
                    </div>
                </div>
                <div class="card p-4">
                    <div class="text-sm text-gray-400">Kehadiran</div>
                    <div class="flex gap-2 mt-1">
                        <span class="text-green-600 font-bold"
                            >{attendanceStats.present}</span
                        >
                        <span class="text-gray-300">|</span>
                        <span class="text-yellow-600 font-bold"
                            >{attendanceStats.sick}</span
                        >
                        <span class="text-gray-300">|</span>
                        <span class="text-blue-600 font-bold"
                            >{attendanceStats.permission}</span
                        >
                        <span class="text-gray-300">|</span>
                        <span class="text-red-600 font-bold"
                            >{attendanceStats.absent}</span
                        >
                    </div>
                </div>
            </div>

            <!-- Tabs -->
            <div class="card">
                <div class="border-b border-gray-200">
                    <div class="flex">
                        <button
                            on:click={() => (activeTab = "iurans")}
                            class="px-6 py-3 font-medium {activeTab === 'iurans'
                                ? 'text-blue-600 border-b-2 border-blue-600'
                                : 'text-gray-500'}"
                        >
                            📊 Nilai
                        </button>
                        <button
                            on:click={() => (activeTab = "attendance")}
                            class="px-6 py-3 font-medium {activeTab ===
                            'attendance'
                                ? 'text-blue-600 border-b-2 border-blue-600'
                                : 'text-gray-500'}"
                        >
                            📋 Kehadiran
                        </button>
                        <button
                            on:click={() => (activeTab = "history")}
                            class="px-6 py-3 font-medium {activeTab ===
                            'history'
                                ? 'text-blue-600 border-b-2 border-blue-600'
                                : 'text-gray-500'}"
                        >
                            📜 Riwayat RT
                        </button>
                    </div>
                </div>

                <div class="p-6">
                    {#if activeTab === "iurans"}
                        {#if iurans.length === 0}
                            <div class="text-center py-8 text-gray-400">
                                Belum ada data nilai
                            </div>
                        {:else}
                            <div class="space-y-3">
                                {#each iurans as iuran}
                                    <div
                                        class="flex items-center justify-between p-4 bg-gray-50 rounded-lg"
                                    >
                                        <div>
                                            <div class="font-medium">
                                                {iuran.kegiatan?.name ||
                                                    "Mata Pelajaran"}
                                            </div>
                                            <div class="text-sm text-gray-500">
                                                {iuran.type}
                                            </div>
                                        </div>
                                        <div
                                            class="text-xl font-bold {iuran.score >=
                                            70
                                                ? 'text-green-600'
                                                : 'text-red-600'}"
                                        >
                                            {iuran.score}/{iuran.max_score}
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        {/if}
                    {:else if activeTab === "attendance"}
                        {#if attendances.length === 0}
                            <div class="text-center py-8 text-gray-400">
                                Belum ada data kehadiran
                            </div>
                        {:else}
                            <div class="space-y-2">
                                {#each attendances.slice(0, 20) as att}
                                    <div
                                        class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
                                    >
                                        <div class="font-medium">
                                            {att.date}
                                        </div>
                                        <span
                                            class="font-medium {getStatusColor(
                                                att.status,
                                            )}"
                                        >
                                            {getStatusLabel(att.status)}
                                        </span>
                                    </div>
                                {/each}
                            </div>
                        {/if}
                    {:else if activeTab === "history"}
                        {#if rtHistory.length === 0}
                            <div class="text-center py-8 text-gray-400">
                                Belum ada riwayat tim
                            </div>
                        {:else}
                            <div class="timeline">
                                {#each rtHistory as history}
                                    <div class="timeline-item">
                                        <div
                                            class="timeline-dot {history.status ===
                                            'active'
                                                ? 'active'
                                                : ''}"
                                        ></div>
                                        <div class="timeline-content">
                                            <div class="timeline-year">
                                                {history.academic_year}
                                            </div>
                                            <div class="timeline-rt">
                                                {history.rt_name}
                                            </div>
                                            <div class="timeline-status">
                                                <span
                                                    class="status-badge {history.status}"
                                                >
                                                    {history.status === "active"
                                                        ? "🟢 Aktif"
                                                        : history.status ===
                                                            "promoted"
                                                          ? "⬆️ Naik RT"
                                                          : history.status ===
                                                              "transferred"
                                                            ? "↔️ Pindah"
                                                            : history.status ===
                                                                "graduated"
                                                              ? "🎓 Lulus"
                                                              : history.status}
                                                </span>
                                            </div>
                                            {#if history.notes}
                                                <div class="timeline-notes">
                                                    {history.notes}
                                                </div>
                                            {/if}
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        {/if}
                    {/if}
                </div>
            </div>
        {/if}
    </div>
</div>

<style>
    .card {
        background: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 1rem;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    }
    .input-field {
        width: 100%;
        padding: 0.5rem 0.75rem;
        background: #ffffff;
        border: 1px solid #d1d5db;
        border-radius: 0.5rem;
        color: #1f2937;
        font-size: 0.875rem;
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
    }
    .btn-secondary {
        padding: 0.5rem 1rem;
        background: #f1f5f9;
        color: #475569;
        font-weight: 500;
        border: 1px solid #e2e8f0;
        border-radius: 0.5rem;
        cursor: pointer;
    }
    .btn-secondary:hover {
        background: #e2e8f0;
    }
    /* Timeline styles */
    .timeline {
        position: relative;
        padding-left: 2rem;
    }
    .timeline::before {
        content: "";
        position: absolute;
        left: 0.5rem;
        top: 0;
        bottom: 0;
        width: 2px;
        background: #e2e8f0;
    }
    .timeline-item {
        position: relative;
        padding-bottom: 1.5rem;
        display: flex;
        gap: 1rem;
    }
    .timeline-item:last-child {
        padding-bottom: 0;
    }
    .timeline-dot {
        position: absolute;
        left: -1.5rem;
        top: 0.25rem;
        width: 1rem;
        height: 1rem;
        border-radius: 50%;
        background: #94a3b8;
        border: 3px solid #ffffff;
        box-shadow: 0 0 0 2px #e2e8f0;
    }
    .timeline-dot.active {
        background: #22c55e;
        box-shadow: 0 0 0 2px #bbf7d0;
    }
    .timeline-content {
        flex: 1;
        background: #f8fafc;
        padding: 1rem;
        border-radius: 0.75rem;
        border: 1px solid #e2e8f0;
    }
    .timeline-year {
        font-size: 0.875rem;
        color: #7c3aed;
        font-weight: 600;
        margin-bottom: 0.25rem;
    }
    .timeline-rt {
        font-size: 1.125rem;
        font-weight: 600;
        color: #1e293b;
        margin-bottom: 0.5rem;
    }
    .status-badge {
        display: inline-block;
        padding: 0.25rem 0.75rem;
        border-radius: 1rem;
        font-size: 0.75rem;
        font-weight: 500;
    }
    .status-badge.active {
        background: #dcfce7;
        color: #166534;
    }
    .status-badge.promoted {
        background: #dbeafe;
        color: #1e40af;
    }
    .status-badge.transferred {
        background: #fef3c7;
        color: #92400e;
    }
    .status-badge.graduated {
        background: #ede9fe;
        color: #5b21b6;
    }
    .timeline-notes {
        margin-top: 0.5rem;
        font-size: 0.8rem;
        color: #64748b;
        font-style: italic;
    }
</style>
