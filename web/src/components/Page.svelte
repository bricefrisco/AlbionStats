<script>
	import Header from './Header.svelte';

	let isDarkMode = typeof localStorage !== 'undefined' && localStorage.getItem('darkMode') === 'true';
	function toggleDarkMode() {
		isDarkMode = !isDarkMode;
		localStorage.setItem('darkMode', isDarkMode.toString());
		document.documentElement.classList.toggle('dark', isDarkMode);
	}
</script>

<div class="min-h-screen bg-white dark:bg-neutral-950 text-gray-900 dark:text-neutral-100">
	<Header {isDarkMode} {toggleDarkMode} />
	<div class="max-w-5xl mx-auto py-8">
		<slot />
	</div>
</div>
