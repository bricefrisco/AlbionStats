<script>
	import { page } from '$app/stores';
	import Page from '../../../components/Page.svelte';

	// Get parameters from URL
	$: region = $page.params.region;
	$: playerName = $page.params.name;
	$: decodedName = playerName ? decodeURIComponent(playerName) : '';

	// Validate region
	$: validRegion = ['americas', 'europe', 'asia'].includes(region);
</script>

<Page>
	{#if !validRegion}
		<div class="text-center">
			<h1 class="mb-4 text-2xl font-bold text-red-600 dark:text-red-400">Invalid Region</h1>
			<p class="text-gray-700 dark:text-gray-300">Valid regions are: americas, europe, asia</p>
		</div>
	{:else if !decodedName}
		<div class="text-center">
			<h1 class="mb-4 text-2xl font-bold text-red-600 dark:text-red-400">Player Not Found</h1>
			<p class="text-gray-700 dark:text-gray-300">Please provide a valid player name</p>
		</div>
	{:else}
		<div class="text-center">
			<h1 class="mb-4 text-3xl font-bold text-gray-900 dark:text-white">
				{decodedName}
			</h1>
			<p class="text-lg text-gray-600 capitalize dark:text-gray-400">
				{region} Server
			</p>
			<div class="mt-8 rounded-lg bg-gray-50 p-6 dark:bg-neutral-800">
				<p class="text-gray-700 dark:text-gray-300">
					This is the player page for <strong>{decodedName}</strong> on the
					<strong>{region}</strong> server.
				</p>
				<p class="mt-4 text-sm text-gray-500 dark:text-gray-400">
					Player details and charts will be added here.
				</p>
			</div>
		</div>
	{/if}
</Page>
