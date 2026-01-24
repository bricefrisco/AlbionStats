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

	function formatAffiliation(player) {
		if (!player) return '';
		if (player.AllianceName && player.GuildName) {
			return `[${player.AllianceName}] ${player.GuildName}`;
		}
		if (player.GuildName) {
			return player.GuildName;
		}
		return '';
	}
</script>

<Table>
	{#snippet header()}
		<TableHeader class="w-16 font-semibold"></TableHeader>
		<TableHeader class="w-64 text-left font-semibold">Player</TableHeader>
		<TableHeader class="text-right font-semibold">IP</TableHeader>
		<TableHeader class="text-right font-semibold">Kills</TableHeader>
		<TableHeader class="text-right font-semibold">Deaths</TableHeader>
		<TableHeader class="text-right font-semibold whitespace-nowrap">Kill Fame</TableHeader>
		<TableHeader class="text-right font-semibold whitespace-nowrap">Death Fame</TableHeader>
		<TableHeader class="text-right font-semibold">Damage</TableHeader>
		<TableHeader class="text-right font-semibold">Heal</TableHeader>
	{/snippet}

	{#each paginatedData as player (player.PlayerName)}
		<TableRow>
			<TableData class="p-2 flex items-center justify-center">
				{#if player.Weapon}
					<img
						src="https://render.albiononline.com/v1/item/{player.Weapon}"
						alt={player.Weapon}
						class="h-10 w-10 min-w-10"
					/>
				{/if}
			</TableData>
			<TableData class="font-medium text-gray-900 dark:text-white">
				<div class="flex flex-col">
					{#if formatAffiliation(player)}
						<span class="text-sm text-gray-600 dark:text-gray-400">
							{formatAffiliation(player)}
						</span>
					{/if}
					<span>{player.PlayerName}</span>
				</div>
			</TableData>
			<TableData class="text-right">
				{formatNumber(player.IP)}
			</TableData>
			<TableData class="text-right text-red-600 dark:text-red-400">
				{formatNumber(player.Kills)}
			</TableData>
			<TableData class="text-right text-gray-600 dark:text-gray-400">
				{formatNumber(player.Deaths)}
			</TableData>
			<TableData class="text-right text-yellow-600 dark:text-yellow-400">
				{formatFame(player.KillFame)}
			</TableData>
			<TableData class="text-right text-gray-500 dark:text-gray-500">
				{formatFame(player.DeathFame)}
			</TableData>
			<TableData class="text-right">
				{formatNumber(player.Damage)}
			</TableData>
			<TableData class="text-right">
				{formatNumber(player.Heal)}
			</TableData>
		</TableRow>
	{/each}
</Table>

<Pagination bind:currentPage {totalPages} />
