<script>
	let { children, class: className = '' } = $props();
</script>

<td class="px-4 py-3 {className}">
	{@render children?.()}
</td>
