<script>
	// Loading component for charts
</script>

<div class="flex h-80 items-center justify-center">
	<div class="flex flex-col items-center space-y-4">
		<!-- Simple loading spinner -->
		<div class="relative">
			<div class="w-12 h-12 border-4 border-gray-200 dark:border-gray-700 rounded-full animate-spin border-t-teal-500"></div>
		</div>
		<div class="text-gray-500 dark:text-gray-400 text-sm font-medium">
			Loading...
		</div>
	</div>
</div>

<style>
	/* Custom spinner animation */
	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}

	.animate-spin {
		animation: spin 1s linear infinite;
	}
</style>
