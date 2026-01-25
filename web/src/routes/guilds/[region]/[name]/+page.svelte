<script>
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import SubHeader from '$components/SubHeader.svelte';
	import Typography from '$components/Typography.svelte';
	import Tabs from '$components/Tabs.svelte';
	import GuildSearchBar from '$components/GuildSearchBar.svelte';
	import { formatNumber, formatRatio } from '$lib/utils';

	let { data } = $props();

	let validRegion = $derived(data.validRegion);
	let guildData = $derived(data.guildData);
	let error = $derived(data.guildError);
	let searchName = $derived(guildData.Name);

	function formatDate(dateString) {
		if (!dateString) return 'Never';
		return new Date(dateString).toLocaleString();
	}

	const statRows = $derived.by(() => {
		const roster = guildData?.RosterStats;
		const battle = guildData?.BattleSummary;

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

	const tabs = [{ id: 'players', label: 'Players' }];
	let activeTab = $state('players');
</script>

<Page>
	<div class="mb-4">
		<GuildSearchBar
			bind:value={searchName}
			placeholder="Guild name"
			links={true}
		/>
	</div>

	{#if !validRegion}
		<PageHeader title="Invalid Region" />
		<Typography>Valid regions are: americas, europe, asia</Typography>
	{:else if !guildData.Name}
		<PageHeader title="Guild Not Found" />
		<Typography>Please provide a valid guild name</Typography>
	{:else if error}
		<PageHeader title="Error" />
		<Typography>{error}</Typography>
	{:else}
		<PageHeader title={guildData.Name} />
		<div class="mt-4">
			<table class="w-full text-sm">
				<tbody class="divide-y divide-gray-200/50 dark:divide-gray-700/50">
					<tr class="bg-gray-50/20 dark:bg-gray-800/20">
						<td class="w-1/4 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white">
							PvP
						</td>
						<td class="w-1/4 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"></td>
						<td class="w-1/4 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white">
							Activity
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

		<SubHeader title="Guild Details" classes="mt-6" />

		<div class="mt-2 mb-6">
			<Tabs {tabs} {activeTab} />
		</div>

		<div class="rounded-lg border border-gray-200 bg-gray-50 p-4 text-sm text-gray-600 dark:border-neutral-800 dark:bg-neutral-900 dark:text-gray-300">
			Players tab content coming soon.
		</div>
	{/if}
</Page>
