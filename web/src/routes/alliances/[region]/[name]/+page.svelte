<script>
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import SubHeader from '$components/SubHeader.svelte';
	import Typography from '$components/Typography.svelte';
	import Tabs from '$components/Tabs.svelte';
	import AllianceSearchBar from '$components/AllianceSearchBar.svelte';
	import AllianceGuilds from '$components/AllianceGuilds.svelte';
	import AllianceGuildPlayers from '$components/AllianceGuildPlayers.svelte';
	import { formatNumber, formatRatio } from '$lib/utils';
	import { page } from '$app/state';
	import { regionState } from '$lib/regionState.svelte.js';

	let { data } = $props();

	let validRegion = $derived(data.validRegion);
	let allianceData = $derived(data.allianceData);
	let error = $derived(data.allianceError);
	let searchName = $derived(allianceData.Name);
	let allianceName = $derived.by(() => allianceData?.Name || 'Alliance');
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
					name: 'Alliances',
					item: `${origin}/alliances/${region}`
				},
				{
					'@type': 'ListItem',
					position: 2,
					name: allianceName,
					item: `${origin}/alliances/${region}/${encodeURIComponent(allianceName)}`
				}
			]
		});
	});

	function formatDate(dateString) {
		if (!dateString) return 'Never';
		return new Date(dateString).toLocaleString();
	}

	const statRows = $derived.by(() => {
		const roster = allianceData?.RosterStats;
		const battle = allianceData?.BattleSummary;

		return [
			{
				label: 'Kill Fame',
				value: battle?.TotalKillFame?.toLocaleString() ?? '-',
				label2: 'Last 7 Days',
				value2: formatNumber(roster?.Active7d)
			},
			{
				label: 'Death Fame',
				value: battle?.TotalDeathFame?.toLocaleString() ?? '-',
				label2: 'Last 30 Days',
				value2: formatNumber(roster?.Active30d)
			},
			{
				label: 'Fame Ratio',
				value: formatRatio(battle?.TotalKillFame, battle?.TotalDeathFame),
				label2: 'Max Participation',
				value2: formatNumber(battle?.MaxPlayers)
			},
			{
				label: 'Kills',
				value: battle?.TotalKills?.toLocaleString() ?? '-',
				label2: 'Last battle at',
				value2: formatDate(battle?.LastBattleAt)
			},
			{
				label: 'Deaths',
				value: battle?.TotalDeaths?.toLocaleString() ?? '-',
				label2: '',
				value2: ''
			}
		];
	});

	const tabs = [
		{ id: 'guilds', label: 'Guilds' },
		{ id: 'players', label: 'Players' }
	];
	let activeTab = $state('guilds');

	const guilds = $derived(allianceData?.Guilds || []);
	const players = $derived(allianceData?.Players || []);
</script>

<svelte:head>
	<title>{allianceName} - AlbionStats - {regionState.label}</title>
	<meta
		name="description"
		content={`Albion Online alliance stats for ${allianceName} in ${regionState.label}. View kills, deaths, fame, and roster activity.`}
	/>
	<link rel="canonical" href={`${page.url.origin}${page.url.pathname}`} />
	<script type="application/ld+json">{breadcrumbJsonLd}</script>
</svelte:head>

<Page>
	<div class="mb-4">
		<AllianceSearchBar
			bind:value={searchName}
			placeholder="Alliance name"
			links={true}
		/>
	</div>

	{#if !validRegion}
		<PageHeader title="Invalid Region" />
		<Typography>Valid regions are: americas, europe, asia</Typography>
	{:else if !allianceData.Name}
		<PageHeader title="Alliance Not Found" />
		<Typography>Please provide a valid alliance name</Typography>
	{:else if error}
		<PageHeader title="Error" />
		<Typography>{error}</Typography>
	{:else}
		<PageHeader title={allianceData.Name} />
		<Typography>
			<h2 class="mb-1">Albion Online Alliance Statistics.</h2>

			<p>Data is based on battle boards from the past 30 days. Collection began on January 19, 2026.
			</p>
		</Typography>
		<div class="mt-4">
			<table class="w-full text-sm">
				<tbody class="divide-y divide-gray-200/50 dark:divide-gray-700/50">
					<tr class="bg-gray-50/20 dark:bg-gray-800/20">
						<td class="w-1/4 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white">
							PvP
						</td>
						<td class="w-1/4 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"></td>
						<td class="w-1/4 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white">
							Active Players
						</td>
						<td class="w-1/4 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"></td>
					</tr>
					{#each statRows as stat (stat.label)}
						<tr>
							<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">{stat.label}</td>
							<td class="py-2 pr-4 text-right font-medium">{stat.value}</td>
							<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">{stat.label2}</td>
							<td class="py-2 pr-4 text-right font-medium">{stat.value2}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<SubHeader title="Alliance Details" classes="mt-6" />

		<div class="mt-2 mb-6">
			<Tabs {tabs} bind:activeTab />
		</div>

		{#if activeTab === 'guilds'}
			<AllianceGuilds data={guilds} />
		{:else if activeTab === 'players'}
			<AllianceGuildPlayers data={players} />
		{/if}
	{/if}
</Page>
