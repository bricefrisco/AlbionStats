<script>
	import Page from '../components/Page.svelte';
	import Chart from '../components/charts/Chart.svelte';
</script>

<Page>
	<h1 class="mb-4 text-2xl font-bold text-gray-900 dark:text-white">Welcome to AlbionStats</h1>
	<p class="mb-4 text-gray-700 dark:text-gray-300">
		We track Albion Online player statistics over time.
	</p>

	<div class="mt-8">
		<h2 class="mb-4 text-xl font-semibold text-gray-900 dark:text-white">Players Tracked</h2>
		<Chart
			timestamps={[1768150800000, 1768154400000, 1768158000000]}
			values={[30153, 36612, 40170]}
			label="Players Tracked"
		/>
	</div>
</Page>
