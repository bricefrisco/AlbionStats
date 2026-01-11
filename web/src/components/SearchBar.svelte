<script>
	import SearchIcon from './icons/SearchIcon.svelte';
	import ChevronDownIcon from './icons/ChevronDownIcon.svelte';

	let selectedRegion = 'americas';
	const regions = [
		{ label: 'Americas', value: 'americas' },
		{ label: 'Europe', value: 'europe' },
		{ label: 'Asia', value: 'asia' }
	];

	let searchQuery = '';
	let players = [];
	let isFetching = false;
	let showDropdown = false;
	let container;

	function handleRegionChange(region) {
		selectedRegion = region;
	}

	async function handleSearch(e) {
		e.preventDefault();
		if (!searchQuery.trim()) {
			players = [];
			showDropdown = false;
			return;
		}
		showDropdown = true;
		await performSearch();
	}

	async function handleInput() {
		if (searchQuery.trim().length < 2) {
			players = [];
			showDropdown = false;
			return;
		}
		showDropdown = true;
		await performSearch();
	}

	async function performSearch() {
		isFetching = true;
		try {
			const response = await fetch(
				`https://albionstats.bricefrisco.com/api/search/${selectedRegion}/${encodeURIComponent(searchQuery)}`
			);
			const data = await response.json();
			players = data.players || [];
		} catch (error) {
			console.error('search failed', error);
			players = [];
		} finally {
			isFetching = false;
			showDropdown = Boolean(players.length || searchQuery.length >= 2);
		}
	}

	const handleClick = (event) => {
		if (container && !container.contains(event.target)) {
			showDropdown = false;
		}
	};

	import { onMount } from 'svelte';

	onMount(() => {
		document.addEventListener('click', handleClick);
		return () => document.removeEventListener('click', handleClick);
	});
</script>

<div class="relative" bind:this={container}>
	<form on:submit={handleSearch} class="flex gap-2">
		<div class="relative flex-1">
			<input
				type="text"
				bind:value={searchQuery}
				on:input={handleInput}
				placeholder="Player name"
				class="w-full rounded border border-gray-300 bg-gray-50 px-3 py-2 pl-9 text-sm text-gray-900 placeholder-gray-500 focus:border-gray-400 focus:outline-none dark:border-neutral-800 dark:bg-neutral-900 dark:text-neutral-100 dark:placeholder-neutral-500 dark:focus:border-neutral-700"
			/>
			<div class="absolute top-[11px] left-3 flex items-center text-gray-500 dark:text-neutral-500">
				<SearchIcon size={16} />
			</div>
		</div>

		<div class="relative">
			<select
				value={selectedRegion}
				on:change={(e) => handleRegionChange(e.target.value)}
				class="cursor-pointer appearance-none rounded border border-gray-300 bg-gray-50 px-3 py-2 pr-8 text-sm text-gray-900 focus:border-gray-400 focus:outline-none dark:border-neutral-800 dark:bg-neutral-900 dark:text-neutral-100 dark:focus:border-neutral-700"
			>
				{#each regions as region}
					<option value={region.value}>
						{region.label}
					</option>
				{/each}
			</select>
			<div
				class="pointer-events-none absolute top-1/2 right-2 flex -translate-y-1/2 items-center text-gray-500 dark:text-neutral-500"
			>
				<ChevronDownIcon size={16} />
			</div>
		</div>

		<button
			type="submit"
			class="rounded bg-gray-100 px-4 py-2 text-sm font-medium text-gray-900 transition-colors hover:bg-gray-200 dark:bg-neutral-100 dark:text-neutral-950 dark:hover:bg-white"
		>
			Search
		</button>
	</form>

	{#if showDropdown && players.length > 0}
		<div
			class="absolute top-full right-0 left-0 z-10 mt-2 max-h-60 overflow-y-auto rounded border border-gray-200 bg-white shadow-lg dark:border-neutral-700 dark:bg-neutral-900"
		>
			{#each players as player}
				<div
					class="border-b border-gray-100 px-3 py-2 text-sm text-gray-900 last:border-none dark:border-neutral-800 dark:text-neutral-100"
				>
					<div class="font-medium">{player.name}</div>
					<div class="text-xs text-gray-500 dark:text-neutral-400">
						{player.guild_name ?? 'No guild'} · {player.alliance_name ?? 'No alliance'}
					</div>
				</div>
			{/each}
		</div>
	{:else if showDropdown && isFetching}
		<div
			class="absolute top-full right-0 left-0 z-10 mt-2 rounded border border-gray-200 bg-white px-3 py-2 text-sm text-gray-500 shadow-lg dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-400"
		>
			Searching...
		</div>
	{:else if showDropdown && searchQuery.length >= 2}
		<div
			class="absolute top-full right-0 left-0 z-10 mt-2 rounded border border-gray-200 bg-white px-3 py-2 text-sm text-gray-500 shadow-lg dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-400"
		>
			No players found
		</div>
	{/if}
</div>
