<script lang="ts">
    import { page } from "$app/stores";
    import { api } from "$lib";
    import { onMount } from "svelte";
    import { showToast } from "$core/components/Toast.svelte";
    import Modal from "$core/components/Modal.svelte";

    $: orgId = $page.params.orgId;
    $: unitId = $page.params.unitId;

    let currentDate = new Date();
    let currentYear = currentDate.getFullYear();
    let currentMonth = currentDate.getMonth() + 1;

    let events: any[] = [];
    let rts: any[] = [];
    let isLoading = true;
    let selectedDay: number | null = null;
    let showEventModal = false;
    let editingEvent: any = null;
    let showDetailModal = false;
    let viewingEvent: any = null;
    let showDayEventsModal = false;
    let selectedDayForEvents: number | null = null;
    let userRole = "";

    // Event form data
    let eventForm = {
        title: "",
        description: "",
        event_type: "activity",
        start_date: "",
        end_date: "",
        start_time: "",
        end_time: "",
        recurrence_type: "once",
        day_of_week: 0,
        is_all_day: false,
        location: "",
        color: "#7c3aed",
        rt_id: "",
    };

    const eventTypes = [
        { value: "schedule", label: "📚 Jadwal", color: "#3b82f6" },
        { value: "exam", label: "📝 Ujian", color: "#dc2626" },
        { value: "holiday", label: "🏖️ Libur", color: "#16a34a" },
        { value: "meeting", label: "👥 Rapat", color: "#2563eb" },
        { value: "activity", label: "🎉 Kegiatan", color: "#7c3aed" },
        { value: "other", label: "📌 Lainnya", color: "#64748b" },
    ];

    const recurrenceTypes = [
        { value: "once", label: "Sekali saja" },
        { value: "daily", label: "Harian" },
        { value: "weekly", label: "Mingguan" },
        { value: "monthly", label: "Bulanan" },
        { value: "semester", label: "Per Semester" },
    ];

    const dayNames = ["Min", "Sen", "Sel", "Rab", "Kam", "Jum", "Sab"];
    const monthNames = [
        "Januari",
        "Februari",
        "Maret",
        "April",
        "Mei",
        "Juni",
        "Juli",
        "Agustus",
        "September",
        "Oktober",
        "November",
        "Desember",
    ];

    onMount(async () => {
        await loadUserRole();
        await loadRTs();
        await loadCalendarEvents();
    });

    async function loadUserRole() {
        try {
            const res = await api.getMyMemberships();
            const perumahanMem = res.data?.unit_memberships?.find(
                (m: any) => m.unit_id === unitId,
            );
            userRole = perumahanMem?.role || "";
        } catch {
            userRole = "";
        }
    }

    async function loadRTs() {
        try {
            const res = await api.getRTs({
                unit_id: unitId,
                limit: 100,
            });
            rts = res.data || [];
        } catch {
            rts = [];
        }
    }

    async function loadCalendarEvents() {
        isLoading = true;
        try {
            const res = await api.getCalendarEntries({
                unit_id: unitId!,
                year: currentYear,
                month: currentMonth,
                limit: 500,
            });
            // Map entries to include date field for display
            events = (res.data || []).map((e: any) => ({
                ...e,
                date: e.date || e.start_date?.split("T")[0] || "",
            }));
        } catch (err) {
            console.error("Failed to load events:", err);
            events = [];
        } finally {
            isLoading = false;
        }
    }

    function getDaysInMonth(year: number, month: number): number {
        return new Date(year, month, 0).getDate();
    }

    function getFirstDayOfMonth(year: number, month: number): number {
        return new Date(year, month - 1, 1).getDay();
    }

    function getEventsForDay(day: number): any[] {
        const dateStr = `${currentYear}-${String(currentMonth).padStart(2, "0")}-${String(day).padStart(2, "0")}`;
        return events.filter((e) => e.date === dateStr);
    }

    function prevMonth() {
        if (currentMonth === 1) {
            currentMonth = 12;
            currentYear--;
        } else {
            currentMonth--;
        }
        loadCalendarEvents();
    }

    function nextMonth() {
        if (currentMonth === 12) {
            currentMonth = 1;
            currentYear++;
        } else {
            currentMonth++;
        }
        loadCalendarEvents();
    }

    function goToToday() {
        const today = new Date();
        currentYear = today.getFullYear();
        currentMonth = today.getMonth() + 1;
        loadCalendarEvents();
    }

    function openAddEventModal(day: number) {
        // Only admin/staff can add events
        if (userRole === "warga" || userRole === "pengurus") return;

        selectedDay = day;
        const dateStr = `${currentYear}-${String(currentMonth).padStart(2, "0")}-${String(day).padStart(2, "0")}`;
        editingEvent = null;
        eventForm = {
            title: "",
            description: "",
            event_type: "activity",
            start_date: dateStr,
            end_date: "",
            start_time: "08:00",
            end_time: "09:00",
            recurrence_type: "once",
            day_of_week: new Date(dateStr).getDay(),
            is_all_day: false,
            location: "",
            color: "#7c3aed",
            rt_id: "",
        };
        showEventModal = true;
    }

    function openEditEventModal(event: any) {
        // Penguruss and wargas only view, admin/staff can edit
        if (userRole === "warga" || userRole === "pengurus") {
            openEventDetailModal(event);
            return;
        }

        editingEvent = event;
        eventForm = {
            title: event.title,
            description: event.description || "",
            event_type: event.event_type,
            start_date: event.date,
            end_date: "",
            start_time: event.start_time || "",
            end_time: event.end_time || "",
            recurrence_type: "once",
            day_of_week: new Date(event.date).getDay(),
            is_all_day: event.is_all_day,
            location: event.location || "",
            color: event.color || "#7c3aed",
            rt_id: "",
        };
        showEventModal = true;
    }

    function openEventDetailModal(event: any) {
        viewingEvent = event;
        showDetailModal = true;
    }

    async function saveEvent() {
        if (!eventForm.title.trim()) {
            showToast("Judul event harus diisi", "error");
            return;
        }

        try {
            if (editingEvent) {
                await api.updateCalendarEntry(editingEvent.id, {
                    title: eventForm.title,
                    description: eventForm.description,
                    event_type: eventForm.event_type,
                    start_time: eventForm.start_time,
                    end_time: eventForm.end_time,
                    location: eventForm.location,
                    color: eventForm.color,
                });
                showToast("Event berhasil diupdate", "success");
            } else {
                await api.createCalendarEntry(
                    {
                        entry_type: "event",
                        title: eventForm.title,
                        description: eventForm.description,
                        event_type: eventForm.event_type,
                        start_date: eventForm.start_date,
                        end_date: eventForm.end_date || undefined,
                        start_time: eventForm.start_time || "00:00",
                        end_time: eventForm.end_time || "00:00",
                        recurrence_type: eventForm.recurrence_type,
                        day_of_week:
                            eventForm.recurrence_type === "weekly"
                                ? eventForm.day_of_week
                                : undefined,
                        is_all_day: eventForm.is_all_day,
                        location: eventForm.location,
                        color: eventForm.color,
                        rt_id: eventForm.rt_id || undefined,
                    },
                    unitId!,
                );
                showToast("Event berhasil dibuat", "success");
            }
            showEventModal = false;
            await loadCalendarEvents();
        } catch (err) {
            console.error("Failed to save event:", err);
            showToast("Gagal menyimpan event", "error");
        }
    }

    async function deleteEvent() {
        if (!editingEvent) return;
        if (!confirm("Hapus event ini?")) return;

        try {
            await api.deleteCalendarEntry(editingEvent.id);
            showToast("Event berhasil dihapus", "success");
            showEventModal = false;
            await loadCalendarEvents();
        } catch (err) {
            showToast("Gagal menghapus event", "error");
        }
    }

    function getEventTypeColor(type: string): string {
        return eventTypes.find((t) => t.value === type)?.color || "#64748b";
    }

    function isToday(day: number): boolean {
        const today = new Date();
        return (
            day === today.getDate() &&
            currentMonth === today.getMonth() + 1 &&
            currentYear === today.getFullYear()
        );
    }

    $: daysInMonth = getDaysInMonth(currentYear, currentMonth);
    $: firstDay = getFirstDayOfMonth(currentYear, currentMonth);
    $: calendarDays = Array.from({ length: daysInMonth }, (_, i) => i + 1);
