<script>
	import { SvelteMap } from 'svelte/reactivity';
	import Table from './Table.svelte';
	import TableHeader from './TableHeader.svelte';
	import TableRow from './TableRow.svelte';
	import TableData from './TableData.svelte';
	import Pagination from './Pagination.svelte';
	import Filter from './Filter.svelte';
	import { formatNumber, formatFame, formatDateUTC } from '$lib/utils';
	import { resolve } from '$app/paths';
	import { regionState } from '$lib/regionState.svelte';

	let { kills = [], players = [] } = $props();

	let currentPage = $state(1);
	const pageSize = 10;
	let search = $state('');

	const filteredData = $derived.by(() => {
		const query = search.trim().toLowerCase();
		if (!query) return kills;
		return kills.filter((kill) => {
			const killer = (kill.KillerName || '').toLowerCase();
			const victim = (kill.VictimName || '').toLowerCase();
			return killer.includes(query) || victim.includes(query);
		});
	});
	let totalPages = $derived(Math.ceil(filteredData.length / pageSize));
	let paginatedData = $derived(
		filteredData.slice((currentPage - 1) * pageSize, currentPage * pageSize)
	);

	const playersByName = $derived.by(() => {
		const map = new SvelteMap();
		for (const player of players || []) {
			if (player?.PlayerName) {
				map.set(player.PlayerName, player);
			}
		}
		return map;
	});

	function getAffiliation(name) {
		const player = playersByName.get(name);
		if (!player) return null;
		return {
			alliance: player.AllianceName,
			guild: player.GuildName
		};
	}

	$effect(() => {
		search;
		currentPage = 1;
	});
</script>

<Filter bind:value={search} placeholder="Filter killer or victim" />

<Table>
	{#snippet header()}
		<TableHeader class="w-32 text-left font-semibold whitespace-nowrap">Time</TableHeader>
		<TableHeader class="w-1/3 text-left font-semibold">Killer</TableHeader>
		<TableHeader class="w-1/3 text-left font-semibold">Victim</TableHeader>
		<TableHeader class="w-24 text-right font-semibold whitespace-nowrap">Killer IP</TableHeader>
		<TableHeader class="w-24 text-right font-semibold whitespace-nowrap">Victim IP</TableHeader>
		<TableHeader class="text-right font-semibold whitespace-nowrap">Fame</TableHeader>
	{/snippet}

	{#each paginatedData as kill, index (`${kill.BattleID}-${kill.TS}-${kill.KillerName}-${kill.VictimName}-${index}`)}
		{@const killerAffiliation = getAffiliation(kill.KillerName)}
		{@const victimAffiliation = getAffiliation(kill.VictimName)}
		<TableRow>
			<TableData class="whitespace-nowrap">{formatDateUTC(kill.TS)}</TableData>
			<TableData>
				<div class="flex items-center gap-3">
					{#if kill.KillerWeapon}
						<img
							src="https://render.albiononline.com/v1/item/{kill.KillerWeapon}"
							alt={kill.KillerWeapon}
							class="h-10 w-10 min-w-10"
						/>
					{/if}
					<div class="flex flex-col">
						<span class="text-sm text-gray-600 dark:text-gray-400">
							{#if killerAffiliation?.alliance && killerAffiliation?.guild}
								[{killerAffiliation.alliance}] {killerAffiliation.guild}
							{:else if killerAffiliation?.guild}
								{killerAffiliation.guild}
							{/if}
						</span>
						<span class="font-medium text-gray-900 dark:text-white">
							{#if kill.KillerName}
								<a
									href={resolve(`/players/${regionState.value}/${encodeURIComponent(kill.KillerName)}`)}
									class="underline hover:text-blue-600 dark:hover:text-blue-400"
								>
									{kill.KillerName}
								</a>
							{:else}
								-
							{/if}
						</span>
					</div>
				</div>
			</TableData>
			<TableData>
				<div class="flex items-center gap-3">
					{#if kill.VictimWeapon}
						<img
							src="https://render.albiononline.com/v1/item/{kill.VictimWeapon}"
							alt={kill.VictimWeapon}
							class="h-10 w-10 min-w-10"
						/>
					{/if}
					<div class="flex flex-col">
						<span class="text-sm text-gray-600 dark:text-gray-400">
							{#if victimAffiliation?.alliance && victimAffiliation?.guild}
								[{victimAffiliation.alliance}] {victimAffiliation.guild}
							{:else if victimAffiliation?.guild}
								{victimAffiliation.guild}
							{/if}
						</span>
						<span class="font-medium text-gray-900 dark:text-white">
							{#if kill.VictimName}
								<a
									href={resolve(`/players/${regionState.value}/${encodeURIComponent(kill.VictimName)}`)}
									class="underline hover:text-blue-600 dark:hover:text-blue-400"
								>
									{kill.VictimName}
								</a>
							{:else}
								-
							{/if}
						</span>
					</div>
				</div>
			</TableData>
			<TableData class="text-right">{formatNumber(kill.KillerIP)}</TableData>
			<TableData class="text-right">{formatNumber(kill.VictimIP)}</TableData>
			<TableData class="text-right text-yellow-600 dark:text-yellow-400">
				{formatFame(kill.Fame)}
			</TableData>
		</TableRow>
	{/each}
</Table>

<Pagination bind:currentPage {totalPages} />
