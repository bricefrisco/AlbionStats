<script>
	import Table from './Table.svelte';
	import TableHeader from './TableHeader.svelte';
	import TableRow from './TableRow.svelte';
	import TableData from './TableData.svelte';
	import Pagination from './Pagination.svelte';
	import { formatNumber, formatFame, formatDateUTC } from '$lib/utils';

	let { data = [] } = $props();

	let currentPage = $state(1);
	const pageSize = 20;

	let totalPages = $derived(Math.ceil(data.length / pageSize));
	let paginatedData = $derived(data.slice((currentPage - 1) * pageSize, currentPage * pageSize));
</script>

<Table>
	{#snippet header()}
		<TableHeader class="w-32 text-left font-semibold whitespace-nowrap">Time</TableHeader>
		<TableHeader class="text-left font-semibold">Killer</TableHeader>
		<TableHeader class="text-left font-semibold">Victim</TableHeader>
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
							IP {formatNumber(kill.KillerIP)}
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
							IP {formatNumber(kill.VictimIP)}
						</span>
						<span class="font-medium text-gray-900 dark:text-white">
							{kill.VictimName || '-'}
						</span>
					</div>
				</div>
			</TableData>
			<TableData class="text-right text-yellow-600 dark:text-yellow-400">
				{formatFame(kill.Fame)}
			</TableData>
		</TableRow>
	{/each}
</Table>

<Pagination bind:currentPage {totalPages} />
