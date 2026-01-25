<script>
	import Chart from './Chart.svelte';
	import ChartLoading from './ChartLoading.svelte';
	import ChartError from './ChartError.svelte';
	import SubHeader from '../SubHeader.svelte';

	let { data = null } = $props();

	let timestamps = $state(data?.timestamps || []);
	let values = $state(data?.values || []);
	let loading = $state(!data);
	let error = $state(data?.error || null);

	$effect(() => {
		if (!data) return;
		timestamps = data.timestamps || [];
		values = data.values || [];
		error = data.error || null;
		loading = false;
	});

</script>

<div>
	<SubHeader title="Players Tracked" classes="mb-4" />
	{#if loading}
		<ChartLoading />
	{:else if error}
		<ChartError {error} />
	{:else}
		<Chart {timestamps} {values} label="Players Tracked" />
	{/if}
</div>
