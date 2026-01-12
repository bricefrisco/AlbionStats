<script>
	import Header from './Header.svelte';

	let isDarkMode =
		typeof localStorage !== 'undefined' && localStorage.getItem('darkMode') === 'true';
	function toggleDarkMode() {
		isDarkMode = !isDarkMode;
		localStorage.setItem('darkMode', isDarkMode.toString());
		document.documentElement.classList.toggle('dark', isDarkMode);
		document.documentElement.style.colorScheme = isDarkMode ? 'dark' : 'light';
	}
</script>

<div class="min-h-screen bg-white text-gray-900 dark:bg-neutral-950 dark:text-neutral-100">
	<Header {isDarkMode} {toggleDarkMode} />
	<div class="mx-auto max-w-5xl px-4 py-8 xl:px-0">
		<slot />
	</div>
</div>
