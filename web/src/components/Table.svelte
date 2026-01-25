<script>
	let { children, header, classes } = $props();
</script>

<div
	class="overflow-x-auto rounded-lg border border-gray-200/60 bg-transparent shadow-sm dark:border-gray-800/60 dark:bg-transparent {classes}"
>
	<table class="w-full table-auto text-sm lg:table-fixed">
		<thead class="bg-gray-50/60 dark:bg-gray-800/40">
			<tr class="text-sm capitalize tracking-wide text-gray-600 dark:text-gray-300">
				{@render header?.()}
			</tr>
		</thead>
		<tbody class="divide-y divide-gray-200/60 dark:divide-gray-700/60">
			{@render children?.()}
		</tbody>
	</table>
</div>
