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

			data = await response.json();
		} catch (err) {
			error = err.message;
			console.error('Failed to fetch PvP data:', err);
		} finally {
			loading = false;
		}
	}

	// Fetch data when props change
	$: if (region && playerId) {
		fetchPvpData();
	}
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
	{:else if data}
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
			<!-- Kill Fame Chart -->
			<div>
				<SubHeader title="Kill Fame" classes="mb-4" />
				{#if data.kill_fame && data.kill_fame.length > 0}
					<Chart
						timestamps={data.timestamps}
						values={data.kill_fame}
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
				{#if data.death_fame && data.death_fame.length > 0}
					<Chart
						timestamps={data.timestamps}
						values={data.death_fame}
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
				{#if data.fame_ratio && data.fame_ratio.filter((r) => r !== null).length > 0}
					<Chart
						timestamps={data.timestamps}
						values={data.fame_ratio}
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
