<script>
	import { untrack } from 'svelte';
	import { regionState } from '$lib/regionState.svelte';

	let { q = '', type = 'alliance', p = '10' } = $props();
	let battles = $state([]);
	let loading = $state(true);
	let error = $state(null);

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
			} else {
				url = new URL(`https://albionstats.bricefrisco.com/api/boards/${regionState.value}`);
				url.searchParams.set('totalPlayers', p || '10');
			}

			const response = await fetch(url.toString());
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = await response.json();
			battles = data.map((battle) => ({
				...battle,
				AllianceEntries: mapEntries(battle.AllianceNames),
				GuildEntries: mapEntries(battle.GuildNames)
			}));
		} catch (err) {
			error = err.message;
			console.error('Failed to fetch battle data:', err);
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		// Fetch battles when region or search parameters change
		regionState.value;
		q;
		type;
		p;
		untrack(() => fetchBattles());
	});
</script>

{#if loading}
	<p class="text-sm text-gray-600 dark:text-gray-300">Loading battles...</p>
{:else if error}
	<p class="text-sm text-red-600 dark:text-red-400">{error}</p>
{:else if battles.length === 0}
	<p class="text-sm text-gray-600 dark:text-gray-300">No battles found.</p>
{:else}
	<div
		class="overflow-hidden rounded-lg border border-gray-200/60 bg-transparent shadow-sm dark:border-gray-800/60 dark:bg-transparent"
	>
		<table class="w-full table-fixed text-sm">
			<thead class="bg-gray-50/60 dark:bg-gray-800/40">
			<tr class="text-sm capitalize tracking-wide text-gray-600 dark:text-gray-300">
				<th class="w-1/6 px-4 py-3 text-left font-semibold">Battle ID</th>
				<th class="w-1/6 px-4 py-3 text-left font-semibold">Start Time</th>
				<th class="w-1/12 px-4 py-3 text-right font-semibold">Players</th>
				<th class="w-1/12 px-4 py-3 text-right font-semibold">Kills</th>
				<th class="w-1/4 px-4 py-3 text-left font-semibold">Alliances</th>
				<th class="w-1/4 px-4 py-3 text-left font-semibold">Guilds</th>
				<th class="w-1/12 px-4 py-3 text-right font-semibold">Fame</th>
			</tr>
			</thead>
			<tbody class="divide-y divide-gray-200/60 dark:divide-gray-700/60">
			{#each battles as battle (battle.BattleID)}
				<tr class="align-top text-gray-700 hover:bg-gray-50/60 dark:text-gray-300 dark:hover:bg-gray-800/30">
					<td class="px-4 py-3 font-medium text-gray-900 dark:text-white break-words">
						{battle.BattleID}
					</td>
					<td class="px-4 py-3 whitespace-nowrap">{formatDate(battle.StartTime)}</td>
					<td class="px-4 py-3 text-right font-medium text-blue-600 dark:text-blue-400">
						{formatNumber(battle.TotalPlayers)}
					</td>
					<td class="px-4 py-3 text-right font-medium text-red-600 dark:text-red-400">
						{formatNumber(battle.TotalKills)}
					</td>
					<td class="px-4 py-3 align-top" style="vertical-align: top;">
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
					</td>
					<td class="px-4 py-3 align-top" style="vertical-align: top;">
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
					</td>
					<td class="px-4 py-3 text-right font-medium text-yellow-600 dark:text-yellow-400">
						{formatFame(battle.TotalFame)}
					</td>
				</tr>
			{/each}
			</tbody>
		</table>
	</div>
{/if}


