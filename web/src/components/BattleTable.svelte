<script>
	import { SvelteSet } from 'svelte/reactivity';
	import { regionState } from '$lib/regionState.svelte';
	import { resolve } from '$app/paths';
	import { formatNumber, formatFame, formatDateUTC } from '$lib/utils';
	import { buildBattleBoardsUrl, mapBattleBoardsData } from '$lib/battleBoards';
	import Table from './Table.svelte';
	import TableHeader from './TableHeader.svelte';
	import TableRow from './TableRow.svelte';
	import TableData from './TableData.svelte';

	let {
		q = '',
		type = 'alliance',
		p = '10',
		offset = 0,
		hasMore = $bindable(true),
		loading = $bindable(true),
		selectedIds = $bindable(new SvelteSet()),
		hasResults = $bindable(false),
		initialBattles = [],
		initialHasMore = true,
		initialError = null
	} = $props();
	let extraBattles = $state([]);
	let error = $state(null);
	let prevOffset = 0;
	let battles = $derived([...initialBattles, ...extraBattles]);

	$effect(() => {
		hasResults = battles.length > 0;
	});

	$effect(() => {
		hasMore = initialHasMore;
		error = initialError;
		loading = false;
		prevOffset = 0;
		extraBattles = [];
		selectedIds.clear();
	});

	function toggleSelection(id) {
		if (selectedIds.has(id)) {
			selectedIds.delete(id);
		} else {
			selectedIds.add(id);
		}
	}

	async function fetchBattles() {
		loading = true;
		error = null;
		try {
			const url = buildBattleBoardsUrl({
				base: 'https://albionstats.bricefrisco.com/api',
				region: regionState.value,
				type,
				q,
				p,
				offset
			});

			const response = await fetch(url.toString());
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = await response.json();

			const newBattles = mapBattleBoardsData(data);
			hasMore = newBattles.length >= 20;

			if (offset > 0 && offset > prevOffset) {
				const existingIds = new Set(
					[...initialBattles, ...extraBattles].map((b) => b.BattleID)
				);
				const uniqueNewBattles = newBattles.filter((b) => !existingIds.has(b.BattleID));
				extraBattles = [...extraBattles, ...uniqueNewBattles];
			}
			prevOffset = offset;
		} catch (err) {
			error = err.message;
			console.error('Failed to fetch battle data:', err);
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		// Only fetch on pagination. Initial data comes from SSR.
		if (offset <= 0) return;
		if (offset <= prevOffset) return;
		fetchBattles();
	});
</script>

{#if loading && battles.length === 0}
	<p class="text-sm text-gray-600 dark:text-gray-300">Loading battles...</p>
{:else if error && battles.length === 0}
	<p class="text-sm text-red-600 dark:text-red-400">{error}</p>
{:else if battles.length === 0 && !loading}
	<p class="text-sm text-gray-600 dark:text-gray-300">No battles found.</p>
{:else}
	<Table>
		{#snippet header()}
			<TableHeader class="w-12 text-left" />
			<TableHeader class="w-1/6 text-left font-semibold">Battle ID</TableHeader>
			<TableHeader class="w-1/6 text-left font-semibold">Start Time</TableHeader>
			<TableHeader class="w-1/12 text-right font-semibold">Players</TableHeader>
			<TableHeader class="w-1/12 text-right font-semibold">Kills</TableHeader>
			<TableHeader class="w-1/4 text-left font-semibold">Alliances</TableHeader>
			<TableHeader class="w-1/4 text-left font-semibold">Guilds</TableHeader>
			<TableHeader class="w-1/12 text-right font-semibold">Fame</TableHeader>
		{/snippet}

		{#each battles as battle (battle.BattleID)}
			<TableRow>
				<TableData>
					<input
						type="checkbox"
						checked={selectedIds.has(battle.BattleID)}
						onchange={() => toggleSelection(battle.BattleID)}
						class="mt-[1.75px] h-[18px] w-[18px] rounded border-gray-300 bg-gray-100 text-blue-600 focus:ring-blue-500 dark:border-neutral-800 dark:bg-neutral-900"
					/>
				</TableData>
				<TableData class="font-medium text-gray-900 dark:text-white break-words">
					<a
						href={resolve(`/battle-boards/${regionState.value}/${battle.BattleID}`)}
						class="underline hover:text-blue-600 dark:hover:text-blue-400"
					>
						{battle.BattleID}
					</a>
				</TableData>
				<TableData class="whitespace-nowrap">{formatDateUTC(battle.StartTime)}</TableData>
				<TableData class="text-right font-medium text-blue-600 dark:text-blue-400">
					{formatNumber(battle.TotalPlayers)}
				</TableData>
				<TableData class="text-right font-medium text-red-600 dark:text-red-400">
					{formatNumber(battle.TotalKills)}
				</TableData>
				<TableData>
					<div class="grid gap-1 text-xs text-gray-700 dark:text-gray-300">
						{#if battle.AllianceEntries?.length}
							{#each battle.AllianceEntries.slice(0, 3) as entry (entry.label)}
								<div class="grid grid-cols-[1fr_auto] items-center gap-2">
									<span class="truncate">{entry.label}</span>
									{#if entry.count}
										<span class="text-right font-semibold">{entry.count}</span>
									{/if}
								</div>
							{/each}
							{#if battle.AllianceEntries.length > 3}
								<span class="text-[11px] text-gray-500 dark:text-gray-400">
									+{battle.AllianceEntries.length - 3} more
								</span>
							{/if}
						{:else}
							<span class="text-[11px] text-gray-500 dark:text-gray-400">-</span>
						{/if}
					</div>
				</TableData>
				<TableData>
					<div class="grid gap-1 text-xs text-gray-700 dark:text-gray-300">
						{#if battle.GuildEntries?.length}
							{#each battle.GuildEntries.slice(0, 3) as entry (entry.label)}
								<div class="grid grid-cols-[1fr_auto] items-center gap-2">
									<span class="truncate">{entry.label}</span>
									{#if entry.count}
										<span class="text-right font-semibold">{entry.count}</span>
									{/if}
								</div>
							{/each}
							{#if battle.GuildEntries.length > 3}
								<span class="text-[11px] text-gray-500 dark:text-gray-400">
									+{battle.GuildEntries.length - 3} more
								</span>
							{/if}
						{:else}
							<span class="text-[11px] text-gray-500 dark:text-gray-400">-</span>
						{/if}
					</div>
				</TableData>
				<TableData class="text-right font-medium text-yellow-600 dark:text-yellow-400">
					{formatFame(battle.TotalFame)}
				</TableData>
			</TableRow>
		{/each}
	</Table>

	{#if loading}
		<p class="mt-4 text-sm text-gray-600 dark:text-gray-300 italic animate-pulse">Loading more battles...</p>
	{/if}

	{#if error}
		<p class="mt-4 text-sm text-red-600 dark:text-red-400">{error}</p>
	{/if}
{/if}


