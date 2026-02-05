<script lang="ts" context="module">
    import { writable } from "svelte/store";

    type Toast = {
        id: number;
        message: string;
        type: "success" | "error" | "info";
    };

    export const toasts = writable<Toast[]>([]);

    let toastId = 0;

    export function showToast(
        message: string,
        type: "success" | "error" | "info" = "info",
    ) {
        const id = toastId++;
        toasts.update((t) => [...t, { id, message, type }]);

        setTimeout(() => {
            toasts.update((t) => t.filter((toast) => toast.id !== id));
        }, 3000);
    }
</script>

<script lang="ts">
    function getToastRT(type: string) {
        const rts = {
            success: "bg-green-100 border-green-500 text-green-800",
            error: "bg-red-100 border-red-500 text-red-800",
            info: "bg-gray-100 border-gray-400 text-gray-800",
        };
        return rts[type as keyof typeof rts] || rts.info;
    }

    function removeToast(id: number) {
        toasts.update((t) => t.filter((toast) => toast.id !== id));
    }
</script>

<div class="fixed top-4 right-4 z-50 space-y-2">
    {#each $toasts as toast (toast.id)}
        <div
            class="border px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 min-w-[300px] {getToastRT(
                toast.type,
            )}"
            role="alert"
        >
            <div class="flex-1">{toast.message}</div>
            <button
                on:click={() => removeToast(toast.id)}
                class="text-current hover:opacity-70"
                aria-label="Close"
            >
                ×
            </button>
        </div>
    {/each}
</div>
