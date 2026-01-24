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
	let players = $state([]);
	let isSearching = $state(false);
	let searchTimeout;
	let activeIndex = $state(0);

	let showNoResults = $derived(value.length >= 3 && players.length === 0);

	$effect(() => {
		activeIndex = 0;
		// Include regionState.value to trigger re-search when region changes
		const currentRegion = regionState.value;
		const currentValue = value;
		if (currentValue.length >= 3 && currentRegion) {
			untrack(() => performSearch());
		} else {
			untrack(() => {
				players = [];
			});
		}
	});

	async function performSearch() {
		if (value.trim().length < 3) return;

		isSearching = true;
		clearTimeout(searchTimeout);

		searchTimeout = setTimeout(async () => {
			try {
				const response = await fetch(
					`https://albionstats.bricefrisco.com/api/players/search/${regionState.value}/${encodeURIComponent(value)}`
				);
				const data = await response.json();
				players = data.players || [];
			} catch (error) {
				console.error('search failed', error);
				players = [];
			} finally {
				isSearching = false;
			}
		}, 200);
	}

	async function handleInput(e) {
		oninput?.(e);
		if (value.trim().length < 3) {
			showDropdown = false;
			return;
		}
		showDropdown = true;
	}

	function selectPlayer(index) {
		const player = players[index];
		if (!player?.name) return;

		if (links) {
			goto(resolve(`/players/${regionState.value}/${encodeURIComponent(player.name)}`), {
				keepFocus: true,
				noScroll: true
			});
		} else {
			value = player.name;
			onselect?.(player);
		}
		showDropdown = false;
	}
</script>

<SearchBar
	bind:value
	bind:showDropdown
	oninput={handleInput}
	onfocus={handleInput}
	itemCount={players.length}
	bind:activeIndex
	onselectIndex={selectPlayer}
	{isSearching}
	{label}
	{placeholder}
>
	{#snippet dropdownContent()}
		{#if players.length > 0}
			{#each players as player, index (player.name)}
				{#if links}
					<a
						href={resolve(`/players/${regionState.value}/${encodeURIComponent(player.name)}`)}
						class="block w-full border-b border-gray-100 px-3 py-2 text-left text-sm text-gray-900 last:border-none hover:bg-gray-50 dark:border-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-800 {index === activeIndex ? 'bg-gray-50 dark:bg-neutral-800' : ''}"
						onclick={() => selectPlayer(index)}
					>
						<div class="font-medium">{player.name}</div>
						<div class="text-xs text-gray-500 dark:text-neutral-400">
							{#if player.guild_name && player.alliance_name}
								[{player.alliance_name}] {player.guild_name}
							{:else if player.guild_name}
								{player.guild_name}
							{/if}
						</div>
					</a>
				{:else}
					<button
						class="block w-full border-b border-gray-100 px-3 py-2 text-left text-sm text-gray-900 last:border-none hover:bg-gray-50 dark:border-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-800 {index === activeIndex ? 'bg-gray-50 dark:bg-neutral-800' : ''}"
						onclick={() => selectPlayer(index)}
					>
						<div class="font-medium">{player.name}</div>
						<div class="text-xs text-gray-500 dark:text-neutral-400">
							{#if player.guild_name && player.alliance_name}
								[{player.alliance_name}] {player.guild_name}
							{:else if player.guild_name}
								{player.guild_name}
							{/if}
						</div>
					</button>
				{/if}
			{/each}
		{:else if showNoResults}
			<div class="px-3 py-2 text-sm text-gray-500 dark:text-neutral-400">No players found</div>
		{/if}
	{/snippet}
</SearchBar>
