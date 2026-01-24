<script>
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import Typography from '$components/Typography.svelte';
	import Tabs from '$components/Tabs.svelte';
	import Card from '$components/Card.svelte';
	import BBAlliances from '$components/BBAlliances.svelte';
	import BBGuilds from '$components/BBGuilds.svelte';
	import BBPlayers from '$components/BBPlayers.svelte';
	import BBKills from '$components/BBKills.svelte';
	import { formatFame, formatNumber } from '$lib/utils';
	import { formatDateUTC } from '$lib/utils';
	import { SvelteMap } from 'svelte/reactivity';

	let { data } = $props();

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
	let kills = $derived(battleData?.Kills || []);

	const statRows = $derived.by(() => [
		{
			leftLabel: 'Battle IDs',
			leftValue: data.battleIds.split(',').join(', '),
			rightLabel: 'Total Alliances',
			rightValue: data.battleData.Alliances.length
		},
		{
			leftLabel: 'Battle Start',
			leftValue: formatDateUTC(data.battleData?.StartTime),
			rightLabel: 'Total Guilds',
			rightValue: data.battleData.Guilds.length
		},
		{
			leftLabel: 'Battle End',
			leftValue: formatDateUTC(data.battleData?.EndTime),
			rightLabel: 'Total Players',
			rightValue: data.battleData.Players.length
		},
		{
			leftLabel: '',
			leftValue: '',
			rightLabel: 'Total Kills',
			rightValue: data.battleData?.TotalKills ?? '-'
		},
		{
			leftLabel: '',
			leftValue: '',
			rightLabel: 'Total Fame',
			rightValue: data.battleData?.TotalFame ? data.battleData.TotalFame.toLocaleString() : '-'
		}
	]);

	const highlightCards = $derived.by(() => {
		const kills = battleData?.Kills || [];
		const playersList = battleData?.Players || [];

		const playersByName = new Map(
			playersList.filter((player) => player?.PlayerName).map((player) => [player.PlayerName, player])
		);

		function formatAffiliation(name) {
			const player = playersByName.get(name);
			if (!player) return '';
			if (player.AllianceName && player.GuildName) {
				return `[${player.AllianceName}] ${player.GuildName}`;
			}
			if (player.GuildName) {
				return player.GuildName;
			}
			return '';
		}

		const killCounts = new SvelteMap();
		for (const kill of kills) {
			if (!kill?.KillerName) continue;
			killCounts.set(kill.KillerName, (killCounts.get(kill.KillerName) || 0) + 1);
		}

		let topKiller = null;
		for (const [name, count] of killCounts.entries()) {
			if (!topKiller || count > topKiller.count) {
				topKiller = { name, count };
			}
		}

		let highestIP = null;
		let mostDeathFame = null;
		let topDamage = null;
		let topHeals = null;

		for (const player of playersList) {
			if (!player?.PlayerName) continue;
			if (!highestIP || player.IP > highestIP.value) {
				highestIP = { name: player.PlayerName, value: player.IP };
			}
			if (!mostDeathFame || player.DeathFame > mostDeathFame.value) {
				mostDeathFame = { name: player.PlayerName, value: player.DeathFame };
			}
			if (!topDamage || player.Damage > topDamage.value) {
				topDamage = { name: player.PlayerName, value: player.Damage };
			}
			if (!topHeals || player.Heal > topHeals.value) {
				topHeals = { name: player.PlayerName, value: player.Heal };
			}
		}

		return [
			{
				title: 'Top Killer',
				value: topKiller?.name || '-',
				meta: topKiller ? formatAffiliation(topKiller.name) : '',
				subtitle: topKiller ? formatNumber(topKiller.count) : '-',
				subtitleClass: 'text-red-600 dark:text-red-400'
			},
			{
				title: 'Highest IP',
				value: highestIP?.name || '-',
				meta: highestIP ? formatAffiliation(highestIP.name) : '',
				subtitle: highestIP ? String(highestIP.value ?? '-') : '-',
				subtitleClass: 'text-white'
			},
			{
				title: 'Most Death Fame',
				value: mostDeathFame?.name || '-',
				meta: mostDeathFame ? formatAffiliation(mostDeathFame.name) : '',
				subtitle: mostDeathFame ? formatFame(mostDeathFame.value) : '-',
				subtitleClass: 'text-gray-600 dark:text-gray-400'
			},
			{
				title: 'Top Damage',
				value: topDamage?.name || '-',
				meta: topDamage ? formatAffiliation(topDamage.name) : '',
				subtitle: topDamage ? formatFame(topDamage.value) : '-',
				subtitleClass: 'text-purple-600 dark:text-purple-400'
			},
			{
				title: 'Top Heals',
				value: topHeals?.name || '-',
				meta: topHeals ? formatAffiliation(topHeals.name) : '',
				subtitle: topHeals ? formatFame(topHeals.value) : '-',
				subtitleClass: 'text-green-600 dark:text-green-400'
			}
		];
	});
</script>

<Page>
	<PageHeader title="Battle Boards" />
	<Typography>
		<h2>Albion Online battle board results</h2>
	</Typography>
	{#if loading}
		<p>Loading battle data...</p>
	{:else if error}
		<p class="text-red-600">{error}</p>
	{:else if battleData}
		<div class="mt-6">
			<table class="w-full text-sm">
				<tbody class="divide-y divide-gray-200/50 dark:divide-gray-700/50">
					<tr class="bg-gray-50/20 dark:bg-gray-800/20">
						<td class="w-1/6 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white">
							Battle
						</td>
						<td class="w-1/3 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"></td>
						<td class="w-1/6 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white">
							Totals
						</td>
						<td class="w-1/3 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"></td>
					</tr>
					{#each statRows as stat (stat.rightLabel)}
						<tr>
							<td class="w-1/6 py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">
								{stat.leftLabel}
							</td>
							<td class="w-1/3 py-2 pr-4 text-right font-medium">{stat.leftValue}</td>
							<td class="w-1/6 py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">
								{stat.rightLabel}
							</td>
							<td class="w-1/3 py-2 pr-4 text-right font-medium">{stat.rightValue}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<div class="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
			{#each highlightCards as card (card.title)}
				<Card
					title={card.title}
					value={card.value}
					subtitle={card.subtitle}
					meta={card.meta}
					subtitleClass={card.subtitleClass}
				/>
			{/each}
		</div>

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
				<BBKills kills={kills} players={players} />
			{/if}
		</div>
	{/if}
</Page>
