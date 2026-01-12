<script>
	import { page } from '$app/stores';
	import Page from '../../../components/Page.svelte';
	import PageHeader from '../../../components/PageHeader.svelte';
	import SubHeader from '../../../components/SubHeader.svelte';
	import Paragraph from '../../../components/Paragraph.svelte';
	import PlayerStats from '../../../components/PlayerStats.svelte';
	import Tabs from '../../../components/Tabs.svelte';

	// Get parameters from URL
	$: region = $page.params.region;
	$: playerName = $page.params.name;
	$: decodedName = playerName ? decodeURIComponent(playerName) : '';

	// Validate region
	$: validRegion = ['americas', 'europe', 'asia'].includes(region);

	// Fetch player data when route parameters change
	$: if (validRegion && decodedName) {
		playerData = null;
		loading = true;
		error = null;
		fetchPlayerData();
	}

	// Player data
	let playerData = null;
	let loading = true;
	let error = null;

	// Active tab state
	let activeTab = 'pvp';

	// Tab configuration
	const tabs = [
		{ id: 'pvp', label: 'PvP' },
		{ id: 'pve', label: 'PvE' },
		{ id: 'gathering', label: 'Gathering' },
		{ id: 'crafting', label: 'Crafting' }
	];

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
			<Paragraph>Valid regions are: americas, europe, asia</Paragraph>
		{:else if !decodedName}
			<PageHeader title="Player Not Found" />
			<Paragraph>Please provide a valid player name</Paragraph>
		{:else if loading}
			<PageHeader title="Loading..." />
		{:else if error}
			<PageHeader title="Error" />
			<Paragraph>{error}</Paragraph>
		{:else if playerData}
			<PageHeader title={playerData.Name} />

			<!-- Guild and Alliance Info -->
			{#if playerData.GuildName || playerData.AllianceName}
				<Paragraph classes="mb-2 text-sm text-gray-600 dark:text-gray-400 mt-[-15px] font-medium">
					{#if playerData.AllianceName && playerData.GuildName}
						[{playerData.AllianceName}] {playerData.GuildName}
					{:else if playerData.AllianceName}
						[{playerData.AllianceName}]
					{:else if playerData.GuildName}
						{playerData.GuildName}
					{/if}
				</Paragraph>
			{/if}

			<PlayerStats {playerData} />

			<SubHeader title="Player Charts" classes="mt-4" />

			<!-- Tab Navigation -->
			<div class="mt-2 mb-6">
				<Tabs {tabs} {activeTab} on:tabChange={(e) => (activeTab = e.detail.tabId)} />
			</div>

			<!-- Tab Content -->
			<div class="mt-4">
				{#if activeTab === 'pvp'}
					<div class="py-12 text-center text-gray-500 dark:text-gray-400">
						<div class="mb-2 text-lg font-medium">PvP Statistics</div>
						<div class="text-sm">Kill fame, death fame, and combat statistics coming soon...</div>
					</div>
				{:else if activeTab === 'pve'}
					<div class="py-12 text-center text-gray-500 dark:text-gray-400">
						<div class="mb-2 text-lg font-medium">PvE Statistics</div>
						<div class="text-sm">
							Monster kills, dungeon clears, and PvE achievements coming soon...
						</div>
					</div>
				{:else if activeTab === 'gathering'}
					<div class="py-12 text-center text-gray-500 dark:text-gray-400">
						<div class="mb-2 text-lg font-medium">Gathering Statistics</div>
						<div class="text-sm">
							Resource collection, gathering efficiency, and mining stats coming soon...
						</div>
					</div>
				{:else if activeTab === 'crafting'}
					<div class="py-12 text-center text-gray-500 dark:text-gray-400">
						<div class="mb-2 text-lg font-medium">Crafting Statistics</div>
						<div class="text-sm">
							Item production, crafting levels, and artisan achievements coming soon...
						</div>
					</div>
				{/if}
			</div>
		{/if}
	</div>
</Page>
