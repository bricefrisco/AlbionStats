<script>
	import Table from './Table.svelte';
	import TableHeader from './TableHeader.svelte';
	import TableRow from './TableRow.svelte';
	import TableData from './TableData.svelte';
	import Pagination from './Pagination.svelte';
	import Filter from './Filter.svelte';
	import { formatNumber } from '$lib/utils';
	import { resolve } from '$app/paths';
	import { regionState } from '$lib/regionState.svelte';

	let { data = [] } = $props();

	let currentPage = $state(1);
	const pageSize = 10;
	let search = $state('');

	const indexedData = $derived.by(() =>
		data.map((player, index) => ({
			...player,
			pos: index + 1
		}))
	);

	const filteredData = $derived.by(() => {
		const query = search.trim().toLowerCase();
		if (!query) return indexedData;
		return indexedData.filter((player) =>
			(player.PlayerName || '').toLowerCase().includes(query)
		);
	});
	let totalPages = $derived(Math.ceil(filteredData.length / pageSize));
	let paginatedData = $derived(
		filteredData.slice((currentPage - 1) * pageSize, currentPage * pageSize)
	);

	$effect(() => {
		search;
		currentPage = 1;
	});

	function formatDate(dateString) {
		if (!dateString) return '-';
		return new Date(dateString).toLocaleString();
	}
</script>

<Filter bind:value={search} placeholder="Filter players" />

<Table>
	{#snippet header()}
		<TableHeader class="w-12 text-right font-semibold">Pos.</TableHeader>
		<TableHeader class="text-left font-semibold">Player</TableHeader>
		<TableHeader class="w-48 text-left font-semibold whitespace-nowrap">Last Battle</TableHeader>
		<TableHeader class="text-right font-semibold">Battles</TableHeader>
		<TableHeader class="text-right font-semibold">Kills</TableHeader>
		<TableHeader class="text-right font-semibold">Deaths</TableHeader>
		<TableHeader class="text-right font-semibold whitespace-nowrap">Kill Fame</TableHeader>
		<TableHeader class="text-right font-semibold whitespace-nowrap">Death Fame</TableHeader>
	{/snippet}

	{#each paginatedData as player (player.PlayerName)}
		<TableRow>
			<TableData class="text-right text-gray-500 dark:text-gray-400">
				{player.pos}
			</TableData>
			<TableData class="font-medium text-gray-900 dark:text-white">
				{#if player.PlayerName}
					<a
						href={resolve(`/players/${regionState.value}/${encodeURIComponent(player.PlayerName)}`)}
						class="underline hover:text-blue-600 dark:hover:text-blue-400"
					>
						{player.PlayerName}
					</a>
				{:else}
					-
				{/if}
			</TableData>
			<TableData class="text-left text-gray-600 dark:text-gray-400">
				{formatDate(player.LastBattle)}
			</TableData>
			<TableData class="text-right text-blue-600 dark:text-blue-400">
				{formatNumber(player.NumBattles)}
			</TableData>
			<TableData class="text-right text-red-600 dark:text-red-400">
				{formatNumber(player.Kills)}
			</TableData>
			<TableData class="text-right text-gray-600 dark:text-gray-400">
				{formatNumber(player.Deaths)}
			</TableData>
			<TableData class="text-right text-yellow-600 dark:text-yellow-400">
				{formatNumber(player.KillFame)}
			</TableData>
			<TableData class="text-right text-gray-500 dark:text-gray-500">
				{formatNumber(player.DeathFame)}
			</TableData>
		</TableRow>
	{/each}
</Table>

<Pagination bind:currentPage {totalPages} />
