<script lang="ts">
    import { page } from "$app/stores";
    import { api } from "$lib";
    import { onMount } from "svelte";

    $: unitId = $page.params.unitId ?? "";

    let rts: { id: string; name: string }[] = [];
    let wargas: any[] = [];
    let iurans: any[] = [];
    let attendances: any[] = [];
    let isLoading = true;

    // Filter
    let filterRTId = "";
    let filterAcademicYear = new Date().getFullYear().toString();
    let filterSemester = 1;
    let reportType = "attendance";

    // For pengurus-kegiatan report
    let kegiatans: { id: string; name: string; code: string }[] = [];
    let pengurussByKegiatan: Map<string, any[]> = new Map();
    let expandedKegiatans: Set<string> = new Set();

    onMount(async () => {
        await Promise.all([loadRTs(), loadKegiatans()]);
        if (rts.length > 0) {
            filterRTId = rts[0].id;
            await loadData();
        }
    });

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

    async function loadData() {
        isLoading = true;
        try {
            const [wargasRes, iuransRes, attendancesRes] = await Promise.all([
                api.getWargas({
                    unit_id: unitId,
                    rt_id: filterRTId,
                    limit: 200,
                }),
                api.getIurans({
                    unit_id: unitId,
                    rt_id: filterRTId,
                    academic_year: filterAcademicYear,
                    semester: filterSemester,
                    limit: 500,
                }),
                api.getAttendances({
                    unit_id: unitId,
                    rt_id: filterRTId,
                    limit: 500,
                }),
            ]);

            wargas = wargasRes.data || [];
            iurans = iuransRes.data || [];
            attendances = attendancesRes.data || [];
        } catch (err) {
            console.error("Failed to load data:", err);
        } finally {
            isLoading = false;
        }
    }

    async function loadKegiatans() {
        try {
            const response = await api.getKegiatans({
                unit_id: unitId,
                limit: 100,
            });
            kegiatans = response.data || [];
            // Load penguruss per kegiatan
            await loadPengurussByKegiatan();
        } catch (err) {
            console.error("Failed to load kegiatans:", err);
        }
    }

    async function loadPengurussByKegiatan() {
        pengurussByKegiatan = new Map();
        for (const kegiatan of kegiatans) {
            try {
                const response = await api.getPengurussByKegiatan(kegiatan.id);
                pengurussByKegiatan.set(kegiatan.id, response.data || []);
            } catch {
                pengurussByKegiatan.set(kegiatan.id, []);
            }
        }
        pengurussByKegiatan = pengurussByKegiatan;
    }

    function toggleKegiatanExpand(kegiatanId: string) {
        if (expandedKegiatans.has(kegiatanId)) {
            expandedKegiatans.delete(kegiatanId);
        } else {
            expandedKegiatans.add(kegiatanId);
        }
        expandedKegiatans = expandedKegiatans;
    }

    // Aggregate attendance by warga
    $: attendanceSummary = wargas.map((warga) => {
        const wargaAtt = attendances.filter(
            (a) => a.warga_id === warga.id,
        );
        return {
            warga,
            present: wargaAtt.filter((a) => a.status === "present").length,
            sick: wargaAtt.filter((a) => a.status === "sick").length,
            permission: wargaAtt.filter((a) => a.status === "permission")
                .length,
            absent: wargaAtt.filter((a) => a.status === "absent").length,
            total: wargaAtt.length,
        };
    });

    // Aggregate iurans by warga
    $: iuranSummary = wargas
        .map((warga) => {
            const wargaIurans = iurans.filter(
                (g) => g.warga_id === warga.id,
            );
            const avgScore =
                wargaIurans.length > 0
                    ? wargaIurans.reduce(
                          (sum, g) => sum + (g.score / g.max_score) * 100,
                          0,
                      ) / wargaIurans.length
                    : 0;
            return {
                warga,
                iuranCount: wargaIurans.length,
                avgScore: avgScore.toFixed(1),
            };
        })
        .sort((a, b) => parseFloat(b.avgScore) - parseFloat(a.avgScore));

    function getRTName(rtId: string): string {
        return rts.find((c) => c.id === rtId)?.name || "-";
    }

    function exportToCSV() {
        let csv = "";
        let filename = "";

        if (reportType === "attendance") {
            csv = "No,Nama Karyawan,NIS,Hadir,Sakit,Izin,Alpha,Total\n";
            attendanceSummary.forEach((row, i) => {
                csv += `${i + 1},${row.warga.name},${row.warga.nis},${row.present},${row.sick},${row.permission},${row.absent},${row.total}\n`;
            });
            filename = `laporan-kehadiran-${filterRTId}.csv`;
        } else {
            csv = "No,Nama Karyawan,NIS,Jumlah Nilai,Rata-rata\n";
            iuranSummary.forEach((row, i) => {
                csv += `${i + 1},${row.warga.name},${row.warga.nis},${row.iuranCount},${row.avgScore}\n`;
            });
            filename = `laporan-nilai-${filterRTId}.csv`;
        }

        const blob = new Blob([csv], { type: "text/csv" });
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = filename;
        a.click();
    }
