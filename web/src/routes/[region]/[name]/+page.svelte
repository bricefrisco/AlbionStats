<script>
	import { page } from '$app/stores';
	import Page from '../../../components/Page.svelte';
	import PageHeader from '../../../components/PageHeader.svelte';

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
		<PageHeader title={`${decodedName}`} />
	{/if}
</Page>
