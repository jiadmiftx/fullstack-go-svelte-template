<script lang="ts">
    import { page } from "$app/stores";
    import { api } from "$lib";
    import { onMount } from "svelte";

    $: orgId = $page.params.orgId;
    $: unitId = $page.params.unitId;

    let stats = {
        totalWarga: 0,
        totalPengurus: 0,
        totalRumah: 0,
        pendingVerifikasi: 0,
    };
    let perumahan: any = null;
    let isLoading = true;

    // Mock data - will be replaced with API calls later
    let kasRT = {
        saldo: 15750000,
        pemasukanBulanIni: 3500000,
        pengeluaranBulanIni: 1250000,
    };

    // Pengumuman from API
    let pengumumanList: any[] = [];

    let aktivitasWarga = [
        {
            id: 1,
            nama: "Ahmad Santoso",
            aktivitas: "Membayar iuran Januari",
            waktu: "2 jam lalu",
            tipe: "iuran",
        },
        {
            id: 2,
            nama: "Siti Rahayu",
            aktivitas: "Mendaftar sebagai warga baru",
            waktu: "5 jam lalu",
            tipe: "pendaftaran",
        },
        {
            id: 3,
            nama: "Budi Prakoso",
            aktivitas: "Menghadiri rapat RT",
            waktu: "1 hari lalu",
            tipe: "kegiatan",
        },
        {
            id: 4,
            nama: "Dewi Lestari",
            aktivitas: "Update data profil",
            waktu: "2 hari lalu",
            tipe: "profil",
        },
    ];

    onMount(async () => {
        await Promise.all([loadPerumahan(), loadStats(), loadPengumuman()]);
        isLoading = false;
    });

    async function loadPerumahan() {
        try {
            const response = await api.getPerumahan(unitId);
            perumahan = response.data;
        } catch (err) {
            console.error("Failed to load perumahan:", err);
        }
    }

    async function loadStats() {
        try {
            const [wargasRes, pendingRes] = await Promise.all([
                api.getWargas({ unit_id: unitId, limit: 1 }),
                api.get(`/units/${unitId}/pending-registrations`),
            ]);

            stats.totalWarga =
                wargasRes.paginate?.total || wargasRes.data?.length || 0;
            stats.pendingVerifikasi = pendingRes.data?.length || 0;
        } catch (err) {
            console.error("Failed to load stats:", err);
        }
    }

    async function loadPengumuman() {
        try {
            const response: any = await api.get(
                `/posts?unit_id=${unitId}&page=1&page_size=5`,
            );
            if (response.data) {
                pengumumanList = response.data;
            }
        } catch (err) {
            console.error("Failed to load pengumuman:", err);
        }
    }

    function formatCurrency(amount: number): string {
        return new Intl.NumberFormat("id-ID", {
            style: "currency",
            currency: "IDR",
            minimumFractionDigits: 0,
        }).format(amount);
    }

    function formatDate(dateStr: string): string {
        return new Date(dateStr).toLocaleDateString("id-ID", {
            day: "numeric",
            month: "short",
            year: "numeric",
        });
    }
</script>

<svelte:head>
    <title>Dashboard - {perumahan?.name || "RT"}</title>
</svelte:head>

