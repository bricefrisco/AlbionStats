<script>
	import { onMount } from 'svelte';
	import { formatCompactDate } from '$lib/utils';

	let {
		timestamps = [],
		values = [],
		label = 'Metric',
		height = 'h-80',
		color = 'rgb(75, 192, 192)',
		datasets = null,
		showLegend = false
	} = $props();

	let canvas = $state();
	let chart = $state();

	// Reactive theme detection with MutationObserver for dynamic updates
	let isDark = $state(
		typeof document !== 'undefined' && document.documentElement.classList.contains('dark')
	);
	const textColor = $derived(isDark ? 'rgb(209, 213, 219)' : 'rgb(55, 65, 81)'); // gray-300 : gray-700
	const gridColor = $derived(isDark ? 'rgb(17, 24, 39)' : 'rgb(229, 231, 235)'); // gray-900 : gray-200

	function updateTheme() {
		const newIsDark =
			typeof document !== 'undefined' && document.documentElement.classList.contains('dark');
		if (newIsDark !== isDark) {
			isDark = newIsDark;
		}
	}

	function buildData() {
		if (Array.isArray(datasets) && datasets.length > 0) {
			return {
				labels: Array.from(timestamps || []),
				datasets: datasets.map((dataset) => ({
					tension: 0.4,
					...dataset
				}))
			};
		}

		return {
			labels: Array.from(timestamps || []),
			datasets: [
				{
					label,
					data: Array.from(values || []),
					borderColor: color,
					backgroundColor: 'rgba(75, 192, 192, 0.2)',
					tension: 0.4
				}
			]
		};
	}

	const options = $derived({
		responsive: true,
		maintainAspectRatio: false,
		plugins: {
			legend: {
				display: showLegend,
				labels: {
					color: textColor
				}
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
	});

	onMount(async () => {
		if (!canvas) {
			return;
		}

		const [{ default: Chart }] = await Promise.all([
			import('chart.js/auto'),
			import('chartjs-adapter-date-fns')
		]);

		chart = new Chart(canvas, {
			type: 'line',
			data: buildData(),
			options
		});

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
	$effect(() => {
		if (!chart) return;
		if (textColor) {
			chart.options.scales.y.ticks.color = textColor;
			chart.options.scales.x.ticks.color = textColor;
			if (chart.options.plugins?.legend?.labels) {
				chart.options.plugins.legend.labels.color = textColor;
			}
		}
		if (gridColor) {
			chart.options.scales.y.grid.color = gridColor;
			chart.options.scales.x.grid.color = gridColor;
		}
		const nextData = buildData();
		if (nextData.labels.length > 0) {
			chart.data = nextData;
		}
		chart.update();
	});
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
