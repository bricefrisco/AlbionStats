<script>
	export let title = '';
	export let value = '';
	export let subtitle = '';
</script>

<div class="rounded-lg border border-gray-200 bg-gray-50 p-4 dark:border-neutral-800 dark:bg-neutral-900">
	<p class="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
		{title}
	</p>
	<p class="mt-2 text-xl font-semibold text-gray-900 dark:text-white">{value}</p>
	{#if subtitle}
		<p class="mt-1 text-sm text-gray-600 dark:text-gray-400">{subtitle}</p>
	{/if}
</div>
