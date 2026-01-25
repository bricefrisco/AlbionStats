<script>
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import Typography from '$components/Typography.svelte';
	import GuildSearchBar from '$components/GuildSearchBar.svelte';
	import BackToTopButton from '$components/BackToTopButton.svelte';
	import Table from '$components/Table.svelte';
	import TableHeader from '$components/TableHeader.svelte';
	import TableRow from '$components/TableRow.svelte';
	import TableData from '$components/TableData.svelte';
	import { page } from '$app/state';
	import { formatFame, formatNumber, formatRatio, getRatioColor } from '$lib/utils.js';
	import { regionState } from '$lib/regionState.svelte.js';
	import { resolve } from '$app/paths';

	let { data } = $props();

	let searchQuery = $derived(page.url.searchParams.get('q') || '');
</script>

<svelte:head>
	<title>Guilds - AlbionStats - {regionState.label}</title>
	<meta
		name="description"
		content={`Top Albion Online guilds in ${regionState.label}. Search guild stats, kills, deaths, and fame.`}
	/>
</svelte:head>

<Page>
	<PageHeader title="Guilds" />
	<Typography>
		<h2>Albion Online Guild Statistics. Search for a guild to view their stats.</h2>
		<p>Below are top 100 guilds based on statistics pulled from battle board data based
			over the past 30 days.</p>
		<p>Collection began on January 19th, 2026.</p>
	</Typography>

	<div class="mb-4">
		<GuildSearchBar
			bind:value={searchQuery}
			links={true}
			placeholder="Guild name"
		/>
	</div>

	{#if data.topGuildsError}
		<Typography>
			<p class="text-red-600 dark:text-red-400">{data.topGuildsError}</p>
		</Typography>
	{:else if data.topGuilds?.length}
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

		{#each data.topGuilds as guild, index (guild.GuildName)}
			{@const ratioColor = getRatioColor(guild.TotalKillFame, guild.TotalDeathFame)}
				<TableRow>
					<TableData class="text-right text-gray-500 dark:text-gray-400">
						{index + 1}
					</TableData>
					<TableData class="font-medium text-gray-900 dark:text-white">
						<a
							href={resolve(`/guilds/${regionState.value}/${encodeURIComponent(guild.GuildName)}`)}
							class="hover:underline hover:text-blue-600 dark:hover:text-blue-400"
						>
							{guild.GuildName}
						</a>
					</TableData>
					<TableData class="text-right text-yellow-600 dark:text-yellow-400">
						{formatFame(guild.TotalKillFame)}
					</TableData>
					<TableData class="text-right text-gray-600 dark:text-gray-400">
						{formatFame(guild.TotalDeathFame)}
					</TableData>
					<TableData class="text-right">
						<span style={ratioColor ? `color: ${ratioColor}` : ''}>
							{formatRatio(guild.TotalKillFame, guild.TotalDeathFame)}
						</span>
					</TableData>
					<TableData class="text-right text-red-600 dark:text-red-400">
						{formatNumber(guild.TotalKills)}
					</TableData>
					<TableData class="text-right text-gray-600 dark:text-gray-400">
						{formatNumber(guild.TotalDeaths)}
					</TableData>
				</TableRow>
			{/each}
		</Table>
		<div class="mt-6 flex justify-center">
			<BackToTopButton />
		</div>
	{:else}
		<Typography>
			<p>No top guilds found.</p>
		</Typography>
	{/if}
</Page>
