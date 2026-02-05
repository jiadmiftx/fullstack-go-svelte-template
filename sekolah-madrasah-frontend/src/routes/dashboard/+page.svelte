<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import {
    isAuthenticated,
    isSuperAdmin,
    currentUser,
    api,
    type Organization,
    type User,
    type Permission,
  } from "$lib";
  import StatsCard from "$core/components/StatsCard.svelte";

  let organizations: Organization[] = [];
  let users: User[] = [];
  let permissions: Permission[] = [];
  let isLoading = true;
  let error = "";

  // Redirect if not authenticated
  $: if (!$isAuthenticated && typeof window !== "undefined") {
    goto("/auth/login");
  }

  // Only Super Admin can access this dashboard
  $: if ($isAuthenticated && !$isSuperAdmin && typeof window !== "undefined") {
    goto("/dashboard");
  }

  onMount(async () => {
    if ($isAuthenticated && $isSuperAdmin) {
      await loadData();
    }
  });

  async function loadData() {
    isLoading = true;
    error = "";
    try {
      const [orgsRes, usersRes, permsRes] = await Promise.all([
        api.getOrganizations({ limit: 100 }),
        api.getUsers({ limit: 100 }),
        api.getPermissions({ limit: 100 }),
      ]);

      organizations = orgsRes.data || [];
      users = usersRes.data || [];
      permissions = permsRes.data || [];
    } catch (err) {
      error =
        err instanceof Error ? err.message : "Failed to load dashboard data";
    } finally {
      isLoading = false;
    }
  }

  // Calculate stats
  $: totalOrgs = organizations.length;
  $: totalUsers = users.length;
  $: totalPermissions = permissions.length;
  $: activeUsers = users.filter((u) => u.is_active).length;
</script>

<svelte:head>
  <title>Super Admin Dashboard - Sekolah Madrasah</title>
</svelte:head>

{#if $isAuthenticated && $isSuperAdmin}
  <div class="py-8 px-6 bg-gray-50 min-h-screen">
    <div class="max-w-7xl mx-auto">
      <!-- Welcome Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">
          Welcome back, {$currentUser?.full_name}! 👋
        </h1>
        <p class="text-gray-600">
          Here's what's happening with your system today.
        </p>
      </div>

      {#if error}
        <div
          class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6"
        >
          {error}
        </div>
      {/if}

      <!-- Stats Grid -->
      <div class="grid md:grid-cols-4 gap-6 mb-8">
        <StatsCard
          title="Total Organizations"
          value={totalOrgs}
          change="+12%"
          changeType="positive"
          iconType="building"
          loading={isLoading}
        />
        <StatsCard
          title="Total Users"
          value={totalUsers}
          change="+8%"
          changeType="positive"
          iconType="users"
          loading={isLoading}
        />
        <StatsCard
          title="Active Users"
          value={activeUsers}
          change={`${activeUsers}/${totalUsers}`}
          changeType="neutral"
          iconType="check"
          loading={isLoading}
        />
        <StatsCard
          title="Permissions"
          value={totalPermissions}
          change="Global"
          changeType="neutral"
          iconType="lock"
          loading={isLoading}
        />
      </div>

      <!-- Recent Organizations -->
      <div class="card">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-xl font-semibold text-gray-900">
            Recent Organizations
          </h2>
          <a
            href="/organizations"
            class="text-gray-700 hover:text-gray-900 text-sm font-medium"
          >
            View all →
          </a>
        </div>

        {#if isLoading}
          <div class="space-y-3">
            {#each Array(3) as _}
              <div class="skeleton h-16 rounded-lg"></div>
            {/each}
          </div>
        {:else if organizations.length === 0}
          <div class="empty-state">
            <p class="text-lg font-medium text-gray-900 mb-1">
              No organizations yet
            </p>
            <p class="text-gray-500 mb-4">
              Get started by creating your first organization
            </p>
            <a href="/organizations" class="btn-primary inline-block">
              Create Organization
            </a>
          </div>
        {:else}
          <div class="space-y-3">
            {#each organizations.slice(0, 5) as org}
              <div
                class="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
              >
                <div class="flex-1">
                  <h3 class="font-medium text-gray-900">{org.name}</h3>
                  <p class="text-sm text-gray-500">{org.code} • {org.type}</p>
                </div>
                <div class="flex items-center gap-3">
                  <span class="badge badge-info">{org.type}</span>
                  <a
                    href="/org/{org.id}/dashboard"
                    class="btn-secondary text-sm"
                  >
                    Enter →
                  </a>
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}
