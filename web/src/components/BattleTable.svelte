<script>
	import { onMount } from 'svelte';

    let battles = [];
	let loading = true;
	let error = null;

	function formatDate(dateString) {
		const date = new Date(dateString);
		const month = date.getMonth() + 1;
		const day = date.getDate();
		const hours = date.getHours() % 12 || 12;
		const minutes = String(date.getMinutes()).padStart(2, '0');
		return `${month}/${day} ${hours}:${minutes}`;
	}

function formatNumber(num) {
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
		try {
			const response = await fetch('https://albionstats.bricefrisco.com/api/boards/americas');
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

	// Initialize data fetching
	onMount(() => {
		fetchBattles();
	});
</script>

{#if loading}
	<p class="text-sm text-gray-600 dark:text-gray-300">Loading battles...</p>
{:else if error}
	<p class="text-sm text-red-600 dark:text-red-400">{error}</p>
{:else if battles.length === 0}
	<p class="text-sm text-gray-600 dark:text-gray-300">No battles found.</p>
{:else}
		<table class="w-full table-fixed text-sm break-words">
			<thead>
				<tr class="bg-gray-50/20 dark:bg-gray-800/20">
					<th class="px-4 py-4 text-left font-semibold text-gray-900 dark:text-white">
						Battle ID
					</th>
					<th class="px-4 py-4 text-left font-semibold text-gray-900 dark:text-white">
						Start Time
					</th>
					<th class="px-4 py-4 text-left font-semibold text-gray-900 dark:text-white">
						Total Players
					</th>
					<th class="px-4 py-4 text-left font-semibold text-gray-900 dark:text-white">
						Total Kills
					</th>
					<th class="px-4 py-4 text-left font-semibold text-gray-900 dark:text-white">
						Total Fame
					</th>
					<th class="px-4 py-4 text-left font-semibold text-gray-900 dark:text-white">
						Alliances
					</th>
					<th class="px-4 py-4 text-left font-semibold text-gray-900 dark:text-white">
						Guilds
					</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-gray-200/50 dark:divide-gray-700/50">
				{#each battles as battle (battle.BattleID)}
					<tr class="align-top">
						<td class="px-4 py-4 text-gray-700 dark:text-gray-300 break-words">
							{battle.BattleID}
						</td>
						<td class="px-4 py-4 text-gray-700 dark:text-gray-300 whitespace-nowrap">
							{formatDate(battle.StartTime)}
						</td>
						<td class="px-4 py-4 text-gray-700 dark:text-gray-300">
							{formatNumber(battle.TotalPlayers)}
						</td>
						<td class="px-4 py-4 text-gray-700 dark:text-gray-300">
							{formatNumber(battle.TotalKills)}
						</td>
						<td class="px-4 py-4 text-gray-700 dark:text-gray-300">
							{formatNumber(battle.TotalFame)}
						</td>
						<td class="px-4 py-4 align-top" style="vertical-align: top;">
							<div class="flex flex-col gap-1.5 items-start">
								{#if battle.AllianceEntries?.length}
									{#each battle.AllianceEntries.slice(0, 3) as entry}
										<div class="flex items-center justify-between text-xs text-blue-700 dark:text-blue-200">
											<span class="truncate mr-2">{entry.label}</span>
											{#if entry.count}
												<span class="font-semibold">{entry.count}</span>
											{/if}
										</div>
									{/each}
									{#if battle.AllianceEntries.length > 3}
										<span class="inline-flex items-center rounded-full bg-gray-100 px-2 py-0.5 text-xs text-gray-600 dark:bg-gray-800 dark:text-gray-300">
											... ({battle.AllianceEntries.length - 3} more)
										</span>
									{/if}
								{:else}
									<span
										class="inline-flex items-center rounded-full bg-gray-100 px-2 py-0.5 text-xs text-gray-600 dark:bg-gray-800 dark:text-gray-300"
										>None</span
									>
								{/if}
							</div>
						</td>
						<td class="px-4 py-4 align-top" style="vertical-align: top;">
							<div class="flex flex-col gap-1.5 items-start">
								{#if battle.GuildEntries?.length}
									{#each battle.GuildEntries.slice(0, 3) as entry}
										<div class="flex items-center justify-between text-xs text-blue-700 dark:text-blue-200">
											<span class="truncate mr-2">{entry.label}</span>
											{#if entry.count}
												<span class="font-semibold">{entry.count}</span>
											{/if}
										</div>
									{/each}
									{#if battle.GuildEntries.length > 3}
										<span class="inline-flex items-center rounded-full bg-gray-100 px-2 py-0.5 text-xs text-gray-600 dark:bg-gray-800 dark:text-gray-300">
											... ({battle.GuildEntries.length - 3} more)
										</span>
									{/if}
								{:else}
									<span
										class="inline-flex items-center rounded-full bg-gray-100 px-2 py-0.5 text-xs text-gray-600 dark:bg-gray-800 dark:text-gray-300"
										>None</span
									>
								{/if}
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
{/if}


