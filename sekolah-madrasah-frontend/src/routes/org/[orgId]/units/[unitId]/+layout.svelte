<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { isAuthenticated, api, selectedOrganization, auth } from "$lib";
    import { onMount } from "svelte";

    // Get IDs from URL
    $: orgId = $page.params.orgId;
    $: unitId = $page.params.unitId;

    // Perumahan data
    let perumahan: any = null;
    let userRole: string = "";
    let userName: string = "";
    let userJabatan: string = "";
    let isLoading = true;

    $: currentPath = $page.url.pathname;

    // Full menu for admin/staff
    const fullMenu = [
        {
            href: "dashboard",
            label: "Dashboard",
            iconType: "dashboard",
            roles: ["admin", "staff", "pengurus", "warga", "parent"],
        },
        {
            href: "wargas",
            label: "Warga",
            iconType: "warga",
            roles: ["admin", "staff"],
        },
        {
            href: "rooms",
            label: "Ruang",
            iconType: "home",
            roles: ["admin", "staff"],
        },
        {
            href: "kegiatans",
            label: "Kegiatan",
            iconType: "kegiatan",
            roles: ["admin", "staff"],
        },
        {
            href: "calendar",
            label: "Kalender",
            iconType: "calendar",
            roles: ["admin", "staff", "pengurus", "warga"],
        },
        {
            href: "iurans",
            label: "Iuran",
            iconType: "iuran",
            roles: ["admin", "staff", "pengurus", "warga"],
        },
        {
            href: "pengumuman",
            label: "Pengumuman",
            iconType: "announcement",
            roles: ["admin", "staff", "pengurus", "warga"],
        },
        {
            href: "reports",
            label: "Laporan",
            iconType: "report",
            roles: ["admin", "staff", "pengurus"],
        },
        {
            href: "pending-registrations",
            label: "Verifikasi Warga",
            iconType: "approval",
            roles: ["admin", "staff"],
        },
        {
            href: "struktur",
            label: "Struktur RT",
            iconType: "struktur",
            roles: ["admin", "staff", "pengurus"],
        },
        {
            href: "jabatan",
            label: "Struktur RT",
            iconType: "struktur",
            roles: ["warga"],
        },
    ];

    // Filter menu based on user role
    $: perumahanMenu =
        perumahan && userRole
            ? fullMenu
                  .filter((item) => item.roles.includes(userRole))
                  .map((item) => ({
                      ...item,
                      href: `/org/${orgId}/units/${unitId}/${item.href}`,
                  }))
            : [];

    // Check if user can go back to org level (admin/staff only)
    $: canGoBackToOrg =
        userRole === "admin" || userRole === "staff" || !userRole;

    onMount(async () => {
        if (!$isAuthenticated) {
            goto("/auth/login");
            return;
        }
        await loadPerumahan();
        await loadUserRole();
        await loadCurrentUser(); // Run last to ensure userName is from current user
    });

    async function loadCurrentUser() {
        try {
            const response = await api.getCurrentUser();
            if (response.data) {
                userName =
                    response.data.full_name || response.data.email || "User";
            }
        } catch (err) {
            console.error("Failed to load current user:", err);
        }
    }

    function handleLogout() {
        auth.logout();
        selectedOrganization.clear();
        goto("/auth/login");
    }

    async function loadPerumahan() {
        if (!unitId) return;
        try {
            const response = await api.getPerumahan(unitId);
            perumahan = response.data;
        } catch (err) {
            console.error("Failed to load perumahan:", err);
        } finally {
            isLoading = false;
        }
    }

    async function loadUserJabatan(perumahanMemberId: string) {
        try {
            console.log(
                "Layout: Fetching warga profile for member:",
                perumahanMemberId,
            );

            // Get user's warga profile by unit_member_id
            const profileRes: any = await api.get(
                `/warga-profiles?unit_member_id=${perumahanMemberId}&limit=1`,
            );
            if (!profileRes.data || profileRes.data.length === 0) {
                console.log("No warga profile found");
                return;
            }
            const wargaProfileId = profileRes.data[0].id;
            console.log("Warga Profile ID:", wargaProfileId);

            // Fetch struktur RT which includes member assignments
            const strukturRes: any = await api.get(
                `/units/${unitId}/struktur`,
            );
            console.log("Struktur Response:", strukturRes);

            // Priority order (higher = more important)
            const jabatanPriority: { [key: string]: number } = {
                "ketua rt": 4,
                "wakil ketua rt": 3,
                sekretaris: 2,
                bendahara: 1,
            };

            let foundJabatan = "";
            let highestPriority = -1;

            if (strukturRes.data) {
                for (const jabatan of strukturRes.data) {
                    console.log(
                        "Checking jabatan:",
                        jabatan.nama,
                        "Members:",
                        jabatan.anggota,
                    );

                    if (jabatan.anggota && jabatan.anggota.length > 0) {
                        // Check if current user is assigned to this jabatan
                        for (const anggota of jabatan.anggota) {
                            console.log(
                                "Checking anggota:",
                                anggota.warga_profile_id,
                                "vs",
                                wargaProfileId,
                            );

                            if (anggota.warga_profile_id === wargaProfileId) {
                                const jabatanName = jabatan.nama || "";
                                const jabatanKey = jabatanName.toLowerCase();
                                const priority =
                                    jabatanPriority[jabatanKey] || 0;

                                console.log(
                                    `Found jabatan: ${jabatanName} with priority ${priority}`,
                                );

                                if (priority > highestPriority) {
                                    highestPriority = priority;
                                    foundJabatan = jabatanName;
                                }
                            }
                        }
                    }
                }
            }

            if (foundJabatan) {
                userJabatan = foundJabatan;
                console.log("Final jabatan:", userJabatan);
            } else {
                console.log("No matching jabatan found");
            }
        } catch (err) {
            console.error("Failed to load user jabatan:", err);
        }
    }

    async function loadUserRole() {
        try {
            const response = await api.getMyMemberships();
            const memberships = response.data;

            // Check if user is super admin (can access any perumahan)
            if (memberships.is_super_admin) {
                userRole = "admin";
                userName = (memberships as any).user?.name || "Super Admin";
                return;
            }

            // Check if user has org-level membership (can access any perumahan in org)
            const hasOrgAccess = memberships.organization_memberships?.some(
                (m: any) => m.org_id === orgId,
            );

            // Check if user has perumahan-level membership
            const perumahanMembership = memberships.unit_memberships?.find(
                (m: any) => m.unit_id === unitId,
            );

            if (perumahanMembership) {
                // User has access to this perumahan
                userRole = perumahanMembership.role;
                userName =
                    (memberships as any).user?.name ||
                    (perumahanMembership as any).user_name ||
                    "User";

                // Load jabatan for warga/pengurus
                if (userRole === "warga" || userRole === "pengurus") {
                    if (perumahanMembership.unit_member_id) {
                        await loadUserJabatan(
                            perumahanMembership.unit_member_id,
                        );
                    }
                }
            } else if (hasOrgAccess) {
                // User has org-level access
                userRole = "staff";
            } else {
                // User is pengurus/warga but trying to access wrong perumahan
                // Find their assigned perumahan and redirect
                const userPerumahans = memberships.unit_memberships || [];
                if (userPerumahans.length > 0) {
                    const assignedPerumahan = userPerumahans[0];
                    console.warn(
                        "Access denied: Redirecting to assigned perumahan",
                    );
                    goto(
                        `/org/${assignedPerumahan.org_id}/units/${assignedPerumahan.unit_id}/dashboard`,
                    );
                    return;
                }
                // No perumahan access at all - redirect to home
                console.warn("Access denied: No perumahan membership found");
                goto("/dashboard");
                return;
            }
        } catch (err) {
            console.error("Failed to load user role:", err);
            userRole = "staff"; // Fallback
        }
    }

    function isActive(href: string): boolean {
        if (href.endsWith("/dashboard")) {
            return (
                currentPath === href ||
                currentPath === `/org/${orgId}/units/${unitId}`
            );
        }
        return currentPath.startsWith(href);
    }
