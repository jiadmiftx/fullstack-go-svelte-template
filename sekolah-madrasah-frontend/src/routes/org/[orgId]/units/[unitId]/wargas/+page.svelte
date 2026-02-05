<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { api } from "$lib";
    import { onMount } from "svelte";
    import Modal from "$core/components/Modal.svelte";
    import DataTable from "$core/components/DataTable.svelte";
    import { showToast } from "$core/components/Toast.svelte";

    interface Warga {
        id: string;
        unit_member_id: string;
        unit_id: string;
        user_id: string;
        // Identitas
        nik?: string;
        gender?: string;
        agama?: string;
        pekerjaan?: string;
        no_whatsapp?: string;
        // Domisili
        blok_rumah?: string;
        nomor_rumah?: string;
        status_kepemilikan?: string;
        status_hunian?: string;
        // System
        status: string;
        user_name: string;
        email?: string;
        phone?: string;
        temp_password?: string;
    }

    $: unitId = $page.params.unitId;
    $: orgId = $page.params.orgId;

    let wargas: Warga[] = [];
    let isLoading = false;
    let error = "";
    let showModal = false;
    let showSuccessModal = false;
    let editingWarga: Warga | null = null;
    let generatedPassword = "";

    // Form fields - Identitas
    let formName = "";
    let formEmail = "";
    let formPhone = "";
    let formNIK = "";
    let formGender = "";
    let formAgama = "";
    let formPekerjaan = "";

    // Form fields - Domisili
    let formBlokRumah = "";
    let formNomorRumah = "";
    let formStatusKepemilikan = "";
    let formStatusHunian = "";

    onMount(async () => {
        await loadWargas();
    });

    async function loadWargas() {
        isLoading = true;
        error = "";
        try {
            const response = await api.getWargas({
                unit_id: unitId,
                limit: 100,
            });
            wargas = response.data || [];
        } catch (err) {
            error =
                err instanceof Error ? err.message : "Gagal memuat data warga";
        } finally {
            isLoading = false;
        }
    }

    function openCreateModal() {
        editingWarga = null;
        formName = "";
        formEmail = "";
        formPhone = "";
        formNIK = "";
        formGender = "";
        formAgama = "";
        formPekerjaan = "";
        formBlokRumah = "";
        formNomorRumah = "";
        formStatusKepemilikan = "";
        formStatusHunian = "";
        generatedPassword = "";
        showModal = true;
    }

    function openEditModal(warga: Warga) {
        editingWarga = warga;
        formName = warga.user_name || "";
        formEmail = warga.email || "";
        formPhone = warga.phone || "";
        formNIK = warga.nik || "";
        formGender = warga.gender || "";
        formAgama = warga.agama || "";
        formPekerjaan = warga.pekerjaan || "";
        formBlokRumah = warga.blok_rumah || "";
        formNomorRumah = warga.nomor_rumah || "";
        formStatusKepemilikan = warga.status_kepemilikan || "";
        formStatusHunian = warga.status_hunian || "";
        generatedPassword = warga.temp_password || "";
        showModal = true;
    }

    async function handleSubmit() {
        error = "";
        try {
            if (editingWarga) {
                await api.updateWarga(editingWarga.id, {
                    nik: formNIK || undefined,
                    gender: formGender || undefined,
                    agama: formAgama || undefined,
                    pekerjaan: formPekerjaan || undefined,
                    blok_rumah: formBlokRumah || undefined,
                    nomor_rumah: formNomorRumah || undefined,
                    status_kepemilikan: formStatusKepemilikan || undefined,
                    status_hunian: formStatusHunian || undefined,
                });
                showToast("Warga berhasil diupdate", "success");
                showModal = false;
            } else {
                const response = await api.registerWarga(unitId, {
                    name: formName,
                    email: formEmail,
                    phone: formPhone || undefined,
                    nik: formNIK || undefined,
                    gender: formGender || undefined,
                    agama: formAgama || undefined,
                    pekerjaan: formPekerjaan || undefined,
                    blok_rumah: formBlokRumah || undefined,
                    nomor_rumah: formNomorRumah || undefined,
                    status_kepemilikan: formStatusKepemilikan || undefined,
                    status_hunian: formStatusHunian || undefined,
                });
                generatedPassword = response.data?.generated_password || "";
                showModal = false;
                showSuccessModal = true;
            }
            await loadWargas();
        } catch (err) {
            const msg =
                err instanceof Error
                    ? err.message
                    : "Gagal menyimpan data warga";
            error = msg;
            showToast(msg, "error");
        }
    }

    async function handleDelete(warga: Warga) {
        if (!confirm(`Yakin ingin menghapus "${warga.user_name}"?`)) return;
        try {
            await api.deleteWarga(warga.id);
            showToast("Warga berhasil dihapus", "success");
            await loadWargas();
        } catch (err) {
            showToast("Gagal menghapus warga", "error");
        }
    }

    function copyPassword() {
        navigator.clipboard.writeText(generatedPassword);
        showToast("Password disalin ke clipboard", "success");
    }

    const columns = [
        { key: "user_name" as keyof Warga, label: "Nama", sortable: true },
        {
            key: "blok_rumah" as keyof Warga,
            label: "Blok",
            render: (v: string) => v || "-",
        },
        {
            key: "nomor_rumah" as keyof Warga,
            label: "No. Rumah",
            render: (v: string) => v || "-",
        },
        {
            key: "phone" as keyof Warga,
            label: "Telepon",
            render: (v: string) => v || "-",
        },
        {
            key: "status_kepemilikan" as keyof Warga,
            label: "Kepemilikan",
            render: (v: string) => v || "-",
        },
        {
            key: "status" as keyof Warga,
            label: "Status",
            render: (status: string) =>
                status === "active"
                    ? '<span class="badge success">Aktif</span>'
                    : '<span class="badge muted">Nonaktif</span>',
        },
    ];

    const actions = [
        {
            label: "Profile",
            variant: "secondary" as const,
            onClick: (s: Warga) =>
                goto(`/org/${orgId}/units/${unitId}/wargas/${s.id}`),
        },
        { label: "Edit", variant: "primary" as const, onClick: openEditModal },
        { label: "Hapus", variant: "danger" as const, onClick: handleDelete },
    ];
