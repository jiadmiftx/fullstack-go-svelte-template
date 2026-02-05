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

    let struktur: JabatanStruktur[] = [];
    let isLoading = true;
    let error = "";

    // Current user info
    let myProfile: any = null;

    onMount(async () => {
        await Promise.all([loadStruktur(), loadMyProfile()]);
    });

    async function loadStruktur() {
        isLoading = true;
        error = "";
        try {
            const res = await api.get(`/units/${unitId}/struktur`);
            struktur = res.data || [];
        } catch (err) {
            error =
                err instanceof Error ? err.message : "Gagal memuat struktur RT";
        } finally {
            isLoading = false;
        }
    }

    async function loadMyProfile() {
        try {
            const membershipRes = await api.getMyMemberships();
            const perumahanMembership =
                membershipRes.data?.unit_memberships?.find(
                    (m: any) => m.unit_id === unitId,
                );

            if (perumahanMembership) {
                myProfile = {
                    role: perumahanMembership.role || "warga",
                    joined_at: perumahanMembership.joined_at,
                    unit_member_id:
                        perumahanMembership.unit_member_id,
                };

                // Get warga profile for ID matching
                try {
                    const wargasRes = await api.getWargas({
                        unit_id: unitId,
                        limit: 100,
                    });
                    const wargas = wargasRes.data || [];
                    const myWarga = wargas.find(
                        (w: any) =>
                            w.unit_member_id ===
                            perumahanMembership.unit_member_id,
                    );
                    if (myWarga) {
                        myProfile = { ...myProfile, warga_id: myWarga.id };
                    }
                } catch (err) {
                    console.log("Could not load warga profile:", err);
                }
            }
        } catch (err) {
            console.error("Failed to load profile:", err);
        }
    }

    function getJabatanIcon(nama: string): string {
        if (nama.includes("Ketua")) return "👤";
        if (nama.includes("Wakil")) return "👥";
        if (nama.includes("Bendahara")) return "💰";
        if (nama.includes("Sekretaris")) return "📝";
        if (nama.includes("Keamanan")) return "🔒";
        if (nama.includes("Kebersihan")) return "🧹";
        return "📋";
    }

    function formatDate(dateStr: string | null): string {
        if (!dateStr) return "-";
        return new Date(dateStr).toLocaleDateString("id-ID", {
            day: "numeric",
            month: "long",
            year: "numeric",
        });
    }

    // Group by tipe and urutan
    $: struktural = struktur
        .filter((s) => s.tipe === "struktural")
        .sort((a, b) => a.urutan - b.urutan);
    $: seksi = struktur
        .filter((s) => s.tipe === "seksi")
        .sort((a, b) => a.urutan - b.urutan);

    // Check if user holds a jabatan
    $: myJabatanList = struktur.filter((j) =>
        j.anggota.some((a) => a.warga_profile_id === myProfile?.warga_id),
    );
</script>

<svelte:head>
    <title>Jabatan & Struktur RT</title>
</svelte:head>

