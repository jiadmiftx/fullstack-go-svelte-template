<script lang="ts">
    import { page } from "$app/stores";
    import { api, selectedOrganization, type Role, type User } from "$lib";
    import { onMount } from "svelte";
    import StatsCard from "$core/components/StatsCard.svelte";

    let orgId: string = "";
    let members: any[] = [];
    let roles: Role[] = [];
    let users: User[] = [];
    let isLoading = true;

    $: orgId = $page.params.orgId ?? "";
    $: organization = $selectedOrganization;

    onMount(async () => {
        if (orgId) {
            await loadDashboardData();
        }
    });

    async function loadDashboardData() {
        isLoading = true;
        try {
            const [membersRes, rolesRes, usersRes] = await Promise.all([
                api.getOrganizationMembers(orgId, { limit: 100 }),
                api.getRoles({ limit: 100, organization_id: orgId }),
                api.getUsers({ limit: 100 }),
            ]);

            members = membersRes.data || [];
            roles = rolesRes.data || [];
            users = usersRes.data || [];
        } catch (err) {
            console.error("Failed to load dashboard data:", err);
        } finally {
            isLoading = false;
        }
    }

    // Calculate stats
    $: totalMembers = members.length;
    $: totalRoles = roles.length;
    $: activeMembers = members.filter((m) => m.is_active).length;
</script>

<svelte:head>
    <title
        >{organization?.name || "Organization"} Dashboard - Sekolah Madrasah</title
    >
</svelte:head>

