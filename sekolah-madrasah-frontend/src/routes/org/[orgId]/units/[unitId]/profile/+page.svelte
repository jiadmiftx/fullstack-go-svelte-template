<script lang="ts">
    import { page } from "$app/stores";
    import { api } from "$lib";
    import { onMount } from "svelte";
    import { showToast } from "$core/components/Toast.svelte";

    $: orgId = $page.params.orgId;
    $: unitId = $page.params.unitId;

    let user: any = null;
    let membership: any = null;
    let wargaProfile: any = null;
    let isLoading = true;
    let isSaving = false;
    let isEditing = false;

    // Form data for user
    let userFormData = {
        full_name: "",
        phone: "",
    };

    // Form data for warga profile (new structure)
    let wargaFormData = {
        // Kategori 1: Identitas
        nik: "",
        gender: "",
        agama: "",
        pekerjaan: "",
        // Kategori 2: Domisili
        blok_rumah: "",
        nomor_rumah: "",
        status_kepemilikan: "",
        status_hunian: "",
        // Kategori 3: Keluarga
        jumlah_anggota_keluarga: 0,
        nama_kontak_darurat: "",
        no_kontak_darurat: "",
        no_plat_mobil: "",
        no_plat_motor: "",
        memiliki_art: false,
        // Kategori 4: Keuangan
        status_iuran: "",
        metode_pembayaran: "",
        keterangan_khusus: "",
    };

    onMount(async () => {
        await loadProfile();
    });

    async function loadProfile() {
        isLoading = true;
        try {
            const userRes = await api.getCurrentUser();
            user = userRes.data;

            const membershipRes = await api.getMyMemberships();
            const perumahanMem =
                membershipRes.data?.unit_memberships?.find(
                    (m: any) => m.unit_id === unitId,
                );
            membership = {
                role: perumahanMem?.role || "",
                perumahan_name: perumahanMem?.perumahan_name || "",
                unit_member_id: perumahanMem?.unit_member_id || "",
            };

            userFormData = {
                full_name: user?.full_name || "",
                phone: user?.phone || "",
            };

            if (membership.role === "warga" && membership.unit_member_id) {
                await loadWargaProfile(membership.unit_member_id);
            }
        } catch (err) {
            console.error("Failed to load profile:", err);
            showToast("Gagal memuat profil", "error");
        } finally {
            isLoading = false;
        }
    }

    async function loadWargaProfile(perumahanMemberId: string) {
        try {
            const wargasRes = await api.getWargas({
                unit_id: unitId,
                limit: 100,
            });
            const wargas = wargasRes.data || [];
            wargaProfile = wargas.find(
                (s: any) => s.unit_member_id === perumahanMemberId,
            );

            if (wargaProfile) {
                wargaFormData = {
                    nik: wargaProfile.nik || "",
                    gender: wargaProfile.gender || "",
                    agama: wargaProfile.agama || "",
                    pekerjaan: wargaProfile.pekerjaan || "",
                    blok_rumah: wargaProfile.blok_rumah || "",
                    nomor_rumah: wargaProfile.nomor_rumah || "",
                    status_kepemilikan: wargaProfile.status_kepemilikan || "",
                    status_hunian: wargaProfile.status_hunian || "",
                    jumlah_anggota_keluarga:
                        wargaProfile.jumlah_anggota_keluarga || 0,
                    nama_kontak_darurat: wargaProfile.nama_kontak_darurat || "",
                    no_kontak_darurat: wargaProfile.no_kontak_darurat || "",
                    no_plat_mobil: wargaProfile.no_plat_mobil || "",
                    no_plat_motor: wargaProfile.no_plat_motor || "",
                    memiliki_art: wargaProfile.memiliki_art || false,
                    status_iuran: wargaProfile.status_iuran || "",
                    metode_pembayaran: wargaProfile.metode_pembayaran || "",
                    keterangan_khusus: wargaProfile.keterangan_khusus || "",
                };
            }
        } catch (err) {
            console.error("Failed to load warga profile:", err);
        }
    }

    async function handleSave() {
        isSaving = true;
        try {
            await api.updateUser(user.id, {
                full_name: userFormData.full_name,
                phone: userFormData.phone,
            });
            user = { ...user, ...userFormData };

            if (membership.role === "warga" && wargaProfile) {
                await api.updateWarga(wargaProfile.id, wargaFormData);
                wargaProfile = { ...wargaProfile, ...wargaFormData };
            }

            isEditing = false;
            showToast("Profil berhasil disimpan", "success");
        } catch (err) {
            console.error("Failed to save profile:", err);
            showToast("Gagal menyimpan profil", "error");
        } finally {
            isSaving = false;
        }
    }

    function handleCancel() {
        userFormData = {
            full_name: user?.full_name || "",
            phone: user?.phone || "",
        };
        if (wargaProfile) {
            wargaFormData = {
                nik: wargaProfile.nik || "",
                gender: wargaProfile.gender || "",
                agama: wargaProfile.agama || "",
                pekerjaan: wargaProfile.pekerjaan || "",
                blok_rumah: wargaProfile.blok_rumah || "",
                nomor_rumah: wargaProfile.nomor_rumah || "",
                status_kepemilikan: wargaProfile.status_kepemilikan || "",
                status_hunian: wargaProfile.status_hunian || "",
                jumlah_anggota_keluarga:
                    wargaProfile.jumlah_anggota_keluarga || 0,
                nama_kontak_darurat: wargaProfile.nama_kontak_darurat || "",
                no_kontak_darurat: wargaProfile.no_kontak_darurat || "",
                no_plat_mobil: wargaProfile.no_plat_mobil || "",
                no_plat_motor: wargaProfile.no_plat_motor || "",
                memiliki_art: wargaProfile.memiliki_art || false,
                status_iuran: wargaProfile.status_iuran || "",
                metode_pembayaran: wargaProfile.metode_pembayaran || "",
                keterangan_khusus: wargaProfile.keterangan_khusus || "",
            };
        }
        isEditing = false;
    }

    function getRoleBadgeColor(role: string): string {
        switch (role) {
            case "admin":
                return "bg-purple";
            case "staff":
                return "bg-blue";
            case "pengurus":
                return "bg-green";
            case "warga":
                return "bg-yellow";
            default:
                return "bg-gray";
        }
    }

    function getGenderLabel(g: string): string {
        return g === "M" ? "Laki-laki" : g === "F" ? "Perempuan" : "-";
    }
