<script>
	// Props
	export let error = '';
</script>

<div class="flex h-80 items-center justify-center">
	<div class="flex flex-col items-center space-y-4">
		<!-- Error icon -->
		<div
			class="flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/20"
		>
			<svg
				class="h-6 w-6 text-red-600 dark:text-red-400"
				fill="none"
				stroke="currentColor"
				viewBox="0 0 24 24"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"
				/>
			</svg>
		</div>
		<div class="text-center">
			<div class="mb-1 font-medium text-red-600 dark:text-red-400">Failed to load data</div>
			{#if error}
				<div class="text-sm text-red-500 dark:text-red-400">
					{error}
				</div>
			{/if}
		</div>
	</div>
</div>
