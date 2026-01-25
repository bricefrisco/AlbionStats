<script>
	import { resolve } from '$app/paths';
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import Typography from '$components/Typography.svelte';
	import PlayersTracked from '$components/charts/PlayersTracked.svelte';
	import TotalDataPoints from '$components/charts/TotalDataPoints.svelte';
	import DailyActiveUsers from '$components/charts/DailyActiveUsers.svelte';
	import SubHeader from '$components/SubHeader.svelte';
	import { page } from '$app/state';

	let { data } = $props();
	const overview = {
		intro: 'AlbionStats helps you explore Albion Online statistics over time.',
		gettingStarted: [
			{
				path: '/battle-boards',
				label: 'Battle Boards',
				description: 'See recent battles and filter by alliance, guild, or player'
			},
			{
				path: '/alliances',
				label: 'Alliances',
				description: 'Compare alliance performance and rosters'
			},
			{
				path: '/guilds',
				label: 'Guilds',
				description: 'Compare guild performance and rosters'
			},
			{
				path: '/players',
				label: 'Players',
				description: 'Look up a player and view their history'
			},
		],
		highlights: [],
		about: [
			'We have been tracking players since January 11th, 2026 and battle boards since January 19th, 2026.',
			'When a player appears in combat activity, we start tracking them automatically. We use battle ' +
			'board data to derive alliance and guild statistics.',
			'The goal is a fast, clear way to explore Albion Online progress and rivalries.'
		]
	};
	let websiteJsonLd = $derived.by(() => {
		const origin = page.url.origin;
		return JSON.stringify({
			'@context': 'https://schema.org',
			'@type': 'WebSite',
			name: 'AlbionStats',
			url: origin
		});
	});
</script>

<svelte:head>
	<title>Home - AlbionStats</title>
	<meta
		name="description"
		content="Track Albion Online player, guild, and alliance statistics over time. Browse battle boards, charts, and activity trends."
	/>
	<link rel="canonical" href={`${page.url.origin}${page.url.pathname}`} />
	<script type="application/ld+json">{websiteJsonLd}</script>
	</svelte:head>

<Page>
	<PageHeader title="Welcome to AlbionStats" />
	<Typography>
		<p>{overview.intro}</p>
		<div class="mt-6">
			<p class="font-medium text-gray-900 dark:text-white">To get started:</p>
			<div class="mt-3 grid grid-cols-1 gap-3 text-sm text-gray-600 dark:text-gray-400 lg:grid-cols-2">
				{#each overview.gettingStarted as item}
					<div
						class="rounded-lg border border-gray-200 bg-gray-50 p-3 dark:border-neutral-800 dark:bg-neutral-900"
					>
						<a
							href={resolve(item.path)}
							class="font-medium text-blue-600 hover:underline dark:text-blue-400"
							>{item.label}</a
						>
						<p class="mt-1 text-xs text-gray-600 dark:text-gray-400">{item.description}.</p>
					</div>
				{/each}
			</div>
		</div>
		<div class="mt-6 grid grid-cols-1 gap-4 text-sm lg:grid-cols-2">
			{#each overview.highlights as highlight}
				<div
					class="rounded-lg border border-gray-200 bg-gray-50 p-3 text-gray-700 dark:border-neutral-800 dark:bg-neutral-900 dark:text-gray-200"
				>
					<span>{highlight}</span>
				</div>
			{/each}
		</div>
	</Typography>

	<div class="mt-8 grid grid-cols-1 gap-8 lg:grid-cols-2">
		<PlayersTracked data={data.playersTracked} />
		<TotalDataPoints data={data.totalDataPoints} />
	</div>

	<div class="mt-8">
		<DailyActiveUsers data={data.dailyActiveUsers} />
	</div>

	<SubHeader title="About" classes="mt-8" />
	<Typography classes="leading-relaxed mt-2">
		{overview.about[0]}
		<br />
		{overview.about[1]}
		<br /><br />
		{overview.about[2]}
	</Typography>
</Page>
