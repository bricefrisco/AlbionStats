<script>
	import { resolve } from '$app/paths';
	import GithubIcon from './icons/GithubIcon.svelte';
	import RegionSelect from './RegionSelect.svelte';
	import { regionState } from '$lib/regionState.svelte';
	import { routes } from '$lib/regionRoutes';

	let isMenuOpen = $state(false);
</script>

<header
	class="sticky top-0 z-50 border-b border-gray-200 bg-white dark:border-neutral-800 dark:bg-neutral-950"
>
	<div class="mx-auto max-w-5xl px-4 py-3 xl:px-0">
		<div class="flex items-center justify-between lg:hidden">
			<a
				href={resolve('/')}
				class="text-lg font-medium text-gray-900 hover:text-gray-700 dark:text-white dark:hover:text-gray-300"
			>
				AlbionStats
			</a>
			<div class="flex items-center gap-3">
				<RegionSelect />
				<a
					href="https://github.com/bricefrisco/AlbionStats"
					target="_blank"
					rel="noopener noreferrer"
					class="flex items-center text-gray-600 transition-colors hover:text-gray-900 dark:text-neutral-400 dark:hover:text-neutral-100"
				>
					<GithubIcon size={20} />
				</a>
				<button
					type="button"
					class="flex items-center text-gray-600 transition-colors hover:text-gray-900 dark:text-neutral-400 dark:hover:text-neutral-100"
					aria-label="Toggle menu"
					aria-expanded={isMenuOpen}
					on:click={() => (isMenuOpen = !isMenuOpen)}
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						class="h-5 w-5"
					>
						<line x1="3" y1="6" x2="21" y2="6" />
						<line x1="3" y1="12" x2="21" y2="12" />
						<line x1="3" y1="18" x2="21" y2="18" />
					</svg>
				</button>
			</div>
		</div>

		{#if isMenuOpen}
			<nav class="mt-4 grid grid-cols-2 gap-2 text-sm lg:hidden">
				{#each routes as route (route.base)}
					<a
						href={resolve(`/${route.base}/${regionState.value}`)}
						class="rounded-md border border-gray-200 px-3 py-2 text-gray-700 hover:bg-gray-50 hover:text-gray-900 dark:border-neutral-800 dark:text-gray-300 dark:hover:bg-neutral-900 dark:hover:text-gray-100"
						on:click={() => (isMenuOpen = false)}
					>
						{route.label}
					</a>
				{/each}
			</nav>
		{/if}

		<div class="hidden items-center justify-between lg:flex">
			<nav class="flex items-center gap-6">
				<a
					href={resolve('/')}
					class="text-lg font-medium text-gray-900 hover:text-gray-700 dark:text-white dark:hover:text-gray-300 pr-4 border-r border-gray-300 dark:border-gray-600"
				>
					AlbionStats
				</a>
				{#each routes as route (route.base)}
					<a
						href={resolve(`/${route.base}/${regionState.value}`)}
						class="text-sm text-gray-700 hover:text-gray-900 dark:text-gray-300 dark:hover:text-gray-100"
					>
						{route.label}
					</a>
				{/each}
			</nav>
			<div class="flex items-center gap-3">
				<RegionSelect />
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