<div class="py-8 px-6 bg-gray-50 min-h-screen">
    <div class="max-w-7xl mx-auto">
        <!-- Header -->
        <div class="mb-8">
            <h1 class="text-3xl font-bold text-gray-900 mb-2">
                {organization?.name || "Organization Dashboard"}
            </h1>
            <p class="text-gray-600">
                Manage your organization's members, roles, and settings
            </p>
        </div>

        <!-- Stats Cards -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <StatsCard
                title="Total Members"
                value={totalMembers}
                change="+{activeMembers} active"
                changeType="neutral"
                iconType="users"
                loading={isLoading}
            />
            <StatsCard
                title="Active Roles"
                value={totalRoles}
                change="Organization roles"
                changeType="neutral"
                iconType="user"
                loading={isLoading}
            />
            <StatsCard
                title="Active Members"
                value={activeMembers}
                change={`${Math.round((activeMembers / totalMembers) * 100) || 0}% of total`}
                changeType="positive"
                iconType="check"
                loading={isLoading}
            />
        </div>

        <!-- Quick Actions & Recent Members -->
        <div class="grid md:grid-cols-2 gap-6">
            <!-- Quick Actions -->
            <div class="card">
                <h2 class="text-xl font-semibold text-gray-900 mb-6">
                    Quick Actions
                </h2>
                <div class="space-y-3">
                    <a
                        href="/org/{orgId}/users"
                        class="block p-4 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors group"
                    >
                        <div class="flex items-center gap-3">
                            <div
                                class="w-10 h-10 bg-gray-800 rounded-lg flex items-center justify-center text-white"
                            >
                                <svg
                                    class="w-5 h-5"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="1.5"
                                    ><path
                                        d="M20.7739 18C21.5232 18 22.1192 17.5285 22.6543 16.8691C23.7498 15.5194 21.9512 14.4408 21.2652 13.9126C20.5679 13.3756 19.7893 13.0714 18.9999 13M17.9999 11C19.3806 11 20.4999 9.88071 20.4999 8.5C20.4999 7.11929 19.3806 6 17.9999 6"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    /><path
                                        d="M3.2259 18C2.47659 18 1.88061 17.5285 1.34548 16.8691C0.250028 15.5194 2.04861 14.4408 2.73458 13.9126C3.43191 13.3756 4.21052 13.0714 4.99994 13M5.99994 11C4.61923 11 3.49994 9.88071 3.49994 8.5C3.49994 7.11929 4.61923 6 5.99994 6"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    /><path
                                        d="M8.08368 15.1112C7.0619 15.743 4.38286 17.0331 6.01458 18.6474C6.81166 19.436 7.6994 20 8.8155 20H15.1843C16.3004 20 17.1881 19.436 17.9852 18.6474C19.6169 17.0331 16.9379 15.743 15.9161 15.1112C13.52 13.6296 10.4798 13.6296 8.08368 15.1112Z"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    /><path
                                        d="M15.4999 7.5C15.4999 9.433 13.9329 11 11.9999 11C10.0669 11 8.49988 9.433 8.49988 7.5C8.49988 5.567 10.0669 4 11.9999 4C13.9329 4 15.4999 5.567 15.4999 7.5Z"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    /></svg
                                >
                            </div>
                            <div class="flex-1">
                                <h3
                                    class="font-medium text-gray-900 group-hover:text-gray-700"
                                >
                                    Manage Members
                                </h3>
                                <p class="text-sm text-gray-600">
                                    Add or manage organization members
                                </p>
                            </div>
                            <span class="text-gray-400">→</span>
                        </div>
                    </a>

                    <a
                        href="/org/{orgId}/roles"
                        class="block p-4 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors group"
                    >
                        <div class="flex items-center gap-3">
                            <div
                                class="w-10 h-10 bg-gray-600 rounded-lg flex items-center justify-center text-white"
                            >
                                <svg
                                    class="w-5 h-5"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="1.5"
                                    ><path
                                        d="M6.57757 15.4816C5.1628 16.324 1.45336 18.0441 3.71266 20.1966C4.81631 21.248 6.04549 22 7.59087 22H16.4091C17.9545 22 19.1837 21.248 20.2873 20.1966C22.5466 18.0441 18.8372 16.324 17.4224 15.4816C14.1048 13.5061 9.89519 13.5061 6.57757 15.4816Z"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    /><path
                                        d="M16.5 6.5C16.5 8.98528 14.4853 11 12 11C9.51472 11 7.5 8.98528 7.5 6.5C7.5 4.01472 9.51472 2 12 2C14.4853 2 16.5 4.01472 16.5 6.5Z"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    /></svg
                                >
                            </div>
                            <div class="flex-1">
                                <h3
                                    class="font-medium text-gray-900 group-hover:text-gray-700"
                                >
                                    Manage Roles
                                </h3>
                                <p class="text-sm text-gray-600">
                                    Define custom roles with permissions
                                </p>
                            </div>
                            <span class="text-gray-400">→</span>
                        </div>
                    </a>
                </div>
            </div>

            <!-- Recent Members -->
            <div class="card">
                <div class="flex items-center justify-between mb-6">
                    <h2 class="text-xl font-semibold text-gray-900">
                        Recent Members
                    </h2>
                    <a
                        href="/org/{orgId}/users"
                        class="text-gray-700 hover:text-gray-900 text-sm font-medium"
                    >
                        View all →
                    </a>
                </div>

                {#if isLoading}
                    <div class="space-y-3">
                        {#each Array(3) as _}
                            <div class="skeleton h-14 rounded-lg"></div>
                        {/each}
                    </div>
                {:else if members.length === 0}
                    <div class="text-center py-8">
                        <p class="text-gray-900 font-medium mb-1">
                            No members yet
                        </p>
                        <p class="text-gray-500 text-sm mb-4">
                            Add your first member to get started
                        </p>
                        <a
                            href="/org/{orgId}/users"
                            class="btn-primary inline-block text-sm"
                        >
                            Add Member
                        </a>
                    </div>
                {:else}
                    <div class="space-y-2">
                        {#each members.slice(0, 5) as member}
                            {@const user = users.find(
                                (u) => u.id === member.user_id,
                            )}
                            {@const role = roles.find(
                                (r) => r.id === member.role_id,
                            )}
                            <div
                                class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg"
                            >
                                <div
                                    class="w-10 h-10 bg-gray-800 rounded-full flex items-center justify-center text-white font-semibold"
                                >
                                    {user?.full_name?.charAt(0) || "?"}
                                </div>
                                <div class="flex-1 min-w-0">
                                    <p
                                        class="font-medium text-gray-900 truncate"
                                    >
                                        {user?.full_name || "Unknown"}
                                    </p>
                                    <p class="text-sm text-gray-500 truncate">
                                        {role?.display_name ||
                                            role?.name ||
                                            "No role"}
                                    </p>
                                </div>
                                {#if member.is_active}
                                    <span class="badge badge-success"
                                        >Active</span
                                    >
                                {:else}
                                    <span class="badge badge-gray"
                                        >Inactive</span
                                    >
                                {/if}
                            </div>
                        {/each}
                    </div>
                {/if}
            </div>
        </div>
    </div>
</div>
