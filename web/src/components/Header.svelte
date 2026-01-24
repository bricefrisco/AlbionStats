<script>
	import { resolve } from '$app/paths';
	import GithubIcon from './icons/GithubIcon.svelte';
	import SunIcon from './icons/SunIcon.svelte';
	import MoonIcon from './icons/MoonIcon.svelte';
	import RegionSelect from './RegionSelect.svelte';

	let { isDarkMode, toggleDarkMode } = $props();
</script>

<header
	class="sticky top-0 z-50 border-b border-gray-200 bg-white dark:border-neutral-800 dark:bg-neutral-950"
>
	<div class="mx-auto max-w-5xl px-4 py-3 xl:px-0">
		<div class="flex items-center justify-between">
			<nav class="flex items-center gap-6">
				<a
					href={resolve('/')}
					class="text-lg font-medium text-gray-900 hover:text-gray-700 dark:text-white dark:hover:text-gray-300 pr-4 border-r border-gray-300 dark:border-gray-600"
				>
					AlbionStats
				</a>
				<a
					href={resolve('/battle-boards')}
					class="text-sm text-gray-700 hover:text-gray-900 dark:text-gray-300 dark:hover:text-gray-100"
				>
					Battle Boards
				</a>
				<a
					href={resolve('/alliances')}
					class="text-sm text-gray-700 hover:text-gray-900 dark:text-gray-300 dark:hover:text-gray-100"
				>
					Alliances
				</a>
				<a
					href={resolve('/guilds')}
					class="text-sm text-gray-700 hover:text-gray-900 dark:text-gray-300 dark:hover:text-gray-100"
				>
					Guilds
				</a>
				<a
					href={resolve('/players')}
					class="text-sm text-gray-700 hover:text-gray-900 dark:text-gray-300 dark:hover:text-gray-100"
				>
					Players
				</a>
			</nav>
			<div class="flex items-center gap-3">
				<RegionSelect />
				<button
					onclick={toggleDarkMode}
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
	</div>
</header>
