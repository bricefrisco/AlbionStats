<script>
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import SubHeader from '$components/SubHeader.svelte';
	import Typography from '$components/Typography.svelte';
	import PlayerStats from '$components/PlayerStats.svelte';
	import PlayerSearchBar from '$components/PlayerSearchBar.svelte';
	import Tabs from '$components/Tabs.svelte';
	import PlayerPvPCharts from '$components/charts/PlayerPvPCharts.svelte';
	import PlayerPvECharts from '$components/charts/PlayerPvECharts.svelte';
	import PlayerGatheringCharts from '$components/charts/PlayerGatheringCharts.svelte';
	import PlayerCraftingCharts from '$components/charts/PlayerCraftingCharts.svelte';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { regionState } from '$lib/regionState.svelte.js';

	let { data } = $props();

	// Get parameters from URL
	let decodedName = $derived(data.decodedName);

	// Validate region
	let validRegion = $derived(data.validRegion);

	// Player data
	let searchName = $derived(decodedName);
	let playerData = $derived(data.playerData);
	let loading = $derived(data.loading);
	let error = $derived(data.playerError);

	// Active tab state
	let activeTab = $state('pvp');

	let playerName = $derived.by(() => playerData?.Name || decodedName || 'Player');
	let breadcrumbJsonLd = $derived.by(() => {
		const origin = page.url.origin;
		const region = page.params.region;
		return JSON.stringify({
			'@context': 'https://schema.org',
			'@type': 'BreadcrumbList',
			itemListElement: [
				{
					'@type': 'ListItem',
					position: 1,
					name: 'Players',
					item: `${origin}/players/${region}`
				},
				{
					'@type': 'ListItem',
					position: 2,
					name: playerName,
					item: `${origin}/players/${region}/${encodeURIComponent(playerName)}`
				}
			]
		});
	});

	// Tab configuration
	const tabs = [
		{ id: 'pvp', label: 'PvP' },
		{ id: 'pve', label: 'PvE' },
		{ id: 'gathering', label: 'Gathering' },
		{ id: 'crafting', label: 'Crafting' }
	];

	function handleTabChange(detail) {
		activeTab = detail.tabId;
	}
</script>

<svelte:head>
	<title>{playerName} - AlbionStats - {regionState.label}</title>
	<meta
		name="description"
		content={`Albion Online stats for ${playerName} in ${regionState.label}. View kills, deaths, fame, and charts.`}
	/>
	<link rel="canonical" href={`${page.url.origin}${page.url.pathname}`} />
	<script type="application/ld+json">{breadcrumbJsonLd}</script>
</svelte:head>

<Page>
	<div class="mb-4">
		<PlayerSearchBar
			bind:value={searchName}
			placeholder="Player name"
			links={true}
		/>
	</div>

	<div>
		{#if !validRegion}
			<PageHeader title="Invalid Region" />
			<Typography>Valid regions are: americas, europe, asia</Typography>
		{:else if !decodedName}
			<PageHeader title="Player Not Found" />
			<Typography>Please provide a valid player name</Typography>
		{:else if loading}
			<PageHeader title="Loading..." />
		{:else if error}
			<PageHeader title="Error" />
			<Typography>{error}</Typography>
		{:else if playerData}
			<PageHeader title={playerData.Name} />
			{#if playerData.GuildName || playerData.AllianceName}
				<Typography classes="mb-2 text-sm text-gray-600 dark:text-gray-400 mt-[-15px] font-medium">
					{#if playerData.AllianceName && playerData.GuildName}
						<a
							href={resolve(`/alliances/${regionState.value}/${encodeURIComponent(playerData.AllianceName)}`)}
							class="underline text-current"
						>
							[{playerData.AllianceName}]
						</a>
						{' '}
						<a
							href={resolve(`/guilds/${regionState.value}/${encodeURIComponent(playerData.GuildName)}`)}
							class="underline text-current"
						>
							{playerData.GuildName}
						</a>
					{:else if playerData.AllianceName}
						<a
							href={resolve(`/alliances/${regionState.value}/${encodeURIComponent(playerData.AllianceName)}`)}
							class="underline text-current"
						>
							[{playerData.AllianceName}]
						</a>
					{:else if playerData.GuildName}
						<a
							href={resolve(`/guilds/${regionState.value}/${encodeURIComponent(playerData.GuildName)}`)}
							class="underline text-current"
						>
							{playerData.GuildName}
						</a>
					{/if}
				</Typography>
			{/if}
			<Typography>
				<h2 class="mb-1">Albion Online Player Statistics.</h2>
			</Typography>
			<!-- Guild and Alliance Info -->

			<PlayerStats {playerData} />

			<SubHeader title="Player Charts" classes="mt-4" />

			<!-- Tab Navigation -->
			<div class="mt-2 mb-6">
				<Tabs {tabs} {activeTab} ontabChange={handleTabChange} />
			</div>

			<Typography>
				<p>Data from the past year. Data collection for charts began on January 11, 2026.</p>
			</Typography>

			<!-- Tab Content -->
			<div class="mt-4">
				{#if activeTab === 'pvp'}
					{#if playerData}
						<PlayerPvPCharts data={data.metrics?.pvp} />
					{:else}
						<div class="py-12 text-center text-gray-500 dark:text-gray-400">
							<div class="mb-2 text-lg font-medium">Loading PvP Data...</div>
						</div>
					{/if}
				{:else if activeTab === 'pve'}
					{#if playerData}
						<PlayerPvECharts data={data.metrics?.pve} />
					{:else}
						<div class="py-12 text-center text-gray-500 dark:text-gray-400">
							<div class="mb-2 text-lg font-medium">Loading PvE Data...</div>
						</div>
					{/if}
				{:else if activeTab === 'gathering'}
					{#if playerData}
						<PlayerGatheringCharts data={data.metrics?.gathering} />
					{:else}
						<div class="py-12 text-center text-gray-500 dark:text-gray-400">
							<div class="mb-2 text-lg font-medium">Loading Gathering Data...</div>
						</div>
					{/if}
				{:else if activeTab === 'crafting'}
					{#if playerData}
						<PlayerCraftingCharts data={data.metrics?.crafting} />
					{:else}
						<div class="py-12 text-center text-gray-500 dark:text-gray-400">
							<div class="mb-2 text-lg font-medium">Loading Crafting Data...</div>
						</div>
					{/if}
				{/if}
			</div>
		{/if}
	</div>
</Page>
