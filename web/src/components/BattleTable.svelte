<script>
	import { SvelteSet } from 'svelte/reactivity';
	import { regionState } from '$lib/regionState.svelte';
	import { resolve } from '$app/paths';
	import Table from './Table.svelte';
	import TableHeader from './TableHeader.svelte';
	import TableRow from './TableRow.svelte';
	import TableData from './TableData.svelte';

	let { q = '', type = 'alliance', p = '10', offset = 0, hasMore = $bindable(true), selectedIds = $bindable(new SvelteSet()), hasResults = $bindable(false) } = $props();
	let battles = $state([]);
	let loading = $state(true);
	let error = $state(null);
	let prevOffset = 0;
	let prevParams = { q, type, p };

	$effect(() => {
		hasResults = battles.length > 0;
	});

	function formatDate(dateString) {
		const date = new Date(dateString);
		const month = date.getUTCMonth() + 1;
		const day = date.getUTCDate();
		const hours = String(date.getUTCHours()).padStart(2, '0');
		const minutes = String(date.getUTCMinutes()).padStart(2, '0');
		return `${month}/${day} ${hours}:${minutes}`;
	}

	function formatNumber(num) {
		return num.toLocaleString();
	}

	function formatFame(num) {
		if (num === 0) return '-';
		if (num >= 1000000) {
			return (num / 1000000).toFixed(2) + 'm';
		}
		if (num >= 1000) {
			return Math.round(num / 1000) + 'k';
		}
		return num.toLocaleString();
	}

	function mapEntries(list = []) {
		return (list || []).map((entry) => {
			const match = entry?.match(/^(.*?)\s*\((\d+)\)$/);
			return {
				label: match ? match[1].trim() : entry,
				count: match ? match[2] : null
			};
		});
	}

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
			let url;
			if (type === 'alliance' && q) {
				url = new URL(
					`https://albionstats.bricefrisco.com/api/boards/alliance/${regionState.value}/${encodeURIComponent(q)}`
				);
				url.searchParams.set('playerCount', p || '10');
			} else if (type === 'guild' && q) {
				url = new URL(
					`https://albionstats.bricefrisco.com/api/boards/guild/${regionState.value}/${encodeURIComponent(q)}`
				);
				url.searchParams.set('playerCount', p || '10');
			} else if (type === 'player' && q) {
				url = new URL(
					`https://albionstats.bricefrisco.com/api/boards/player/${regionState.value}/${encodeURIComponent(q)}`
				);
				url.searchParams.set('playerCount', p || '10');
			} else {
				url = new URL(`https://albionstats.bricefrisco.com/api/boards/${regionState.value}`);
				url.searchParams.set('totalPlayers', p || '10');
			}

			if (offset > 0) {
				url.searchParams.set('offset', offset.toString());
			}

			const response = await fetch(url.toString());
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = await response.json();
			if (Array.isArray(data)) {
				const newBattles = data.map((battle) => ({
					...battle,
					AllianceEntries: mapEntries(battle.AllianceNames),
					GuildEntries: mapEntries(battle.GuildNames)
				}));

				hasMore = newBattles.length >= 20;

				if (offset > 0 && offset > prevOffset) {
					const existingIds = new Set(battles.map((b) => b.BattleID));
					const uniqueNewBattles = newBattles.filter((b) => !existingIds.has(b.BattleID));
					battles = [...battles, ...uniqueNewBattles];
				} else {
					battles = newBattles;
				}
				prevOffset = offset;
			} else {
				hasMore = false;
				if (offset === 0) {
					battles = [];
				}
			}
		} catch (err) {
			error = err.message;
			console.error('Failed to fetch battle data:', err);
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		// Fetch battles when region or search parameters change
		const currentParams = { q, type, p };
		if (currentParams.q !== prevParams.q || currentParams.type !== prevParams.type || currentParams.p !== prevParams.p) {
			battles = [];
			selectedIds.clear();
			prevParams = currentParams;
		}

		// Ensure we track these dependencies for the effect
		regionState.value;
		offset;

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
				<TableData class="whitespace-nowrap">{formatDate(battle.StartTime)}</TableData>
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