</script>

<svelte:head>
    <title>Profil Saya - Perumahan</title>
</svelte:head>

<div class="py-8 px-4">
    <div class="max-w-3xl mx-auto">
        <div class="page-header">
            <h1 class="page-title">👤 Profil Saya</h1>
            <p class="page-subtitle">Kelola informasi akun Anda</p>
        </div>

        {#if isLoading}
            <div class="loading-state">
                <div class="spinner"></div>
                <p>Memuat profil...</p>
            </div>
        {:else if user}
            <div class="profile-card">
                <!-- Avatar Section -->
                <div class="avatar-section">
                    <div class="avatar-large">
                        {user.full_name?.charAt(0).toUpperCase() || "U"}
                    </div>
                    <div class="avatar-info">
                        <h2>{user.full_name || "User"}</h2>
                        <span
                            class="role-badge {getRoleBadgeColor(
                                membership?.role,
                            )}">{membership?.role || "Member"}</span
                        >
                        {#if wargaProfile?.blok_rumah && wargaProfile?.nomor_rumah}
                            <span class="location-badge"
                                >🏠 Blok {wargaProfile.blok_rumah}-{wargaProfile.nomor_rumah}</span
                            >
                        {/if}
                    </div>
                </div>

                <!-- Account Info Section -->
                <div class="section-title">📧 Informasi Akun</div>
                <div class="info-section">
                    <div class="info-row">
                        <div class="info-group">
                            <label>Email</label>
                            <div class="info-value readonly">{user.email}</div>
                        </div>
                        <div class="info-group">
                            <label>Perumahan</label>
                            <div class="info-value readonly">
                                {membership?.perumahan_name || "-"}
                            </div>
                        </div>
                    </div>
                    <div class="info-row">
                        <div class="info-group">
                            <label>Nama Lengkap</label>
                            {#if isEditing}
                                <input
                                    type="text"
                                    bind:value={userFormData.full_name}
                                    class="form-input"
                                    placeholder="Nama lengkap"
                                />
                            {:else}
                                <div class="info-value">
                                    {user.full_name || "-"}
                                </div>
                            {/if}
                        </div>
                        <div class="info-group">
                            <label>No. WhatsApp</label>
                            {#if isEditing}
                                <input
                                    type="tel"
                                    bind:value={userFormData.phone}
                                    class="form-input"
                                    placeholder="08xxxxxxxxxx"
                                />
                            {:else}
                                <div class="info-value">
                                    {user.phone || "-"}
                                </div>
                            {/if}
                        </div>
                    </div>
                </div>

                <!-- Warga Profile Section -->
                {#if membership.role === "warga" && wargaProfile}
                    <!-- Kategori 1: Identitas -->
                    <div class="section-title">📋 Data Identitas</div>
                    <div class="info-section">
                        <div class="info-row">
                            <div class="info-group">
                                <label>NIK</label>
                                {#if isEditing}
                                    <input
                                        type="text"
                                        bind:value={wargaFormData.nik}
                                        class="form-input"
                                        placeholder="Nomor Induk Kependudukan"
                                        maxlength="16"
                                    />
                                {:else}
                                    <div class="info-value">
                                        {wargaProfile.nik || "-"}
                                    </div>
                                {/if}
                            </div>
                            <div class="info-group">
                                <label>Jenis Kelamin</label>
                                {#if isEditing}
                                    <select
                                        bind:value={wargaFormData.gender}
                                        class="form-input"
                                    >
                                        <option value="">Pilih</option>
                                        <option value="M">Laki-laki</option>
                                        <option value="F">Perempuan</option>
                                    </select>
                                {:else}
                                    <div class="info-value">
                                        {getGenderLabel(wargaProfile.gender)}
                                    </div>
                                {/if}
                            </div>
                        </div>
                        <div class="info-row">
                            <div class="info-group">
                                <label>Agama</label>
                                {#if isEditing}
                                    <select
                                        bind:value={wargaFormData.agama}
                                        class="form-input"
                                    >
                                        <option value="">Pilih</option>
                                        <option value="Islam">Islam</option>
                                        <option value="Kristen">Kristen</option>
                                        <option value="Katolik">Katolik</option>
                                        <option value="Hindu">Hindu</option>
                                        <option value="Buddha">Buddha</option>
                                        <option value="Konghucu"
                                            >Konghucu</option
                                        >
                                    </select>
                                {:else}
                                    <div class="info-value">
                                        {wargaProfile.agama || "-"}
                                    </div>
                                {/if}
                            </div>
                            <div class="info-group">
                                <label>Pekerjaan</label>
                                {#if isEditing}
                                    <input
                                        type="text"
                                        bind:value={wargaFormData.pekerjaan}
                                        class="form-input"
                                        placeholder="Pekerjaan"
                                    />
                                {:else}
                                    <div class="info-value">
                                        {wargaProfile.pekerjaan || "-"}
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>

                    <!-- Kategori 2: Domisili -->
                    <div class="section-title">🏠 Data Domisili</div>
                    <div class="info-section">
                        <div class="info-row">
                            <div class="info-group">
                                <label>Blok Rumah</label>
                                {#if isEditing}
                                    <input
                                        type="text"
                                        bind:value={wargaFormData.blok_rumah}
                                        class="form-input"
                                        placeholder="A, B, C"
                                    />
                                {:else}
                                    <div class="info-value">
                                        {wargaProfile.blok_rumah || "-"}
                                    </div>
                                {/if}
                            </div>
                            <div class="info-group">
                                <label>Nomor Rumah</label>
                                {#if isEditing}
                                    <input
                                        type="text"
                                        bind:value={wargaFormData.nomor_rumah}
                                        class="form-input"
                                        placeholder="12, 15A"
                                    />
                                {:else}
                                    <div class="info-value">
                                        {wargaProfile.nomor_rumah || "-"}
                                    </div>
                                {/if}
                            </div>
                        </div>
                        <div class="info-row">
                            <div class="info-group">
                                <label>Status Kepemilikan</label>
                                {#if isEditing}
                                    <select
                                        bind:value={
                                            wargaFormData.status_kepemilikan
                                        }
                                        class="form-input"
                                    >
                                        <option value="">Pilih</option>
                                        <option value="milik"
                                            >Milik Sendiri</option
                                        >
                                        <option value="kontrak">Kontrak</option>
                                        <option value="sewa">Sewa</option>
                                    </select>
                                {:else}
                                    <div class="info-value">
                                        {wargaProfile.status_kepemilikan || "-"}
                                    </div>
                                {/if}
                            </div>
                            <div class="info-group">
                                <label>Status Hunian</label>
                                {#if isEditing}
                                    <select
                                        bind:value={wargaFormData.status_hunian}
                                        class="form-input"
                                    >
                                        <option value="">Pilih</option>
                                        <option value="ditempati"
                                            >Ditempati</option
                                        >
                                        <option value="kosong">Kosong</option>
                                        <option value="renovasi"
                                            >Renovasi</option
                                        >
                                    </select>
                                {:else}
                                    <div class="info-value">
                                        {wargaProfile.status_hunian || "-"}
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>

                    <!-- Kategori 3: Keluarga & Keamanan -->
                    <div class="section-title">👨‍👩‍👧 Keluarga & Keamanan</div>
                    <div class="info-section">
                        <div class="info-row">
                            <div class="info-group">
                                <label>Jumlah Anggota Keluarga</label>
                                {#if isEditing}
                                    <input
                                        type="number"
                                        bind:value={
                                            wargaFormData.jumlah_anggota_keluarga
                                        }
                                        class="form-input"
                                        min="1"
                                    />
                                {:else}
                                    <div class="info-value">
                                        {wargaProfile.jumlah_anggota_keluarga ||
                                            "-"} orang
                                    </div>
                                {/if}
                            </div>
                            <div class="info-group">
                                <label>Memiliki ART</label>
                                {#if isEditing}
                                    <select
                                        bind:value={wargaFormData.memiliki_art}
                                        class="form-input"
                                    >
                                        <option value={false}>Tidak</option>
                                        <option value={true}>Ya</option>
                                    </select>
                                {:else}
                                    <div class="info-value">
                                        {wargaProfile.memiliki_art
                                            ? "Ya"
                                            : "Tidak"}
                                    </div>
                                {/if}
                            </div>
                        </div>
                        <div class="info-row">
                            <div class="info-group">
                                <label>Nama Kontak Darurat</label>
                                {#if isEditing}
                                    <input
                                        type="text"
                                        bind:value={
                                            wargaFormData.nama_kontak_darurat
                                        }
                                        class="form-input"
                                        placeholder="Nama"
                                    />
                                {:else}
                                    <div class="info-value">
                                        {wargaProfile.nama_kontak_darurat ||
                                            "-"}
                                    </div>
                                {/if}
                            </div>
                            <div class="info-group">
                                <label>No. Kontak Darurat</label>
                                {#if isEditing}
                                    <input
                                        type="tel"
                                        bind:value={
                                            wargaFormData.no_kontak_darurat
                                        }
                                        class="form-input"
                                        placeholder="08xxx"
                                    />
                                {:else}
                                    <div class="info-value">
                                        {wargaProfile.no_kontak_darurat || "-"}
                                    </div>
                                {/if}
                            </div>
                        </div>
                        <div class="info-row">
                            <div class="info-group">
                                <label>No. Plat Mobil</label>
                                {#if isEditing}
                                    <input
                                        type="text"
                                        bind:value={wargaFormData.no_plat_mobil}
                                        class="form-input"
                                        placeholder="B 1234 ABC"
                                    />
                                {:else}
                                    <div class="info-value">
                                        {wargaProfile.no_plat_mobil || "-"}
                                    </div>
                                {/if}
                            </div>
                            <div class="info-group">
                                <label>No. Plat Motor</label>
                                {#if isEditing}
                                    <input
                                        type="text"
                                        bind:value={wargaFormData.no_plat_motor}
                                        class="form-input"
                                        placeholder="B 5678 XYZ"
                                    />
                                {:else}
                                    <div class="info-value">
                                        {wargaProfile.no_plat_motor || "-"}
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>

                    <!-- Kategori 4: Keuangan -->
                    <div class="section-title">💰 Data Keuangan</div>
                    <div class="info-section">
                        <div class="info-row">
                            <div class="info-group">
                                <label>Status Iuran</label>
                                {#if isEditing}
                                    <select
                                        bind:value={wargaFormData.status_iuran}
                                        class="form-input"
                                    >
                                        <option value="">Pilih</option>
                                        <option value="aktif">Aktif</option>
                                        <option value="non-aktif"
                                            >Non-Aktif</option
                                        >
                                    </select>
                                {:else}
                                    <div
                                        class="info-value status-badge {wargaProfile.status_iuran ===
                                        'aktif'
                                            ? 'active'
                                            : 'inactive'}"
                                    >
                                        {wargaProfile.status_iuran || "-"}
                                    </div>
                                {/if}
                            </div>
                            <div class="info-group">
                                <label>Metode Pembayaran</label>
                                {#if isEditing}
                                    <select
                                        bind:value={
                                            wargaFormData.metode_pembayaran
                                        }
                                        class="form-input"
                                    >
                                        <option value="">Pilih</option>
                                        <option value="transfer"
                                            >Transfer</option
                                        >
                                        <option value="cash">Cash</option>
                                    </select>
                                {:else}
                                    <div class="info-value">
                                        {wargaProfile.metode_pembayaran || "-"}
                                    </div>
                                {/if}
                            </div>
                        </div>
                        <div class="info-group full-width">
                            <label>Keterangan Khusus</label>
                            {#if isEditing}
                                <textarea
                                    bind:value={wargaFormData.keterangan_khusus}
                                    class="form-input"
                                    placeholder="Catatan khusus"
                                    rows="2"
                                ></textarea>
                            {:else}
                                <div class="info-value">
                                    {wargaProfile.keterangan_khusus || "-"}
                                </div>
                            {/if}
                        </div>
                    </div>
                {/if}

                <!-- Actions -->
                <div class="actions-section">
                    {#if isEditing}
                        <button
                            class="btn-secondary"
                            on:click={handleCancel}
                            disabled={isSaving}>Batal</button
                        >
                        <button
                            class="btn-primary"
                            on:click={handleSave}
                            disabled={isSaving}
                        >
                            {isSaving ? "Menyimpan..." : "💾 Simpan Perubahan"}
                        </button>
                    {:else}
                        <button
                            class="btn-primary"
                            on:click={() => (isEditing = true)}
                            >✏️ Edit Profil</button
                        >
                    {/if}
                </div>
            </div>
        {:else}
            <div class="error-state"><p>Gagal memuat profil</p></div>
        {/if}
    </div>
</div>

<style>
    .page-header {
        margin-bottom: 2rem;
    }
    .page-title {
        font-size: 1.5rem;
        font-weight: 700;
        color: #1e293b;
        margin: 0 0 0.25rem;
    }
    .page-subtitle {
        font-size: 0.875rem;
        color: #64748b;
        margin: 0;
    }
    .loading-state {
        text-align: center;
        padding: 4rem;
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
    .profile-card {
        background: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 1rem;
        padding: 2rem;
    }
    .avatar-section {
        display: flex;
        align-items: center;
        gap: 1.5rem;
        padding-bottom: 1.5rem;
        border-bottom: 1px solid #e2e8f0;
        margin-bottom: 1.5rem;
    }
    .avatar-large {
        width: 80px;
        height: 80px;
        border-radius: 50%;
        background: linear-gradient(135deg, #7c3aed 0%, #a855f7 100%);
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        font-weight: 700;
        font-size: 2rem;
        flex-shrink: 0;
    }
    .avatar-info h2 {
        font-size: 1.25rem;
        font-weight: 600;
        color: #1e293b;
        margin: 0 0 0.5rem;
    }
    .role-badge {
        display: inline-block;
        padding: 0.25rem 0.75rem;
        border-radius: 1rem;
        font-size: 0.75rem;
        font-weight: 600;
        text-transform: uppercase;
        margin-right: 0.5rem;
    }
    .role-badge.bg-purple {
        background: #ede9fe;
        color: #7c3aed;
    }
    .role-badge.bg-blue {
        background: #dbeafe;
        color: #2563eb;
    }
    .role-badge.bg-green {
        background: #dcfce7;
        color: #16a34a;
    }
    .role-badge.bg-yellow {
        background: #fef3c7;
        color: #d97706;
    }
    .role-badge.bg-gray {
        background: #f1f5f9;
        color: #475569;
    }
    .location-badge {
        display: inline-block;
        padding: 0.25rem 0.75rem;
        border-radius: 0.5rem;
        font-size: 0.75rem;
        font-weight: 500;
        background: #f0f9ff;
        color: #0284c7;
        border: 1px solid #bae6fd;
    }
    .section-title {
        font-size: 0.875rem;
        font-weight: 700;
        color: #1e293b;
        margin: 1.5rem 0 1rem;
        padding-top: 1rem;
        border-top: 1px solid #e2e8f0;
    }
    .info-section {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }
    .info-row {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1rem;
    }
    .info-group {
        display: flex;
        flex-direction: column;
        gap: 0.375rem;
    }
    .info-group.full-width {
        grid-column: span 2;
    }
    .info-group label {
        font-size: 0.75rem;
        font-weight: 600;
        color: #64748b;
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }
    .info-value {
        padding: 0.625rem 0.875rem;
        background: #f8fafc;
        border: 1px solid #e2e8f0;
        border-radius: 0.5rem;
        font-size: 0.875rem;
        color: #1e293b;
        min-height: 40px;
        display: flex;
        align-items: center;
    }
    .info-value.readonly {
        color: #64748b;
    }
    .info-value.status-badge.active {
        background: #dcfce7;
        color: #16a34a;
        border-color: #86efac;
    }
    .info-value.status-badge.inactive {
        background: #fee2e2;
        color: #dc2626;
        border-color: #fca5a5;
    }
    .form-input {
        padding: 0.625rem 0.875rem;
        border: 1px solid #d1d5db;
        border-radius: 0.5rem;
        font-size: 0.875rem;
        color: #1e293b;
        transition: all 0.2s;
        width: 100%;
        background: #fff;
    }
    .form-input:focus {
        outline: none;
        border-color: #7c3aed;
        box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.1);
    }
    textarea.form-input {
        resize: vertical;
        min-height: 60px;
    }
    .actions-section {
        display: flex;
        justify-content: flex-end;
        gap: 0.75rem;
        margin-top: 2rem;
        padding-top: 1.5rem;
        border-top: 1px solid #e2e8f0;
    }
    .btn-primary {
        padding: 0.75rem 1.5rem;
        background: linear-gradient(135deg, #7c3aed 0%, #a855f7 100%);
        color: white;
        border: none;
        border-radius: 0.5rem;
        font-size: 0.875rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s;
    }
    .btn-primary:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(124, 58, 237, 0.3);
    }
    .btn-primary:disabled {
        opacity: 0.6;
        cursor: not-allowed;
        transform: none;
    }
    .btn-secondary {
        padding: 0.75rem 1.5rem;
        background: #ffffff;
        color: #475569;
        border: 1px solid #e2e8f0;
        border-radius: 0.5rem;
        font-size: 0.875rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s;
    }
    .btn-secondary:hover {
        background: #f8fafc;
        border-color: #d1d5db;
    }
    .error-state {
        text-align: center;
        padding: 4rem;
        color: #dc2626;
    }
    @media (max-width: 640px) {
        .info-row {
            grid-template-columns: 1fr;
        }
        .info-group.full-width {
            grid-column: span 1;
        }
    }
</style>
