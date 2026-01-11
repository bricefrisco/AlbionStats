<script>
	import { onMount } from 'svelte';
	import Chart from 'chart.js/auto';
	import 'chartjs-adapter-date-fns';

	// Props
	export let timestamps = [];
	export let values = [];
	export let label = 'Metric';

	let canvas;
	let chart;

	// Reactive theme detection with MutationObserver for dynamic updates
	let isDark =
		typeof document !== 'undefined' && document.documentElement.classList.contains('dark');
	$: textColor = isDark ? 'rgb(209, 213, 219)' : 'rgb(55, 65, 81)'; // gray-300 : gray-700

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
					color: textColor
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
		if (data && data.labels && data.labels.length > 0) {
			chart.data = data;
		}
		chart.update();
	}
</script>

<div class="h-80">
	<canvas bind:this={canvas}></canvas>
</div>

<style>
	canvas {
		width: 100% !important;
		height: 100% !important;
	}
</style>
