<script>
	import { onMount } from 'svelte';
	import Chart from './Chart.svelte';
	import ChartLoading from './ChartLoading.svelte';
	import ChartError from './ChartError.svelte';
	import SubHeader from '../SubHeader.svelte';

	// Props
	export let region = '';
	export let playerId = '';

	// Data state
	let data = null;
	let loading = true;
	let error = null;

	async function fetchCraftingData() {
		try {
			loading = true;
			error = null;

			const response = await fetch(
				`https://albionstats.bricefrisco.com/api/metrics/crafting/${region}/${playerId}`
			);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			data = await response.json();
		} catch (err) {
			error = err.message;
			console.error('Failed to fetch crafting data:', err);
		} finally {
			loading = false;
		}
	}

	// Fetch data when props change
	$: if (region && playerId) {
		fetchCraftingData();
	}
</script>

<div class="space-y-6">
	{#if loading}
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
			<div><ChartLoading /></div>
		</div>
	{:else if error}
		<div
			class="rounded-lg border border-red-200 bg-red-50 p-4 dark:border-red-800 dark:bg-red-900/20"
		>
			<ChartError {error} />
		</div>
	{:else if data}
		<div>
			<SubHeader title="Total Crafting Fame" classes="mb-4" />
			{#if data.total && data.total.length > 0}
				<Chart
					timestamps={data.timestamps}
					values={data.total}
					label="Total Crafting Fame"
					height="h-40"
				/>
			{:else}
				<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
					No total crafting fame data available
				</div>
			{/if}
		</div>
	{/if}
</div>
