<script>
	import Chart from './Chart.svelte';
	import ChartLoading from './ChartLoading.svelte';
	import ChartError from './ChartError.svelte';
	import SubHeader from '../SubHeader.svelte';

	let { data = null } = $props();

	let timestamps = $derived(data?.timestamps ?? []);
	let americas = $derived(data?.americas ?? []);
	let europe = $derived(data?.europe ?? []);
	let asia = $derived(data?.asia ?? []);
	let loading = $derived(data == null);
	let error = $derived(data?.error ?? null);

	const totals = $derived(
		americas.map((value, index) => value + (europe[index] ?? 0) + (asia[index] ?? 0))
	);

	const datasets = $derived([
		{
			label: 'Total',
			data: Array.from(totals || []),
			borderColor: 'rgb(148, 163, 184)',
			backgroundColor: 'rgba(148, 163, 184, 0.18)'
		},
		{
			label: 'Americas',
			data: Array.from(americas || []),
			borderColor: 'rgb(59, 130, 246)',
			backgroundColor: 'rgba(59, 130, 246, 0.18)'
		},
		{
			label: 'Europe',
			data: Array.from(europe || []),
			borderColor: 'rgb(16, 185, 129)',
			backgroundColor: 'rgba(16, 185, 129, 0.18)'
		},
		{
			label: 'Asia',
			data: Array.from(asia || []),
			borderColor: 'rgb(234, 179, 8)',
			backgroundColor: 'rgba(234, 179, 8, 0.18)'
		}
	]);
</script>

<div>
	<SubHeader title="Daily Active Users" classes="mb-4" />
	{#if loading}
		<ChartLoading />
	{:else if error}
		<ChartError {error} />
	{:else}
		<Chart {timestamps} {datasets} showLegend={true} />
	{/if}
</div>
