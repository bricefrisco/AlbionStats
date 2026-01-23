<script>
	let { currentPage = $bindable(1), totalPages } = $props();
</script>

{#if totalPages > 1}
	<div class="mt-4 flex items-center justify-between">
		<div class="text-sm text-gray-600 dark:text-gray-400">
			Page {currentPage} of {totalPages}
		</div>
		<div class="flex gap-2">
			<button
				type="button"
				onclick={() => (currentPage = Math.max(1, currentPage - 1))}
				disabled={currentPage === 1}
				class="cursor-pointer rounded border border-gray-300 bg-gray-100 px-3 py-1 text-sm font-medium text-gray-700 hover:bg-gray-200 focus:border-gray-400 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 dark:border-neutral-800 dark:bg-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-700 dark:focus:border-neutral-700"
			>
				Previous
			</button>
			<button
				type="button"
				onclick={() => (currentPage = Math.min(totalPages, currentPage + 1))}
				disabled={currentPage === totalPages}
				class="cursor-pointer rounded border border-gray-300 bg-gray-100 px-3 py-1 text-sm font-medium text-gray-700 hover:bg-gray-200 focus:border-gray-400 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 dark:border-neutral-800 dark:bg-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-700 dark:focus:border-neutral-700"
			>
				Next
			</button>
		</div>
	</div>
{/if}
