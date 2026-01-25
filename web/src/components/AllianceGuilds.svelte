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
		data.map((guild, index) => ({
			...guild,
			pos: index + 1
		}))
	);

	const filteredData = $derived.by(() => {
		const query = search.trim().toLowerCase();
		if (!query) return indexedData;
		return indexedData.filter((guild) => (guild.Name || '').toLowerCase().includes(query));
	});
	let totalPages = $derived(Math.ceil(filteredData.length / pageSize));
	let paginatedData = $derived(
		filteredData.slice((currentPage - 1) * pageSize, currentPage * pageSize)
	);

	$effect(() => {
		search;
		currentPage = 1;
	});
</script>

<Filter bind:value={search} placeholder="Filter guilds" />

<Table>
	{#snippet header()}
		<TableHeader class="w-12 text-right font-semibold">Pos.</TableHeader>
		<TableHeader class="w-64 text-left font-semibold">Guild</TableHeader>
		<TableHeader class="text-right font-semibold whitespace-nowrap">Battles</TableHeader>
		<TableHeader class="text-right font-semibold whitespace-nowrap">Max Players</TableHeader>
		<TableHeader class="text-right font-semibold">Kills</TableHeader>
		<TableHeader class="text-right font-semibold">Deaths</TableHeader>
		<TableHeader class="text-right font-semibold whitespace-nowrap">Kill Fame</TableHeader>
		<TableHeader class="text-right font-semibold whitespace-nowrap">Death Fame</TableHeader>
	{/snippet}

	{#each paginatedData as guild (guild.Name)}
		<TableRow>
			<TableData class="text-right text-gray-500 dark:text-gray-400">
				{guild.pos}
			</TableData>
			<TableData class="font-medium text-gray-900 dark:text-white">
				{#if guild.Name}
					<a
						href={resolve(`/guilds/${regionState.value}/${encodeURIComponent(guild.Name)}`)}
						class="underline hover:text-blue-600 dark:hover:text-blue-400"
					>
						{guild.Name}
					</a>
				{:else}
					-
				{/if}
			</TableData>
			<TableData class="text-right text-gray-600 dark:text-gray-400">
				{formatNumber(guild.NumBattles)}
			</TableData>
			<TableData class="text-right text-blue-600 dark:text-blue-400">
				{formatNumber(guild.MaxPlayerCount)}
			</TableData>
			<TableData class="text-right text-red-600 dark:text-red-400">
				{formatNumber(guild.Kills)}
			</TableData>
			<TableData class="text-right text-gray-600 dark:text-gray-400">
				{formatNumber(guild.Deaths)}
			</TableData>
			<TableData class="text-right text-yellow-600 dark:text-yellow-400">
				{formatNumber(guild.KillFame)}
			</TableData>
			<TableData class="text-right text-gray-500 dark:text-gray-500">
				{formatNumber(guild.DeathFame)}
			</TableData>
		</TableRow>
	{/each}
</Table>

<Pagination bind:currentPage {totalPages} />
