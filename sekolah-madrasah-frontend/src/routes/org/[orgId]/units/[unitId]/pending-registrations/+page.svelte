<script lang="ts">
    import { page } from "$app/stores";
    import { api } from "$lib";
    import { onMount } from "svelte";
    import { showToast } from "$core/components/Toast.svelte";
    import Modal from "$core/components/Modal.svelte";

    interface PendingRegistration {
        id: string;
        user_name: string;
        email: string;
        phone: string;
        blok_rumah: string;
        nomor_rumah: string;
        approval_status: string;
        created_at: string;
    }

    interface HistoryItem {
        id: string;
        user_name: string;
        email: string;
        phone: string;
        blok_rumah: string;
        nomor_rumah: string;
        approval_status: string;
        rejected_reason?: string;
        created_at: string;
    }

    $: unitId = $page.params.unitId;
    $: orgId = $page.params.orgId;

    let registrations: PendingRegistration[] = [];
    let history: HistoryItem[] = [];
    let isLoading = false;
    let isLoadingHistory = false;
    let error = "";
    let showRejectModal = false;
    let selectedReg: PendingRegistration | null = null;
    let rejectReason = "";
    let showHistory = false;

    onMount(async () => {
        await Promise.all([loadPendingRegistrations(), loadHistory()]);
    });

    async function loadPendingRegistrations() {
        isLoading = true;
        error = "";
        try {
            const response = await api.get(
                `/units/${unitId}/pending-registrations`,
            );
            registrations = response.data || [];
        } catch (err) {
            error = err instanceof Error ? err.message : "Gagal memuat data";
        } finally {
            isLoading = false;
        }
    }

    async function handleApprove(reg: PendingRegistration) {
        if (!confirm(`Setujui pendaftaran ${reg.user_name}?`)) return;
        try {
            await api.post(
                `/units/${unitId}/registrations/${reg.id}/approve`,
                {},
            );
            showToast("Pendaftaran disetujui", "success");
            await Promise.all([loadPendingRegistrations(), loadHistory()]);
        } catch (err) {
            showToast(
                err instanceof Error ? err.message : "Gagal menyetujui",
                "error",
            );
        }
    }

    function openRejectModal(reg: PendingRegistration) {
        selectedReg = reg;
        rejectReason = "";
        showRejectModal = true;
    }

    async function handleReject() {
        if (!selectedReg) return;
        try {
            await api.post(
                `/units/${unitId}/registrations/${selectedReg.id}/reject`,
                { reason: rejectReason },
            );
            showToast("Pendaftaran ditolak", "success");
            showRejectModal = false;
            await Promise.all([loadPendingRegistrations(), loadHistory()]);
        } catch (err) {
            showToast(
                err instanceof Error ? err.message : "Gagal menolak",
                "error",
            );
        }
    }

    async function loadHistory() {
        isLoadingHistory = true;
        try {
            const response = await api.get(
                `/units/${unitId}/verification-history`,
            );
            history = response.data || [];
        } catch (err) {
            console.error("Failed to load history:", err);
        } finally {
            isLoadingHistory = false;
        }
    }

    function formatDate(dateStr: string): string {
        return new Date(dateStr).toLocaleDateString("id-ID", {
            day: "numeric",
            month: "short",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    }

    async function handleReReview(item: HistoryItem) {
        if (!confirm(`Kembalikan ${item.user_name} ke antrian review?`)) return;
        try {
            await api.post(
                `/units/${unitId}/registrations/${item.id}/re-review`,
                {},
            );
            showToast("Dipindahkan ke antrian review", "success");
            await Promise.all([loadPendingRegistrations(), loadHistory()]);
        } catch (err) {
            showToast(
                err instanceof Error ? err.message : "Gagal memproses",
                "error",
            );
        }
    }

    async function handleApproveFromHistory(item: HistoryItem) {
        if (!confirm(`Setujui pendaftaran ${item.user_name}?`)) return;
        try {
            await api.post(
                `/units/${unitId}/registrations/${item.id}/approve-from-history`,
                {},
            );
            showToast("Pendaftaran disetujui", "success");
            await Promise.all([loadPendingRegistrations(), loadHistory()]);
        } catch (err) {
            showToast(
                err instanceof Error ? err.message : "Gagal menyetujui",
                "error",
            );
        }
    }
</script>

<svelte:head>
    <title>Verifikasi Warga</title>
</svelte:head>

<div class="page-container">
    <div class="page-header">
        <div class="header-content">
            <svg
                xmlns="http://www.w3.org/2000/svg"
                width="28"
                height="28"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="header-icon"
            >
                <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
                <circle cx="9" cy="7" r="4"></circle>
                <path d="M22 21v-2a4 4 0 0 0-3-3.87"></path>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
            </svg>
            <div>
                <h1>Verifikasi Warga</h1>
                <p>Setujui atau tolak pendaftaran warga baru</p>
            </div>
        </div>
        <span class="badge-count">{registrations.length} Menunggu</span>
    </div>

    {#if error}
        <div class="error-alert">{error}</div>
    {/if}

    <div class="card">
        {#if isLoading}
            <div class="loading">Memuat data...</div>
        {:else if registrations.length === 0}
            <div class="empty-state">
                <div class="empty-icon">✅</div>
                <h3>Tidak Ada Pendaftaran</h3>
                <p>Semua pendaftaran sudah diverifikasi</p>
            </div>
        {:else}
            <div class="reg-list">
                {#each registrations as reg}
                    <div class="reg-item">
                        <div class="reg-avatar">
                            {reg.user_name?.charAt(0).toUpperCase() || "?"}
                        </div>
                        <div class="reg-info">
                            <div class="reg-name">{reg.user_name}</div>
                            <div class="reg-details">
                                <span>📧 {reg.email}</span>
                                {#if reg.phone}
                                    <span>📱 {reg.phone}</span>
                                {/if}
                            </div>
                            <div class="reg-location">
                                🏠 Blok {reg.blok_rumah} No. {reg.nomor_rumah}
                            </div>
                            <div class="reg-date">
                                Mendaftar: {formatDate(reg.created_at)}
                            </div>
                        </div>
                        <div class="reg-actions">
                            <button
                                class="btn-approve"
                                on:click={() => handleApprove(reg)}
                            >
                                ✓ Setujui
                            </button>
                            <button
                                class="btn-reject"
                                on:click={() => openRejectModal(reg)}
                            >
                                ✗ Tolak
                            </button>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </div>

    <!-- Verification History Section -->
    <div class="history-section">
        <button
            class="btn-toggle-history"
            on:click={() => (showHistory = !showHistory)}
        >
            {showHistory ? "▼" : "▶"} Riwayat Verifikasi ({history.length})
        </button>

        {#if showHistory}
            <div class="card history-card">
                {#if isLoadingHistory}
                    <div class="loading">Memuat riwayat...</div>
                {:else if history.length === 0}
                    <div class="empty-state" style="padding: 2rem;">
                        <p>Belum ada riwayat verifikasi</p>
                    </div>
                {:else}
                    <table class="history-table">
                        <thead>
                            <tr>
                                <th>Nama</th>
                                <th>Email</th>
                                <th>Alamat</th>
                                <th>Status</th>
                                <th>Keterangan</th>
                                <th>Aksi</th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each history as item}
                                <tr>
                                    <td class="name-cell">{item.user_name}</td>
                                    <td class="email-cell">{item.email}</td>
                                    <td
                                        >Blok {item.blok_rumah} No. {item.nomor_rumah}</td
                                    >
                                    <td>
                                        <span
                                            class="status-badge {item.approval_status}"
                                        >
                                            {item.approval_status === "approved"
                                                ? "✓ Disetujui"
                                                : "✗ Ditolak"}
                                        </span>
                                    </td>
                                    <td class="reason-cell"
                                        >{item.rejected_reason || "-"}</td
                                    >
                                    <td class="action-cell">
                                        {#if item.approval_status === "rejected"}
                                            <button
                                                class="btn-re-review"
                                                on:click={() =>
                                                    handleReReview(item)}
                                            >
                                                ↻ Review Ulang
                                            </button>
                                            <button
                                                class="btn-approve-small"
                                                on:click={() =>
                                                    handleApproveFromHistory(
                                                        item,
                                                    )}
                                            >
                                                ✓ Setujui
                                            </button>
                                        {:else}
                                            <span class="no-action">-</span>
                                        {/if}
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                {/if}
            </div>
        {/if}
    </div>
</div>

<Modal bind:isOpen={showRejectModal} title="Tolak Pendaftaran" size="sm">
    <div class="reject-content">
        <p>Tolak pendaftaran <strong>{selectedReg?.user_name}</strong>?</p>
        <div class="form-group">
            <label for="reason">Alasan Penolakan</label>
            <textarea
                id="reason"
                bind:value={rejectReason}
                placeholder="Opsional - berikan alasan penolakan"
                rows="3"
            ></textarea>
        </div>
    </div>
    <div slot="footer">
        <button on:click={() => (showRejectModal = false)} class="btn-secondary"
            >Batal</button
        >
        <button on:click={handleReject} class="btn-danger"
            >Tolak Pendaftaran</button
        >
    </div>
</Modal>

<style>
    .page-container {
        padding: 2rem;
        max-width: 900px;
        margin: 0 auto;
    }
    .page-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 2rem;
    }
    .header-content {
        display: flex;
        align-items: center;
        gap: 0.75rem;
    }
    .header-icon {
        color: #7c3aed;
    }
    .page-header h1 {
        font-size: 1.5rem;
        font-weight: 700;
        color: #1e293b;
        margin: 0;
    }
    .page-header p {
        color: #475569;
        margin: 0.25rem 0 0;
    }
    .badge-count {
        background: #fef3c7;
        color: #d97706;
        padding: 0.5rem 1rem;
        border-radius: 2rem;
        font-weight: 600;
        font-size: 0.875rem;
    }
    .card {
        background: #ffffff;
        border-radius: 1rem;
        border: 1px solid #e2e8f0;
        overflow: hidden;
    }
    .loading {
        padding: 3rem;
        text-align: center;
        color: #475569;
    }
    .empty-state {
        padding: 4rem 2rem;
        text-align: center;
    }
    .empty-icon {
        font-size: 3rem;
        margin-bottom: 1rem;
    }
    .empty-state h3 {
        font-size: 1.125rem;
        font-weight: 600;
        color: #1e293b;
        margin: 0 0 0.5rem;
    }
    .empty-state p {
        color: #64748b;
        margin: 0;
    }
    .error-alert {
        background: #fef2f2;
        color: #dc2626;
        padding: 1rem;
        border-radius: 0.5rem;
        margin-bottom: 1rem;
        border: 1px solid #fecaca;
    }

    .reg-list {
        display: flex;
        flex-direction: column;
    }
    .reg-item {
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 1.25rem 1.5rem;
        border-bottom: 1px solid #e2e8f0;
    }
    .reg-item:last-child {
        border-bottom: none;
    }
    .reg-avatar {
        width: 48px;
        height: 48px;
        background: linear-gradient(135deg, #7c3aed 0%, #a855f7 100%);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        font-weight: 700;
        font-size: 1.25rem;
        flex-shrink: 0;
    }
    .reg-info {
        flex: 1;
        min-width: 0;
    }
    .reg-name {
        font-weight: 600;
        color: #1e293b;
        font-size: 1rem;
        margin-bottom: 0.25rem;
    }
    .reg-details {
        display: flex;
        flex-wrap: wrap;
        gap: 0.75rem;
        font-size: 0.8rem;
        color: #475569;
        margin-bottom: 0.25rem;
    }
    .reg-location {
        font-size: 0.875rem;
        color: #7c3aed;
        font-weight: 500;
    }
    .reg-date {
        font-size: 0.75rem;
        color: #94a3b8;
        margin-top: 0.25rem;
    }
    .reg-actions {
        display: flex;
        gap: 0.5rem;
        flex-shrink: 0;
    }
    .btn-approve {
        padding: 0.5rem 1rem;
        background: #dcfce7;
        color: #16a34a;
        font-weight: 600;
        border: none;
        border-radius: 0.5rem;
        cursor: pointer;
        font-size: 0.8rem;
        transition: all 0.2s;
    }
    .btn-approve:hover {
        background: #bbf7d0;
    }
    .btn-reject {
        padding: 0.5rem 1rem;
        background: #fee2e2;
        color: #dc2626;
        font-weight: 600;
        border: none;
        border-radius: 0.5rem;
        cursor: pointer;
        font-size: 0.8rem;
        transition: all 0.2s;
    }
    .btn-reject:hover {
        background: #fecaca;
    }

    .reject-content p {
        margin: 0 0 1rem;
        color: #475569;
    }
    .form-group {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }
    .form-group label {
        font-size: 0.875rem;
        font-weight: 500;
        color: #374151;
    }
    .form-group textarea {
        padding: 0.75rem;
        border: 1px solid #d1d5db;
        border-radius: 0.5rem;
        font-size: 0.875rem;
        resize: vertical;
    }
    .form-group textarea:focus {
        outline: none;
        border-color: #7c3aed;
        box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.15);
    }
    .btn-secondary {
        padding: 0.625rem 1.25rem;
        background: #f1f5f9;
        color: #475569;
        font-weight: 500;
        border: none;
        border-radius: 0.5rem;
        cursor: pointer;
    }
    .btn-danger {
        padding: 0.625rem 1.25rem;
        background: #dc2626;
        color: white;
        font-weight: 600;
        border: none;
        border-radius: 0.5rem;
        cursor: pointer;
    }

    @media (max-width: 640px) {
        .reg-item {
            flex-wrap: wrap;
        }
        .reg-actions {
            width: 100%;
            margin-top: 0.75rem;
        }
    }

    /* History Section Styles */
    .history-section {
        margin-top: 2rem;
    }
    .btn-toggle-history {
        background: none;
        border: none;
        font-size: 0.9rem;
        font-weight: 600;
        color: #64748b;
        cursor: pointer;
        padding: 0.5rem 0;
        margin-bottom: 0.5rem;
    }
    .btn-toggle-history:hover {
        color: #7c3aed;
    }
    .history-card {
        margin-top: 0.5rem;
        overflow-x: auto;
    }
    .history-table {
        width: 100%;
        border-collapse: collapse;
        font-size: 0.875rem;
    }
    .history-table th {
        text-align: left;
        padding: 0.75rem 1rem;
        background: #f8fafc;
        border-bottom: 2px solid #e2e8f0;
        color: #475569;
        font-weight: 600;
    }
    .history-table td {
        padding: 0.75rem 1rem;
        border-bottom: 1px solid #e2e8f0;
        color: #334155;
    }
    .history-table tr:hover {
        background: #f8fafc;
    }
    .name-cell {
        font-weight: 500;
    }
    .email-cell {
        color: #64748b;
        font-size: 0.8rem;
    }
    .reason-cell {
        color: #94a3b8;
        font-size: 0.8rem;
        max-width: 150px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }
    .status-badge {
        padding: 0.25rem 0.625rem;
        border-radius: 1rem;
        font-size: 0.75rem;
        font-weight: 600;
    }
    .status-badge.approved {
        background: #dcfce7;
        color: #16a34a;
    }
    .status-badge.rejected {
        background: #fee2e2;
        color: #dc2626;
    }
    .action-cell {
        display: flex;
        gap: 0.5rem;
        flex-wrap: wrap;
    }
    .btn-re-review {
        padding: 0.25rem 0.5rem;
        background: #fef3c7;
        color: #d97706;
        font-weight: 500;
        border: none;
        border-radius: 0.375rem;
        cursor: pointer;
        font-size: 0.7rem;
        transition: all 0.2s;
        white-space: nowrap;
    }
    .btn-re-review:hover {
        background: #fde68a;
    }
    .btn-approve-small {
        padding: 0.25rem 0.5rem;
        background: #dcfce7;
        color: #16a34a;
        font-weight: 500;
        border: none;
        border-radius: 0.375rem;
        cursor: pointer;
        font-size: 0.7rem;
        transition: all 0.2s;
        white-space: nowrap;
    }
    .btn-approve-small:hover {
        background: #bbf7d0;
    }
    .no-action {
        color: #94a3b8;
    }
</style>
