<script>
	import Page from '../components/Page.svelte';
	import PageHeader from '../components/PageHeader.svelte';
	import Typography from '../components/Typography.svelte';
	import PlayersTracked from '../components/charts/PlayersTracked.svelte';
	import TotalDataPoints from '../components/charts/TotalDataPoints.svelte';
	import SubHeader from '../components/SubHeader.svelte';
</script>

<Page>
	<PageHeader title="Welcome to AlbionStats" />
	<Typography>
		<p>We track Albion Online player statistics over time, including:</p>
		<div class="mt-3 grid grid-cols-2 gap-2 text-sm md:grid-cols-3">
			<div class="flex items-center">
				<span class="mr-2 h-1.5 w-1.5 rounded-full bg-gray-400"></span>
				Kill fame
			</div>
			<div class="flex items-center">
				<span class="mr-2 h-1.5 w-1.5 rounded-full bg-gray-400"></span>
				Death fame
			</div>
			<div class="flex items-center">
				<span class="mr-2 h-1.5 w-1.5 rounded-full bg-gray-400"></span>
				PvE Fame
			</div>
			<div class="flex items-center">
				<span class="mr-2 h-1.5 w-1.5 rounded-full bg-gray-400"></span>
				Gathering Fame
			</div>
			<div class="flex items-center">
				<span class="mr-2 h-1.5 w-1.5 rounded-full bg-gray-400"></span>
				Crafting Fame
			</div>
			<div class="flex items-center">
				<span class="mr-2 h-1.5 w-1.5 rounded-full bg-gray-400"></span>
				Guild History
			</div>
		</div>

		<div class="mt-3 font-medium">To get started, search for a player above.</div>
	</Typography>

	<div class="mt-8 grid grid-cols-1 gap-8 lg:grid-cols-2">
		<PlayersTracked />
		<TotalDataPoints />
	</div>

	<SubHeader title="About" classes="mt-8" />
	<Typography classes="leading-relaxed mt-2">
		We've been tracking players since January 11th, 2026.
		<br />
		When a player shows up on the killboards, we automatically start tracking their statistics.
		<br /><br />
		The main motivation for this project is to be able to look players up <i>quickly</i> and
		<i>easily</i>. Other projects like this exist and work, but they are often slow or have a clunky
		UI. This site aims to be fast and easy to use. All of the code for this project is open source
		and available on
		<a
			href="https://github.com/bricefrisco/AlbionStats"
			target="_blank"
			rel="noopener noreferrer"
			class="text-blue-600 underline transition-colors duration-200 hover:text-blue-800 hover:no-underline dark:text-blue-400 dark:hover:text-blue-300"
			>GitHub</a
		>.
	</Typography>
</Page>
