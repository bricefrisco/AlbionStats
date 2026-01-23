<script>
	import Table from './Table.svelte';
	import TableHeader from './TableHeader.svelte';
	import TableRow from './TableRow.svelte';
	import TableData from './TableData.svelte';
	import Pagination from './Pagination.svelte';
	import { formatNumber, formatFame } from '$lib/utils';

	let { data = [] } = $props();

	let currentPage = $state(1);
	const pageSize = 20;

	let totalPages = $derived(Math.ceil(data.length / pageSize));
	let paginatedData = $derived(data.slice((currentPage - 1) * pageSize, currentPage * pageSize));
</script>

<Table>
	{#snippet header()}
		<TableHeader class="text-left font-semibold">Alliance</TableHeader>
		<TableHeader class="text-right font-semibold">Players</TableHeader>
		<TableHeader class="text-right font-semibold">Kills</TableHeader>
		<TableHeader class="text-right font-semibold">Deaths</TableHeader>
		<TableHeader class="text-right font-semibold whitespace-nowrap">Kill Fame</TableHeader>
		<TableHeader class="text-right font-semibold whitespace-nowrap">Death Fame</TableHeader>
		<TableHeader class="text-right font-semibold">IP</TableHeader>
	{/snippet}

	{#each paginatedData as alliance (alliance.AllianceName)}
		<TableRow>
			<TableData class="font-medium text-gray-900 dark:text-white">
				{alliance.AllianceName || '-'}
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
			<TableData class="text-right">
				{formatNumber(alliance.IP)}
			</TableData>
		</TableRow>
	{/each}
</Table>

<Pagination bind:currentPage {totalPages} />
