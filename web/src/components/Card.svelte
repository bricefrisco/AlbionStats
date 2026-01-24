<script>
	export let title = '';
	export let value = '';
	export let subtitle = '';
	export let meta = '';
	export let subtitleClass = '';
</script>

<div class="rounded-lg border border-gray-200 bg-gray-50 p-3 dark:border-neutral-800 dark:bg-neutral-900">
	<div class="flex items-start justify-between gap-4">
		<div class="flex-1">
			<p class="text-[11px] font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
				{title}
			</p>
			{#if meta}
				<p class="mt-2 text-xs text-gray-600 dark:text-gray-400">{meta}</p>
			{/if}
			<p class="mt-[-2px] text-xl font-semibold text-gray-900 dark:text-white">{value}</p>
		</div>
		{#if subtitle}
			<p class="text-base font-semibold {subtitleClass}">{subtitle}</p>
		{/if}
	</div>
</div>
