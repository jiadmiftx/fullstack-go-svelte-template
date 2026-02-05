<script lang="ts">
    import { goto } from "$app/navigation";
    import { isAuthenticated, isSuperAdmin } from "$lib";
    import { onMount } from "svelte";

    export let requireSuperAdmin = false;
    export let redirectTo = "/auth/login";

    let isAllowed = false;
    let isChecking = true;

    onMount(() => {
        const unsubAuth = isAuthenticated.subscribe((authed) => {
            if (!authed) {
                goto(redirectTo);
                return;
            }

            if (requireSuperAdmin) {
                const unsubSuper = isSuperAdmin.subscribe((isSuper) => {
                    if (!isSuper) {
                        goto("/dashboard");
                    } else {
                        isAllowed = true;
                    }
                    isChecking = false;
                });
                return () => unsubSuper();
            } else {
                isAllowed = true;
                isChecking = false;
            }
        });

        return () => unsubAuth();
    });
</script>

{#if isChecking}
    <div class="flex items-center justify-center min-h-[200px]">
        <div class="text-slate-400">Checking access...</div>
    </div>
{:else if isAllowed}
    <slot />
{:else}
    <div class="flex items-center justify-center min-h-[200px]">
        <div class="text-red-400">Access denied</div>
    </div>
{/if}
