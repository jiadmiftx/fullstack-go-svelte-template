<script lang="ts">
    import { page } from "$app/stores";
    import { api } from "$lib";
    import { onMount } from "svelte";
    import Modal from "$core/components/Modal.svelte";
    import DataTable from "$core/components/DataTable.svelte";
    import { showToast } from "$core/components/Toast.svelte";

    interface Room {
        id: string;
        unit_id: string;
        code: string;
        name: string;
        type: string;
        building: string;
        floor: number;
        capacity: number;
        facilities: string;
        is_active: boolean;
    }

    $: unitId = $page.params.unitId;

    let rooms: Room[] = [];
    let isLoading = false;
    let error = "";
    let showModal = false;
    let editingRoom: Room | null = null;

    // Form fields
    let formCode = "";
    let formName = "";
    let formType = "rtroom";
    let formBuilding = "";
    let formFloor = 1;
    let formCapacity = 30;
    let formFacilities = "";

    const roomTypes = [
        { value: "rtroom", label: "Ruang RT" },
        { value: "laboratory", label: "Laboratorium" },
        { value: "library", label: "Perpustakaan" },
        { value: "hall", label: "Aula" },
        { value: "office", label: "Ruang Pengurus/TU" },
        { value: "other", label: "Lainnya" },
    ];

    onMount(async () => {
        await loadRooms();
    });

    async function loadRooms() {
        isLoading = true;
        error = "";
        try {
            const response = await api.getRooms({
                unit_id: unitId,
                limit: 100,
            });
            rooms = response.data || [];
        } catch (err) {
            error =
                err instanceof Error ? err.message : "Gagal memuat data ruang";
        } finally {
            isLoading = false;
        }
    }

    function openCreateModal() {
        editingRoom = null;
        formCode = "";
        formName = "";
        formType = "rtroom";
        formBuilding = "";
        formFloor = 1;
        formCapacity = 30;
        formFacilities = "";
        showModal = true;
    }

    function openEditModal(room: Room) {
        editingRoom = room;
        formCode = room.code;
        formName = room.name;
        formType = room.type || "rtroom";
        formBuilding = room.building || "";
        formFloor = room.floor || 1;
        formCapacity = room.capacity || 30;
        formFacilities = room.facilities || "";
        showModal = true;
    }

    async function handleSubmit() {
        error = "";
        try {
            if (editingRoom) {
                await api.updateRoom(editingRoom.id, {
                    code: formCode,
                    name: formName,
                    type: formType,
                    building: formBuilding,
                    floor: formFloor,
                    capacity: formCapacity,
                    facilities: formFacilities,
                });
                showToast("Ruang berhasil diupdate", "success");
            } else {
                await api.createRoom({
                    unit_id: unitId,
                    code: formCode,
                    name: formName,
                    type: formType,
                    building: formBuilding,
                    floor: formFloor,
                    capacity: formCapacity,
                    facilities: formFacilities,
                });
                showToast("Ruang berhasil ditambahkan", "success");
            }
            showModal = false;
            await loadRooms();
        } catch (err) {
            const msg =
                err instanceof Error
                    ? err.message
                    : "Gagal menyimpan data ruang";
            error = msg;
            showToast(msg, "error");
        }
    }

    async function handleDelete(room: Room) {
        if (!confirm(`Yakin ingin menghapus "${room.name}"?`)) return;
        try {
            await api.deleteRoom(room.id);
            showToast("Ruang berhasil dihapus", "success");
            await loadRooms();
        } catch (err) {
            showToast("Gagal menghapus ruang", "error");
        }
    }

    function getTypeLabel(type: string): string {
        const found = roomTypes.find((t) => t.value === type);
        return found ? found.label : type;
    }

    const columns = [
        { key: "code" as keyof Room, label: "Kode", sortable: true },
        { key: "name" as keyof Room, label: "Nama Ruang", sortable: true },
        {
            key: "type" as keyof Room,
            label: "Jenis",
            render: (v: string) => getTypeLabel(v),
        },
        {
            key: "building" as keyof Room,
            label: "Gedung",
            render: (v: string) => v || "-",
        },
        {
            key: "floor" as keyof Room,
            label: "Lantai",
            render: (v: number) => `Lt. ${v}`,
        },
        {
            key: "capacity" as keyof Room,
            label: "Kapasitas",
            render: (v: number) => `${v} orang`,
        },
    ];

    const actions = [
        { label: "Edit", variant: "primary" as const, onClick: openEditModal },
        { label: "Hapus", variant: "danger" as const, onClick: handleDelete },
    ];
</script>