<div class="dashboard-container">
    <!-- Header -->
    <div class="dashboard-header">
        <div class="header-icon">
            <svg
                xmlns="http://www.w3.org/2000/svg"
                width="28"
                height="28"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
            >
                <rect x="3" y="3" width="7" height="9" rx="1"></rect>
                <rect x="14" y="3" width="7" height="5" rx="1"></rect>
                <rect x="14" y="12" width="7" height="9" rx="1"></rect>
                <rect x="3" y="16" width="7" height="5" rx="1"></rect>
            </svg>
        </div>
        <div>
            <h1>Dashboard</h1>
            <p class="subtitle">Selamat datang di {perumahan?.name || "RT"}</p>
        </div>
    </div>

    <!-- Summary Cards -->
    <div class="summary-grid">
        <div class="summary-card warga">
            <div class="summary-icon">👥</div>
            <div class="summary-content">
                <div class="summary-value">{stats.totalWarga}</div>
                <div class="summary-label">Total Warga</div>
            </div>
        </div>
        <div class="summary-card pending">
            <div class="summary-icon">⏳</div>
            <div class="summary-content">
                <div class="summary-value">{stats.pendingVerifikasi}</div>
                <div class="summary-label">Menunggu Verifikasi</div>
            </div>
        </div>
        <div class="summary-card kas">
            <div class="summary-icon">💰</div>
            <div class="summary-content">
                <div class="summary-value">{formatCurrency(kasRT.saldo)}</div>
                <div class="summary-label">Saldo Kas RT</div>
            </div>
        </div>
        <div class="summary-card iuran">
            <div class="summary-icon">📊</div>
            <div class="summary-content">
                <div class="summary-value positive">
                    +{formatCurrency(kasRT.pemasukanBulanIni)}
                </div>
                <div class="summary-label">Pemasukan Bulan Ini</div>
            </div>
        </div>
    </div>

    <!-- Main Grid -->
    <div class="main-grid">
        <!-- Left Column -->
        <div class="left-column">
            <!-- Pengumuman -->
            <div class="card">
                <div class="card-header">
                    <h2>📢 Pengumuman</h2>
                    <a
                        href="/org/{orgId}/units/{unitId}/pengumuman"
                        class="view-all">Lihat Semua</a
                    >
                </div>
                <div class="card-content">
                    {#if pengumumanList.length === 0}
                        <div class="empty-state">Belum ada pengumuman</div>
                    {:else}
                        {#each pengumumanList as item}
                            <a
                                href="/org/{orgId}/units/{unitId}/pengumuman"
                                class="announcement-item"
                                class:penting={item.is_important}
                            >
                                <div class="announcement-badges">
                                    {#if item.is_important}
                                        <span class="badge-penting"
                                            >PENTING</span
                                        >
                                    {/if}
                                    {#if item.is_pinned}
                                        <span class="badge-pinned">📌</span>
                                    {/if}
                                    {#if item.is_org_wide}
                                        <span class="badge-org">Organisasi</span
                                        >
                                    {:else}
                                        <span class="badge-unit">Unit</span>
                                    {/if}
                                </div>
                                {#if item.title}
                                    <h3>{item.title}</h3>
                                {/if}
                                <p>
                                    {item.content.length > 100
                                        ? item.content.substring(0, 100) + "..."
                                        : item.content}
                                </p>
                                <div class="announcement-meta">
                                    <span class="author"
                                        >👤 {item.author_name}</span
                                    >
                                    <span class="date"
                                        >{formatDate(item.created_at)}</span
                                    >
                                    <span class="comments"
                                        >💬 {item.comment_count}</span
                                    >
                                </div>
                            </a>
                        {/each}
                    {/if}
                </div>
            </div>
        </div>

        <!-- Right Column -->
        <div class="right-column">
            <!-- Kas RT Summary -->
            <div class="card kas-card">
                <div class="card-header">
                    <h2>💰 Kas RT</h2>
                    <a href="/org/{orgId}/units/{unitId}/kas" class="view-all"
                        >Detail</a
                    >
                </div>
                <div class="card-content">
                    <div class="kas-summary">
                        <div class="kas-item">
                            <span class="kas-label">Saldo Saat Ini</span>
                            <span class="kas-value saldo"
                                >{formatCurrency(kasRT.saldo)}</span
                            >
                        </div>
                        <div class="kas-divider"></div>
                        <div class="kas-row">
                            <div class="kas-item small">
                                <span class="kas-label">Pemasukan</span>
                                <span class="kas-value positive"
                                    >+{formatCurrency(
                                        kasRT.pemasukanBulanIni,
                                    )}</span
                                >
                            </div>
                            <div class="kas-item small">
                                <span class="kas-label">Pengeluaran</span>
                                <span class="kas-value negative"
                                    >-{formatCurrency(
                                        kasRT.pengeluaranBulanIni,
                                    )}</span
                                >
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Aktivitas Warga -->
            <div class="card">
                <div class="card-header">
                    <h2>📋 Aktivitas Warga</h2>
                </div>
                <div class="card-content">
                    {#if aktivitasWarga.length === 0}
                        <div class="empty-state">Belum ada aktivitas</div>
                    {:else}
                        <div class="activity-list">
                            {#each aktivitasWarga as item}
                                <div class="activity-item">
                                    <div class="activity-icon {item.tipe}">
                                        {#if item.tipe === "iuran"}💳
                                        {:else if item.tipe === "pendaftaran"}🆕
                                        {:else if item.tipe === "kegiatan"}📅
                                        {:else}👤{/if}
                                    </div>
                                    <div class="activity-content">
                                        <span class="activity-name"
                                            >{item.nama}</span
                                        >
                                        <span class="activity-desc"
                                            >{item.aktivitas}</span
                                        >
                                    </div>
                                    <span class="activity-time"
                                        >{item.waktu}</span
                                    >
                                </div>
                            {/each}
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    .dashboard-container {
        padding: 1.5rem;
        max-width: 1400px;
        margin: 0 auto;
    }

    .dashboard-header {
        display: flex;
        align-items: center;
        gap: 1rem;
        margin-bottom: 1.5rem;
    }

    .header-icon {
        width: 56px;
        height: 56px;
        border-radius: 1rem;
        background: linear-gradient(135deg, #7c3aed 0%, #a855f7 100%);
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        box-shadow: 0 8px 20px rgba(124, 58, 237, 0.3);
    }

    .dashboard-header h1 {
        font-size: 1.75rem;
        font-weight: 700;
        color: #1e293b;
        margin: 0;
    }

    .subtitle {
        color: #64748b;
        margin: 0.25rem 0 0;
    }

    /* Summary Cards */
    .summary-grid {
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        gap: 1rem;
        margin-bottom: 1.5rem;
    }

    @media (max-width: 1024px) {
        .summary-grid {
            grid-template-columns: repeat(2, 1fr);
        }
    }

    @media (max-width: 640px) {
        .summary-grid {
            grid-template-columns: 1fr;
        }
    }

    .summary-card {
        background: white;
        border-radius: 1rem;
        padding: 1.25rem;
        display: flex;
        align-items: center;
        gap: 1rem;
        border: 1px solid #e2e8f0;
        transition: all 0.2s;
    }

    .summary-card:hover {
        transform: translateY(-2px);
        box-shadow: 0 8px 20px rgba(0, 0, 0, 0.08);
    }

    .summary-icon {
        font-size: 2rem;
        width: 56px;
        height: 56px;
        border-radius: 0.75rem;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .summary-card.warga .summary-icon {
        background: #dbeafe;
    }
    .summary-card.pending .summary-icon {
        background: #fef3c7;
    }
    .summary-card.kas .summary-icon {
        background: #dcfce7;
    }
    .summary-card.iuran .summary-icon {
        background: #f3e8ff;
    }

    .summary-value {
        font-size: 1.5rem;
        font-weight: 700;
        color: #1e293b;
    }

    .summary-value.positive {
        color: #16a34a;
    }

    .summary-label {
        font-size: 0.8rem;
        color: #64748b;
    }

    /* Main Grid */
    .main-grid {
        display: grid;
        grid-template-columns: 1fr 400px;
        gap: 1.5rem;
    }

    @media (max-width: 1024px) {
        .main-grid {
            grid-template-columns: 1fr;
        }
    }

    .left-column,
    .right-column {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
    }

    /* Cards */
    .card {
        background: white;
        border-radius: 1rem;
        border: 1px solid #e2e8f0;
        overflow: hidden;
    }

    .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem 1.25rem;
        border-bottom: 1px solid #f1f5f9;
    }

    .card-header h2 {
        font-size: 1rem;
        font-weight: 600;
        color: #1e293b;
        margin: 0;
    }

    .view-all {
        font-size: 0.8rem;
        color: #7c3aed;
        text-decoration: none;
        font-weight: 500;
    }

    .view-all:hover {
        text-decoration: underline;
    }

    .card-content {
        padding: 1rem 1.25rem;
    }

    .empty-state {
        text-align: center;
        color: #94a3b8;
        padding: 2rem;
    }

    /* Announcements */
    .announcement-item {
        padding: 1rem;
        border-radius: 0.75rem;
        background: #f8fafc;
        margin-bottom: 0.75rem;
        position: relative;
    }

    .announcement-item:last-child {
        margin-bottom: 0;
    }

    .announcement-item.penting {
        background: #fef2f2;
        border-left: 3px solid #ef4444;
    }

    .announcement-badges {
        display: flex;
        flex-wrap: wrap;
        gap: 0.375rem;
        margin-bottom: 0.5rem;
    }

    .badge-penting {
        display: inline-block;
        background: #ef4444;
        color: white;
        font-size: 0.65rem;
        font-weight: 600;
        padding: 0.2rem 0.5rem;
        border-radius: 0.25rem;
    }

    .badge-pinned {
        display: inline-block;
        font-size: 0.75rem;
    }

    .badge-unit {
        display: inline-block;
        background: #dcfce7;
        color: #166534;
        font-size: 0.65rem;
        font-weight: 600;
        padding: 0.2rem 0.5rem;
        border-radius: 0.25rem;
    }

    .badge-org {
        display: inline-block;
        background: #e0e7ff;
        color: #4338ca;
        font-size: 0.65rem;
        font-weight: 600;
        padding: 0.2rem 0.5rem;
        border-radius: 0.25rem;
    }

    .announcement-meta {
        display: flex;
        gap: 0.75rem;
        margin-top: 0.5rem;
        font-size: 0.75rem;
        color: #94a3b8;
    }

    .announcement-meta .author {
        color: #64748b;
    }

    .announcement-item h3 {
        font-size: 0.9rem;
        font-weight: 600;
        color: #1e293b;
        margin: 0 0 0.5rem;
    }

    .announcement-item p {
        font-size: 0.85rem;
        color: #475569;
        margin: 0 0 0.5rem;
        line-height: 1.5;
    }

    .date {
        font-size: 0.75rem;
        color: #94a3b8;
    }

    /* Kas Card */
    .kas-summary {
        padding: 0.5rem 0;
    }

    .kas-item {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .kas-item.small {
        flex: 1;
    }

    .kas-label {
        font-size: 0.8rem;
        color: #64748b;
    }

    .kas-value {
        font-size: 1.25rem;
        font-weight: 700;
        color: #1e293b;
    }

    .kas-value.saldo {
        font-size: 1.75rem;
        color: #7c3aed;
    }

    .kas-value.positive {
        color: #16a34a;
        font-size: 1rem;
    }

    .kas-value.negative {
        color: #ef4444;
        font-size: 1rem;
    }

    .kas-divider {
        height: 1px;
        background: #e2e8f0;
        margin: 1rem 0;
    }

    .kas-row {
        display: flex;
        gap: 1rem;
    }

    /* Activity List */
    .activity-list {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }

    .activity-item {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.75rem;
        border-radius: 0.5rem;
        background: #f8fafc;
    }

    .activity-icon {
        width: 36px;
        height: 36px;
        border-radius: 0.5rem;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 1rem;
    }

    .activity-icon.iuran {
        background: #dcfce7;
    }
    .activity-icon.pendaftaran {
        background: #dbeafe;
    }
    .activity-icon.kegiatan {
        background: #fef3c7;
    }
    .activity-icon.profil {
        background: #f3e8ff;
    }

    .activity-content {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 0.125rem;
    }

    .activity-name {
        font-size: 0.85rem;
        font-weight: 600;
        color: #1e293b;
    }

    .activity-desc {
        font-size: 0.8rem;
        color: #64748b;
    }

    .activity-time {
        font-size: 0.75rem;
        color: #94a3b8;
        white-space: nowrap;
    }
</style>
