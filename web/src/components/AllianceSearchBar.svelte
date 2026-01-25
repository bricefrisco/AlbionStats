<script>
	import SearchBar from './SearchBar.svelte';
	import { untrack } from 'svelte';
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';
	import { regionState } from '$lib/regionState.svelte';
	import { getApiBase } from '$lib/apiBase.js';

	let {
		links = true,
		onselect,
		oninput,
		label = '',
		placeholder = '',
		value = $bindable('')
	} = $props();

	let showDropdown = $state(false);
	let alliances = $state([]);
	let isSearching = $state(false);
	let searchTimeout;
	let activeIndex = $state(0);

	let showNoResults = $derived(value.length >= 1 && alliances.length === 0);

	$effect(() => {
		activeIndex = 0;
		// Include regionState.value to trigger re-search when region changes
		const currentRegion = regionState.value;
		const currentValue = value;
		if (currentValue.length >= 1 && currentRegion) {
			untrack(() => performSearch());
		} else {
			untrack(() => {
				alliances = [];
			});
		}
	});

	async function performSearch() {
		if (value.trim().length < 1) return;

		isSearching = true;
		clearTimeout(searchTimeout);

		searchTimeout = setTimeout(async () => {
			try {
				const response = await fetch(
					`https://albionstats.com/api/alliances/search/${regionState.value}/${encodeURIComponent(value)}`
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

	function selectAlliance(index) {
		const alliance = alliances[index];
		if (!alliance) return;

		if (links) {
			goto(resolve(`/alliances/${regionState.value}/${encodeURIComponent(alliance)}`), {
				keepFocus: true,
				noScroll: true
			});
		} else {
			value = alliance;
			onselect?.(alliance);
		}
		showDropdown = false;
	}
</script>

<SearchBar
	bind:value
	bind:showDropdown
	oninput={handleInput}
	onfocus={handleInput}
	itemCount={alliances.length}
	bind:activeIndex
	onselectIndex={selectAlliance}
	allowSubmitOnEnter={!links}
	{isSearching}
	{label}
	{placeholder}
>
	{#snippet dropdownContent()}
		{#if alliances.length > 0}
			{#each alliances as alliance, index (alliance)}
				{#if links}
					<a
						href={resolve(`/alliances/${regionState.value}/${encodeURIComponent(alliance)}`)}
						class="block w-full border-b border-gray-100 px-3 py-2 text-left text-sm text-gray-900 last:border-none hover:bg-gray-50 dark:border-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-800 {index === activeIndex ? 'bg-gray-50 dark:bg-neutral-800' : ''}"
						onclick={() => selectAlliance(index)}
					>
						<div class="font-medium">{alliance}</div>
					</a>
				{:else}
					<button
						class="block w-full border-b border-gray-100 px-3 py-2 text-left text-sm text-gray-900 last:border-none hover:bg-gray-50 dark:border-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-800 {index === activeIndex ? 'bg-gray-50 dark:bg-neutral-800' : ''}"
						onclick={() => selectAlliance(index)}
					>
						<div class="font-medium">{alliance}</div>
					</button>
				{/if}
			{/each}
		{:else if showNoResults}
			<div class="px-3 py-2 text-sm text-gray-500 dark:text-neutral-400">No alliances found</div>
		{/if}
	{/snippet}
</SearchBar>
