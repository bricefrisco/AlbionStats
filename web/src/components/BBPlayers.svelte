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
		return data.filter((player) => {
			const playerName = (player.PlayerName || '').toLowerCase();
			const guildName = (player.GuildName || '').toLowerCase();
			const allianceName = (player.AllianceName || '').toLowerCase();
			return (
				playerName.includes(query) || guildName.includes(query) || allianceName.includes(query)
			);
		});
	});
	let totalPages = $derived(Math.ceil(filteredData.length / pageSize));
	let paginatedData = $derived(
		filteredData.slice((currentPage - 1) * pageSize, currentPage * pageSize)
	);

	function buildAffiliation(player) {
		if (player.AllianceName && player.GuildName) {
			return {
				text: `[${player.AllianceName}] ${player.GuildName}`
			};
		}
		if (player.GuildName) {
			return {
				text: player.GuildName
			};
		}
		return null;
	}

	$effect(() => {
		search;
		currentPage = 1;
	});
</script>

<Filter bind:value={search} placeholder="Filter players, guilds, or alliances" />

<Table>
	{#snippet header()}
		<TableHeader class="w-16 font-semibold"></TableHeader>
		<TableHeader class="w-64 text-left font-semibold">Player</TableHeader>
		<TableHeader class="text-right font-semibold">IP</TableHeader>
		<TableHeader class="text-right font-semibold">Kills</TableHeader>
		<TableHeader class="text-right font-semibold">Deaths</TableHeader>
		<TableHeader class="hidden text-right font-semibold whitespace-nowrap lg:table-cell">
			Kill Fame
		</TableHeader>
		<TableHeader class="hidden text-right font-semibold whitespace-nowrap lg:table-cell">
			Death Fame
		</TableHeader>
		<TableHeader class="hidden text-right font-semibold lg:table-cell">Damage</TableHeader>
		<TableHeader class="hidden text-right font-semibold lg:table-cell">Heal</TableHeader>
	{/snippet}

	{#each paginatedData as player (player.PlayerName)}
		{@const affiliation = buildAffiliation(player)}
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
			<TableData class="text-gray-900 dark:text-white">
				<div class="flex flex-col">
					{#if affiliation}
						<span class="text-sm text-gray-600 dark:text-gray-400">
							{affiliation.text}
						</span>
					{/if}
					<span>
						{#if player.PlayerName}
							<a
								href={resolve(`/players/${regionState.value}/${encodeURIComponent(player.PlayerName)}`)}
								class="font-medium underline hover:text-blue-600 dark:hover:text-blue-400"
							>
								{player.PlayerName}
							</a>
						{:else}
							-
						{/if}
					</span>
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
			<TableData class="hidden text-right text-yellow-600 dark:text-yellow-400 lg:table-cell">
				{formatFame(player.KillFame)}
			</TableData>
			<TableData class="hidden text-right text-gray-500 dark:text-gray-500 lg:table-cell">
				{formatFame(player.DeathFame)}
			</TableData>
			<TableData class="hidden text-right lg:table-cell">
				{formatNumber(player.Damage)}
			</TableData>
			<TableData class="hidden text-right lg:table-cell">
				{formatNumber(player.Heal)}
			</TableData>
		</TableRow>
	{/each}
</Table>

<Pagination bind:currentPage {totalPages} />
