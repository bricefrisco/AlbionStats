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
	import { untrack } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	let searchQuery = $derived(page.url.searchParams.get('q') || '');
	let searchType = $state(page.url.searchParams.get('type') || 'alliance');

	function updateUrl(q) {
		const url = new URL(page.url);
		if (q) {
			url.searchParams.set('q', q);
		} else {
			url.searchParams.delete('q');
		}

		if (searchType !== 'alliance') {
			url.searchParams.set('type', searchType);
		} else {
			url.searchParams.delete('type');
		}

		goto(resolve(url.pathname + url.search), {
			keepFocus: true,
			noScroll: true
		});
	}

	$effect(() => {
		// Update URL when search type changes without a query change
		const url = new URL(page.url);
		const currentType = url.searchParams.get('type') || 'alliance';

		if (searchType !== currentType) {
			untrack(() => updateUrl(searchQuery));
		}
	});
</script>

<Page>
	<PageHeader title="Battle Boards" />
	<Typography>
		<p>Recent battles on the {regionState.label} server.</p>
	</Typography>

	<div class="mb-4 flex items-center gap-2">
		<Select
			bind:value={searchType}
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
					onselect={(p) => updateUrl(p.name)}
					links={false}
				/>
			{:else}
				<AllianceSearchBar
					bind:value={searchQuery}
					onselect={updateUrl}
				/>
			{/if}
		</div>
	</div>

	<BattleTable />
</Page>