<script>
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import BattleTable from '$components/BattleTable.svelte';
	import Typography from '$components/Typography.svelte';
	import Select from '$components/Select.svelte';
	import AllianceSearchBar from '$components/AllianceSearchBar.svelte';
	import PlayerSearchBar from '$components/PlayerSearchBar.svelte';
	import { regionState } from '$lib/regionState.svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	let searchQuery = $state(page.url.searchParams.get('q') || '');
	let searchType = $state(page.url.searchParams.get('type') || 'alliance');
	let minPlayers = $state(page.url.searchParams.get('p') || '10');

	// "Snapshot" of search values to pass to BattleTable, updated only on submission or URL change
	let activeQ = $state(searchQuery);
	let activeType = $state(searchType);
	let activeP = $state(minPlayers);

	async function updateUrl(q = searchQuery, p = minPlayers) {
		const url = new URL(page.url);
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

			activeQ = urlQ;
			activeType = urlType;
			activeP = urlP;
		}
	});
</script>

<Page>
	<PageHeader title="Battle Boards" />
	<Typography>
		<p>Recent battles on the {regionState.label} server.</p>
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
				class="w-24 rounded border border-gray-300 bg-gray-50 px-3 py-2 text-sm text-gray-900 placeholder-gray-500 focus:border-gray-400 focus:outline-none dark:border-neutral-800 dark:bg-neutral-900 dark:text-neutral-100 dark:placeholder-neutral-500 dark:focus:border-neutral-700"
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
			{:else}
				<AllianceSearchBar
					bind:value={searchQuery}
					onselect={(a) => searchQuery = a}
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
	</form>

	<BattleTable q={activeQ} type={activeType} p={activeP} />
</Page>