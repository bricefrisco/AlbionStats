<script>
	import Table from './Table.svelte';
	import TableHeader from './TableHeader.svelte';
	import TableRow from './TableRow.svelte';
	import TableData from './TableData.svelte';
	import Pagination from './Pagination.svelte';
	import Filter from './Filter.svelte';
	import { formatNumber, formatFame } from '$lib/utils';
	import { resolve } from '$app/paths';
	import { regionState } from '$lib/regionState.svelte';

	let { data = [] } = $props();

	let currentPage = $state(1);
	const pageSize = 10;
	let search = $state('');

	const filteredData = $derived.by(() => {
		const query = search.trim().toLowerCase();
		if (!query) return data;
		return data.filter((alliance) =>
			(alliance.AllianceName || '').toLowerCase().includes(query)
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
</script>

<Filter bind:value={search} placeholder="Filter alliances" />

<Table>
	{#snippet header()}
		<TableHeader class="text-left font-semibold">Alliance</TableHeader>
		<TableHeader class="text-right font-semibold">Players</TableHeader>
		<TableHeader class="text-right font-semibold">Kills</TableHeader>
		<TableHeader class="text-right font-semibold">Deaths</TableHeader>
		<TableHeader class="text-right font-semibold whitespace-nowrap">Kill Fame</TableHeader>
		<TableHeader class="text-right font-semibold whitespace-nowrap">Death Fame</TableHeader>
		<TableHeader class="hidden text-right font-semibold lg:table-cell">IP</TableHeader>
	{/snippet}

	{#each paginatedData as alliance (alliance.AllianceName)}
		<TableRow>
			<TableData class="font-medium text-gray-900 dark:text-white">
				{#if alliance.AllianceName}
					<a
						href={resolve(`/alliances/${regionState.value}/${encodeURIComponent(alliance.AllianceName)}`)}
						class="underline hover:text-blue-600 dark:hover:text-blue-400"
					>
						{alliance.AllianceName}
					</a>
				{:else}
					-
				{/if}
			</TableData>
			<TableData class="text-right text-blue-600 dark:text-blue-400">
				{formatNumber(alliance.PlayerCount)}
			</TableData>
			<TableData class="text-right text-red-600 dark:text-red-400">
				{formatNumber(alliance.Kills)}
			</TableData>
			<TableData class="text-right text-gray-600 dark:text-gray-400">
				{formatNumber(alliance.Deaths)}
			</TableData>
			<TableData class="text-right text-yellow-600 dark:text-yellow-400">
				{formatFame(alliance.KillFame)}
			</TableData>
			<TableData class="text-right text-gray-500 dark:text-gray-500">
				{formatFame(alliance.DeathFame)}
			</TableData>
			<TableData class="hidden text-right lg:table-cell">
				{formatNumber(alliance.IP)}
			</TableData>
		</TableRow>
	{/each}
</Table>

<Pagination bind:currentPage {totalPages} />
