<script>
	export let playerData = null;

	function formatNumber(num) {
		return num ? num.toLocaleString() : '0';
	}

	function formatDate(dateString) {
		if (!dateString) return 'Never';
		return new Date(dateString).toLocaleString();
	}

	$: lastActivity = playerData
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
		: null;
</script>

<div class="mt-4">
	<table class="w-full text-sm">
		<tbody class="divide-y divide-gray-200/50 dark:divide-gray-700/50">
			<tr class="bg-gray-50/20 dark:bg-gray-800/20">
				<td class="w-1/4 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white">PvP</td>
				<td class="w-1/4 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"></td>
				<td class="w-1/4 px-4 py-3 text-left font-semibold text-gray-900 dark:text-white">PvE</td>
				<td class="w-1/4 px-4 py-3 text-right font-semibold text-gray-900 dark:text-white"></td>
			</tr>
			<tr>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Kill Fame</td>
				<td class="py-2 pr-4 text-right font-medium">{formatNumber(playerData.KillFame)}</td>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Total</td>
				<td class="py-2 pr-4 text-right font-medium">{formatNumber(playerData.PveTotal)}</td>
			</tr>
			<tr>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Death Fame</td>
				<td class="py-2 pr-4 text-right font-medium">{formatNumber(playerData.DeathFame)}</td>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Royal</td>
				<td class="py-2 pr-4 text-right font-medium">{formatNumber(playerData.PveRoyal)}</td>
			</tr>
			<tr>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Fame Ratio</td>
				<td class="py-2 pr-4 text-right font-medium"
					>{playerData.FameRatio?.toFixed(2) || '0.00'}</td
				>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Outlands</td>
				<td class="py-2 pr-4 text-right font-medium">{formatNumber(playerData.PveOutlands)}</td>
			</tr>
			<tr>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">-</td>
				<td class="py-2 pr-4 text-right font-medium">-</td>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Avalon</td>
				<td class="py-2 pr-4 text-right font-medium">{formatNumber(playerData.PveAvalon)}</td>
			</tr>

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
			<tr>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Total</td>
				<td class="py-2 pr-4 text-right font-medium">{formatNumber(playerData.GatherAllTotal)}</td>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Total</td>
				<td class="py-2 pr-4 text-right font-medium">{formatNumber(playerData.CraftingTotal)}</td>
			</tr>
			<tr>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Royal</td>
				<td class="py-2 pr-4 text-right font-medium">{formatNumber(playerData.GatherAllRoyal)}</td>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">-</td>
				<td class="py-2 pr-4 text-right font-medium">-</td>
			</tr>
			<tr>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">Outlands</td>
				<td class="py-2 pr-4 text-right font-medium"
					>{formatNumber(playerData.GatherAllOutlands)}</td
				>
				<td class="py-2 pr-4 pl-4 text-gray-600 dark:text-gray-400">-</td>
				<td class="py-2 pr-4 text-right font-medium">-</td>
			</tr>

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
