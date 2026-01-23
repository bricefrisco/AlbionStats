<script>
	import { untrack } from 'svelte';
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import Typography from '$components/Typography.svelte';
	import Tabs from '$components/Tabs.svelte';
	import BBAlliances from '$components/BBAlliances.svelte';
	import BBGuilds from '$components/BBGuilds.svelte';
	import BBPlayers from '$components/BBPlayers.svelte';
	import BBKills from '$components/BBKills.svelte';
	import { page } from '$app/state';

	let region = $derived(page.params.region);
	let battleIds = $derived(page.params.battleIds);

	let data = $state(null);
	let loading = $state(true);
	let error = $state(null);

	// Tabs state
	const tabs = [
		{ id: 'alliances', label: 'Alliances' },
		{ id: 'guilds', label: 'Guilds' },
		{ id: 'players', label: 'Players' },
		{ id: 'kills', label: 'Kills' }
	];
	let activeTab = $state('alliances');

	// Pagination state
	let alliances = $derived(data?.Alliances || []);
	let guilds = $derived(data?.Guilds || []);
	let players = $derived(data?.Players || []);

	async function fetchBattles() {
		loading = true;
		error = null;
		try {
			const response = await fetch(
				`https://albionstats.bricefrisco.com/api/battles/${region}/${battleIds}`
			);
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			data = await response.json();
		} catch (err) {
			error = err.message;
			console.error('Failed to fetch battle data:', err);
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (region && battleIds) {
			untrack(() => fetchBattles());
		}
	});
</script>

<Page>
	<PageHeader title="Battle Boards" />
	<Typography>
		<p>Battle board results for the {region} server: {battleIds}</p>
	</Typography>
	{#if loading}
		<p>Loading battle data...</p>
	{:else if error}
		<p class="text-red-600">{error}</p>
	{:else if data}
		<div class="mt-8">
			<Tabs {tabs} bind:activeTab />
		</div>

		<div class="mt-4">
			{#if activeTab === 'alliances'}
				<BBAlliances data={alliances} />
			{:else if activeTab === 'guilds'}
				<BBGuilds data={guilds} />
			{:else if activeTab === 'players'}
				<BBPlayers data={players} />
			{:else if activeTab === 'kills'}
				<BBKills />
			{/if}
		</div>
	{/if}
</Page>
