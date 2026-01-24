<script>
	import SearchBar from './SearchBar.svelte';
	import { untrack } from 'svelte';
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';
	import { regionState } from '$lib/regionState.svelte';

	let {
		links = true,
		onselect,
		oninput,
		label = '',
		placeholder = '',
		value = $bindable('')
	} = $props();

	let showDropdown = $state(false);
	let guilds = $state([]);
	let isSearching = $state(false);
	let searchTimeout;
	let activeIndex = $state(0);

	let showNoResults = $derived(value.length >= 1 && guilds.length === 0);

	$effect(() => {
		activeIndex = 0;
		// Include regionState.value to trigger re-search when region changes
		const currentRegion = regionState.value;
		const currentValue = value;
		if (currentValue.length >= 1 && currentRegion) {
			untrack(() => performSearch());
		} else {
			untrack(() => {
				guilds = [];
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
					`https://albionstats.bricefrisco.com/api/guilds/search/${regionState.value}/${encodeURIComponent(value)}`
				);
				const data = await response.json();
				guilds = data.guilds || [];
			} catch (error) {
				console.error('search failed', error);
				guilds = [];
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

	function selectGuild(index) {
		const guild = guilds[index];
		if (!guild) return;

		if (links) {
			goto(resolve(`/guilds/${regionState.value}/${encodeURIComponent(guild)}`), {
				keepFocus: true,
				noScroll: true
			});
		} else {
			value = guild;
			onselect?.(guild);
		}
		showDropdown = false;
	}
</script>

<SearchBar
	bind:value
	bind:showDropdown
	oninput={handleInput}
	onfocus={handleInput}
	itemCount={guilds.length}
	bind:activeIndex
	onselectIndex={selectGuild}
	{isSearching}
	{label}
	{placeholder}
>
	{#snippet dropdownContent()}
		{#if guilds.length > 0}
			{#each guilds as guild, index (guild)}
				{#if links}
					<a
						href={resolve(`/guilds/${regionState.value}/${encodeURIComponent(guild)}`)}
						class="block w-full border-b border-gray-100 px-3 py-2 text-left text-sm text-gray-900 last:border-none hover:bg-gray-50 dark:border-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-800 {index === activeIndex ? 'bg-gray-50 dark:bg-neutral-800' : ''}"
						onclick={() => selectGuild(index)}
					>
						<div class="font-medium">{guild}</div>
					</a>
				{:else}
					<button
						class="block w-full border-b border-gray-100 px-3 py-2 text-left text-sm text-gray-900 last:border-none hover:bg-gray-50 dark:border-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-800 {index === activeIndex ? 'bg-gray-50 dark:bg-neutral-800' : ''}"
						onclick={() => selectGuild(index)}
					>
						<div class="font-medium">{guild}</div>
					</button>
				{/if}
			{/each}
		{:else if showNoResults}
			<div class="px-3 py-2 text-sm text-gray-500 dark:text-neutral-400">No guilds found</div>
		{/if}
	{/snippet}
</SearchBar>
