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

			data = await response.json();
		} catch (err) {
			error = err.message;
			console.error('Failed to fetch PvE data:', err);
		} finally {
			loading = false;
		}
	}

	// Fetch data when props change
	$: if (region && playerId) {
		fetchPveData();
	}
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
	{:else if data}
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
			<div>
				<SubHeader title="Total PvE Fame" classes="mb-4" />
				{#if data.total && data.total.length > 0}
					<Chart
						timestamps={data.timestamps}
						values={data.total}
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
				{#if data.royal && data.royal.length > 0}
					<Chart
						timestamps={data.timestamps}
						values={data.royal}
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
				{#if data.outlands && data.outlands.length > 0}
					<Chart
						timestamps={data.timestamps}
						values={data.outlands}
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
				{#if data.avalon && data.avalon.length > 0}
					<Chart
						timestamps={data.timestamps}
						values={data.avalon}
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
				{#if data.hellgate && data.hellgate.length > 0}
					<Chart
						timestamps={data.timestamps}
						values={data.hellgate}
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
				{#if data.corrupted && data.corrupted.length > 0}
					<Chart
						timestamps={data.timestamps}
						values={data.corrupted}
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
				{#if data.mists && data.mists.length > 0}
					<Chart
						timestamps={data.timestamps}
						values={data.mists}
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