</script>

<svelte:head>
    <title>Manajemen Warga</title>
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
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
                <circle cx="9" cy="7" r="4"></circle>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
            </svg>
            <div>
                <h1>Manajemen Warga</h1>
                <p>Kelola data warga di perumahan ini</p>
            </div>
        </div>
        <button on:click={openCreateModal} class="btn-add">
            <svg
                xmlns="http://www.w3.org/2000/svg"
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
                <circle cx="9" cy="7" r="4"></circle>
                <line x1="19" y1="8" x2="19" y2="14"></line>
                <line x1="22" y1="11" x2="16" y2="11"></line>
            </svg>
            Tambah Warga
        </button>
    </div>

    {#if error}
        <div class="error-alert">{error}</div>
    {/if}

    <div class="table-card">
        {#if isLoading}
            <div class="loading">Memuat data warga...</div>
        {:else}
            <DataTable
                data={wargas}
                {columns}
                {actions}
                emptyMessage="Belum ada data warga. Klik 'Tambah Warga' untuk menambahkan."
            />
        {/if}
    </div>
</div>

<Modal
    bind:isOpen={showModal}
    title={editingWarga ? "Edit Warga" : "Tambah Warga Baru"}
    size="lg"
>
    <form on:submit|preventDefault={handleSubmit} class="form-sections">
        <!-- Kategori 1: Identitas -->
        <div class="form-section">
            <h3 class="section-title">📋 Data Identitas</h3>
            <div class="form-grid">
                <div class="form-group">
                    <label for="name">Nama Lengkap *</label>
                    <input
                        id="name"
                        type="text"
                        bind:value={formName}
                        class="form-input"
                        placeholder="Nama lengkap"
                        required
                        disabled={!!editingWarga}
                    />
                </div>
                <div class="form-group">
                    <label for="nik">NIK</label>
                    <input
                        id="nik"
                        type="text"
                        bind:value={formNIK}
                        class="form-input"
                        placeholder="Nomor Induk Kependudukan"
                        maxlength="16"
                    />
                </div>
                <div class="form-group">
                    <label for="email">Email {!editingWarga ? "*" : ""}</label>
                    <input
                        id="email"
                        type="email"
                        bind:value={formEmail}
                        class="form-input"
                        placeholder="email@contoh.com"
                        required={!editingWarga}
                        disabled={!!editingWarga}
                    />
                </div>
                <div class="form-group">
                    <label for="phone">No. WhatsApp</label>
                    <input
                        id="phone"
                        type="text"
                        bind:value={formPhone}
                        class="form-input"
                        placeholder="08xxxxxxxxxx"
                        disabled={!!editingWarga}
                    />
                </div>
                <div class="form-group">
                    <label for="gender">Jenis Kelamin</label>
                    <select
                        id="gender"
                        bind:value={formGender}
                        class="form-input"
                    >
                        <option value="">-- Pilih --</option>
                        <option value="M">Laki-laki</option>
                        <option value="F">Perempuan</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="agama">Agama</label>
                    <select
                        id="agama"
                        bind:value={formAgama}
                        class="form-input"
                    >
                        <option value="">-- Pilih --</option>
                        <option value="Islam">Islam</option>
                        <option value="Kristen">Kristen</option>
                        <option value="Katolik">Katolik</option>
                        <option value="Hindu">Hindu</option>
                        <option value="Buddha">Buddha</option>
                        <option value="Konghucu">Konghucu</option>
                    </select>
                </div>
                <div class="form-group full">
                    <label for="pekerjaan">Pekerjaan</label>
                    <input
                        id="pekerjaan"
                        type="text"
                        bind:value={formPekerjaan}
                        class="form-input"
                        placeholder="Contoh: PNS, Wiraswasta, dll"
                    />
                </div>
            </div>
        </div>

        <!-- Kategori 2: Domisili -->
        <div class="form-section">
            <h3 class="section-title">🏠 Data Domisili</h3>
            <div class="form-grid">
                <div class="form-group">
                    <label for="blok">Blok Rumah</label>
                    <input
                        id="blok"
                        type="text"
                        bind:value={formBlokRumah}
                        class="form-input"
                        placeholder="Contoh: A, B, C"
                    />
                </div>
                <div class="form-group">
                    <label for="nomor">Nomor Rumah</label>
                    <input
                        id="nomor"
                        type="text"
                        bind:value={formNomorRumah}
                        class="form-input"
                        placeholder="Contoh: 12, 15A"
                    />
                </div>
                <div class="form-group">
                    <label for="kepemilikan">Status Kepemilikan</label>
                    <select
                        id="kepemilikan"
                        bind:value={formStatusKepemilikan}
                        class="form-input"
                    >
                        <option value="">-- Pilih --</option>
                        <option value="milik">Milik Sendiri</option>
                        <option value="kontrak">Kontrak</option>
                        <option value="sewa">Sewa</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="hunian">Status Hunian</label>
                    <select
                        id="hunian"
                        bind:value={formStatusHunian}
                        class="form-input"
                    >
                        <option value="">-- Pilih --</option>
                        <option value="ditempati">Ditempati</option>
                        <option value="kosong">Kosong</option>
                        <option value="renovasi">Renovasi</option>
                    </select>
                </div>
            </div>
        </div>

        {#if !editingWarga}
            <div class="info-box">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <circle cx="12" cy="12" r="10"></circle>
                    <path d="M12 16v-4"></path>
                    <path d="M12 8h.01"></path>
                </svg>
                <span
                    >Password akan dibuat otomatis dan ditampilkan setelah
                    berhasil.</span
                >
            </div>
        {/if}

        {#if editingWarga}
            <div class="password-box">
                <label>🔑 Password Login (Dev Mode)</label>
                {#if generatedPassword}
                    <div class="password-row">
                        <input
                            type="text"
                            value={generatedPassword}
                            readonly
                            class="form-input password-field"
                        />
                        <button
                            type="button"
                            on:click={copyPassword}
                            class="btn-copy">📋 Copy</button
                        >
                    </div>
                {:else}
                    <div class="password-unavailable">
                        <span>⚠️ Password tidak tersedia</span>
                        <small>Password belum tersimpan di database</small>
                    </div>
                {/if}
            </div>
        {/if}
    </form>
    <div slot="footer">
        <button on:click={() => (showModal = false)} class="btn-secondary"
            >Batal</button
        >
        <button on:click={handleSubmit} class="btn-primary"
            >{editingWarga ? "Simpan" : "Daftar & Buat Akun"}</button
        >
    </div>
</Modal>

<Modal
    bind:isOpen={showSuccessModal}
    title="✅ Warga Berhasil Didaftarkan"
    size="md"
>
    <div class="success-content">
        <div class="success-icon">🎉</div>
        <p>
            <strong>{formName}</strong> telah terdaftar sebagai warga dan dapat login
            ke sistem.
        </p>
        <div class="credential-box">
            <div class="credential-item">
                <label>Email</label>
                <code>{formEmail}</code>
            </div>
            <div class="credential-item">
                <label>Password</label>
                <div class="password-display">
                    <code>{generatedPassword}</code>
                    <button on:click={copyPassword} class="btn-copy"
                        >📋 Copy</button
                    >
                </div>
            </div>
        </div>
        <div class="warning-box">
            ⚠️ Password ini juga tersedia di menu Edit.
        </div>
    </div>
    <div slot="footer">
        <button on:click={() => (showSuccessModal = false)} class="btn-primary"
            >Tutup</button
        >
    </div>
</Modal>

<style>
    .page-container {
        padding: 2rem;
        max-width: 1200px;
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
    .table-card {
        background: #ffffff;
        border-radius: 1rem;
        border: 1px solid #e2e8f0;
        overflow: hidden;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    }
    .btn-add {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1.25rem;
        background: linear-gradient(135deg, #7c3aed 0%, #6d28d9 100%);
        color: white;
        border: none;
        border-radius: 0.75rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s;
    }
    .btn-add:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 15px rgba(124, 58, 237, 0.4);
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
    .error-alert {
        background: #fef2f2;
        color: #dc2626;
        padding: 1rem;
        border-radius: 0.5rem;
        margin-bottom: 1rem;
        border: 1px solid #fecaca;
    }
    .loading {
        padding: 3rem;
        text-align: center;
        color: #475569;
    }

    .form-sections {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
    }
    .form-section {
        background: #f8fafc;
        border-radius: 0.75rem;
        padding: 1rem;
        border: 1px solid #e2e8f0;
    }
    .section-title {
        font-size: 0.95rem;
        font-weight: 600;
        color: #374151;
        margin: 0 0 1rem;
    }
    .form-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1rem;
    }
    .form-group {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }
    .form-group.full {
        grid-column: 1 / -1;
    }
    .form-group label {
        font-size: 0.875rem;
        font-weight: 500;
        color: #374151;
    }
    .form-input {
        width: 100%;
        padding: 0.625rem 0.875rem;
        background-color: #ffffff !important;
        border: 1px solid #d1d5db;
        border-radius: 0.5rem;
        color: #1f2937 !important;
        font-size: 0.875rem;
        transition: all 0.2s;
    }
    .form-input:focus {
        outline: none;
        border-color: #7c3aed;
        box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.15);
    }
    .form-input:disabled {
        background-color: #f3f4f6 !important;
        color: #6b7280 !important;
        cursor: not-allowed;
    }
    .form-input::placeholder {
        color: #9ca3af;
    }
    select.form-input {
        cursor: pointer;
    }

    .info-box {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        background: #ede9fe;
        padding: 0.75rem 1rem;
        border-radius: 0.5rem;
        color: #6d28d9;
        font-size: 0.875rem;
        border: 1px solid #ddd6fe;
    }
    .password-box {
        background: #f0fdf4;
        padding: 1rem;
        border-radius: 0.5rem;
        border: 1px solid #a7f3d0;
    }
    .password-box > label {
        color: #065f46 !important;
        font-weight: 600;
        font-size: 0.9rem;
        display: block;
        margin-bottom: 0.5rem;
    }
    .password-row {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }
    .password-field {
        flex: 1;
        font-family: monospace;
        font-size: 1rem !important;
        letter-spacing: 0.05em;
    }
    .password-unavailable {
        padding: 0.75rem;
        background: #fef3c7;
        border-radius: 0.375rem;
        border: 1px solid #fcd34d;
    }
    .password-unavailable span {
        display: block;
        color: #92400e;
        font-weight: 500;
    }
    .password-unavailable small {
        display: block;
        color: #a16207;
        font-size: 0.75rem;
        margin-top: 0.25rem;
    }
    .btn-copy {
        padding: 0.5rem 0.75rem;
        background: #7c3aed;
        border: none;
        border-radius: 0.375rem;
        color: white;
        cursor: pointer;
        font-size: 0.875rem;
        font-weight: 500;
        transition: all 0.2s;
        white-space: nowrap;
    }
    .btn-copy:hover {
        background: #6d28d9;
    }

    .success-content {
        text-align: center;
    }
    .success-icon {
        font-size: 3rem;
        margin-bottom: 1rem;
    }
    .credential-box {
        background: #f8fafc;
        border-radius: 0.75rem;
        padding: 1.25rem;
        margin: 1.5rem 0;
        text-align: left;
        border: 1px solid #e2e8f0;
    }
    .credential-item {
        margin-bottom: 1rem;
    }
    .credential-item:last-child {
        margin-bottom: 0;
    }
    .credential-item label {
        display: block;
        font-size: 0.75rem;
        color: #475569;
        margin-bottom: 0.25rem;
    }
    .credential-item code {
        display: block;
        background: #ffffff;
        padding: 0.5rem 0.75rem;
        border-radius: 0.375rem;
        color: #1e293b;
        font-family: monospace;
        font-size: 1rem;
        border: 1px solid #e2e8f0;
    }
    .password-display {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }
    .password-display code {
        flex: 1;
    }
    .warning-box {
        background: rgba(34, 197, 94, 0.1);
        padding: 0.75rem 1rem;
        border-radius: 0.5rem;
        color: #166534;
        font-size: 0.875rem;
        text-align: left;
    }

    :global(.badge.success) {
        background: rgba(16, 185, 129, 0.2);
        color: #10b981;
        padding: 0.25rem 0.5rem;
        border-radius: 0.25rem;
        font-size: 0.75rem;
    }
    :global(.badge.muted) {
        background: rgba(148, 163, 184, 0.2);
        color: #475569;
        padding: 0.25rem 0.5rem;
        border-radius: 0.25rem;
        font-size: 0.75rem;
    }
</style>