</script>

<svelte:head>
    <title>Laporan - SeekOlah</title>
</svelte:head>

<div class="py-8 px-4">
    <div class="max-w-7xl mx-auto">
        <!-- Header -->
        <div class="flex justify-between items-center mb-8">
            <div>
                <h1 class="text-3xl font-bold mb-2">📑 Laporan</h1>
                <p class="text-gray-600">Rekap kehadiran dan nilai warga</p>
            </div>
            <button on:click={exportToCSV} class="btn-primary">
                📥 Export CSV
            </button>
        </div>

        <!-- Filters -->
        <div class="card mb-6">
            <div class="p-4 flex gap-4 flex-wrap items-end">
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1"
                        >Tipe Laporan</label
                    >
                    <select
                        bind:value={reportType}
                        class="input-field min-w-[180px]"
                    >
                        <option value="attendance">Kehadiran</option>
                        <option value="iurans">Nilai</option>
                        <option value="pengurus-kegiatan">Pengurus per Proyek</option>
                    </select>
                </div>
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1"
                        >RT</label
                    >
                    <select
                        bind:value={filterRTId}
                        on:change={loadData}
                        class="input-field min-w-[150px]"
                    >
                        {#each rts as cls}
                            <option value={cls.id}>{cls.name}</option>
                        {/each}
                    </select>
                </div>
                {#if reportType === "iurans"}
                    <div>
                        <label
                            class="block text-sm font-medium text-gray-700 mb-1"
                            >Tahun Ajaran</label
                        >
                        <input
                            type="text"
                            bind:value={filterAcademicYear}
                            on:change={loadData}
                            class="input-field w-[100px]"
                        />
                    </div>
                    <div>
                        <label
                            class="block text-sm font-medium text-gray-700 mb-1"
                            >Semester</label
                        >
                        <select
                            bind:value={filterSemester}
                            on:change={loadData}
                            class="input-field"
                        >
                            <option value={1}>Ganjil</option>
                            <option value={2}>Genap</option>
                        </select>
                    </div>
                {/if}
            </div>
        </div>

        <!-- Report Table -->
        <div
            class="card bg-white border border-gray-200 rounded-xl overflow-hidden"
        >
            {#if isLoading}
                <div class="text-center py-12 text-gray-600">
                    Memuat data...
                </div>
            {:else}
                <div class="px-6 py-4 border-b border-gray-200 bg-gray-100/50">
                    <h2 class="text-lg font-semibold text-gray-900">
                        {reportType === "attendance"
                            ? "Rekap Kehadiran"
                            : reportType === "iurans"
                              ? "Rekap Nilai"
                              : "Pengurus per Mata Pelajaran"}
                        {#if reportType !== "pengurus-kegiatan"}
                            - RT {getRTName(filterRTId)}
                        {/if}
                    </h2>
                </div>

                <div class="overflow-x-auto">
                    {#if reportType === "attendance"}
                        <table class="w-full">
                            <thead style="background-color: #f3f4f6;">
                                <tr>
                                    <th
                                        class="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase"
                                        >No</th
                                    >
                                    <th
                                        class="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase"
                                        >Nama</th
                                    >
                                    <th
                                        class="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase"
                                        >NIS</th
                                    >
                                    <th
                                        class="px-4 py-3 text-center text-xs font-semibold text-green-700 uppercase"
                                        >Hadir</th
                                    >
                                    <th
                                        class="px-4 py-3 text-center text-xs font-semibold text-amber-600 uppercase"
                                        >Sakit</th
                                    >
                                    <th
                                        class="px-4 py-3 text-center text-xs font-semibold text-blue-700 uppercase"
                                        >Izin</th
                                    >
                                    <th
                                        class="px-4 py-3 text-center text-xs font-semibold text-red-700 uppercase"
                                        >Alpha</th
                                    >
                                    <th
                                        class="px-4 py-3 text-center text-xs font-medium text-gray-700 uppercase"
                                        >Total</th
                                    >
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-gray-200">
                                {#each attendanceSummary as row, i}
                                    <tr class="hover:bg-gray-100/30">
                                        <td
                                            class="px-4 py-3 text-sm text-gray-600"
                                            >{i + 1}</td
                                        >
                                        <td
                                            class="px-4 py-3 text-sm font-medium text-gray-900"
                                            >{row.warga.name}</td
                                        >
                                        <td
                                            class="px-4 py-3 text-sm text-gray-600"
                                            >{row.warga.nis}</td
                                        >
                                        <td
                                            class="px-4 py-3 text-center text-sm font-medium text-green-600"
                                            >{row.present}</td
                                        >
                                        <td
                                            class="px-4 py-3 text-center text-sm font-medium text-yellow-600"
                                            >{row.sick}</td
                                        >
                                        <td
                                            class="px-4 py-3 text-center text-sm font-medium text-blue-600"
                                            >{row.permission}</td
                                        >
                                        <td
                                            class="px-4 py-3 text-center text-sm font-medium text-red-600"
                                            >{row.absent}</td
                                        >
                                        <td
                                            class="px-4 py-3 text-center text-sm text-gray-600"
                                            >{row.total}</td
                                        >
                                    </tr>
                                {/each}
                            </tbody>
                        </table>
                    {:else if reportType === "iurans"}
                        <table class="w-full">
                            <thead class="bg-gray-100/50">
                                <tr>
                                    <th
                                        class="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase"
                                        >Rank</th
                                    >
                                    <th
                                        class="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase"
                                        >Nama</th
                                    >
                                    <th
                                        class="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase"
                                        >NIS</th
                                    >
                                    <th
                                        class="px-4 py-3 text-center text-xs font-medium text-gray-700 uppercase"
                                        >Jumlah Nilai</th
                                    >
                                    <th
                                        class="px-4 py-3 text-center text-xs font-medium text-gray-700 uppercase"
                                        >Rata-rata</th
                                    >
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-slate-700">
                                {#each iuranSummary as row, i}
                                    <tr class="hover:bg-gray-100/30">
                                        <td class="px-4 py-3 text-sm">
                                            {#if i === 0}
                                                <span class="text-2xl">🥇</span>
                                            {:else if i === 1}
                                                <span class="text-2xl">🥈</span>
                                            {:else if i === 2}
                                                <span class="text-2xl">🥉</span>
                                            {:else}
                                                <span
                                                    class="font-medium text-gray-600"
                                                    >{i + 1}</span
                                                >
                                            {/if}
                                        </td>
                                        <td
                                            class="px-4 py-3 text-sm font-medium text-gray-900"
                                            >{row.warga.name}</td
                                        >
                                        <td
                                            class="px-4 py-3 text-sm text-gray-600"
                                            >{row.warga.nis}</td
                                        >
                                        <td
                                            class="px-4 py-3 text-center text-sm text-gray-600"
                                            >{row.iuranCount}</td
                                        >
                                        <td class="px-4 py-3 text-center">
                                            <span
                                                class="text-lg font-bold {parseFloat(
                                                    row.avgScore,
                                                ) >= 70
                                                    ? 'text-green-600'
                                                    : 'text-red-600'}"
                                            >
                                                {row.avgScore}
                                            </span>
                                        </td>
                                    </tr>
                                {/each}
                            </tbody>
                        </table>
                    {:else}
                        <!-- Pengurus-Kegiatan Report -->
                        <div class="p-6">
                            <p class="text-sm text-gray-600 mb-4">
                                Total {kegiatans.length} kegiatan terdaftar
                            </p>
                            <div class="space-y-3">
                                {#each kegiatans as kegiatan}
                                    {@const penguruss =
                                        pengurussByKegiatan.get(kegiatan.id) || []}
                                    <div
                                        class="bg-gray-100/50 rounded-lg border border-gray-300 overflow-hidden"
                                    >
                                        <button
                                            on:click={() =>
                                                toggleKegiatanExpand(kegiatan.id)}
                                            class="w-full flex items-center justify-between px-4 py-3 hover:bg-gray-100/70 transition-colors"
                                        >
                                            <div
                                                class="flex items-center gap-3"
                                            >
                                                <span
                                                    class="text-xs font-mono bg-blue-500/20 text-blue-600 px-2 py-1 rounded"
                                                    >{kegiatan.code}</span
                                                >
                                                <span
                                                    class="font-medium text-gray-900"
                                                    >{kegiatan.name}</span
                                                >
                                            </div>
                                            <div
                                                class="flex items-center gap-3"
                                            >
                                                <span
                                                    class="text-sm {penguruss.length >
                                                    0
                                                        ? 'text-green-600'
                                                        : 'text-yellow-600'}"
                                                >
                                                    {penguruss.length} pengurus
                                                </span>
                                                <span
                                                    class="text-gray-600 text-lg"
                                                >
                                                    {expandedKegiatans.has(
                                                        kegiatan.id,
                                                    )
                                                        ? "▼"
                                                        : "▶"}
                                                </span>
                                            </div>
                                        </button>
                                        {#if expandedKegiatans.has(kegiatan.id)}
                                            <div
                                                class="border-t border-gray-300 px-4 py-3"
                                            >
                                                {#if penguruss.length === 0}
                                                    <p
                                                        class="text-sm text-gray-600 italic"
                                                    >
                                                        Belum ada pengurus yang
                                                        mengajar proyek ini
                                                    </p>
                                                {:else}
                                                    <div
                                                        class="flex flex-wrap gap-2"
                                                    >
                                                        {#each penguruss as ts}
                                                            <span
                                                                class="inline-flex items-center gap-1 px-3 py-1 bg-gray-200 rounded-full text-sm"
                                                            >
                                                                <span
                                                                    class="w-2 h-2 bg-green-400 rounded-full"
                                                                ></span>
                                                                {ts.pengurus_name}
                                                            </span>
                                                        {/each}
                                                    </div>
                                                {/if}
                                            </div>
                                        {/if}
                                    </div>
                                {/each}
                            </div>
                        </div>
                    {/if}
                </div>

                {#if wargas.length === 0}
                    <div class="text-center py-12 text-gray-600">
                        Tidak ada data warga untuk tim ini
                    </div>
                {/if}
            {/if}
        </div>
    </div>
</div>

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
    :global(table thead) {
        background: #f3f4f6;
    }
    :global(table th) {
        color: #374151;
    }
    :global(table td) {
        color: #1f2937;
    }
</style>
