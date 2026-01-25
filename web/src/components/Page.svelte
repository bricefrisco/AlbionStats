<script>
	import Header from './Header.svelte';
	import Footer from './Footer.svelte';
	import { page } from '$app/state';
	import { browser } from '$app/environment';
	import { regionState } from '$lib/regionState.svelte';
	import { validRegions } from '$lib/utils';

	let isDarkMode = $state(
		typeof localStorage !== 'undefined' && localStorage.getItem('darkMode') === 'true'
	);
	function toggleDarkMode() {
		isDarkMode = !isDarkMode;
		localStorage.setItem('darkMode', isDarkMode.toString());
		document.documentElement.classList.toggle('dark', isDarkMode);
		document.documentElement.style.colorScheme = isDarkMode ? 'dark' : 'light';
	}

	let { children } = $props();

	$effect(() => {
		if (!browser) return;
		const segments = page.url.pathname.split('/').filter(Boolean);
		const region = segments.length > 1 ? segments[1] : null;
		if (region && validRegions.has(region) && regionState.value !== region) {
			regionState.value = region;
		}
	});
</script>

<div class="min-h-screen bg-white text-gray-900 dark:bg-neutral-950 dark:text-neutral-100">
	<Header {isDarkMode} {toggleDarkMode} />
	<div class="mx-auto max-w-5xl px-4 py-8 xl:px-0">
		{@render children()}
	</div>
	<Footer />
</div>
