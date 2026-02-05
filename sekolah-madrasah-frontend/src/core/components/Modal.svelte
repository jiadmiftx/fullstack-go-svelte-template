<script lang="ts">
    import { createEventDispatcher } from "svelte";

    export let isOpen = false;
    export let title = "";
    export let size: "sm" | "md" | "lg" = "md";

    const dispatch = createEventDispatcher();

    const sizeRTs = {
        sm: "max-w-md",
        md: "max-w-2xl",
        lg: "max-w-4xl",
    };

    function close() {
        isOpen = false;
        dispatch("close");
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape" && isOpen) {
            close();
        }
    }

    function handleBackdropClick(e: MouseEvent) {
        if (e.target === e.currentTarget) {
            close();
        }
    }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if isOpen}
    <div
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
        on:click={handleBackdropClick}
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal-title"
    >
        <div
            class="bg-white border border-gray-200 rounded-xl p-6 shadow-xl {sizeRTs[
                size
            ]} w-full max-h-[90vh] overflow-y-auto"
        >
            <!-- Header -->
            <div
                class="flex justify-between items-center mb-6 pb-4 border-b border-gray-200"
            >
                <h2 id="modal-title" class="text-xl font-bold text-gray-900">
                    {title}
                </h2>
                <button
                    on:click={close}
                    class="text-gray-400 hover:text-gray-600 transition-colors p-1"
                    aria-label="Close modal"
                >
                    <svg
                        class="w-6 h-6"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M6 18L18 6M6 6l12 12"
                        />
                    </svg>
                </button>
            </div>

            <!-- Content -->
            <div class="mb-6 text-gray-700">
                <slot />
            </div>

            <!-- Footer -->
            {#if $$slots.footer}
                <div
                    class="flex justify-end gap-3 pt-4 border-t border-gray-200"
                >
                    <slot name="footer" />
                </div>
            {/if}
        </div>
    </div>
{/if}
