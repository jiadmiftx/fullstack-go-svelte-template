<script lang="ts">
  import "../app.css";
  import { onMount } from "svelte";
  import { page } from "$app/stores";
  import {
    auth,
    isAuthenticated,
    currentUser,
    isSuperAdmin,
    api,
    selectedOrganization,
    type Organization,
  } from "$lib";
  import { goto } from "$app/navigation";
  import Toast from "$core/components/Toast.svelte";

  let organizations: Organization[] = [];
  let sidebarCollapsed = false;
  let showOrgDropdown = false;

  // Super Admin menu - global access
  const superAdminMenu = [
    { href: "/dashboard", label: "Dashboard", iconType: "dashboard" },
    { href: "/organizations", label: "Organizations", iconType: "building" },
    { href: "/users", label: "Users", iconType: "users" },
    { href: "/permissions", label: "Permissions", iconType: "lock" },
  ];

  // Org Admin menu - org-scoped access
  $: orgAdminMenu = $selectedOrganization
    ? [
        {
          href: `/org/${$selectedOrganization.id}/dashboard`,
          label: "Dashboard",
          iconType: "dashboard",
        },
        {
          href: `/org/${$selectedOrganization.id}/units`,
          label: "Units",
          iconType: "perumahan",
        },
        {
          href: `/org/${$selectedOrganization.id}/users`,
          label: "Users",
          iconType: "users",
        },
        {
          href: `/org/${$selectedOrganization.id}/roles`,
          label: "Roles",
          iconType: "user",
        },
      ]
    : [];

  $: currentPath = $page.url.pathname;

  // Check if currently on auth page (login/register)
  $: isAuthPage = currentPath.startsWith("/auth/");

  // Check if currently on public registration page
  $: isPublicPage = currentPath.startsWith("/daftar/");

  // Check if currently viewing a perumahan route (has its own sidebar)
  $: isPerumahanRoute = /^\/org\/[^/]+\/units\/[^/]+/.test(currentPath);

  // Check if currently viewing an org route
  $: isOrgRoute = currentPath.startsWith("/org/");

  // Dynamic menu based on user role and current route
  // Super Admin in org route: show ONLY org menu (like org admin)
  // Super Admin in global route: show super admin menu
  // Org Admin: show org menu only
  $: menuItems = $isSuperAdmin
    ? isOrgRoute && $selectedOrganization
      ? orgAdminMenu // Only org menu when in org route
      : superAdminMenu
    : orgAdminMenu;

  onMount(async () => {
    if ($isAuthenticated) {
      await loadOrganizations();
    }
  });

  async function loadOrganizations() {
    try {
      const response = await api.getOrganizations({ limit: 100 });
      organizations = response.data || [];

      // Auto-select first org if none selected (for non-super-admin)
      if (
        !$isSuperAdmin &&
        !$selectedOrganization &&
        organizations.length > 0
      ) {
        selectedOrganization.set(organizations[0]);
      }
    } catch (err) {
      console.error("Failed to load organizations:", err);
    }
  }

  function selectOrganization(org: Organization) {
    selectedOrganization.set(org);
    showOrgDropdown = false;
    goto(`/org/${org.id}/dashboard`);
  }

  function exitOrgMode() {
    goto("/dashboard");
  }

  function handleLogout() {
    auth.logout();
    selectedOrganization.clear();
    goto("/auth/login");
  }
</script>

