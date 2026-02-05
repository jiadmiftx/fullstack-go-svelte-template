<script lang="ts" generics="T">
    type Column<T> = {
        key: keyof T;
        label: string;
        sortable?: boolean;
        render?: (value: any, item: T) => string;
    };

    export let data: T[] = [];
    export let columns: Column<T>[] = [];
    export let actions: Array<{
        label: string;
        icon?: string;
        variant?: "primary" | "danger" | "secondary";
        onClick: (item: T) => void;
        isHidden?: (item: T) => boolean;
    }> = [];
    export let emptyMessage = "No data available";
    export let searchPlaceholder = "Search...";
    export let searchable = true;
    export let pageSize = 10;
    export let theme: "light" | "dark" = "light";

    let searchTerm = "";
    let currentPage = 1;
    let sortKey: keyof T | null = null;
    let sortDirection: "asc" | "desc" = "asc";

    // Theme rts
    $: containerRT =
        theme === "dark"
            ? "bg-slate-800 rounded-xl border border-slate-700 overflow-hidden"
            : "bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden";
    $: searchBarBgRT =
        theme === "dark"
            ? "border-b border-slate-700 bg-slate-800/50"
            : "border-b border-gray-200 bg-gray-50";
    $: inputRT =
        theme === "dark"
            ? "w-full pl-10 pr-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            : "w-full pl-10 pr-4 py-2 bg-white border border-gray-300 rounded-lg text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent";
    $: iconRT = theme === "dark" ? "text-slate-400" : "text-gray-500";
    $: clearBtnRT =
        theme === "dark"
            ? "text-slate-400 hover:text-white"
            : "text-gray-500 hover:text-gray-700";
    $: emptyIconRT = theme === "dark" ? "text-slate-500" : "text-gray-500";
    $: emptyTitleRT = theme === "dark" ? "text-slate-300" : "text-gray-700";
    $: emptySubRT = theme === "dark" ? "text-slate-500" : "text-gray-500";
    $: theadRT =
        theme === "dark"
            ? "bg-slate-700/50 border-b border-slate-600"
            : "bg-gray-100 border-b border-gray-200";
    $: theadStyle = theme === "dark" ? "" : "background-color: #f3f4f6;";
    $: thRT = theme === "dark" ? "text-slate-300" : "text-gray-800";
    $: thStyle = theme === "dark" ? "" : "color: #1f2937;";
    $: tbodyDividerRT =
        theme === "dark"
            ? "divide-y divide-slate-700"
            : "divide-y divide-gray-200";
    $: trHoverRT =
        theme === "dark" ? "hover:bg-slate-700/30" : "hover:bg-gray-50";
    $: tdRT = theme === "dark" ? "text-slate-200" : "text-gray-900";
    $: tdStyle = theme === "dark" ? "" : "color: #111827;";
    $: sortArrowRT = theme === "dark" ? "text-slate-500" : "text-gray-500";
    $: paginationBorderRT =
        theme === "dark"
            ? "border-t border-slate-700 bg-slate-800/50"
            : "border-t border-gray-200 bg-gray-50";
    $: paginationTextRT =
        theme === "dark" ? "text-slate-400" : "text-gray-700";
    $: paginationHighlightRT =
        theme === "dark" ? "text-slate-200" : "text-gray-900";
    $: paginationBtnRT =
        theme === "dark"
            ? "bg-slate-700 border border-slate-600 text-slate-300 hover:bg-slate-600"
            : "bg-white border border-gray-300 text-gray-700 hover:bg-gray-50";
    $: pageNumRT =
        theme === "dark"
            ? "bg-slate-700 text-slate-300 border border-slate-600 hover:bg-slate-600"
            : "bg-white text-gray-700 border border-gray-300 hover:bg-gray-50";

    // Search filtering
    $: filteredData = searchTerm
        ? data.filter((item) => {
              return columns.some((column) => {
                  const value = item[column.key];
                  return value
                      ?.toString()
                      .toLowerCase()
                      .includes(searchTerm.toLowerCase());
              });
          })
        : data;

    // Sorting
    $: sortedData = sortKey
        ? [...filteredData].sort((a, b) => {
              const aVal = a[sortKey!];
              const bVal = b[sortKey!];
              if (aVal < bVal) return sortDirection === "asc" ? -1 : 1;
              if (aVal > bVal) return sortDirection === "asc" ? 1 : -1;
              return 0;
          })
        : filteredData;

    // Pagination
    $: totalPages = Math.ceil(sortedData.length / pageSize);
    $: paginatedData = sortedData.slice(
        (currentPage - 1) * pageSize,
        currentPage * pageSize,
    );

    // Reset page when search changes
    $: if (searchTerm) {
        currentPage = 1;
    }

    function handleSort(column: Column<T>) {
        if (!column.sortable) return;

        if (sortKey === column.key) {
            sortDirection = sortDirection === "asc" ? "desc" : "asc";
        } else {
            sortKey = column.key;
            sortDirection = "asc";
        }
    }

    function goToPage(page: number) {
        currentPage = Math.max(1, Math.min(page, totalPages));
    }

    function getActionButtonRT(variant: string = "secondary") {
        if (theme === "dark") {
            const rts = {
                primary: "bg-blue-600 hover:bg-blue-700 text-white",
                danger: "bg-red-600 hover:bg-red-700 text-white",
                secondary: "bg-slate-600 hover:bg-slate-500 text-white",
            };
            return (
                rts[variant as keyof typeof rts] || rts.secondary
            );
        } else {
            const rts = {
                primary: "bg-purple-600 hover:bg-purple-700 text-white",
                danger: "bg-red-50 hover:bg-red-100 text-red-600 border border-red-200",
                secondary:
                    "bg-white hover:bg-gray-50 text-gray-700 border border-gray-300",
            };
            return (
                rts[variant as keyof typeof rts] || rts.secondary
            );
        }
    }
