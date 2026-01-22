<script>
	import Select from './Select.svelte';
	import { regionState } from '$lib/regionState.svelte';

	const regions = [
		{ label: 'Americas', value: 'americas' },
		{ label: 'Europe', value: 'europe' },
		{ label: 'Asia', value: 'asia' }
	];
</script>

<Select bind:value={regionState.value} options={regions} />