<svelte:head>
    <title>Ruang RT</title>
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
                <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
                <polyline points="9 22 9 12 15 12 15 22"></polyline>
            </svg>
            <div>
                <h1>Ruang RT</h1>
                <p>Kelola data ruangan dan fasilitas RT</p>
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
                <line x1="12" y1="5" x2="12" y2="19"></line>
                <line x1="5" y1="12" x2="19" y2="12"></line>
            </svg>
            Tambah Ruang
        </button>
    </div>

    {#if error}
        <div class="error-alert">{error}</div>
    {/if}

    <div
        class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden"
    >
        {#if isLoading}
            <div class="loading">Memuat data ruang...</div>
        {:else}
            <DataTable
                data={rooms}
                {columns}
                {actions}
                emptyMessage="Belum ada data ruang. Klik 'Tambah Ruang' untuk menambahkan."
            />
        {/if}
    </div>
</div>

<Modal
    bind:isOpen={showModal}
    title={editingRoom ? "Edit Ruang" : "Tambah Ruang Baru"}
    size="lg"
>
    <form on:submit|preventDefault={handleSubmit} class="form-grid">
        <div class="form-group">
            <label for="code">Kode Ruang *</label>
            <input
                id="code"
                type="text"
                bind:value={formCode}
                class="input-field"
                placeholder="R101"
                required
            />
        </div>
        <div class="form-group">
            <label for="name">Nama Ruang *</label>
            <input
                id="name"
                type="text"
                bind:value={formName}
                class="input-field"
                placeholder="Ruang RT 1A"
                required
            />
        </div>
        <div class="form-group">
            <label for="type">Jenis Ruang</label>
            <select id="type" bind:value={formType} class="input-field">
                {#each roomTypes as t}
                    <option value={t.value}>{t.label}</option>
                {/each}
            </select>
        </div>
        <div class="form-group">
            <label for="building">Gedung</label>
            <input
                id="building"
                type="text"
                bind:value={formBuilding}
                class="input-field"
                placeholder="Gedung A"
            />
        </div>
        <div class="form-group">
            <label for="floor">Lantai</label>
            <input
                id="floor"
                type="number"
                bind:value={formFloor}
                class="input-field"
                min="1"
            />
        </div>
        <div class="form-group">
            <label for="capacity">Kapasitas</label>
            <input
                id="capacity"
                type="number"
                bind:value={formCapacity}
                class="input-field"
                min="1"
            />
        </div>
        <div class="form-group full">
            <label for="facilities">Fasilitas</label>
            <textarea
                id="facilities"
                bind:value={formFacilities}
                class="input-field"
                rows="2"
                placeholder="AC, Kegiatanor, Whiteboard"
            ></textarea>
        </div>
        <div class="form-actions full">
            <button
                type="button"
                on:click={() => (showModal = false)}
                class="btn-cancel">Batal</button
            >
            <button type="submit" class="btn-primary"
                >{editingRoom ? "Update" : "Simpan"}</button
            >
        </div>
    </form>
</Modal>

<style>
    .page-container {
        padding: 2rem;
        max-width: 1200px;
        margin: 0 auto;
    }
    .page-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 1.5rem;
    }
    .header-content {
        display: flex;
        align-items: center;
        gap: 1rem;
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
        margin: 0;
        font-size: 0.875rem;
    }
    .btn-add {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1.25rem;
        background: linear-gradient(135deg, #7c3aed, #6d28d9);
        color: white;
        border: none;
        border-radius: 0.5rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s;
    }
    .btn-add:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(124, 58, 237, 0.4);
    }
    .error-alert {
        background: rgba(239, 68, 68, 0.1);
        border: 1px solid rgba(239, 68, 68, 0.3);
        color: #dc2626;
        padding: 1rem;
        border-radius: 0.5rem;
        margin-bottom: 1rem;
    }
    .loading {
        padding: 3rem;
        text-align: center;
        color: #475569;
    }
    .form-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 1rem;
    }
    .form-group {
        display: flex;
        flex-direction: column;
        gap: 0.375rem;
    }
    .form-group.full {
        grid-column: span 2;
    }
    .form-group label {
        font-size: 0.875rem;
        font-weight: 500;
        color: #475569;
    }
    .input-field {
        padding: 0.625rem 0.875rem;
        background: #e2e8f0;
        border: 1px solid #e2e8f0;
        border-radius: 0.375rem;
        color: #1e293b;
        font-size: 0.875rem;
        transition: border-color 0.2s;
    }
    .input-field:focus {
        outline: none;
        border-color: #7c3aed;
    }
    textarea.input-field {
        resize: vertical;
    }
    .form-actions {
        display: flex;
        justify-content: flex-end;
        gap: 0.75rem;
        margin-top: 0.5rem;
    }
    .btn-cancel {
        padding: 0.625rem 1.25rem;
        background: #e2e8f0;
        border: none;
        border-radius: 0.375rem;
        color: #1e293b;
        font-weight: 500;
        cursor: pointer;
    }
    .btn-cancel:hover {
        background: #cbd5e1;
    }
    .btn-primary {
        padding: 0.625rem 1.5rem;
        background: linear-gradient(135deg, #7c3aed, #6d28d9);
        border: none;
        border-radius: 0.375rem;
        color: white;
        font-weight: 600;
        cursor: pointer;
    }
    .btn-primary:hover {
        background: linear-gradient(135deg, #6d28d9, #5b21b6);
    }
</style>
