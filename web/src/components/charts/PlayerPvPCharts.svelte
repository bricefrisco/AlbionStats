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

	async function fetchPvpData() {
		try {
			loading = true;
			error = null;

			const response = await fetch(
				`https://albionstats.bricefrisco.com/api/metrics/pvp/${region}/${playerId}`
			);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			chartData = await response.json();
		} catch (err) {
			error = err.message;
			console.error('Failed to fetch PvP data:', err);
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
		fetchPvpData();
	});
</script>

<div class="space-y-6">
	{#if loading}
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
			<div>
				<ChartLoading />
			</div>
			<div>
				<ChartLoading />
			</div>
			<div class="lg:col-span-2">
				<ChartLoading />
			</div>
		</div>
	{:else if error}
		<div
			class="rounded-lg border border-red-200 bg-red-50 p-4 dark:border-red-800 dark:bg-red-900/20"
		>
			<ChartError {error} />
		</div>
	{:else if chartData}
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
			<!-- Kill Fame Chart -->
			<div>
				<SubHeader title="Kill Fame" classes="mb-4" />
				{#if chartData.kill_fame && chartData.kill_fame.length > 0}
					<Chart
						timestamps={chartData.timestamps}
						values={chartData.kill_fame}
						label="Kill Fame"
						height="h-40"
					/>
				{:else}
					<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
						No kill fame data available
					</div>
				{/if}
			</div>

			<!-- Death Fame Chart -->
			<div>
				<SubHeader title="Death Fame" classes="mb-4" />
				{#if chartData.death_fame && chartData.death_fame.length > 0}
					<Chart
						timestamps={chartData.timestamps}
						values={chartData.death_fame}
						label="Death Fame"
						height="h-40"
					/>
				{:else}
					<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
						No death fame data available
					</div>
				{/if}
			</div>

			<!-- Fame Ratio Chart -->
			<div class="lg:col-span-2">
				<SubHeader title="Fame Ratio" classes="mb-4" />
				{#if chartData.fame_ratio && chartData.fame_ratio.filter((r) => r !== null).length > 0}
					<Chart
						timestamps={chartData.timestamps}
						values={chartData.fame_ratio}
						label="Fame Ratio"
						height="h-40"
					/>
				{:else}
					<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
						No fame ratio data available
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
