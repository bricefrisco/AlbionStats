<script>
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import SubHeader from '$components/SubHeader.svelte';
	import Typography from '$components/Typography.svelte';
	import PlayerStats from '$components/PlayerStats.svelte';
	import Tabs from '$components/Tabs.svelte';
	import PlayerPvPCharts from '$components/charts/PlayerPvPCharts.svelte';
	import PlayerPvECharts from '$components/charts/PlayerPvECharts.svelte';
	import PlayerGatheringCharts from '$components/charts/PlayerGatheringCharts.svelte';
	import PlayerCraftingCharts from '$components/charts/PlayerCraftingCharts.svelte';

	// Get parameters from URL
	let region = $derived($page.params.region);
	let playerName = $derived($page.params.name);
	let decodedName = $derived(playerName ? decodeURIComponent(playerName) : '');

	// Validate region
	let validRegion = $derived(['americas', 'europe', 'asia'].includes(region));

	// Fetch player data when route parameters change
	$effect(() => {
		if (validRegion && decodedName) {
			playerData = null;
			loading = true;
			error = null;
			fetchPlayerData();
		}
	});

	// Player data
	let playerData = $state(null);
	let loading = $state(true);
	let error = $state(null);

	// Active tab state
	let activeTab = $derived($page.url.searchParams.get('tab') || 'pvp');

	// Tab configuration
	const tabs = [
		{ id: 'pvp', label: 'PvP' },
		{ id: 'pve', label: 'PvE' },
		{ id: 'gathering', label: 'Gathering' },
		{ id: 'crafting', label: 'Crafting' }
	];

	function handleTabChange(detail) {
		const newTabId = detail.tabId;
		const url = new URL($page.url);
		url.searchParams.set('tab', newTabId);
		goto(resolve(url.pathname + url.search), { replaceState: true, noScroll: true });
	}

	async function fetchPlayerData() {
		try {
			const response = await fetch(
				`https://albionstats.bricefrisco.com/api/players/${region}/${encodeURIComponent(decodedName)}`
			);

			if (response.status === 404) {
				error = 'Player not found';
				return;
			}

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			playerData = await response.json();
		} catch (err) {
			error = err.message;
			console.error('Failed to fetch player data:', err);
		} finally {
			loading = false;
		}
	}
</script>

<Page>
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

			<!-- Guild and Alliance Info -->
			{#if playerData.GuildName || playerData.AllianceName}
				<Typography classes="mb-2 text-sm text-gray-600 dark:text-gray-400 mt-[-15px] font-medium">
					{#if playerData.AllianceName && playerData.GuildName}
						[{playerData.AllianceName}] {playerData.GuildName}
					{:else if playerData.AllianceName}
						[{playerData.AllianceName}]
					{:else if playerData.GuildName}
						{playerData.GuildName}
					{/if}
				</Typography>
			{/if}

			<PlayerStats {playerData} />

			<SubHeader title="Player Charts" classes="mt-4" />

			<!-- Tab Navigation -->
			<div class="mt-2 mb-6">
				<Tabs {tabs} {activeTab} ontabChange={handleTabChange} />
			</div>

			<!-- Tab Content -->
			<div class="mt-4">
				{#if activeTab === 'pvp'}
					{#if playerData}
						<PlayerPvPCharts {region} playerId={playerData.PlayerID} />
					{:else}
						<div class="py-12 text-center text-gray-500 dark:text-gray-400">
							<div class="mb-2 text-lg font-medium">Loading PvP Data...</div>
						</div>
					{/if}
				{:else if activeTab === 'pve'}
					{#if playerData}
						<PlayerPvECharts {region} playerId={playerData.PlayerID} />
					{:else}
						<div class="py-12 text-center text-gray-500 dark:text-gray-400">
							<div class="mb-2 text-lg font-medium">Loading PvE Data...</div>
						</div>
					{/if}
				{:else if activeTab === 'gathering'}
					{#if playerData}
						<PlayerGatheringCharts {region} playerId={playerData.PlayerID} />
					{:else}
						<div class="py-12 text-center text-gray-500 dark:text-gray-400">
							<div class="mb-2 text-lg font-medium">Loading Gathering Data...</div>
						</div>
					{/if}
				{:else if activeTab === 'crafting'}
					{#if playerData}
						<PlayerCraftingCharts {region} playerId={playerData.PlayerID} />
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
