<script>
	import Table from './Table.svelte';
	import TableHeader from './TableHeader.svelte';
	import TableRow from './TableRow.svelte';
	import TableData from './TableData.svelte';
	import Pagination from './Pagination.svelte';
	import { formatNumber, formatFame, formatDateUTC } from '$lib/utils';

	let { kills = [], players = [] } = $props();

	let currentPage = $state(1);
	const pageSize = 20;

	let totalPages = $derived(Math.ceil(kills.length / pageSize));
	let paginatedData = $derived(kills.slice((currentPage - 1) * pageSize, currentPage * pageSize));

	const playersByName = $derived.by(() => {
		const map = new Map();
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
</script>

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
