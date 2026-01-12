<script>
	import GithubIcon from './icons/GithubIcon.svelte';
	import SunIcon from './icons/SunIcon.svelte';
	import MoonIcon from './icons/MoonIcon.svelte';
	import SearchBar from './SearchBar.svelte';

	export let isDarkMode;
	export let toggleDarkMode;
</script>

<header class="border-b border-gray-200 dark:border-neutral-800">
	<div class="mx-auto max-w-5xl px-4 py-4 xl:px-0">
		<div class="mb-4 flex items-center justify-between">
			<a
				href="/"
				class="text-lg font-medium text-gray-900 hover:text-gray-700 dark:text-white dark:hover:text-gray-300"
			>
				AlbionStats
			</a>
			<div class="flex items-center gap-3">
				<button
					on:click={toggleDarkMode}
					class="flex items-center text-gray-600 transition-colors hover:text-gray-900 dark:text-neutral-400 dark:hover:text-neutral-100"
					aria-label="Toggle dark mode"
				>
					{#if isDarkMode}
						<SunIcon size={20} />
					{:else}
						<MoonIcon size={20} />
					{/if}
				</button>
				<a
					href="https://github.com/bricefrisco/AlbionStats"
					target="_blank"
					rel="noopener noreferrer"
					class="flex items-center text-gray-600 transition-colors hover:text-gray-900 dark:text-neutral-400 dark:hover:text-neutral-100"
				>
					<GithubIcon size={20} />
				</a>
			</div>
		</div>

		<SearchBar />
	</div>
</header>
