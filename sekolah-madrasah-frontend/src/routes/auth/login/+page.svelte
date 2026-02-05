<script lang="ts">
  import { auth, api } from "$lib";
  import { goto } from "$app/navigation";

  let email = "";
  let password = "";
  let error = "";
  let isLoading = false;

  async function handleLogin() {
    error = "";
    isLoading = true;

    try {
      const response = await auth.login(email, password);
      const user = response.data.user;

      // Try to get user's memberships for role-based redirect
      try {
        // Small delay to ensure token is propagated to API client
        await new Promise((r) => setTimeout(r, 100));
        const memberships = await api.getMyMemberships();
        const data = memberships.data;

        if (data.is_super_admin) {
          // Super Admin goes to global dashboard (pick organization)
          goto("/dashboard");
          return;
        }

        // Check perumahan memberships first
        const perumahanMemberships = data.unit_memberships || [];

        if (perumahanMemberships.length > 0) {
          // User has perumahan membership(s)
          const firstPerumahan = perumahanMemberships[0];

          if (
            firstPerumahan.role === "pengurus" ||
            firstPerumahan.role === "warga"
          ) {
            // Pengurus/Warga: Go directly to perumahan dashboard
            goto(
              `/org/${firstPerumahan.org_id}/units/${firstPerumahan.unit_id}/dashboard`,
            );
            return;
          }

          // Admin/Staff/Parent: Can navigate via org dashboard
          goto(`/org/${firstPerumahan.org_id}/dashboard`);
          return;
        }

        // Check organization memberships
        const orgMemberships = data.organization_memberships || [];

        if (orgMemberships.length > 0) {
          // User has org membership - go to org dashboard
          goto(`/org/${orgMemberships[0].org_id}/dashboard`);
          return;
        }
      } catch (membershipErr) {
        console.warn(
          "Failed to get memberships, using fallback:",
          membershipErr,
        );
        // Fallback for super admin based on user data
        if (user.is_super_admin) {
          goto("/dashboard");
          return;
        }
      }

      // Fallback: Check if user owns any organization
      try {
        const orgsResponse = await api.getOrganizations({ limit: 100 });
        const ownedOrgs =
          orgsResponse.data?.filter((org) => org.owner_id === user.id) || [];

        if (ownedOrgs.length > 0) {
          goto(`/org/${ownedOrgs[0].id}/dashboard`);
          return;
        }
      } catch {
        // Ignore errors
      }

      // No memberships found - go to default dashboard
      goto("/dashboard");
    } catch (err) {
      error = err instanceof Error ? err.message : "Login failed";
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>Login - Sekolah Madrasah</title>
</svelte:head>

<div class="min-h-[calc(100vh-4rem)] flex items-center justify-center px-4">
  <div class="w-full max-w-md">
    <div class="card">
      <h1 class="text-2xl font-bold text-center mb-6">Login</h1>

      {#if error}
        {#if error.includes("menunggu verifikasi")}
          <div
            class="bg-amber-500/20 border border-amber-500 text-amber-300 px-4 py-3 rounded-lg mb-4"
          >
            ⏳ {error}
          </div>
        {:else}
          <div
            class="bg-red-500/20 border border-red-500 text-red-400 px-4 py-3 rounded-lg mb-4"
          >
            {error}
          </div>
        {/if}
      {/if}

      <form on:submit|preventDefault={handleLogin} class="space-y-4">
        <div>
          <label
            for="email"
            class="block text-sm font-medium text-slate-300 mb-1"
          >
            Email
          </label>
          <input
            id="email"
            type="email"
            bind:value={email}
            class="input-field"
            placeholder="you@example.com"
            required
          />
        </div>

        <div>
          <label
            for="password"
            class="block text-sm font-medium text-slate-300 mb-1"
          >
            Password
          </label>
          <input
            id="password"
            type="password"
            bind:value={password}
            class="input-field"
            placeholder="••••••••"
            required
          />
        </div>

        <button
          type="submit"
          class="btn-primary w-full py-3"
          disabled={isLoading}
        >
          {isLoading ? "Logging in..." : "Login"}
        </button>
      </form>

      <p class="text-center text-slate-400 mt-6">
        Don't have an account?
        <a href="/auth/register" class="text-primary-500 hover:underline"
          >Register</a
        >
      </p>
    </div>
  </div>
</div>
