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
	import { resolve } from '$app/paths';
	import { formatFame, formatNumber, formatRatio, getRatioColor } from '$lib/utils.js';
	import { regionState } from '$lib/regionState.svelte.js';
	import Tooltip from '$components/Tooltip.svelte';
	import HelpIcon from '$components/icons/HelpIcon.svelte';

	let { data } = $props();

	let searchQuery = $derived('');
	let websiteJsonLd = $derived.by(() => {
		const origin = page.url.origin;
		const region = page.params.region;
		return JSON.stringify({
			'@context': 'https://schema.org',
			'@type': 'WebSite',
			name: 'AlbionStats',
			url: origin,
			potentialAction: {
				'@type': 'SearchAction',
				target: `${origin}/players/${region}/{search_term_string}`,
				'query-input': 'required name=search_term_string'
			}
		});
	});
	let itemListJsonLd = $derived.by(() => {
		if (!data.topPlayers?.length) return '';
		const origin = page.url.origin;
		const region = page.params.region;
		return JSON.stringify({
			'@context': 'https://schema.org',
			'@type': 'ItemList',
			name: `Top Albion Online players in ${regionState.label}`,
			itemListOrder: 'https://schema.org/ItemListOrderDescending',
			itemListElement: data.topPlayers.map((player, index) => ({
				'@type': 'ListItem',
				position: index + 1,
				name: player.PlayerName,
				url: `${origin}/players/${region}/${encodeURIComponent(player.PlayerName)}`
			}))
		});
	});
</script>

<svelte:head>
	<title>Players - AlbionStats - {regionState.label}</title>
	<meta
		name="description"
		content={`Top Albion Online players in ${regionState.label}. Search player stats, kills, deaths, and fame.`}
	/>
	<link rel="canonical" href={`${page.url.origin}${page.url.pathname}`} />
	<script type="application/ld+json">{websiteJsonLd}</script>
	{#if itemListJsonLd}
		<script type="application/ld+json">{itemListJsonLd}</script>
	{/if}
</svelte:head>

<Page>
	<PageHeader title="Players" />
	<Typography>
		<h2>Albion Online Player Statistics. Search for a player to view their stats.</h2>
		<p>Below are top 100 players based on statistics pulled from battle board data
			based over the past 30 days.</p>
		<p>Collection began on January 19th, 2026.
			<Tooltip content="Small fights (1v1's) not show up on the battle boards, so some kills from hellgates or mists do not count on this leaderboard at this time.">
				<button
					type="button"
					class="flex items-center text-gray-400 transition-colors hover:text-gray-600 dark:text-neutral-500 dark:hover:text-neutral-300"
					aria-label="Participants info"
				>
					<HelpIcon size={14} />
				</button>
			</Tooltip>
		</p>
	</Typography>

	<div class="mb-4">
		<PlayerSearchBar links={true} bind:value={searchQuery} placeholder="Player name" />
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
						<a
							href={resolve(`/players/${regionState.value}/${encodeURIComponent(player.PlayerName)}`)}
							class="hover:underline hover:text-blue-600 dark:hover:text-blue-400"
						>
							{player.PlayerName}
						</a>
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