<div class="min-h-screen bg-gray-50 flex">
  <!-- Sidebar (hide on auth pages, public pages, and perumahan routes which have their own sidebar) -->
  {#if $isAuthenticated && !isAuthPage && !isPublicPage && !isPerumahanRoute}
    <aside
      class="fixed left-0 top-0 h-full bg-white border-r border-gray-200 flex flex-col transition-all duration-300 z-40 shadow-sm {sidebarCollapsed
        ? 'w-20'
        : 'w-64'}"
    >
      <!-- Logo / Header -->
      <div class="h-16 flex items-center px-4 border-b border-gray-200">
        <a href="/" class="flex items-center gap-3">
          <!-- Shield Icon -->
          <svg
            class="w-8 h-8 text-gray-800"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.5"
          >
            <path
              d="M11.9982 2C8.99043 2 7.04018 4.01899 4.73371 4.7549C3.79589 5.05413 3.32697 5.20374 3.1372 5.41465C2.94743 5.62556 2.89186 5.93375 2.78072 6.55013C1.59143 13.146 4.1909 19.244 10.3903 21.6175C11.0564 21.8725 11.3894 22 12.0015 22C12.6136 22 12.9466 21.8725 13.6127 21.6175C19.8116 19.2439 22.4086 13.146 21.2194 6.55013C21.1082 5.93375 21.0526 5.62556 20.8629 5.41465C20.6731 5.20374 20.2042 5.05413 19.2664 4.7549C16.9595 4.01899 15.0092 2 11.9982 2Z"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            <path
              d="M9 13L10.7528 14.4023C11.1707 14.7366 11.7777 14.6826 12.1301 14.2799L15 11"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
          {#if !sidebarCollapsed}
            <span class="text-xl font-bold text-gray-900">RBAC</span>
          {/if}
        </a>
        <button
          on:click={() => (sidebarCollapsed = !sidebarCollapsed)}
          class="ml-auto p-2 hover:bg-gray-100 rounded-lg transition-colors text-gray-600"
        >
          <svg
            class="w-5 h-5 transition-transform {sidebarCollapsed
              ? 'rotate-180'
              : ''}"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M11 19l-7-7 7-7m8 14l-7-7 7-7"
            />
          </svg>
        </button>
      </div>

      <!-- Super Admin in Org Mode: Show org banner + exit button -->
      {#if $isSuperAdmin && isOrgRoute && $selectedOrganization}
        <div class="p-3 bg-gray-100 border-b border-gray-200">
          {#if !sidebarCollapsed}
            <div class="flex items-center gap-2 mb-2">
              <!-- Building Icon -->
              <svg
                class="w-5 h-5 text-gray-700"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  d="M14 22V8C14 5.17157 14 3.75736 13.1213 2.87868C12.2426 2 10.8284 2 8 2C5.17157 2 3.75736 2 2.87868 2.87868C2 3.75736 2 5.17157 2 8V16C2 18.8284 2 20.2426 2.87868 21.1213C3.75736 22 5.17157 22 8 22H14Z"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M14 22H16C18.8284 22 20.2426 22 21.1213 21.1213C22 20.2426 22 18.8284 22 16V13C22 11.1144 22 10.1716 21.4142 9.58579C20.8284 9 19.8856 9 18 9H14"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M5 6H6M5 9.5H6M5 13H6"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M17.5 13H18M17.5 16H18"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M9 22V19C9 18.0572 9 17.5858 9.29289 17.2929C9.58579 17 10.0572 17 11 17C11.9428 17 12.4142 17 12.7071 17.2929C13 17.5858 13 18.0572 13 19V22"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
              <div class="flex-1 min-w-0">
                <p class="text-xs text-gray-600 font-medium">
                  Viewing as Org Admin
                </p>
                <p class="text-sm font-semibold text-gray-900 truncate">
                  {$selectedOrganization.name}
                </p>
              </div>
            </div>
            <button
              on:click={exitOrgMode}
              class="w-full flex items-center justify-center gap-2 px-3 py-2 bg-gray-800 hover:bg-gray-900 text-white rounded-lg transition-colors text-sm font-medium"
            >
              <!-- Arrow Left Icon -->
              <svg
                class="w-4 h-4"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  d="M15 6L9 12L15 18"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
              <span>Back to Super Admin</span>
            </button>
          {:else}
            <button
              on:click={exitOrgMode}
              class="w-full flex items-center justify-center p-2 bg-gray-800 hover:bg-gray-900 text-white rounded-lg transition-colors"
              title="Back to Super Admin"
            >
              <svg
                class="w-4 h-4"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  d="M15 6L9 12L15 18"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </button>
          {/if}
        </div>
      {/if}

      <!-- Organization Display (Org Admin only - No dropdown, just display) -->
      {#if !$isSuperAdmin && $selectedOrganization}
        <div class="p-4 border-b border-gray-200">
          <div
            class="flex items-center gap-2 px-3 py-2 bg-gray-100 rounded-lg text-gray-900"
          >
            <!-- Building Icon -->
            <svg
              class="w-5 h-5 text-gray-700"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
            >
              <path
                d="M14 22V8C14 5.17157 14 3.75736 13.1213 2.87868C12.2426 2 10.8284 2 8 2C5.17157 2 3.75736 2 2.87868 2.87868C2 3.75736 2 5.17157 2 8V16C2 18.8284 2 20.2426 2.87868 21.1213C3.75736 22 5.17157 22 8 22H14Z"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M14 22H16C18.8284 22 20.2426 22 21.1213 21.1213C22 20.2426 22 18.8284 22 16V13C22 11.1144 22 10.1716 21.4142 9.58579C20.8284 9 19.8856 9 18 9H14"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M5 6H6M5 9.5H6M5 13H6"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M17.5 13H18M17.5 16H18"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M9 22V19C9 18.0572 9 17.5858 9.29289 17.2929C9.58579 17 10.0572 17 11 17C11.9428 17 12.4142 17 12.7071 17.2929C13 17.5858 13 18.0572 13 19V22"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
            {#if !sidebarCollapsed}
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-gray-900 truncate">
                  {$selectedOrganization.name}
                </p>
                <p class="text-xs text-gray-500">
                  {$selectedOrganization.code}
                </p>
              </div>
            {/if}
          </div>
        </div>
      {/if}

      <!-- Navigation Menu -->
      <nav class="flex-1 p-4 space-y-1 overflow-y-auto">
        {#each menuItems as item}
          <a
            href={item.href}
            class="flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all duration-200 font-medium
              {currentPath === item.href ||
            currentPath.startsWith(item.href + '/')
              ? 'bg-gray-100 text-gray-900 shadow-sm'
              : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
          >
            <!-- Dynamic Icon based on iconType -->
            {#if item.iconType === "dashboard"}
              <svg
                class="w-5 h-5"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  d="M2 6C2 4.11438 2 3.17157 2.58579 2.58579C3.17157 2 4.11438 2 6 2C7.88562 2 8.82843 2 9.41421 2.58579C10 3.17157 10 4.11438 10 6C10 7.88562 10 8.82843 9.41421 9.41421C8.82843 10 7.88562 10 6 10C4.11438 10 3.17157 10 2.58579 9.41421C2 8.82843 2 7.88562 2 6Z"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M2 18C2 16.1144 2 15.1716 2.58579 14.5858C3.17157 14 4.11438 14 6 14C7.88562 14 8.82843 14 9.41421 14.5858C10 15.1716 10 16.1144 10 18C10 19.8856 10 20.8284 9.41421 21.4142C8.82843 22 7.88562 22 6 22C4.11438 22 3.17157 22 2.58579 21.4142C2 20.8284 2 19.8856 2 18Z"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M14 6C14 4.11438 14 3.17157 14.5858 2.58579C15.1716 2 16.1144 2 18 2C19.8856 2 20.8284 2 21.4142 2.58579C22 3.17157 22 4.11438 22 6C22 7.88562 22 8.82843 21.4142 9.41421C20.8284 10 19.8856 10 18 10C16.1144 10 15.1716 10 14.5858 9.41421C14 8.82843 14 7.88562 14 6Z"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M14 18C14 16.1144 14 15.1716 14.5858 14.5858C15.1716 14 16.1144 14 18 14C19.8856 14 20.8284 14 21.4142 14.5858C22 15.1716 22 16.1144 22 18C22 19.8856 22 20.8284 21.4142 21.4142C20.8284 22 19.8856 22 18 22C16.1144 22 15.1716 22 14.5858 21.4142C14 20.8284 14 19.8856 14 18Z"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            {:else if item.iconType === "building"}
              <svg
                class="w-5 h-5"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  d="M14 22V8C14 5.17157 14 3.75736 13.1213 2.87868C12.2426 2 10.8284 2 8 2C5.17157 2 3.75736 2 2.87868 2.87868C2 3.75736 2 5.17157 2 8V16C2 18.8284 2 20.2426 2.87868 21.1213C3.75736 22 5.17157 22 8 22H14Z"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M14 22H16C18.8284 22 20.2426 22 21.1213 21.1213C22 20.2426 22 18.8284 22 16V13C22 11.1144 22 10.1716 21.4142 9.58579C20.8284 9 19.8856 9 18 9H14"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M5 6H6M5 9.5H6M5 13H6"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M17.5 13H18M17.5 16H18"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M9 22V19C9 18.0572 9 17.5858 9.29289 17.2929C9.58579 17 10.0572 17 11 17C11.9428 17 12.4142 17 12.7071 17.2929C13 17.5858 13 18.0572 13 19V22"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            {:else if item.iconType === "users"}
              <svg
                class="w-5 h-5"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  d="M20.7739 18C21.5232 18 22.1192 17.5285 22.6543 16.8691C23.7498 15.5194 21.9512 14.4408 21.2652 13.9126C20.5679 13.3756 19.7893 13.0714 18.9999 13M17.9999 11C19.3806 11 20.4999 9.88071 20.4999 8.5C20.4999 7.11929 19.3806 6 17.9999 6"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M3.2259 18C2.47659 18 1.88061 17.5285 1.34548 16.8691C0.250028 15.5194 2.04861 14.4408 2.73458 13.9126C3.43191 13.3756 4.21052 13.0714 4.99994 13M5.99994 11C4.61923 11 3.49994 9.88071 3.49994 8.5C3.49994 7.11929 4.61923 6 5.99994 6"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M8.08368 15.1112C7.0619 15.743 4.38286 17.0331 6.01458 18.6474C6.81166 19.436 7.6994 20 8.8155 20H15.1843C16.3004 20 17.1881 19.436 17.9852 18.6474C19.6169 17.0331 16.9379 15.743 15.9161 15.1112C13.52 13.6296 10.4798 13.6296 8.08368 15.1112Z"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M15.4999 7.5C15.4999 9.433 13.9329 11 11.9999 11C10.0669 11 8.49988 9.433 8.49988 7.5C8.49988 5.567 10.0669 4 11.9999 4C13.9329 4 15.4999 5.567 15.4999 7.5Z"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            {:else if item.iconType === "user"}
              <svg
                class="w-5 h-5"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  d="M6.57757 15.4816C5.1628 16.324 1.45336 18.0441 3.71266 20.1966C4.81631 21.248 6.04549 22 7.59087 22H16.4091C17.9545 22 19.1837 21.248 20.2873 20.1966C22.5466 18.0441 18.8372 16.324 17.4224 15.4816C14.1048 13.5061 9.89519 13.5061 6.57757 15.4816Z"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M16.5 6.5C16.5 8.98528 14.4853 11 12 11C9.51472 11 7.5 8.98528 7.5 6.5C7.5 4.01472 9.51472 2 12 2C14.4853 2 16.5 4.01472 16.5 6.5Z"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            {:else if item.iconType === "lock"}
              <svg
                class="w-5 h-5"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  d="M12 16.5V14.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M4.26781 18.8447C4.49269 20.515 5.87613 21.8235 7.55966 21.9009C8.97627 21.966 10.4153 22 12 22C13.5847 22 15.0237 21.966 16.4403 21.9009C18.1239 21.8235 19.5073 20.515 19.7322 18.8447C19.879 17.7547 20 16.6376 20 15.5C20 14.3624 19.879 13.2453 19.7322 12.1553C19.5073 10.485 18.1239 9.17649 16.4403 9.09909C15.0237 9.03397 13.5847 9 12 9C10.4153 9 8.97627 9.03397 7.55966 9.09909C5.87613 9.17649 4.49269 10.485 4.26781 12.1553C4.12105 13.2453 4 14.3624 4 15.5C4 16.6376 4.12105 17.7547 4.26781 18.8447Z"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M7.5 9V6.5C7.5 4.01472 9.51472 2 12 2C14.4853 2 16.5 4.01472 16.5 6.5V9"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            {:else if item.iconType === "perumahan"}
              <svg
                class="w-5 h-5"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  d="M2 22H22"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M3 22V11L12 2L21 11V22"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M9 22V17C9 15.8954 9.89543 15 11 15H13C14.1046 15 15 15.8954 15 17V22"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M12 2V6"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            {:else}
              <svg
                class="w-5 h-5"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <circle
                  cx="12"
                  cy="12"
                  r="10"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            {/if}
            {#if !sidebarCollapsed}
              <span>{item.label}</span>
            {/if}
          </a>
        {/each}
      </nav>

      <!-- User Section -->
      <div class="p-4 border-t border-gray-200 bg-gray-50">
        <div class="flex items-center gap-3 mb-3">
          <div
            class="w-10 h-10 rounded-full bg-gray-800 flex items-center justify-center text-lg font-bold text-white"
          >
            {$currentUser?.full_name?.charAt(0) || "U"}
          </div>
          {#if !sidebarCollapsed}
            <div class="flex-1 min-w-0">
              <div class="font-medium truncate text-gray-900">
                {$currentUser?.full_name}
              </div>
              <div class="text-xs text-gray-500 truncate">
                {$isSuperAdmin ? "Super Admin" : $currentUser?.email}
              </div>
            </div>
          {/if}
        </div>
        <button
          on:click={handleLogout}
          class="w-full flex items-center justify-center gap-2 px-3 py-2 bg-red-50 text-red-600 hover:bg-red-100 rounded-lg transition-colors font-medium"
        >
          <!-- Logout Icon -->
          <svg
            class="w-5 h-5"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.5"
          >
            <path
              d="M14 3.09502C13.543 3.03241 13.0755 3 12.6 3C7.29807 3 3 7.02944 3 12C3 16.9706 7.29807 21 12.6 21C13.0755 21 13.543 20.9676 14 20.905"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            <path
              d="M21 12L11 12M21 12C21 11.2998 19.0057 9.99153 18.5 9.5M21 12C21 12.7002 19.0057 14.0085 18.5 14.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
          {#if !sidebarCollapsed}
            <span>Logout</span>
          {/if}
        </button>
      </div>
    </aside>
  {/if}

  <!-- Main Content Area -->
  <div
    class="flex-1 transition-all duration-300 {$isAuthenticated &&
    !isAuthPage &&
    !isPublicPage &&
    !isPerumahanRoute
      ? sidebarCollapsed
        ? 'ml-20'
        : 'ml-64'
      : ''}"
  >
    <!-- Top Bar (For non-authenticated users) -->
    {#if !$isAuthenticated}
      <nav class="bg-white border-b border-gray-200 shadow-sm">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16 items-center">
            <a href="/" class="flex items-center gap-3">
              <!-- Shield Icon -->
              <svg
                class="w-8 h-8 text-gray-800"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  d="M11.9982 2C8.99043 2 7.04018 4.01899 4.73371 4.7549C3.79589 5.05413 3.32697 5.20374 3.1372 5.41465C2.94743 5.62556 2.89186 5.93375 2.78072 6.55013C1.59143 13.146 4.1909 19.244 10.3903 21.6175C11.0564 21.8725 11.3894 22 12.0015 22C12.6136 22 12.9466 21.8725 13.6127 21.6175C19.8116 19.2439 22.4086 13.146 21.2194 6.55013C21.1082 5.93375 21.0526 5.62556 20.8629 5.41465C20.6731 5.20374 20.2042 5.05413 19.2664 4.7549C16.9595 4.01899 15.0092 2 11.9982 2Z"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M9 13L10.7528 14.4023C11.1707 14.7366 11.7777 14.6826 12.1301 14.2799L15 11"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
              <span class="text-xl font-bold text-gray-900">RBAC</span>
            </a>
            <div class="flex items-center space-x-4">
              <a
                href="/auth/login"
                class="text-gray-700 hover:text-gray-900 font-medium transition-colors"
              >
                Login
              </a>
              <a href="/auth/register" class="btn-primary text-sm">Register</a>
            </div>
          </div>
        </div>
      </nav>
    {/if}

    <!-- Main Content -->
    <main class="min-h-screen">
      <slot />
    </main>
  </div>

  <!-- Toast Notifications -->
  <Toast />
</div>

<style>
  .rotate-180 {
    transform: rotate(180deg);
  }
</style>
