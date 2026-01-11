<script>
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import Page from '../../../components/Page.svelte';
	import PageHeader from '../../../components/PageHeader.svelte';
	import SubHeader from '../../../components/SubHeader.svelte';
	import Paragraph from '../../../components/Paragraph.svelte';
	import StatSection from '../../../components/StatSection.svelte';
	import StatRow from '../../../components/StatRow.svelte';

	// Get parameters from URL
	$: region = $page.params.region;
	$: playerName = $page.params.name;
	$: decodedName = playerName ? decodeURIComponent(playerName) : '';

	// Validate region
	$: validRegion = ['americas', 'europe', 'asia'].includes(region);

	// Player data
	let playerData = null;
	let loading = true;
	let error = null;

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

	onMount(() => {
		if (validRegion && decodedName) {
			fetchPlayerData();
		} else {
			loading = false;
		}
	});

	// Format numbers with commas
	function formatNumber(num) {
		return num ? num.toLocaleString() : '0';
	}

	// Format date and time
	function formatDate(dateString) {
		if (!dateString) return 'Never';
		return new Date(dateString).toLocaleString();
	}

	// Get the most recent activity date
	$: lastActivity = playerData
		? (() => {
				const killboardDate = playerData.KillboardLastActivity
					? new Date(playerData.KillboardLastActivity)
					: null;
				const otherDate = playerData.OtherLastActivity
					? new Date(playerData.OtherLastActivity)
					: null;

				if (!killboardDate && !otherDate) return null;
				if (!killboardDate) return otherDate;
				if (!otherDate) return killboardDate;

				return killboardDate > otherDate ? killboardDate : otherDate;
			})()
		: null;
</script>

<Page>
	<div>
		{#if !validRegion}
			<div class="text-center">
				<PageHeader title="Invalid Region" />
				<Paragraph>Valid regions are: americas, europe, asia</Paragraph>
			</div>
		{:else if !decodedName}
			<div class="text-center">
				<PageHeader title="Player Not Found" />
				<Paragraph>Please provide a valid player name</Paragraph>
			</div>
		{:else if loading}
			<PageHeader title="Loading..." />
		{:else if error}
			<div class="text-center">
				<PageHeader title="Error" />
				<Paragraph classes="text-red-600 dark:text-red-400">{error}</Paragraph>
			</div>
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

			<div class="mt-4">
				<table class="w-full text-sm">
					<tbody class="divide-y divide-gray-200/50 dark:divide-gray-700/50">
						<tr class="bg-gray-50/20 dark:bg-gray-800/20">
							<td class="w-1/4 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white"
								>PvP</td
							>
							<td class="w-1/4 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"
							></td>
							<td class="w-1/4 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white"
								>PvE</td
							>
							<td class="w-1/4 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"
							></td>
						</tr>
						<StatRow
							label="Kill Fame"
							value={formatNumber(playerData.KillFame)}
							label2="Total"
							value2={formatNumber(playerData.PveTotal)}
						/>
						<StatRow
							label="Death Fame"
							value={formatNumber(playerData.DeathFame)}
							label2="Royal"
							value2={formatNumber(playerData.PveRoyal)}
						/>
						<StatRow
							label="Fame Ratio"
							value={playerData.FameRatio?.toFixed(2) || '0.00'}
							label2="Outlands"
							value2={formatNumber(playerData.PveOutlands)}
						/>
						<StatRow
							label="-"
							value="-"
							label2="Avalon"
							value2={formatNumber(playerData.PveAvalon)}
						/>

						<tr class="bg-gray-50/20 dark:bg-gray-800/20">
							<td class="w-1/4 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white"
								>Gathering</td
							>
							<td class="w-1/4 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"
							></td>
							<td class="w-1/4 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white"
								>Crafting</td
							>
							<td class="w-1/4 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"
							></td>
						</tr>
						<StatRow
							label="Total"
							value={formatNumber(playerData.GatherAllTotal)}
							label2="Total"
							value2={formatNumber(playerData.CraftingTotal)}
						/>
						<StatRow
							label="Royal"
							value={formatNumber(playerData.GatherAllRoyal)}
							label2="-"
							value2="-"
						/>
						<StatRow
							label="Outlands"
							value={formatNumber(playerData.GatherAllOutlands)}
							label2="-"
							value2="-"
						/>

						<StatSection title="Activity" />
						<tr>
							<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Last polled</td>
							<td class="py-2 pr-4 text-right font-medium" colspan="3"
								>{formatDate(playerData.TS)}</td
							>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Last activity</td>
							<td class="py-2 pr-4 text-right font-medium" colspan="3"
								>{formatDate(lastActivity)}</td
							>
						</tr>
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</Page>
