<script>
	import { formatNumber } from '$lib/utils';
	let { playerData = null } = $props();

	function formatDate(dateString) {
		if (!dateString) return 'Never';
		return new Date(dateString).toLocaleString();
	}

	let lastActivity = $derived(playerData
		? (() => {
				const killboardDate = playerData.KillboardLastActivity
					? new Date(playerData.KillboardLastActivity)
					: null;
				const otherDate = playerData.OtherLastActivity
					? new Date(playerData.OtherLastActivity)
					: null;

				if (!killboardDate && !otherDate) return null;
				if (!killboardDate) return otherDate;
				if (!otherDate) return killboardDate;

				return killboardDate > otherDate ? killboardDate : otherDate;
			})()
		: null);

	let pvpStats = $derived([
		{
			label: 'Kill Fame',
			value: formatNumber(playerData.KillFame),
			label2: 'Total',
			value2: formatNumber(playerData.PveTotal)
		},
		{
			label: 'Death Fame',
			value: formatNumber(playerData.DeathFame),
			label2: 'Royal',
			value2: formatNumber(playerData.PveRoyal)
		},
		{
			label: 'Fame Ratio',
			value: playerData.FameRatio?.toFixed(2) || '0.00',
			label2: 'Outlands',
			value2: formatNumber(playerData.PveOutlands)
		},
		{ label: '-', value: '-', label2: 'Avalon', value2: formatNumber(playerData.PveAvalon) }
	]);

	let gatheringStats = $derived([
		{
			label: 'Total',
			value: formatNumber(playerData.GatherAllTotal),
			label2: 'Total',
			value2: formatNumber(playerData.CraftingTotal)
		},
		{ label: 'Royal', value: formatNumber(playerData.GatherAllRoyal), label2: '-', value2: '-' },
		{
			label: 'Outlands',
			value: formatNumber(playerData.GatherAllOutlands),
			label2: '-',
			value2: '-'
		}
	]);
</script>

<div class="mt-4">
	<table class="w-full text-sm">
		<tbody class="divide-y divide-gray-200/50 dark:divide-gray-700/50">
			<!-- PvP/PvE Section -->
			<tr class="bg-gray-50/20 dark:bg-gray-800/20">
				<td class="w-1/4 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white">PvP</td>
				<td class="w-1/4 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"></td>
				<td class="w-1/4 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white">PvE</td>
				<td class="w-1/4 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"></td>
			</tr>
			{#each pvpStats as stat (stat.label)}
				<tr>
					<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">{stat.label}</td>
					<td class="py-2 pr-4 text-right font-medium">{stat.value}</td>
					<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">{stat.label2}</td>
					<td class="py-2 pr-4 text-right font-medium">{stat.value2}</td>
				</tr>
			{/each}

			<!-- Gathering/Crafting Section -->
			<tr class="bg-gray-50/20 dark:bg-gray-800/20">
				<td class="w-1/4 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white"
					>Gathering</td
				>
				<td class="w-1/4 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"></td>
				<td class="w-1/4 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white"
					>Crafting</td
				>
				<td class="w-1/4 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"></td>
			</tr>
			{#each gatheringStats as stat (stat.label)}
				<tr>
					<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">{stat.label}</td>
					<td class="py-2 pr-4 text-right font-medium">{stat.value}</td>
					<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">{stat.label2}</td>
					<td class="py-2 pr-4 text-right font-medium">{stat.value2}</td>
				</tr>
			{/each}

			<!-- Activity Section -->
			<tr class="bg-gray-50/20 dark:bg-gray-800/20">
				<td colspan="4" class="px-4 py-3 font-semibold text-gray-900 dark:text-white">Activity</td>
			</tr>
			<tr>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Last polled</td>
				<td class="py-2 pr-4 text-right font-medium" colspan="3">{formatDate(playerData.TS)}</td>
			</tr>
			<tr>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Last activity</td>
				<td class="py-2 pr-4 text-right font-medium" colspan="3">{formatDate(lastActivity)}</td>
			</tr>
		</tbody>
	</table>
</div>
