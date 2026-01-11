<script>
	// Props
	export let label = '';
	export let value = '';
	export let label2 = '';
	export let value2 = '';
</script>

<tr>
	<td class="py-2 pl-4 pr-4 text-gray-600 dark:text-gray-400">{label}</td>
	<td class="py-2 pr-4 text-right font-medium">{value}</td>
	<td class="py-2 pl-4 pr-4 text-gray-600 dark:text-gray-400">{label2 || '-'}</td>
	<td class="py-2 pr-4 text-right font-medium">{value2 || '-'}</td>
</tr>