<div class="page-container">
    <div class="page-header">
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
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
                <circle cx="9" cy="7" r="4"></circle>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
            </svg>
        </div>
        <div>
            <h1>Jabatan & Struktur RT</h1>
            <p class="subtitle">Struktur kepengurusan RT</p>
        </div>
    </div>

    {#if error}
        <div class="alert alert-error">{error}</div>
    {/if}

    {#if isLoading}
        <div class="loading">
            <div class="spinner"></div>
            <p>Memuat struktur...</p>
        </div>
    {:else}
        <!-- Jabatan Saya Summary -->
        {#if myJabatanList.length > 0}
            <div class="my-jabatan-banner">
                <span class="banner-label">Jabatan Anda:</span>
                {#each myJabatanList as jab}
                    <span class="banner-item"
                        >{getJabatanIcon(jab.nama)} {jab.nama}</span
                    >
                {/each}
            </div>
        {/if}

        <!-- Org Chart -->
        <div class="org-chart">
            <!-- Struktural Level (Top to Bottom) -->
            {#each struktural as jabatan, index}
                <div class="org-level">
                    <div class="level-title">{jabatan.nama}</div>
                    <div class="level-members">
                        {#if jabatan.anggota.length === 0}
                            <div class="member-card empty">
                                <div class="member-avatar">?</div>
                                <div class="member-info">
                                    <span class="member-name">Belum ada</span>
                                </div>
                            </div>
                        {:else}
                            {#each jabatan.anggota as anggota}
                                <div
                                    class="member-card"
                                    class:is-me={anggota.warga_profile_id ===
                                        myProfile?.warga_id}
                                >
                                    <div class="member-avatar">
                                        {(anggota.nama_warga || "X")
                                            .charAt(0)
                                            .toUpperCase()}
                                    </div>
                                    <div class="member-info">
                                        <span class="member-name"
                                            >{anggota.nama_warga ||
                                                "Nama tidak tersedia"}</span
                                        >
                                        <span class="member-role"
                                            >{jabatan.nama}</span
                                        >
                                    </div>
                                    {#if anggota.warga_profile_id === myProfile?.warga_id}
                                        <span class="me-badge">Anda</span>
                                    {/if}
                                </div>
                            {/each}
                        {/if}
                    </div>
                </div>
                {#if index < struktural.length - 1 || seksi.length > 0}
                    <div class="org-connector"></div>
                {/if}
            {/each}

            <!-- Seksi Level (Side by Side) -->
            {#if seksi.length > 0}
                <div class="org-level seksi-level">
                    <div class="level-title">Seksi-Seksi</div>
                    <div class="seksi-grid">
                        {#each seksi as jabatan}
                            <div class="seksi-card">
                                <div class="seksi-header">
                                    <span class="seksi-icon"
                                        >{getJabatanIcon(jabatan.nama)}</span
                                    >
                                    <span class="seksi-name"
                                        >{jabatan.nama}</span
                                    >
                                </div>
                                <div class="seksi-members">
                                    {#if jabatan.anggota.length === 0}
                                        <span class="empty-text"
                                            >Belum ada anggota</span
                                        >
                                    {:else}
                                        {#each jabatan.anggota as anggota}
                                            <div
                                                class="seksi-member"
                                                class:is-me={anggota.warga_profile_id ===
                                                    myProfile?.warga_id}
                                            >
                                                {#if anggota.is_ketua}
                                                    <span class="ketua-badge"
                                                        >K</span
                                                    >
                                                {/if}
                                                <span class="member-name"
                                                    >{anggota.nama_warga}</span
                                                >
                                                {#if anggota.warga_profile_id === myProfile?.warga_id}
                                                    <span class="me-badge"
                                                        >Anda</span
                                                    >
                                                {/if}
                                            </div>
                                        {/each}
                                    {/if}
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            {/if}
        </div>

        <div class="note-section">
            <p>
                ℹ️ Data ini bersifat <strong>read-only</strong>. Hubungi
                pengurus RT untuk perubahan.
            </p>
        </div>
    {/if}
</div>

<style>
    .page-container {
        padding: 1.5rem;
        max-width: 1000px;
        margin: 0 auto;
    }

    .page-header {
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

    .my-jabatan-banner {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        background: linear-gradient(135deg, #f0fdf4 0%, #dcfce7 100%);
        border: 2px solid #86efac;
        border-radius: 0.75rem;
        padding: 1rem 1.25rem;
        margin-bottom: 1.5rem;
        flex-wrap: wrap;
    }

    .banner-label {
        font-weight: 600;
        color: #16a34a;
    }

    .banner-item {
        background: white;
        padding: 0.375rem 0.75rem;
        border-radius: 0.5rem;
        font-weight: 600;
        color: #1e293b;
        border: 1px solid #86efac;
    }

    /* Org Chart */
    .org-chart {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 2rem 0;
    }

    .org-level {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 100%;
    }

    .level-title {
        font-size: 0.75rem;
        color: #64748b;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        font-weight: 600;
        margin-bottom: 0.75rem;
    }

    .level-members {
        display: flex;
        justify-content: center;
        gap: 1rem;
        flex-wrap: wrap;
    }

    .member-card {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        background: white;
        border: 2px solid #e2e8f0;
        border-radius: 1rem;
        padding: 1rem 1.25rem;
        min-width: 220px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
        position: relative;
    }

    .member-card.empty {
        border-style: dashed;
        opacity: 0.6;
    }

    .member-card.is-me {
        border-color: #7c3aed;
        background: linear-gradient(135deg, #faf5ff 0%, #f3e8ff 100%);
        box-shadow: 0 4px 12px rgba(124, 58, 237, 0.2);
    }

    .member-avatar {
        width: 48px;
        height: 48px;
        border-radius: 50%;
        background: linear-gradient(135deg, #7c3aed 0%, #a855f7 100%);
        color: white;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: 700;
        font-size: 1.25rem;
    }

    .member-info {
        display: flex;
        flex-direction: column;
    }

    .member-name {
        font-weight: 600;
        color: #1e293b;
        font-size: 0.95rem;
    }

    .member-role {
        color: #64748b;
        font-size: 0.8rem;
    }

    .me-badge {
        position: absolute;
        top: -8px;
        right: -8px;
        font-size: 0.65rem;
        background: #7c3aed;
        color: white;
        padding: 0.2rem 0.5rem;
        border-radius: 1rem;
        font-weight: 600;
    }

    .org-connector {
        width: 2px;
        height: 32px;
        background: linear-gradient(to bottom, #cbd5e1, #e2e8f0);
        margin: 0.5rem 0;
    }

    /* Seksi Level */
    .seksi-level {
        width: 100%;
        margin-top: 0.5rem;
    }

    .seksi-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 1rem;
        width: 100%;
    }

    .seksi-card {
        background: white;
        border: 1px solid #e2e8f0;
        border-left: 4px solid #8b5cf6;
        border-radius: 0.75rem;
        padding: 1rem;
    }

    .seksi-header {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin-bottom: 0.75rem;
        padding-bottom: 0.5rem;
        border-bottom: 1px solid #f1f5f9;
    }

    .seksi-icon {
        font-size: 1.25rem;
    }

    .seksi-name {
        font-weight: 600;
        color: #1e293b;
        font-size: 0.9rem;
    }

    .seksi-members {
        display: flex;
        flex-direction: column;
        gap: 0.375rem;
    }

    .empty-text {
        color: #94a3b8;
        font-size: 0.8rem;
        font-style: italic;
    }

    .seksi-member {
        display: flex;
        align-items: center;
        gap: 0.375rem;
        font-size: 0.85rem;
        color: #334155;
        padding: 0.25rem 0.5rem;
        border-radius: 0.25rem;
    }

    .seksi-member.is-me {
        background: #f0fdf4;
        color: #16a34a;
        font-weight: 600;
    }

    .ketua-badge {
        width: 18px;
        height: 18px;
        background: #fef3c7;
        color: #92400e;
        font-size: 0.6rem;
        font-weight: 700;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .note-section {
        background: #fef3c7;
        border: 1px solid #f59e0b;
        border-radius: 0.75rem;
        padding: 1rem;
        margin-top: 1.5rem;
    }

    .note-section p {
        margin: 0;
        color: #92400e;
        font-size: 0.9rem;
    }
</style>
