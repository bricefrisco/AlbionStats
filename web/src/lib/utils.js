export function formatNumber(num) {
	return num?.toLocaleString() ?? '0';
}

export function formatFame(num) {
	if (!num) return '-';
	if (num >= 1000000) {
		return (num / 1000000).toFixed(2) + 'm';
	}
	if (num >= 1000) {
		return Math.round(num / 1000) + 'k';
	}
	return num.toLocaleString();
}