</script>

{#if $isAuthenticated}
    <div class="perumahan-layout dark">
        <!-- Perumahan Sidebar -->
        <aside class="perumahan-sidebar">
            <!-- Back to Perumahans (only for admin/staff) -->
            {#if canGoBackToOrg}
                <a href="/org/{orgId}/units" class="back-link">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        <line x1="19" y1="12" x2="5" y2="12"></line>
                        <polyline points="12 19 5 12 12 5"></polyline>
                    </svg>
                    Kembali ke Daftar RT
                </a>
            {/if}

            <!-- Perumahan Header -->
            {#if perumahan}
                <div class="perumahan-header">
                    <div class="perumahan-logo">
                        {#if perumahan.logo}
                            <img src={perumahan.logo} alt={perumahan.name} />
                        {:else}
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                width="24"
                                height="24"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="1.5"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <path d="M22 10v6M2 10l10-5 10 5-10 5z"></path>
                                <path d="M6 12v5c3 3 9 3 12 0v-5"></path>
                            </svg>
                        {/if}
                    </div>
                    <div class="perumahan-info">
                        <h2>{perumahan.name}</h2>
                        <span class="perumahan-type">{perumahan.type}</span>
                    </div>
                </div>
            {:else if isLoading}
                <div class="perumahan-header loading">
                    <div class="skeleton-avatar"></div>
                    <div class="skeleton-text"></div>
                </div>
            {/if}

            <!-- Navigation -->
            <nav class="perumahan-nav">
                {#each perumahanMenu as item}
                    <a
                        href={item.href}
                        class="nav-item"
                        rt:active={isActive(item.href)}
                    >
                        <span class="nav-icon">
                            {#if item.iconType === "profile"}
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
                                    <path
                                        d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"
                                    ></path>
                                    <circle cx="12" cy="7" r="4"></circle>
                                </svg>
                            {:else if item.iconType === "dashboard"}
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
                                    <rect
                                        x="3"
                                        y="3"
                                        width="7"
                                        height="9"
                                        rx="1"
                                    ></rect>
                                    <rect
                                        x="14"
                                        y="3"
                                        width="7"
                                        height="5"
                                        rx="1"
                                    ></rect>
                                    <rect
                                        x="14"
                                        y="12"
                                        width="7"
                                        height="9"
                                        rx="1"
                                    ></rect>
                                    <rect
                                        x="3"
                                        y="16"
                                        width="7"
                                        height="5"
                                        rx="1"
                                    ></rect>
                                </svg>
                            {:else if item.iconType === "pengurus"}
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
                                    <path
                                        d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"
                                    ></path>
                                    <circle cx="9" cy="7" r="4"></circle>
                                    <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
                                    <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
                                </svg>
                            {:else if item.iconType === "warga"}
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
                                    <path d="M22 10v6M2 10l10-5 10 5-10 5z"
                                    ></path>
                                    <path d="M6 12v5c3 3 9 3 12 0v-5"></path>
                                </svg>
                            {:else if item.iconType === "rt"}
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
                                    <path
                                        d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"
                                    ></path>
                                    <polyline points="9 22 9 12 15 12 15 22"
                                    ></polyline>
                                </svg>
                            {:else if item.iconType === "home"}
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
                                    <rect x="3" y="3" width="7" height="9"
                                    ></rect>
                                    <rect x="14" y="3" width="7" height="9"
                                    ></rect>
                                    <rect x="3" y="14" width="7" height="7"
                                    ></rect>
                                    <rect x="14" y="14" width="7" height="7"
                                    ></rect>
                                </svg>
                            {:else if item.iconType === "members"}
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
                                    <path
                                        d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"
                                    ></path>
                                    <circle cx="9" cy="7" r="4"></circle>
                                    <path d="M22 21v-2a4 4 0 0 0-3-3.87"></path>
                                    <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
                                </svg>
                            {:else if item.iconType === "kegiatan"}
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
                                    <path
                                        d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"
                                    ></path>
                                </svg>
                            {:else if item.iconType === "schedule"}
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
                                    <rect
                                        width="18"
                                        height="18"
                                        x="3"
                                        y="4"
                                        rx="2"
                                        ry="2"
                                    ></rect>
                                    <line x1="16" x2="16" y1="2" y2="6"></line>
                                    <line x1="8" x2="8" y1="2" y2="6"></line>
                                    <line x1="3" x2="21" y1="10" y2="10"></line>
                                </svg>
                            {:else if item.iconType === "iuran"}
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
                                    <line x1="18" x2="18" y1="20" y2="10"
                                    ></line>
                                    <line x1="12" x2="12" y1="20" y2="4"></line>
                                    <line x1="6" x2="6" y1="20" y2="14"></line>
                                </svg>
                            {:else if item.iconType === "calendar"}
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
                                    <rect
                                        x="3"
                                        y="4"
                                        width="18"
                                        height="18"
                                        rx="2"
                                        ry="2"
                                    ></rect>
                                    <line x1="16" x2="16" y1="2" y2="6"></line>
                                    <line x1="8" x2="8" y1="2" y2="6"></line>
                                    <line x1="3" x2="21" y1="10" y2="10"></line>
                                    <path d="M8 14h.01"></path>
                                    <path d="M12 14h.01"></path>
                                    <path d="M16 14h.01"></path>
                                    <path d="M8 18h.01"></path>
                                    <path d="M12 18h.01"></path>
                                    <path d="M16 18h.01"></path>
                                </svg>
                            {:else if item.iconType === "attendance"}
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
                                    <path d="M9 11l3 3L22 4"></path>
                                    <path
                                        d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"
                                    ></path>
                                </svg>
                            {:else if item.iconType === "report"}
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
                                    <path
                                        d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
                                    ></path>
                                    <polyline points="14 2 14 8 20 8"
                                    ></polyline>
                                    <line x1="16" x2="8" y1="13" y2="13"></line>
                                    <line x1="16" x2="8" y1="17" y2="17"></line>
                                    <polyline points="10 9 9 9 8 9"></polyline>
                                </svg>
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
                                    <circle cx="12" cy="12" r="3"></circle>
                                    <path
                                        d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"
                                    ></path>
                                </svg>
                            {:else if item.iconType === "approval"}
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
                                    <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"
                                    ></path>
                                    <polyline points="22 4 12 14.01 9 11.01"
                                    ></polyline>
                                </svg>
                            {:else if item.iconType === "struktur"}
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
                                    <circle cx="12" cy="5" r="3"></circle>
                                    <line x1="12" y1="8" x2="12" y2="12"></line>
                                    <line x1="8" y1="14" x2="16" y2="14"></line>
                                    <line x1="8" y1="14" x2="8" y2="16"></line>
                                    <line x1="16" y1="14" x2="16" y2="16"
                                    ></line>
                                    <circle cx="8" cy="19" r="3"></circle>
                                    <circle cx="16" cy="19" r="3"></circle>
                                </svg>
                            {:else if item.iconType === "announcement"}
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
                                    <path d="m3 11 18-5v12L3 14v-3z"></path>
                                    <path d="M11.6 16.8a3 3 0 1 1-5.8-1.6"
                                    ></path>
                                </svg>
                            {/if}
                        </span>
                        <span class="nav-label">{item.label}</span>
                    </a>
                {/each}
            </nav>

            <!-- User Profile Section -->
            <div class="sidebar-divider"></div>
            <div class="user-profile-section">
                <a
                    href="/org/{orgId}/units/{unitId}/profile"
                    class="user-info-link"
                >
                    <div class="user-avatar">
                        {userName ? userName.charAt(0).toUpperCase() : "U"}
                    </div>
                    <div class="user-details">
                        <span class="user-name">{userName || "User"}</span>
                        <span class="user-role">
                            {#if userJabatan}
                                {userRole} | {userJabatan}
                            {:else}
                                {userRole}
                            {/if}
                        </span>
                    </div>
                </a>
                <div class="profile-actions">
                    <button class="logout-btn" on:click={handleLogout}>
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="16"
                            height="16"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        >
                            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"
                            ></path>
                            <polyline points="16 17 21 12 16 7"></polyline>
                            <line x1="21" x2="9" y1="12" y2="12"></line>
                        </svg>
                        Logout
                    </button>
                </div>
            </div>
        </aside>

        <!-- Main Content -->
        <main class="perumahan-content">
            <slot />
        </main>
    </div>
{/if}

<style>
    .perumahan-layout {
        display: flex;
        min-height: 100vh;
        background: #f8fafc;
    }

    .perumahan-sidebar {
        width: 280px;
        background: #ffffff;
        border-right: 1px solid #e2e8f0;
        padding: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 1.25rem;
        box-shadow: 2px 0 8px rgba(0, 0, 0, 0.03);
    }

    .back-link {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        color: #475569;
        text-decoration: none;
        font-size: 0.875rem;
        padding: 0.625rem 0.875rem;
        border-radius: 0.75rem;
        transition: all 0.2s ease;
        border: 1px solid transparent;
    }

    .back-link:hover {
        color: #1e293b;
        background: #f1f5f9;
        border-color: #e2e8f0;
    }

    .perumahan-header {
        padding: 1.25rem;
        background: linear-gradient(135deg, #ede9fe 0%, #fae8ff 100%);
        border-radius: 1rem;
        display: flex;
        gap: 1rem;
        align-items: center;
        border: 1px solid #ddd6fe;
    }

    .perumahan-logo {
        width: 56px;
        height: 56px;
        border-radius: 0.875rem;
        overflow: hidden;
        background: linear-gradient(135deg, #7c3aed 0%, #a855f7 100%);
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        flex-shrink: 0;
    }

    .perumahan-logo img {
        width: 100%;
        height: 100%;
        object-fit: cover;
    }

    .perumahan-info h2 {
        font-size: 1.0625rem;
        font-weight: 600;
        color: #1e293b;
        margin: 0;
        line-height: 1.4;
    }

    .perumahan-type {
        font-size: 0.75rem;
        color: #7c3aed;
        background: rgba(124, 58, 237, 0.1);
        padding: 0.25rem 0.625rem;
        border-radius: 0.375rem;
        display: inline-block;
        margin-top: 0.375rem;
        font-weight: 500;
    }

    .perumahan-nav {
        display: flex;
        flex-direction: column;
        gap: 0.375rem;
        margin-top: 0.5rem;
    }

    .nav-item {
        display: flex;
        align-items: center;
        gap: 0.875rem;
        padding: 0.875rem 1rem;
        color: #475569;
        text-decoration: none;
        border-radius: 0.75rem;
        transition: all 0.2s ease;
        border: 1px solid transparent;
    }

    .nav-item:hover {
        color: #1e293b;
        background: #f8fafc;
        border-color: #e2e8f0;
    }

    .nav-item.active {
        color: white;
        background: linear-gradient(135deg, #7c3aed 0%, #a855f7 100%);
        border-color: transparent;
        box-shadow: 0 4px 12px rgba(124, 58, 237, 0.3);
    }

    .nav-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 24px;
        height: 24px;
    }

    .nav-label {
        font-size: 0.9375rem;
        font-weight: 500;
    }

    .sidebar-divider {
        height: 1px;
        background: linear-gradient(
            90deg,
            transparent 0%,
            #e2e8f0 50%,
            transparent 100%
        );
        margin: 0.5rem 0;
    }

    .quick-stats {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .stats-label {
        font-size: 0.75rem;
        color: #475569;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        font-weight: 600;
        padding: 0 0.5rem;
    }

    .quick-link {
        display: flex;
        align-items: center;
        gap: 0.625rem;
        padding: 0.625rem 0.875rem;
        color: #475569;
        text-decoration: none;
        font-size: 0.875rem;
        border-radius: 0.625rem;
        transition: all 0.2s ease;
    }

    .quick-link:hover {
        color: #1e293b;
        background: #f8fafc;
    }

    .perumahan-content {
        flex: 1;
        overflow-y: auto;
        background: #f8fafc;
    }

    .loading {
        justify-content: center;
    }

    .skeleton-avatar {
        width: 56px;
        height: 56px;
        border-radius: 0.875rem;
        background: linear-gradient(
            90deg,
            #e2e8f0 0%,
            #f1f5f9 50%,
            #e2e8f0 100%
        );
        background-size: 200% 100%;
        animation: shimmer 1.5s infinite;
    }

    .skeleton-text {
        flex: 1;
        height: 40px;
        border-radius: 0.5rem;
        background: linear-gradient(
            90deg,
            #e2e8f0 0%,
            #f1f5f9 50%,
            #e2e8f0 100%
        );
        background-size: 200% 100%;
        animation: shimmer 1.5s infinite;
    }

    @keyframes shimmer {
        0% {
            background-position: 200% 0;
        }
        100% {
            background-position: -200% 0;
        }
    }

    /* User Profile Section */
    .user-profile-section {
        padding: 0.75rem;
        background: #f8fafc;
        border-radius: 0.75rem;
        border: 1px solid #e2e8f0;
        margin-top: auto;
    }

    .user-info {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        margin-bottom: 0.75rem;
    }

    .user-avatar {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        background: linear-gradient(135deg, #7c3aed 0%, #a855f7 100%);
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        font-weight: 600;
        font-size: 1rem;
    }

    .user-details {
        display: flex;
        flex-direction: column;
        flex: 1;
        overflow: hidden;
    }

    .user-name {
        font-size: 0.875rem;
        font-weight: 600;
        color: #1e293b;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .user-role {
        font-size: 0.75rem;
        color: #475569;
        text-transform: capitalize;
    }

    .logout-btn {
        width: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        padding: 0.625rem;
        background: #fef2f2;
        border: 1px solid #fecaca;
        border-radius: 0.5rem;
        color: #dc2626;
        font-size: 0.875rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s ease;
    }

    .logout-btn:hover {
        background: #fee2e2;
        border-color: #fca5a5;
    }

    .user-info-link {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.5rem;
        margin: -0.5rem;
        margin-bottom: 0.5rem;
        border-radius: 0.5rem;
        text-decoration: none;
        transition: all 0.2s ease;
    }

    .user-info-link:hover {
        background: #ede9fe;
    }

    .profile-actions {
        display: flex;
        gap: 0.5rem;
    }

    .profile-btn {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.375rem;
        padding: 0.5rem;
        background: #f0fdf4;
        border: 1px solid #86efac;
        border-radius: 0.5rem;
        color: #16a34a;
        font-size: 0.8125rem;
        font-weight: 500;
        text-decoration: none;
        transition: all 0.2s ease;
    }

    .profile-btn:hover {
        background: #dcfce7;
        border-color: #4ade80;
    }

    /* Main Content Wrapper */
    .main-wrapper {
        flex: 1;
        display: flex;
        flex-direction: column;
        overflow: hidden;
    }

    /* Top Header Bar */
    .top-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.75rem 1.5rem;
        background: white;
        border-bottom: 1px solid #e2e8f0;
        position: sticky;
        top: 0;
        z-index: 100;
    }

    .header-spacer {
        flex: 1;
    }

    .header-right {
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    /* Profile Dropdown */
    .profile-dropdown {
        position: relative;
    }

    .profile-avatar-btn {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.375rem 0.75rem;
        background: transparent;
        border: 1px solid #e2e8f0;
        border-radius: 2rem;
        cursor: pointer;
        transition: all 0.2s;
    }

    .profile-avatar-btn:hover {
        background: #f8fafc;
        border-color: #cbd5e1;
    }

    .header-avatar {
        width: 32px;
        height: 32px;
        border-radius: 50%;
        background: linear-gradient(135deg, #7c3aed 0%, #a855f7 100%);
        color: white;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: 600;
        font-size: 0.875rem;
    }

    .header-username {
        font-size: 0.875rem;
        font-weight: 500;
        color: #1e293b;
    }

    .profile-dropdown-menu {
        display: none;
        position: absolute;
        top: 100%;
        right: 0;
        margin-top: 0.5rem;
        min-width: 200px;
        background: white;
        border: 1px solid #e2e8f0;
        border-radius: 0.75rem;
        box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
        z-index: 200;
        overflow: hidden;
    }

    .profile-dropdown-menu.show {
        display: block;
    }

    .dropdown-user-info {
        padding: 1rem;
        display: flex;
        flex-direction: column;
        gap: 0.125rem;
    }

    .dropdown-name {
        font-weight: 600;
        color: #1e293b;
    }

    .dropdown-role {
        font-size: 0.75rem;
        color: #64748b;
        text-transform: capitalize;
    }

    .dropdown-divider {
        height: 1px;
        background: #e2e8f0;
    }

    .dropdown-item {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.75rem 1rem;
        color: #475569;
        text-decoration: none;
        font-size: 0.875rem;
        transition: background 0.15s;
        background: none;
        border: none;
        width: 100%;
        cursor: pointer;
    }

    .dropdown-item:hover {
        background: #f8fafc;
        color: #1e293b;
    }

    .dropdown-item.logout {
        color: #dc2626;
    }

    .dropdown-item.logout:hover {
        background: #fef2f2;
        color: #b91c1c;
    }
</style>
