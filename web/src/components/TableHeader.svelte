<script>
	let { children, class: className = '' } = $props();
</script>

<th class="px-4 py-3 {className}">
	{@render children?.()}
</th>
