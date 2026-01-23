<script>
	let { children, class: className = '' } = $props();
</script>

<td class="align-middle px-4 py-3 {className}">
	{@render children?.()}
</td>