</script>

<svelte:head>
    <title>Kalender - SeekOlah</title>
</svelte:head>

<div class="calendar-page">
    <div class="page-header">
        <div class="header-left">
            <h1>📅 Kalender RT</h1>
            <p>Lihat event dan jadwal kegiatan</p>
        </div>
        {#if userRole !== "warga" && userRole !== "pengurus"}
            <button
                class="btn-add"
                on:click={() => openAddEventModal(new Date().getDate())}
            >
                + Tambah Event
            </button>
        {/if}
    </div>

    <div class="calendar-nav">
        <button class="nav-btn" on:click={prevMonth}>‹</button>
        <div class="nav-center">
            <h2>{monthNames[currentMonth - 1]} {currentYear}</h2>
            <button class="today-btn" on:click={goToToday}>Hari Ini</button>
        </div>
        <button class="nav-btn" on:click={nextMonth}>›</button>
    </div>

    {#if isLoading}
        <div class="loading">Memuat kalender...</div>
    {:else}
        <div class="calendar-grid">
            {#each dayNames as dayName}
                <div class="day-header">{dayName}</div>
            {/each}

            {#each Array(firstDay) as _}
                <div class="day-cell empty"></div>
            {/each}

            {#each calendarDays as day}
                <div
                    class="day-cell"
                    rt:today={isToday(day)}
                    rt:clickable={userRole !== "warga" &&
                        userRole !== "pengurus"}
                    on:click={() =>
                        userRole !== "warga" &&
                        userRole !== "pengurus" &&
                        openAddEventModal(day)}
                    on:keypress={(e) =>
                        e.key === "Enter" &&
                        userRole !== "warga" &&
                        userRole !== "pengurus" &&
                        openAddEventModal(day)}
                    role={userRole !== "warga" && userRole !== "pengurus"
                        ? "button"
                        : undefined}
                    tabindex={userRole !== "warga" && userRole !== "pengurus"
                        ? 0
                        : -1}
                >
                    <span class="day-number">{day}</span>
                    <div class="day-events">
                        {#each getEventsForDay(day).slice(0, 3) as event}
                            <div
                                class="event-dot"
                                style="background-color: {getEventTypeColor(
                                    event.event_type,
                                )}"
                                title={event.title}
                                on:click|stopPropagation={() =>
                                    openEditEventModal(event)}
                                on:keypress|stopPropagation={(e) =>
                                    e.key === "Enter" &&
                                    openEditEventModal(event)}
                                role="button"
                                tabindex="0"
                            >
                                <span class="event-title">{event.title}</span>
                            </div>
                        {/each}
                        {#if getEventsForDay(day).length > 3}
                            <div
                                class="more-events"
                                on:click|stopPropagation={() => {
                                    selectedDayForEvents = day;
                                    showDayEventsModal = true;
                                }}
                                on:keypress|stopPropagation={(e) => {
                                    if (e.key === "Enter") {
                                        selectedDayForEvents = day;
                                        showDayEventsModal = true;
                                    }
                                }}
                                role="button"
                                tabindex="0"
                            >
                                +{getEventsForDay(day).length - 3} lagi
                            </div>
                        {/if}
                    </div>
                </div>
            {/each}
        </div>

        <div class="legend">
            {#each eventTypes as type}
                <span class="legend-item">
                    <span
                        class="legend-dot"
                        style="background-color: {type.color}"
                    ></span>
                    {type.label}
                </span>
            {/each}
        </div>
    {/if}
</div>

<!-- Event Modal -->
<Modal
    bind:isOpen={showEventModal}
    title={editingEvent ? "Edit Event" : "Tambah Event Baru"}
    size="md"
>
    <form on:submit|preventDefault={saveEvent}>
        <div class="form-group">
            <label>Judul Event *</label>
            <input
                type="text"
                bind:value={eventForm.title}
                placeholder="Contoh: Ujian Tengah Semester"
                required
            />
        </div>

        <div class="form-row">
            <div class="form-group">
                <label>Tipe Event</label>
                <select bind:value={eventForm.event_type}>
                    {#each eventTypes as type}
                        <option value={type.value}>{type.label}</option>
                    {/each}
                </select>
            </div>
            <div class="form-group">
                <label>Tim (opsional)</label>
                <select bind:value={eventForm.rt_id}>
                    <option value="">Semua RT</option>
                    {#each rts as cls}
                        <option value={cls.id}>{cls.name}</option>
                    {/each}
                </select>
            </div>
        </div>

        <div class="form-row">
            <div class="form-group">
                <label>Tanggal Mulai</label>
                <input type="date" bind:value={eventForm.start_date} required />
            </div>
            <div class="form-group">
                <label>Perulangan</label>
                <select bind:value={eventForm.recurrence_type}>
                    {#each recurrenceTypes as type}
                        <option value={type.value}>{type.label}</option>
                    {/each}
                </select>
            </div>
        </div>

        {#if eventForm.recurrence_type !== "once"}
            <div class="form-group">
                <label>Tanggal Berakhir</label>
                <input type="date" bind:value={eventForm.end_date} />
            </div>
        {/if}

        <div class="form-group">
            <label class="checkbox-label">
                <input type="checkbox" bind:checked={eventForm.is_all_day} />
                Seharian
            </label>
        </div>

        {#if !eventForm.is_all_day}
            <div class="form-row">
                <div class="form-group">
                    <label>Waktu Mulai</label>
                    <input type="time" bind:value={eventForm.start_time} />
                </div>
                <div class="form-group">
                    <label>Waktu Selesai</label>
                    <input type="time" bind:value={eventForm.end_time} />
                </div>
            </div>
        {/if}

        <div class="form-group">
            <label>Lokasi</label>
            <input
                type="text"
                bind:value={eventForm.location}
                placeholder="Ruang, Gedung, dll"
            />
        </div>

        <div class="form-group">
            <label>Deskripsi</label>
            <textarea
                bind:value={eventForm.description}
                rows="2"
                placeholder="Keterangan tambahan"
            ></textarea>
        </div>

        <div class="modal-actions">
            {#if editingEvent}
                <button type="button" class="btn-danger" on:click={deleteEvent}
                    >Hapus</button
                >
            {/if}
            <button
                type="button"
                class="btn-cancel"
                on:click={() => (showEventModal = false)}>Batal</button
            >
            <button type="submit" class="btn-save">💾 Simpan</button>
        </div>
    </form>
</Modal>

<!-- View Event Detail Modal (for penguruss/wargas) -->
<Modal bind:isOpen={showDetailModal} title="📅 Detail Event" size="sm">
    {#if viewingEvent}
        <div class="detail-content">
            <div class="detail-row">
                <span class="detail-label">Judul</span>
                <span class="detail-value">{viewingEvent.title}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Tipe</span>
                <span class="detail-value"
                    >{eventTypes.find(
                        (t) => t.value === viewingEvent.event_type,
                    )?.label || viewingEvent.event_type}</span
                >
            </div>
            <div class="detail-row">
                <span class="detail-label">Tanggal</span>
                <span class="detail-value"
                    >{new Date(viewingEvent.date).toLocaleDateString("id-ID", {
                        weekday: "long",
                        year: "numeric",
                        month: "long",
                        day: "numeric",
                    })}</span
                >
            </div>
            {#if viewingEvent.start_time}
                <div class="detail-row">
                    <span class="detail-label">Waktu</span>
                    <span class="detail-value"
                        >{viewingEvent.start_time}{viewingEvent.end_time
                            ? ` - ${viewingEvent.end_time}`
                            : ""}</span
                    >
                </div>
            {/if}
            {#if viewingEvent.location}
                <div class="detail-row">
                    <span class="detail-label">Lokasi</span>
                    <span class="detail-value">{viewingEvent.location}</span>
                </div>
            {/if}
            {#if viewingEvent.description}
                <div class="detail-row">
                    <span class="detail-label">Deskripsi</span>
                    <span class="detail-value">{viewingEvent.description}</span>
                </div>
            {/if}
            {#if viewingEvent.rt_name}
                <div class="detail-row">
                    <span class="detail-label">RT</span>
                    <span class="detail-value">{viewingEvent.rt_name}</span>
                </div>
            {/if}
        </div>
    {/if}
    <div class="modal-actions" style="justify-content: center;">
        <button class="btn-save" on:click={() => (showDetailModal = false)}
            >Tutup</button
        >
    </div>
</Modal>

<!-- Day Events List Modal -->
<Modal
    bind:isOpen={showDayEventsModal}
    title="📅 Kegiatan Tanggal {selectedDayForEvents} {monthNames[
        currentMonth - 1
    ]}"
    size="md"
>
    {#if selectedDayForEvents}
        <div class="day-events-list">
            {#each getEventsForDay(selectedDayForEvents) as event}
                <div
                    class="day-event-item"
                    on:click={() => {
                        showDayEventsModal = false;
                        openEditEventModal(event);
                    }}
                    on:keypress={(e) => {
                        if (e.key === "Enter") {
                            showDayEventsModal = false;
                            openEditEventModal(event);
                        }
                    }}
                    role="button"
                    tabindex="0"
                >
                    <span
                        class="event-color-dot"
                        style="background-color: {getEventTypeColor(
                            event.event_type,
                        )}"
                    ></span>
                    <div class="event-info">
                        <div class="event-item-title">{event.title}</div>
                        <div class="event-item-meta">
                            {eventTypes.find(
                                (t) => t.value === event.event_type,
                            )?.label || event.event_type}
                            {#if event.start_time}
                                • {event.start_time}{event.end_time
                                    ? ` - ${event.end_time}`
                                    : ""}{/if}
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
    <div class="modal-actions" style="justify-content: center;">
        <button class="btn-save" on:click={() => (showDayEventsModal = false)}
            >Tutup</button
        >
    </div>
</Modal>

<style>
    .calendar-page {
        padding: 1.5rem;
        max-width: 1200px;
        margin: 0 auto;
    }
    .page-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 1.5rem;
    }
    .header-left h1 {
        font-size: 1.5rem;
        font-weight: 700;
        color: #1e293b;
        margin: 0 0 0.25rem;
    }
    .header-left p {
        font-size: 0.875rem;
        color: #64748b;
        margin: 0;
    }

    .btn-add {
        padding: 0.625rem 1.25rem;
        background: linear-gradient(135deg, #7c3aed 0%, #a855f7 100%);
        color: white;
        border: none;
        border-radius: 0.5rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s;
    }
    .btn-add:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(124, 58, 237, 0.3);
    }

    .calendar-nav {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 1rem;
        padding: 1rem;
        background: white;
        border-radius: 0.75rem;
        border: 1px solid #e2e8f0;
    }
    .nav-btn {
        width: 40px;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: center;
        border: none;
        background: #f8fafc;
        border-radius: 0.5rem;
        font-size: 1.5rem;
        cursor: pointer;
        transition: background 0.2s;
    }
    .nav-btn:hover {
        background: #ede9fe;
    }
    .nav-center {
        display: flex;
        align-items: center;
        gap: 1rem;
    }
    .nav-center h2 {
        font-size: 1.25rem;
        font-weight: 600;
        color: #1e293b;
        margin: 0;
    }
    .today-btn {
        padding: 0.375rem 0.75rem;
        background: #ede9fe;
        color: #7c3aed;
        border: none;
        border-radius: 0.375rem;
        font-size: 0.75rem;
        font-weight: 600;
        cursor: pointer;
    }

    .loading {
        text-align: center;
        padding: 4rem;
        color: #64748b;
    }

    .calendar-grid {
        display: grid;
        grid-template-columns: repeat(7, 1fr);
        gap: 1px;
        background: #e2e8f0;
        border: 1px solid #e2e8f0;
        border-radius: 0.75rem;
        overflow: hidden;
    }
    .day-header {
        padding: 0.75rem;
        background: #f8fafc;
        text-align: center;
        font-size: 0.75rem;
        font-weight: 600;
        color: #64748b;
        text-transform: uppercase;
    }
    .day-cell {
        min-height: 100px;
        padding: 0.5rem;
        background: white;
        cursor: default;
        transition: background 0.2s;
    }
    .day-cell.clickable {
        cursor: pointer;
    }
    .day-cell.clickable:hover {
        background: #faf5ff;
    }
    .day-cell.empty {
        background: #f8fafc;
        cursor: default;
    }
    .day-cell.today {
        background: #ede9fe;
    }
    .day-cell.today .day-number {
        background: #7c3aed;
        color: white;
    }
    .day-number {
        display: inline-block;
        width: 24px;
        height: 24px;
        line-height: 24px;
        text-align: center;
        font-size: 0.75rem;
        font-weight: 600;
        color: #374151;
        border-radius: 50%;
    }
    .day-events {
        margin-top: 0.25rem;
        display: flex;
        flex-direction: column;
        gap: 2px;
    }
    .event-dot {
        padding: 2px 4px;
        border-radius: 3px;
        font-size: 0.625rem;
        color: white;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        cursor: pointer;
    }
    .event-dot:hover {
        opacity: 0.9;
    }
    .event-title {
        font-size: 0.625rem;
    }
    .more-events {
        font-size: 0.625rem;
        color: #64748b;
        padding: 2px 0;
    }

    .legend {
        display: flex;
        gap: 1.5rem;
        justify-content: center;
        padding: 1rem;
        margin-top: 1rem;
    }
    .legend-item {
        display: flex;
        align-items: center;
        gap: 0.375rem;
        font-size: 0.75rem;
        color: #64748b;
    }
    .legend-dot {
        width: 10px;
        height: 10px;
        border-radius: 50%;
    }

    /* Modal Form */
    .form-group {
        margin-bottom: 1rem;
    }
    .form-group label {
        display: block;
        font-size: 0.75rem;
        font-weight: 600;
        color: #374151;
        margin-bottom: 0.375rem;
        text-transform: uppercase;
    }
    .form-group input,
    .form-group select,
    .form-group textarea {
        width: 100%;
        padding: 0.625rem;
        border: 1px solid #d1d5db;
        border-radius: 0.5rem;
        font-size: 0.875rem;
        transition: border-color 0.2s;
        background: white;
        color: #1e293b;
    }
    .form-group input:focus,
    .form-group select:focus,
    .form-group textarea:focus {
        outline: none;
        border-color: #7c3aed;
        box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.1);
    }
    .form-row {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1rem;
    }
    .checkbox-label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        cursor: pointer;
    }
    .checkbox-label input {
        width: auto;
    }

    .modal-actions {
        display: flex;
        justify-content: flex-end;
        gap: 0.75rem;
        margin-top: 1.5rem;
        padding-top: 1rem;
        border-top: 1px solid #e2e8f0;
    }
    .btn-cancel {
        padding: 0.625rem 1rem;
        background: white;
        color: #64748b;
        border: 1px solid #e2e8f0;
        border-radius: 0.5rem;
        cursor: pointer;
    }
    .btn-save {
        padding: 0.625rem 1.25rem;
        background: linear-gradient(135deg, #7c3aed 0%, #a855f7 100%);
        color: white;
        border: none;
        border-radius: 0.5rem;
        font-weight: 600;
        cursor: pointer;
    }
    .btn-danger {
        padding: 0.625rem 1rem;
        background: #fef2f2;
        color: #dc2626;
        border: 1px solid #fecaca;
        border-radius: 0.5rem;
        cursor: pointer;
        margin-right: auto;
    }

    /* Tabs */
    .view-tabs {
        display: flex;
        gap: 0.5rem;
        margin-bottom: 1rem;
        background: #f1f5f9;
        padding: 0.25rem;
        border-radius: 0.5rem;
        width: fit-content;
    }
    .tab-btn {
        padding: 0.5rem 1rem;
        background: transparent;
        border: none;
        border-radius: 0.375rem;
        font-weight: 500;
        color: #64748b;
        cursor: pointer;
        transition: all 0.15s;
    }
    .tab-btn.active {
        background: white;
        color: #1e293b;
        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
    }
    .tab-btn:hover:not(.active) {
        color: #1e293b;
    }

    /* Schedule Grid */
    .schedule-grid-container {
        background: white;
        border: 1px solid #e2e8f0;
        border-radius: 0.75rem;
        overflow: hidden;
    }
    .schedule-info {
        padding: 0.75rem 1rem;
        background: #f8fafc;
        border-bottom: 1px solid #e2e8f0;
        font-size: 0.875rem;
        color: #64748b;
        display: flex;
        align-items: center;
        gap: 1rem;
        flex-wrap: wrap;
    }
    .rt-select {
        padding: 0.375rem 0.75rem;
        border: 1px solid #e2e8f0;
        border-radius: 0.375rem;
        font-size: 0.875rem;
        background: white;
        color: #1e293b;
    }
    .schedule-count {
        color: #94a3b8;
        font-size: 0.75rem;
    }
    .empty-state {
        padding: 3rem;
        text-align: center;
        color: #64748b;
    }
    .empty-state a {
        color: #7c3aed;
    }
    .schedule-table {
        display: flex;
        flex-direction: column;
    }
    .schedule-header {
        display: grid;
        grid-template-columns: 100px repeat(6, 1fr);
        background: #f1f5f9;
        border-bottom: 1px solid #e2e8f0;
    }
    .header-cell {
        padding: 0.75rem 0.5rem;
        font-weight: 600;
        font-size: 0.875rem;
        text-align: center;
        color: #475569;
    }
    .header-cell.time-col {
        text-align: left;
        padding-left: 1rem;
    }
    .schedule-row {
        display: grid;
        grid-template-columns: 100px repeat(6, 1fr);
        border-bottom: 1px solid #f1f5f9;
    }
    .schedule-row:last-child {
        border-bottom: none;
    }
    .time-cell {
        display: flex;
        flex-direction: column;
        padding: 0.5rem 0.75rem;
        background: #fafafa;
        border-right: 1px solid #e2e8f0;
    }
    .period-num {
        font-weight: 600;
        font-size: 0.75rem;
        color: #475569;
    }
    .period-time {
        font-size: 0.7rem;
        color: #94a3b8;
    }
    .schedule-cell {
        padding: 0.375rem;
        min-height: 60px;
        border-right: 1px solid #f1f5f9;
    }
    .schedule-cell:last-child {
        border-right: none;
    }
    .schedule-cell.has-schedule {
        background: #eff6ff;
    }
    .schedule-content {
        display: flex;
        flex-direction: column;
        gap: 0.125rem;
        padding: 0.375rem;
        background: #3b82f6;
        border-radius: 0.375rem;
        color: white;
        font-size: 0.75rem;
        height: 100%;
    }
    .kegiatan-name {
        font-weight: 600;
    }
    .rt-name {
        font-size: 0.65rem;
        opacity: 0.9;
    }
    .pengurus-name {
        font-size: 0.65rem;
        opacity: 0.8;
    }

    /* View Detail Modal */
    .detail-content {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
        margin-bottom: 1.5rem;
    }
    .detail-row {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }
    .detail-label {
        font-size: 0.75rem;
        font-weight: 600;
        color: #64748b;
        text-transform: uppercase;
    }
    .detail-value {
        font-size: 1rem;
        color: #1e293b;
        font-weight: 500;
    }

    /* More Events Button */
    .more-events {
        cursor: pointer;
    }
    .more-events:hover {
        text-decoration: underline;
        color: #6d28d9;
    }

    /* Day Events List Modal */
    .day-events-list {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        margin-bottom: 1rem;
        max-height: 400px;
        overflow-y: auto;
    }
    .day-event-item {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.75rem;
        background: #f8fafc;
        border-radius: 0.5rem;
        cursor: pointer;
        transition: background 0.2s;
    }
    .day-event-item:hover {
        background: #f1f5f9;
    }
    .event-color-dot {
        width: 12px;
        height: 12px;
        border-radius: 50%;
        flex-shrink: 0;
    }
    .event-info {
        flex: 1;
    }
    .event-item-title {
        font-weight: 600;
        color: #1e293b;
    }
    .event-item-meta {
        font-size: 0.75rem;
        color: #64748b;
    }
</style>
