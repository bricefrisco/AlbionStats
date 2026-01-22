<script>
	import SearchBar from './SearchBar.svelte';
	import { regionState } from '$lib/regionState.svelte';

	let {
		onselect,
		oninput,
		value = $bindable('')
	} = $props();

	let showDropdown = $state(false);
	let alliances = $state([]);
	let isSearching = $state(false);
	let searchTimeout;

	let showNoResults = $derived(value.length >= 1 && alliances.length === 0);

	$effect(() => {
		// Include regionState.value to trigger re-search when region changes
		const currentRegion = regionState.value;
		if (value.length >= 1 && currentRegion) {
			performSearch();
		} else {
			alliances = [];
		}
	});

	async function performSearch() {
		if (value.trim().length < 1) return;

		isSearching = true;
		clearTimeout(searchTimeout);

		searchTimeout = setTimeout(async () => {
			try {
				const response = await fetch(
					`https://albionstats.bricefrisco.com/api/alliances/search/${regionState.value}/${encodeURIComponent(value)}`
				);
				const data = await response.json();
				alliances = data.alliances || [];
			} catch (error) {
				console.error('search failed', error);
				alliances = [];
			} finally {
				isSearching = false;
			}
		}, 200);
	}

	async function handleInput(e) {
		oninput?.(e);
		if (value.trim().length < 1) {
			showDropdown = false;
			return;
		}
		showDropdown = true;
	}
</script>

<SearchBar
	bind:value
	bind:showDropdown
	oninput={handleInput}
	onfocus={handleInput}
	{isSearching}
	placeholder="Alliance name"
>
	{#snippet dropdownContent()}
		{#if alliances.length > 0}
			{#each alliances as alliance (alliance)}
				<button
					class="block w-full border-b border-gray-100 px-3 py-2 text-left text-sm text-gray-900 last:border-none hover:bg-gray-50 dark:border-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-800"
					onclick={() => {
						value = alliance;
						showDropdown = false;
						onselect?.(alliance);
					}}
				>
					<div class="font-medium">{alliance}</div>
				</button>
			{/each}
		{:else if showNoResults}
			<div class="px-3 py-2 text-sm text-gray-500 dark:text-neutral-400">No alliances found</div>
		{/if}
	{/snippet}
</SearchBar>
