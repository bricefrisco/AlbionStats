<script>
	import { SvelteMap } from 'svelte/reactivity';
	import Table from './Table.svelte';
	import TableHeader from './TableHeader.svelte';
	import TableRow from './TableRow.svelte';
	import TableData from './TableData.svelte';
	import Pagination from './Pagination.svelte';
	import BBFilter from './BBFilter.svelte';
	import { formatNumber, formatFame, formatDateUTC } from '$lib/utils';

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

	function formatAffiliation(name) {
		const player = playersByName.get(name);
		const alliance = player?.AllianceName;
		const guild = player?.GuildName;
		const parts = [];
		if (alliance && guild) {
			parts.push(`[${alliance}] ${guild}`);
		} else if (guild) {
			parts.push(guild);
		}
		return parts.join(' - ');
	}

	$effect(() => {
		search;
		currentPage = 1;
	});
</script>

<BBFilter bind:value={search} placeholder="Filter killer or victim" />

<Table>
	{#snippet header()}
		<TableHeader class="w-32 text-left font-semibold whitespace-nowrap">Time</TableHeader>
		<TableHeader class="text-left font-semibold">Killer</TableHeader>
		<TableHeader class="text-left font-semibold">Victim</TableHeader>
		<TableHeader class="w-24 text-right font-semibold whitespace-nowrap">Killer IP</TableHeader>
		<TableHeader class="w-24 text-right font-semibold whitespace-nowrap">Victim IP</TableHeader>
		<TableHeader class="text-right font-semibold whitespace-nowrap">Fame</TableHeader>
	{/snippet}

	{#each paginatedData as kill, index (`${kill.BattleID}-${kill.TS}-${kill.KillerName}-${kill.VictimName}-${index}`)}
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
							{formatAffiliation(kill.KillerName)}
						</span>
						<span class="font-medium text-gray-900 dark:text-white">
							{kill.KillerName || '-'}
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
							{formatAffiliation(kill.VictimName)}
						</span>
						<span class="font-medium text-gray-900 dark:text-white">
							{kill.VictimName || '-'}
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
