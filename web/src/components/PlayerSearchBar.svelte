<script>
	let { onSelect } = $props();

	import SearchIcon from './icons/SearchIcon.svelte';

	import { regionState } from '$lib/regionState.svelte';

	let container = $state();

	let showDropdown = $state(false);
	let searchQuery = $state('');
	let players = $state([]);

	let searchTimeout;
	let isSearching = $state(false);

	let showNoResults = $derived(searchQuery.length >= 2 && players.length === 0);

	async function handleInput() {
		if (searchQuery.trim().length < 2) {
			players = [];
			showDropdown = false;
			clearTimeout(searchTimeout);
			isSearching = false;
			return;
		}
		showDropdown = true;
		isSearching = true;

		clearTimeout(searchTimeout);

		searchTimeout = setTimeout(async () => {
			await performSearch();
		}, 200);
	}

	async function performSearch() {
		try {
			const response = await fetch(
				`https://albionstats.bricefrisco.com/api/search/${regionState.value}/${encodeURIComponent(searchQuery)}`
			);
			const data = await response.json();
			players = data.players || [];
		} catch (error) {
			console.error('search failed', error);
			players = [];
		} finally {
			isSearching = false;
			showDropdown = Boolean(players.length || searchQuery.length >= 2);
		}
	}

	const handleDocumentClick = (event) => {
		if (container && !container.contains(event.target)) {
			showDropdown = false;
			clearTimeout(searchTimeout);
			isSearching = false;
		}
	};

	const handleOptionClick = (player) => {
		showDropdown = false;
		searchQuery = player.name;
		onSelect?.(player);
	};

	import { onMount } from 'svelte';

	onMount(() => {
		document.addEventListener('click', handleDocumentClick);
		return () => document.removeEventListener('click', handleDocumentClick);
	});
</script>

<div class="relative" bind:this={container}>
	<div class="relative flex-1">
		<input
			type="text"
			bind:value={searchQuery}
			oninput={handleInput}
			onfocus={handleInput}
			placeholder="Player name"
			class="w-full rounded border border-gray-300 bg-gray-50 px-3 py-2 pl-9 text-sm text-gray-900 placeholder-gray-500 focus:border-gray-400 focus:outline-none dark:border-neutral-800 dark:bg-neutral-900 dark:text-neutral-100 dark:placeholder-neutral-500 dark:focus:border-neutral-700"
		/>
		<div class="absolute top-[11px] left-3 flex items-center text-gray-500 dark:text-neutral-500">
			<SearchIcon size={16} />
		</div>
	</div>

	{#if showDropdown}
		<div
			class="absolute top-full right-0 left-0 z-10 mt-2 max-h-60 overflow-y-auto rounded border border-gray-200 bg-white shadow-lg dark:border-neutral-700 dark:bg-neutral-900"
		>
			{#if isSearching}
				<div class="px-3 py-2 text-sm text-gray-500 dark:text-neutral-400">Searching...</div>
			{:else if players.length > 0}
				{#each players as player (player.name)}
					<button
						type="button"
						class="block w-full border-b border-gray-100 px-3 py-2 text-left text-sm text-gray-900 last:border-none hover:bg-gray-50 dark:border-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-800"
						onclick={() => handleOptionClick(player)}
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
				{/each}
			{:else if showNoResults}
				<div class="px-3 py-2 text-sm text-gray-500 dark:text-neutral-400">No players found</div>
			{/if}
		</div>
	{/if}
</div>
