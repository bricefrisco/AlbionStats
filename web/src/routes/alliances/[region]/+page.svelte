<script>
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import Typography from '$components/Typography.svelte';
	import AllianceSearchBar from '$components/AllianceSearchBar.svelte';
	import Table from '$components/Table.svelte';
	import TableHeader from '$components/TableHeader.svelte';
	import TableRow from '$components/TableRow.svelte';
	import TableData from '$components/TableData.svelte';
	import { formatFame, formatNumber, formatRatio, getRatioColor } from '$lib/utils.js';

	let { data } = $props();
	let searchQuery = $state('');
</script>

<Page>
	<PageHeader title="Alliances" />
	<Typography>
		<h2>Albion Online Alliances. Search for an alliance to view it's statistics. Below are top 100 alliances
			based on statistics pulled from battle board data based over the past 30 days.</h2>
	</Typography>

	<div class="mb-4">
		<AllianceSearchBar
			bind:value={searchQuery}
			links={true}
			placeholder="Alliance name"
		/>
	</div>

	{#if data.topAlliancesError}
		<Typography>
			<p class="text-red-600 dark:text-red-400">{data.topAlliancesError}</p>
		</Typography>
	{:else if data.topAlliances?.length}
		<Table>
			{#snippet header()}
				<TableHeader class="w-12 text-right font-semibold">#</TableHeader>
				<TableHeader class="text-left font-semibold">Name</TableHeader>
				<TableHeader class="text-right font-semibold whitespace-nowrap">Kill Fame</TableHeader>
				<TableHeader class="text-right font-semibold whitespace-nowrap">Death Fame</TableHeader>
				<TableHeader class="text-right font-semibold whitespace-nowrap">K/D Ratio</TableHeader>
				<TableHeader class="text-right font-semibold">Kills</TableHeader>
				<TableHeader class="text-right font-semibold">Deaths</TableHeader>
			{/snippet}

			{#each data.topAlliances as alliance, index (alliance.AllianceName)}
				{@const ratioColor = getRatioColor(alliance.TotalKillFame, alliance.TotalDeathFame)}
				<TableRow>
					<TableData class="text-right text-gray-500 dark:text-gray-400">
						{index + 1}
					</TableData>
					<TableData class="font-medium text-gray-900 dark:text-white">
						{alliance.AllianceName}
					</TableData>
					<TableData class="text-right text-yellow-600 dark:text-yellow-400">
						{formatFame(alliance.TotalKillFame)}
					</TableData>
					<TableData class="text-right text-gray-600 dark:text-gray-400">
						{formatFame(alliance.TotalDeathFame)}
					</TableData>
					<TableData class="text-right">
						<span style={ratioColor ? `color: ${ratioColor}` : ''}>
							{formatRatio(alliance.TotalKillFame, alliance.TotalDeathFame)}
						</span>
					</TableData>
					<TableData class="text-right text-red-600 dark:text-red-400">
						{formatNumber(alliance.TotalKills)}
					</TableData>
					<TableData class="text-right text-gray-600 dark:text-gray-400">
						{formatNumber(alliance.TotalDeaths)}
					</TableData>
				</TableRow>
			{/each}
		</Table>
	{:else}
		<Typography>
			<p>No top alliances found.</p>
		</Typography>
	{/if}
</Page>
