<script>
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import PlayerSearchBar from '$components/PlayerSearchBar.svelte';
	import Typography from '$components/Typography.svelte';
	import BackToTopButton from '$components/BackToTopButton.svelte';
	import Table from '$components/Table.svelte';
	import TableHeader from '$components/TableHeader.svelte';
	import TableRow from '$components/TableRow.svelte';
	import TableData from '$components/TableData.svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { formatFame, formatNumber, formatRatio, getRatioColor } from '$lib/utils.js';

	let { data } = $props();

	let searchQuery = $derived(page.url.searchParams.get('q') || '');

	function updateUrl(q) {
		searchQuery = q;
		const url = new URL(page.url);
		if (q) {
			url.searchParams.set('q', q);
		} else {
			url.searchParams.delete('q');
		}
		goto(resolve(url.pathname + url.search), {
			keepFocus: true,
			noScroll: true
		});
	}
</script>

<Page>
	<PageHeader title="Players" />
	<Typography>
		<h2>Albion Online Players. Search for a player to view their statistics. Below are top 100 players
			based on statistics pulled from battle board data based over the past 30 days.</h2>
	</Typography>

	<div class="mb-4">
		<PlayerSearchBar bind:value={searchQuery} onselect={updateUrl} placeholder="Player name" />
	</div>

	{#if data.topPlayersError}
		<Typography>
			<p class="text-red-600 dark:text-red-400">{data.topPlayersError}</p>
		</Typography>
	{:else if data.topPlayers?.length}
		<Table>
			{#snippet header()}
				<TableHeader class="w-12 text-right font-semibold">#</TableHeader>
				<TableHeader class="text-left font-semibold">Name</TableHeader>
				<TableHeader class="text-right font-semibold whitespace-nowrap">Kill Fame</TableHeader>
				<TableHeader class="text-right font-semibold whitespace-nowrap">Death Fame</TableHeader>
				<TableHeader class="text-right font-semibold whitespace-nowrap">K/D Ratio</TableHeader>
				<TableHeader class="text-right font-semibold">Kills</TableHeader>
				<TableHeader class="text-right font-semibold">Deaths</TableHeader>
			{/snippet}

			{#each data.topPlayers as player, index (player.PlayerName)}
				{@const ratioColor = getRatioColor(player.TotalKillFame, player.TotalDeathFame)}
				<TableRow>
					<TableData class="text-right text-gray-500 dark:text-gray-400">
						{index + 1}
					</TableData>
					<TableData class="font-medium text-gray-900 dark:text-white">
						{player.PlayerName}
					</TableData>
					<TableData class="text-right text-yellow-600 dark:text-yellow-400">
						{formatFame(player.TotalKillFame)}
					</TableData>
					<TableData class="text-right text-gray-600 dark:text-gray-400">
						{formatFame(player.TotalDeathFame)}
					</TableData>
					<TableData class="text-right">
						<span style={ratioColor ? `color: ${ratioColor}` : ''}>
							{formatRatio(player.TotalKillFame, player.TotalDeathFame)}
						</span>
					</TableData>
					<TableData class="text-right text-red-600 dark:text-red-400">
						{formatNumber(player.TotalKills)}
					</TableData>
					<TableData class="text-right text-gray-600 dark:text-gray-400">
						{formatNumber(player.TotalDeaths)}
					</TableData>
				</TableRow>
			{/each}
		</Table>
		<div class="mt-6 flex justify-center">
			<BackToTopButton />
		</div>
	{:else}
		<Typography>
			<p>No top players found.</p>
		</Typography>
	{/if}
</Page>
