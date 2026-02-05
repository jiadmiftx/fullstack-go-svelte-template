<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";

    $: orgId = $page.params.orgId;
    $: unitId = $page.params.unitId;

    onMount(() => {
        // Redirect to dashboard
        goto(`/org/${orgId}/units/${unitId}/dashboard`, {
            replaceState: true,
        });
    });
</script>

<div class="loading">
    <p>Redirecting to dashboard...</p>
</div>

<style>
    .loading {
        display: flex;
        align-items: center;
        justify-content: center;
        height: 50vh;
        color: #64748b;
    }
</style>
