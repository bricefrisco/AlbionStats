export function formatNumber(num) {
	if (num === 0) return '-';
	return num?.toLocaleString() ?? '0';
}

export function formatFame(num) {
	if (!num) return '-';
	if (num >= 1000000000) {
		return (num / 1000000000).toFixed(2) + 'b';
	}
	if (num >= 1000000) {
		return (num / 1000000).toFixed(2) + 'm';
	}
	if (num >= 1000) {
		return Math.round(num / 1000) + 'k';
	}
	return num.toLocaleString();
}

export function formatRatio(numerator, denominator) {
	if (!denominator) return '-';
	return (numerator / denominator).toFixed(2);
}

export function getRatioColor(numerator, denominator, min = 0.5, max = 1.5) {
	if (!denominator) return null;
	const ratio = numerator / denominator;
	const clamped = Math.max(min, Math.min(max, ratio));
	const t = (clamped - min) / (max - min);
	const hue = Math.round(t * 120);
	return `hsl(${hue} 70% 45%)`;
}

const compactDateFormatter =
	typeof Intl !== 'undefined'
		? new Intl.DateTimeFormat(undefined, { month: 'numeric', day: 'numeric' })
		: null;
const compactTimeFormatter =
	typeof Intl !== 'undefined'
		? new Intl.DateTimeFormat(undefined, { hour: 'numeric', hour12: true })
		: null;

export function formatCompactDate(value) {
	if (!compactDateFormatter || !compactTimeFormatter) return '';
	const date = new Date(value);
	return `${compactDateFormatter.format(date)}, ${compactTimeFormatter.format(date)}`;
}
