<script>
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import Page from '../../../components/Page.svelte';
	import PageHeader from '../../../components/PageHeader.svelte';
	import SubHeader from '../../../components/SubHeader.svelte';
	import Paragraph from '../../../components/Paragraph.svelte';

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

	// Format date
	function formatDate(dateString) {
		if (!dateString) return 'Never';
		return new Date(dateString).toLocaleDateString();
	}
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

			<!-- All Stats in Single Table -->
			<div class="mt-4">
				<table class="w-full text-sm">
					<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
						<!-- PvP Section -->
						<tr class="bg-gray-50/30 dark:bg-gray-800/30">
							<td colspan="2" class="px-4 py-3 font-semibold text-gray-900 dark:text-white"
								>PvP Statistics</td
							>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-8 text-gray-600 dark:text-gray-400">Kill Fame</td>
							<td class="py-2 pr-4 text-right font-medium">{formatNumber(playerData.KillFame)}</td>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-8 text-gray-600 dark:text-gray-400">Death Fame</td>
							<td class="py-2 pr-4 text-right font-medium">{formatNumber(playerData.DeathFame)}</td>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-8 text-gray-600 dark:text-gray-400">Fame Ratio</td>
							<td class="py-2 pr-4 text-right font-medium"
								>{playerData.FameRatio?.toFixed(2) || '0.00'}</td
							>
						</tr>

						<!-- PvE Section -->
						<tr class="bg-gray-50/30 dark:bg-gray-800/30">
							<td colspan="2" class="px-4 py-3 font-semibold text-gray-900 dark:text-white"
								>PvE Statistics</td
							>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-8 text-gray-600 dark:text-gray-400">Total</td>
							<td class="py-2 pr-4 text-right font-medium">{formatNumber(playerData.PveTotal)}</td>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-8 text-gray-600 dark:text-gray-400">Royal</td>
							<td class="py-2 pr-4 text-right font-medium">{formatNumber(playerData.PveRoyal)}</td>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-8 text-gray-600 dark:text-gray-400">Outlands</td>
							<td class="py-2 pr-4 text-right font-medium"
								>{formatNumber(playerData.PveOutlands)}</td
							>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-8 text-gray-600 dark:text-gray-400">Avalon</td>
							<td class="py-2 pr-4 text-right font-medium">{formatNumber(playerData.PveAvalon)}</td>
						</tr>

						<!-- Gathering Section -->
						<tr class="bg-gray-50/30 dark:bg-gray-800/30">
							<td colspan="2" class="px-4 py-3 font-semibold text-gray-900 dark:text-white"
								>Gathering</td
							>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-8 text-gray-600 dark:text-gray-400">Total</td>
							<td class="py-2 pr-4 text-right font-medium"
								>{formatNumber(playerData.GatherAllTotal)}</td
							>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-8 text-gray-600 dark:text-gray-400">Royal</td>
							<td class="py-2 pr-4 text-right font-medium"
								>{formatNumber(playerData.GatherAllRoyal)}</td
							>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-8 text-gray-600 dark:text-gray-400">Outlands</td>
							<td class="py-2 pr-4 text-right font-medium"
								>{formatNumber(playerData.GatherAllOutlands)}</td
							>
						</tr>

						<!-- Crafting Section -->
						<tr class="bg-gray-50/30 dark:bg-gray-800/30">
							<td colspan="2" class="px-4 py-3 font-semibold text-gray-900 dark:text-white"
								>Crafting</td
							>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-8 text-gray-600 dark:text-gray-400">Total</td>
							<td class="py-2 pr-4 text-right font-medium"
								>{formatNumber(playerData.CraftingTotal)}</td
							>
						</tr>

						<!-- Activity Section -->
						<tr class="bg-gray-50/30 dark:bg-gray-800/30">
							<td colspan="2" class="px-4 py-3 font-semibold text-gray-900 dark:text-white"
								>Activity</td
							>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-8 text-gray-600 dark:text-gray-400"
								>Last Killboard Activity</td
							>
							<td class="py-2 pr-4 text-right font-medium"
								>{formatDate(playerData.KillboardLastActivity)}</td
							>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-8 text-gray-600 dark:text-gray-400">Last Other Activity</td>
							<td class="py-2 pr-4 text-right font-medium"
								>{formatDate(playerData.OtherLastActivity)}</td
							>
						</tr>
						<tr>
							<td class="py-2 pr-4 pl-8 text-gray-600 dark:text-gray-400">Last Encountered</td>
							<td class="py-2 pr-4 text-right font-medium"
								>{formatDate(playerData.LastEncountered)}</td
							>
						</tr>
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</Page>
