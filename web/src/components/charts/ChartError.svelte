<script>
	// Props
	export let error = '';
</script>

<div class="flex h-80 items-center justify-center">
	<div class="flex flex-col items-center space-y-4">
		<!-- Error icon -->
		<div class="flex items-center justify-center w-12 h-12 bg-red-100 dark:bg-red-900/20 rounded-full">
			<svg class="w-6 h-6 text-red-600 dark:text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"/>
			</svg>
		</div>
		<div class="text-center">
			<div class="text-red-600 dark:text-red-400 font-medium mb-1">
				Failed to load data
			</div>
			{#if error}
				<div class="text-red-500 dark:text-red-400 text-sm">
					{error}
				</div>
			{/if}
		</div>
	</div>
</div>
