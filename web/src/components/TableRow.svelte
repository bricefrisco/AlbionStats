<script>
	let { children, class: className = '' } = $props();
</script>

<tr class="align-top text-gray-700 hover:bg-gray-50/60 dark:text-gray-300 dark:hover:bg-gray-800/30 {className}">
	{@render children?.()}
</tr>
