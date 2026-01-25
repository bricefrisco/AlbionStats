<script>
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import BattleTable from '$components/BattleTable.svelte';
	import Typography from '$components/Typography.svelte';
	import Select from '$components/Select.svelte';
	import AllianceSearchBar from '$components/AllianceSearchBar.svelte';
	import BackToTopButton from '$components/BackToTopButton.svelte';
	import GuildSearchBar from '$components/GuildSearchBar.svelte';
	import PlayerSearchBar from '$components/PlayerSearchBar.svelte';
	import Tooltip from '$components/Tooltip.svelte';
	import HelpIcon from '$components/icons/HelpIcon.svelte';
	import { regionState } from '$lib/regionState.svelte.js';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { SvelteSet } from 'svelte/reactivity';

	let { data } = $props();

	let searchQuery = $state(data.q);
	let searchType = $state(data.type);
	let minPlayers = $state(data.p);
	let offset = $state(0);
	let hasMore = $state(data.initialHasMore);
	let loading = $state(false);
	let hasResults = $state(data.initialBattles.length > 0);
	let selectedIds = new SvelteSet();

	// "Snapshot" of search values to pass to BattleTable, updated only on submission or URL change
	let activeQ = $state(searchQuery);
	let activeType = $state(searchType);
	let activeP = $state(minPlayers);
	let hasFilters = $derived.by(() => {
		const params = page.url.searchParams;
		return Boolean(params.get('q') || params.get('p') || params.get('type'));
	});

	async function updateUrl(q = searchQuery, p = minPlayers) {
		const url = new URL(page.url);
		offset = 0;
		if (q) {
			url.searchParams.set('q', q);
		} else {
			url.searchParams.delete('q');
		}

		if (p || p === 0) {
			url.searchParams.set('p', p);
		} else {
			url.searchParams.delete('p');
		}

		url.searchParams.set('type', searchType);

		await goto(resolve(url.pathname + url.search), {
			keepFocus: true,
			noScroll: true
		});

		activeQ = q;
		activeType = searchType;
		activeP = p;
	}


	$effect(() => {
		const urlQ = page.url.searchParams.get('q') || '';
		const urlType = page.url.searchParams.get('type') || 'alliance';
		const urlP = page.url.searchParams.get('p') || '10';

		if (urlQ !== activeQ || urlType !== activeType || urlP !== activeP) {
			searchQuery = urlQ;
			searchType = urlType;
			minPlayers = urlP;
			offset = 0;

			activeQ = urlQ;
			activeType = urlType;
			activeP = urlP;
		}
	});
</script>

<svelte:head>
	<title>Battle Boards - AlbionStats - {regionState.label}</title>
	<meta
		name="description"
		content={`Recent Albion Online battles in ${regionState.label}. Filter by alliance, guild, or player and participants.`}
	/>
	<link rel="canonical" href={`${page.url.origin}${page.url.pathname}`} />
	{#if hasFilters}
		<meta name="robots" content="noindex,follow" />
	{/if}
</svelte:head>

<Page>
	<PageHeader title="Battle Boards" />
	<Typography>
		<h2>Albion Online Battle Boards. Filter by alliance, guild, and number of participants.</h2>
		<p>Below are the most recent battles occurring over the past 30 days.</p>
		<p>Collection began on January 19th, 2026.</p>
	</Typography>

	<form class="mb-4 flex items-end gap-2" onsubmit={(e) => { e.preventDefault(); updateUrl(searchQuery); }}>
		<div class="flex flex-col gap-1">
			<div class="flex w-full items-center justify-between">
				<label for="min-players" class="text-xs font-medium text-gray-600 dark:text-gray-400">
					Participants
				</label>
				<Tooltip content="If filtering by alliance or guild, the participant count is for that alliance or guild specifically. Otherwise, it is the total amount of players in the battle">
					<button
						type="button"
						class="flex items-center text-gray-400 transition-colors hover:text-gray-600 dark:text-neutral-500 dark:hover:text-neutral-300"
						aria-label="Participants info"
					>
						<HelpIcon size={14} />
					</button>
				</Tooltip>
			</div>
			<input
				id="min-players"
				type="number"
				bind:value={minPlayers}
				min="1"
				max="300"
				onkeydown={(e) => {
					if (e.key !== 'Enter') return;
					e.preventDefault();
					updateUrl(searchQuery);
				}}
				class="w-32 rounded border border-gray-300 bg-gray-50 px-3 py-2 text-sm text-gray-900 placeholder-gray-500 focus:border-gray-400 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 dark:border-neutral-800 dark:bg-neutral-900 dark:text-neutral-100 dark:placeholder-neutral-500 dark:focus:border-neutral-700"
			/>
		</div>

		<Select
			bind:value={searchType}
			classes="w-32"
			label="Type"
			options={[
				{ label: 'Alliance', value: 'alliance' },
				{ label: 'Guild', value: 'guild' },
				{ label: 'Player', value: 'player' }
			]}
		/>

		<div class="flex-1">
			{#if searchType === 'player'}
				<PlayerSearchBar
					bind:value={searchQuery}
					onselect={(p) => searchQuery = p.name}
					links={false}
					label="Name"
				/>
			{:else if searchType === 'guild'}
				<GuildSearchBar
					bind:value={searchQuery}
					onselect={(g) => searchQuery = g}
					links={false}
					label="Name"
				/>
			{:else}
				<AllianceSearchBar
					bind:value={searchQuery}
					onselect={(a) => searchQuery = a}
					links={false}
					label="Name"
				/>
			{/if}
		</div>

		<button
			type="submit"
			class="h-[38px] cursor-pointer rounded border border-gray-300 bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-200 focus:border-gray-400 focus:outline-none dark:border-neutral-800 dark:bg-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-700 dark:focus:border-neutral-700"
		>
			Apply filter
		</button>

		{#if selectedIds.size > 0}
			<a
				href={resolve(`/battle-boards/${regionState.value}/${Array.from(selectedIds).join(',')}`)}
				class="inline-flex h-[38px] items-center cursor-pointer rounded border border-gray-300 bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-200 focus:border-gray-400 focus:outline-none dark:border-neutral-800 dark:bg-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-700 dark:focus:border-neutral-700"
			>
				View battles
			</a>
		{/if}
	</form>

	<BattleTable
		q={activeQ}
		type={activeType}
		p={activeP}
		{offset}
		initialBattles={data.initialBattles}
		initialHasMore={data.initialHasMore}
		initialError={data.initialError}
		bind:hasMore
		bind:loading
		bind:selectedIds
		bind:hasResults
	/>

	<div class="mt-8 flex justify-center gap-4">
		{#if hasMore && !loading && searchQuery === activeQ && searchType === activeType && minPlayers === activeP}
			<button
				type="button"
				onclick={() => offset += 20}
				class="cursor-pointer rounded border border-gray-300 bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-200 focus:border-gray-400 focus:outline-none dark:border-neutral-800 dark:bg-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-700 dark:focus:border-neutral-700"
			>
				Load more
			</button>
		{/if}

		{#if hasResults}
			<BackToTopButton />
		{/if}
	</div>
</Page>
