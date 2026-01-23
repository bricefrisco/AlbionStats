<script>
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import Typography from '$components/Typography.svelte';
	import Tabs from '$components/Tabs.svelte';
	import BBAlliances from '$components/BBAlliances.svelte';
	import BBGuilds from '$components/BBGuilds.svelte';
	import BBPlayers from '$components/BBPlayers.svelte';
	import BBKills from '$components/BBKills.svelte';

	let { data } = $props();

	let region = $derived(data.region);
	let battleIds = $derived(data.battleIds);
	let battleData = $derived(data.battleData);
	let loading = $derived(data.loading);
	let error = $derived(data.error);

	// Tabs state
	const tabs = [
		{ id: 'alliances', label: 'Alliances' },
		{ id: 'guilds', label: 'Guilds' },
		{ id: 'players', label: 'Players' },
		{ id: 'kills', label: 'Kills' }
	];
	let activeTab = $state('alliances');

	// Pagination state
	let alliances = $derived(battleData?.Alliances || []);
	let guilds = $derived(battleData?.Guilds || []);
	let players = $derived(battleData?.Players || []);
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
	{:else if battleData}
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
