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

<Page>
	<PageHeader title="Battle Boards" />
	<Typography>
		<h2>Recent battles on the {regionState.label} server.</h2>
	</Typography>

	<form class="mb-4 flex items-end gap-2" onsubmit={(e) => { e.preventDefault(); updateUrl(searchQuery); }}>
		<div class="flex flex-col gap-1">
			<label for="min-players" class="text-xs font-medium text-gray-600 dark:text-gray-400">
				Participants
			</label>
			<input
				id="min-players"
				type="number"
				bind:value={minPlayers}
				onkeydown={(e) => {
					if (e.key !== 'Enter') return;
					e.preventDefault();
					updateUrl(searchQuery);
				}}
				class="w-24 rounded border border-gray-300 bg-gray-50 px-3 py-2 text-sm text-gray-900 placeholder-gray-500 focus:border-gray-400 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 dark:border-neutral-800 dark:bg-neutral-900 dark:text-neutral-100 dark:placeholder-neutral-500 dark:focus:border-neutral-700"
			/>
		</div>

		<Select
			bind:value={searchType}
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
