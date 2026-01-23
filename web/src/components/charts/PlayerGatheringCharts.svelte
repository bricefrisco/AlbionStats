<script>
	import Chart from './Chart.svelte';
	import ChartLoading from './ChartLoading.svelte';
	import ChartError from './ChartError.svelte';
	import SubHeader from '../SubHeader.svelte';

	// Props
	let { region = '', playerId = '', data = null } = $props();

	// Data state
	let chartData = $state(data?.data || null);
	let loading = $state(!data);
	let error = $state(data?.error || null);
	let lastKey = $state('');

	async function fetchGatheringData() {
		try {
			loading = true;
			error = null;

			const response = await fetch(
				`https://albionstats.bricefrisco.com/api/metrics/gathering/${region}/${playerId}`
			);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			chartData = await response.json();
		} catch (err) {
			error = err.message;
			console.error('Failed to fetch gathering data:', err);
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (data) {
			chartData = data.data || null;
			error = data.error || null;
			loading = false;
			return;
		}

		if (!region || !playerId) return;
		const key = `${region}:${playerId}`;
		if (key === lastKey) return;
		lastKey = key;
		fetchGatheringData();
	});
</script>

<div class="space-y-6">
	{#if loading}
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
			<div><ChartLoading /></div>
			<div><ChartLoading /></div>
			<div><ChartLoading /></div>
			<div><ChartLoading /></div>
		</div>
	{:else if error}
		<div
			class="rounded-lg border border-red-200 bg-red-50 p-4 dark:border-red-800 dark:bg-red-900/20"
		>
			<ChartError {error} />
		</div>
	{:else if chartData}
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
			<div>
				<SubHeader title="Total Gathering Fame" classes="mb-4" />
				{#if chartData.total && chartData.total.length > 0}
					<Chart
						timestamps={chartData.timestamps}
						values={chartData.total}
						label="Total Gathering Fame"
						height="h-40"
					/>
				{:else}
					<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
						No total gathering fame data available
					</div>
				{/if}
			</div>

			<div>
				<SubHeader title="Royal Gathering Fame" classes="mb-4" />
				{#if chartData.royal && chartData.royal.length > 0}
					<Chart
						timestamps={chartData.timestamps}
						values={chartData.royal}
						label="Royal Gathering Fame"
						height="h-40"
					/>
				{:else}
					<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
						No royal gathering fame data available
					</div>
				{/if}
			</div>

			<div>
				<SubHeader title="Outlands Gathering Fame" classes="mb-4" />
				{#if chartData.outlands && chartData.outlands.length > 0}
					<Chart
						timestamps={chartData.timestamps}
						values={chartData.outlands}
						label="Outlands Gathering Fame"
						height="h-40"
					/>
				{:else}
					<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
						No outlands gathering fame data available
					</div>
				{/if}
			</div>

			<div>
				<SubHeader title="Avalon Gathering Fame" classes="mb-4" />
				{#if chartData.avalon && chartData.avalon.length > 0}
					<Chart
						timestamps={chartData.timestamps}
						values={chartData.avalon}
						label="Avalon Gathering Fame"
						height="h-40"
					/>
				{:else}
					<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
						No avalon gathering fame data available
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
