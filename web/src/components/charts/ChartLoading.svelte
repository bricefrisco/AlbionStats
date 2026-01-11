<script>
	// Loading component for charts
</script>

<div class="flex h-80 items-center justify-center">
	<div class="flex flex-col items-center space-y-4">
		<!-- Simple loading spinner -->
		<div class="relative">
			<div
				class="h-12 w-12 animate-spin rounded-full border-4 border-gray-200 border-t-teal-500 dark:border-gray-700"
			></div>
		</div>
		<div class="text-sm font-medium text-gray-500 dark:text-gray-400">Loading...</div>
	</div>
</div>

<style>
	/* Custom spinner animation */
	@keyframes spin {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}

	.animate-spin {
		animation: spin 1s linear infinite;
	}
</style>
