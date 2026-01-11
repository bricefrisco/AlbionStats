<script>
	import { onMount } from 'svelte';
	import Chart from './Chart.svelte';
	import ChartLoading from './ChartLoading.svelte';
	import ChartError from './ChartError.svelte';
	import SubHeader from '../SubHeader.svelte';

	let timestamps = [];
	let values = [];
	let loading = true;
	let error = null;

	async function fetchData() {
		try {
			const response = await fetch(
				'https://albionstats.bricefrisco.com/api/metrics/snapshots?granularity=1w'
			);
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = await response.json();
			timestamps = data.timestamps || [];
			values = (data.values || []).map((value) => value * 40);
		} catch (err) {
			error = err.message;
			console.error('Failed to fetch snapshots data:', err);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchData();
	});
</script>

<div>
	<SubHeader title="Total Data Points" classes="mb-4" />
	{#if loading}
		<ChartLoading />
	{:else if error}
		<ChartError {error} />
	{:else}
		<Chart {timestamps} {values} label="Total Data Points" />
	{/if}
</div>
