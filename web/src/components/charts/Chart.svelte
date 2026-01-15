<script>
	import { onMount } from 'svelte';
	import Chart from 'chart.js/auto';
	import 'chartjs-adapter-date-fns';

	// Props
	export let timestamps = [];
	export let values = [];
	export let label = 'Metric';
	export let height = 'h-80';

	let canvas;
	let chart;

	// Reactive theme detection with MutationObserver for dynamic updates
	let isDark =
		typeof document !== 'undefined' && document.documentElement.classList.contains('dark');
	$: textColor = isDark ? 'rgb(209, 213, 219)' : 'rgb(55, 65, 81)'; // gray-300 : gray-700
	$: gridColor = isDark ? 'rgb(17, 24, 39)' : 'rgb(229, 231, 235)'; // gray-900 : gray-200

	const compactDateFormatter =
		typeof Intl !== 'undefined'
			? new Intl.DateTimeFormat(undefined, { month: 'numeric', day: 'numeric' })
			: null;
	const compactTimeFormatter =
		typeof Intl !== 'undefined'
			? new Intl.DateTimeFormat(undefined, {
					hour: 'numeric',
					hour12: true
				})
			: null;

	const formatCompactDate = (value) => {
		if (!compactDateFormatter || !compactTimeFormatter) return '';
		const date = new Date(value);
		return `${compactDateFormatter.format(date)}, ${compactTimeFormatter.format(date)}`;
	};

	function updateTheme() {
		const newIsDark =
			typeof document !== 'undefined' && document.documentElement.classList.contains('dark');
		if (newIsDark !== isDark) {
			isDark = newIsDark;
		}
	}

	// Chart data reactive to props
	$: data = {
		labels: timestamps,
		datasets: [
			{
				label,
				data: values,
				borderColor: 'rgb(75, 192, 192)',
				backgroundColor: 'rgba(75, 192, 192, 0.2)',
				tension: 0.4
			}
		]
	};

	$: options = {
		responsive: true,
		maintainAspectRatio: false,
		plugins: {
			legend: {
				display: false
			},
			title: {
				display: false
			},
			tooltip: {
				callbacks: {
					label: (context) => {
						const dateStr = formatCompactDate(context.parsed.x);
						return `${dateStr}: ${context.formattedValue}`;
					}
				}
			}
		},
		scales: {
			y: {
				beginAtZero: false,
				title: {
					display: false
				},
				ticks: {
					color: textColor
				},
				grid: {
					color: gridColor
				}
			},
			x: {
				type: 'time',
				time: {
					unit: 'hour'
				},
				title: {
					display: false
				},
				ticks: {
					color: textColor,
					callback: (value) => formatCompactDate(value)
				},
				grid: {
					color: gridColor
				}
			}
		}
	};

	onMount(() => {
		if (canvas) {
			chart = new Chart(canvas, {
				type: 'line',
				data,
				options
			});
		}

		// Watch for theme changes
		const observer = new MutationObserver(() => {
			updateTheme();
		});

		if (typeof document !== 'undefined') {
			observer.observe(document.documentElement, {
				attributes: true,
				attributeFilter: ['class']
			});
		}

		return () => {
			if (chart) {
				chart.destroy();
			}
			observer.disconnect();
		};
	});

	// Update chart when data or theme changes
	$: if (chart) {
		if (textColor) {
			chart.options.scales.y.ticks.color = textColor;
			chart.options.scales.x.ticks.color = textColor;
		}
		if (gridColor) {
			chart.options.scales.y.grid.color = gridColor;
			chart.options.scales.x.grid.color = gridColor;
		}
		if (data && data.labels && data.labels.length > 0) {
			chart.data = data;
		}
		chart.update();
	}
</script>

<div class={height}>
	<canvas bind:this={canvas}></canvas>
</div>

<style>
	canvas {
		width: 100% !important;
		height: 100% !important;
	}
</style>