</script>

<div class={containerRT}>
    <!-- Search Bar -->
    {#if searchable && data.length > 0}
        <div class="px-6 py-4 {searchBarBgRT}">
            <div class="relative">
                <input
                    type="text"
                    bind:value={searchTerm}
                    placeholder={searchPlaceholder}
                    class={inputRT}
                />
                <span
                    class="absolute left-3 top-1/2 -translate-y-1/2 {iconRT}"
                >
                    <svg
                        class="w-4 h-4"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="1.5"
                        ><path
                            d="M17.5 17.5L22 22"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        /><path
                            d="M20 11C20 6.02944 15.9706 2 11 2C6.02944 2 2 6.02944 2 11C2 15.9706 6.02944 20 11 20C15.9706 20 20 15.9706 20 11Z"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        /></svg
                    >
                </span>
                {#if searchTerm}
                    <button
                        on:click={() => (searchTerm = "")}
                        class="absolute right-3 top-1/2 -translate-y-1/2 {clearBtnRT}"
                    >
                        ✕
                    </button>
                {/if}
            </div>
        </div>
    {/if}

    <!-- Table -->
    {#if paginatedData.length === 0}
        <div class="empty-state py-12 text-center">
            <div class="mb-4">
                <svg
                    class="w-16 h-16 mx-auto {emptyIconRT}"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                    ><path
                        d="M7 18V15M12 18V12M17 18V9"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    /><path
                        d="M2.5 12C2.5 7.52166 2.5 5.28249 3.89124 3.89124C5.28249 2.5 7.52166 2.5 12 2.5C16.4783 2.5 18.7175 2.5 20.1088 3.89124C21.5 5.28249 21.5 7.52166 21.5 12C21.5 16.4783 21.5 18.7175 20.1088 20.1088C18.7175 21.5 16.4783 21.5 12 21.5C7.52166 21.5 5.28249 21.5 3.89124 20.1088C2.5 18.7175 2.5 16.4783 2.5 12Z"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    /></svg
                >
            </div>
            <p class="text-lg font-medium {emptyTitleRT} mb-1">
                {emptyMessage}
            </p>
            {#if searchTerm}
                <p class={emptySubRT}>Try adjusting your search</p>
            {/if}
        </div>
    {:else}
        <div class="overflow-x-auto">
            <table class="w-full">
                <thead class={theadRT} style={theadStyle}>
                    <tr>
                        {#each columns as column}
                            <th
                                class="px-6 py-3 text-left text-xs font-semibold {thRT} uppercase tracking-wider"
                                style={thStyle}
                                rt:cursor-pointer={column.sortable}
                                on:click={() => handleSort(column)}
                            >
                                <div class="flex items-center gap-2">
                                    {column.label}
                                    {#if column.sortable && sortKey === column.key}
                                        <span class="text-xs">
                                            {sortDirection === "asc"
                                                ? "↑"
                                                : "↓"}
                                        </span>
                                    {:else if column.sortable}
                                        <span class="text-xs {sortArrowRT}"
                                            >↕</span
                                        >
                                    {/if}
                                </div>
                            </th>
                        {/each}
                        {#if actions.length > 0}
                            <th
                                class="px-6 py-3 text-right text-xs font-semibold {thRT} uppercase tracking-wider"
                                style={thStyle}>Actions</th
                            >
                        {/if}
                    </tr>
                </thead>
                <tbody class={tbodyDividerRT}>
                    {#each paginatedData as item}
                        <tr class="{trHoverRT} transition-colors">
                            {#each columns as column}
                                <td
                                    class="px-6 py-4 text-sm {tdRT}"
                                    style={tdStyle}
                                >
                                    {#if column.render}
                                        {@html column.render(
                                            item[column.key],
                                            item,
                                        )}
                                    {:else}
                                        {item[column.key]}
                                    {/if}
                                </td>
                            {/each}
                            {#if actions.length > 0}
                                <td class="px-6 py-4">
                                    <div class="flex justify-end gap-2">
                                        {#each actions as action}
                                            {#if !action.isHidden || !action.isHidden(item)}
                                                <button
                                                    on:click={() =>
                                                        action.onClick(item)}
                                                    class="{getActionButtonRT(
                                                        action.variant,
                                                    )} px-3 py-1.5 rounded-lg text-sm font-medium transition-all duration-200"
                                                >
                                                    {action.label}
                                                </button>
                                            {/if}
                                        {/each}
                                    </div>
                                </td>
                            {/if}
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>

        <!-- Pagination -->
        {#if totalPages > 1}
            <div class="px-6 py-4 {paginationBorderRT}">
                <div class="flex items-center justify-between">
                    <div class="text-sm {paginationTextRT}">
                        Showing <span
                            class="font-medium {paginationHighlightRT}"
                            >{(currentPage - 1) * pageSize + 1}</span
                        >
                        to
                        <span class="font-medium {paginationHighlightRT}"
                            >{Math.min(
                                currentPage * pageSize,
                                sortedData.length,
                            )}</span
                        >
                        of
                        <span class="font-medium {paginationHighlightRT}"
                            >{sortedData.length}</span
                        > results
                    </div>

                    <div class="flex items-center gap-2">
                        <button
                            on:click={() => goToPage(currentPage - 1)}
                            disabled={currentPage === 1}
                            class="px-3 py-1.5 {paginationBtnRT} rounded-lg text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                        >
                            ← Previous
                        </button>

                        <div class="flex gap-1">
                            {#each Array(totalPages) as _, i}
                                {@const pageNum = i + 1}
                                {#if pageNum === 1 || pageNum === totalPages || (pageNum >= currentPage - 1 && pageNum <= currentPage + 1)}
                                    <button
                                        on:click={() => goToPage(pageNum)}
                                        class="w-10 h-10 flex items-center justify-center rounded-lg text-sm font-medium transition-colors
                                            {pageNum === currentPage
                                            ? 'bg-blue-600 text-white'
                                            : pageNumRT}"
                                    >
                                        {pageNum}
                                    </button>
                                {:else if pageNum === currentPage - 2 || pageNum === currentPage + 2}
                                    <span
                                        class="w-10 h-10 flex items-center justify-center {paginationTextRT}"
                                        >...</span
                                    >
                                {/if}
                            {/each}
                        </div>

                        <button
                            on:click={() => goToPage(currentPage + 1)}
                            disabled={currentPage === totalPages}
                            class="px-3 py-1.5 {paginationBtnRT} rounded-lg text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                        >
                            Next →
                        </button>
                    </div>
                </div>
            </div>
        {/if}
    {/if}
</div>
