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

	async function fetchPveData() {
		try {
			loading = true;
			error = null;

			const response = await fetch(
				`https://albionstats.bricefrisco.com/api/metrics/pve/${region}/${playerId}`
			);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			chartData = await response.json();
		} catch (err) {
			error = err.message;
			console.error('Failed to fetch PvE data:', err);
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
		fetchPveData();
	});
</script>

<div class="space-y-6">
	{#if loading}
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
			<div><ChartLoading /></div>
			<div><ChartLoading /></div>
			<div><ChartLoading /></div>
			<div><ChartLoading /></div>
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
				<SubHeader title="Total PvE Fame" classes="mb-4" />
				{#if chartData.total && chartData.total.length > 0}
					<Chart
						timestamps={chartData.timestamps}
						values={chartData.total}
						label="Total PvE Fame"
						height="h-40"
					/>
				{:else}
					<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
						No total PvE fame data available
					</div>
				{/if}
			</div>

			<div>
				<SubHeader title="Royal PvE Fame" classes="mb-4" />
				{#if chartData.royal && chartData.royal.length > 0}
					<Chart
						timestamps={chartData.timestamps}
						values={chartData.royal}
						label="Royal PvE Fame"
						height="h-40"
					/>
				{:else}
					<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
						No royal PvE fame data available
					</div>
				{/if}
			</div>

			<div>
				<SubHeader title="Outlands PvE Fame" classes="mb-4" />
				{#if chartData.outlands && chartData.outlands.length > 0}
					<Chart
						timestamps={chartData.timestamps}
						values={chartData.outlands}
						label="Outlands PvE Fame"
						height="h-40"
					/>
				{:else}
					<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
						No outlands PvE fame data available
					</div>
				{/if}
			</div>

			<div>
				<SubHeader title="Avalon PvE Fame" classes="mb-4" />
				{#if chartData.avalon && chartData.avalon.length > 0}
					<Chart
						timestamps={chartData.timestamps}
						values={chartData.avalon}
						label="Avalon PvE Fame"
						height="h-40"
					/>
				{:else}
					<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
						No avalon PvE fame data available
					</div>
				{/if}
			</div>

			<div>
				<SubHeader title="Hellgate PvE Fame" classes="mb-4" />
				{#if chartData.hellgate && chartData.hellgate.length > 0}
					<Chart
						timestamps={chartData.timestamps}
						values={chartData.hellgate}
						label="Hellgate PvE Fame"
						height="h-40"
					/>
				{:else}
					<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
						No hellgate PvE fame data available
					</div>
				{/if}
			</div>

			<div>
				<SubHeader title="Corrupted PvE Fame" classes="mb-4" />
				{#if chartData.corrupted && chartData.corrupted.length > 0}
					<Chart
						timestamps={chartData.timestamps}
						values={chartData.corrupted}
						label="Corrupted PvE Fame"
						height="h-40"
					/>
				{:else}
					<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
						No corrupted PvE fame data available
					</div>
				{/if}
			</div>

			<div>
				<SubHeader title="Mists PvE Fame" classes="mb-4" />
				{#if chartData.mists && chartData.mists.length > 0}
					<Chart
						timestamps={chartData.timestamps}
						values={chartData.mists}
						label="Mists PvE Fame"
						height="h-40"
					/>
				{:else}
					<div class="flex h-40 items-center justify-center text-gray-500 dark:text-gray-400">
						No mists PvE fame data available
					</div>
				{/if}
			</div>

			<div></div>
		</div>
	{/if}
</div>
