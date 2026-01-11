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

	function handleRegionChange(region) {
		selectedRegion = region;
	}

	function handleSearch(e) {
		e.preventDefault();
		console.log('Searching for:', searchQuery, 'in region:', selectedRegion);
		// TODO: Implement actual search logic
	}
</script>

<form on:submit={handleSearch} class="flex gap-2">
	<div class="flex-1 relative">
		<input
			type="text"
			bind:value={searchQuery}
			placeholder="Player name"
			class="w-full px-3 py-2 pl-9 bg-gray-50 dark:bg-neutral-900 border border-gray-300 dark:border-neutral-800 rounded text-sm text-gray-900 dark:text-neutral-100 placeholder-gray-500 dark:placeholder-neutral-500 focus:outline-none focus:border-gray-400 dark:focus:border-neutral-700"
		/>
		<div class="absolute left-3 top-[11px] flex items-center text-gray-500 dark:text-neutral-500">
			<SearchIcon size={16} />
		</div>
	</div>

	<div class="relative">
		<select
			value={selectedRegion}
			on:change={(e) => handleRegionChange(e.target.value)}
			class="px-3 pr-8 py-2 bg-gray-50 dark:bg-neutral-900 border border-gray-300 dark:border-neutral-800 rounded text-sm text-gray-900 dark:text-neutral-100 focus:outline-none focus:border-gray-400 dark:focus:border-neutral-700 cursor-pointer appearance-none"
		>
			{#each regions as region}
				<option value={region.value}>
					{region.label}
				</option>
			{/each}
		</select>
		<div class="absolute right-2 top-1/2 -translate-y-1/2 flex items-center text-gray-500 dark:text-neutral-500 pointer-events-none">
			<ChevronDownIcon size={16} />
		</div>
	</div>

	<button
		type="submit"
		class="px-4 py-2 bg-gray-100 hover:bg-gray-200 dark:bg-neutral-100 dark:hover:bg-white text-gray-900 dark:text-neutral-950 rounded text-sm font-medium transition-colors"
	>
		Search
	</button>
</form>
