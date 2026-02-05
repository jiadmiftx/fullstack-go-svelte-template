<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import {
        isAuthenticated,
        api,
        selectedOrganization,
        isSuperAdmin,
        type Organization,
    } from "$lib";

    let orgId: string = "";
    let organization: Organization | null = null;
    let isLoading = true;
    let error = "";

    $: orgId = $page.params.orgId ?? "";

    // Redirect unauthenticated users
    $: if (!$isAuthenticated && typeof window !== "undefined") {
        goto("/auth/login");
    }

    onMount(async () => {
        if (orgId && $isAuthenticated) {
            await loadOrganization();
            await checkUserAccess();
        }
    });

    async function loadOrganization() {
        isLoading = true;
        error = "";
        try {
            const response = await api.getOrganization(orgId);
            organization = response.data;
            // Update selected organization in store
            if (organization) {
                selectedOrganization.set(organization);
            }
        } catch (err) {
            error =
                err instanceof Error ? err.message : "Organization not found";
            // Redirect if org not accessible
            goto("/dashboard");
        } finally {
            isLoading = false;
        }
    }

    async function checkUserAccess() {
        try {
            const response = await api.getMyMemberships();
            const memberships = response.data;

            // Super admin can access anything
            if (memberships.is_super_admin) return;

            // Check if user has org-level membership
            const hasOrgAccess = memberships.organization_memberships?.some(
                (m: any) => m.org_id === orgId,
            );

            if (hasOrgAccess) return; // Org members can access org pages

            // Check if user only has perumahan-level membership (pengurus/warga)
            const perumahansInOrg =
                memberships.unit_memberships?.filter(
                    (m: any) => m.org_id === orgId,
                ) || [];

            if (perumahansInOrg.length > 0) {
                const firstPerumahan = perumahansInOrg[0];
                // Pengurus/warga/parent - redirect to their perumahan dashboard
                if (
                    ["pengurus", "warga", "parent"].includes(firstPerumahan.role)
                ) {
                    goto(
                        `/org/${orgId}/units/${firstPerumahan.unit_id}/dashboard`,
                    );
                    return;
                }
            }

            // No access at all - redirect to main dashboard
            goto("/dashboard");
        } catch (err) {
            console.error("Failed to check user access:", err);
            // If access check fails, redirect non-super-admin to dashboard
            if (!$isSuperAdmin) {
                goto("/dashboard");
            }
        }
    }
</script>

{#if isLoading}
    <div class="flex items-center justify-center min-h-[400px]">
        <div class="text-slate-400">Loading organization...</div>
    </div>
{:else if error}
    <div class="flex items-center justify-center min-h-[400px]">
        <div class="text-red-400">{error}</div>
    </div>
{:else}
    <slot />
{/if}
