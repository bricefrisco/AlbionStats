<script>
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import BattleTable from '$components/BattleTable.svelte';
	import Typography from '$components/Typography.svelte';
	import AllianceSearchBar from '$components/AllianceSearchBar.svelte';
	import { regionState } from '$lib/regionState.svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	let searchQuery = $derived(page.url.searchParams.get('q') || '');

	function updateUrl(q) {
		searchQuery = q;
		const url = new URL(page.url);
		if (q) {
			url.searchParams.set('q', q);
		} else {
			url.searchParams.delete('q');
		}
		goto(resolve(url.pathname + url.search), {
			keepFocus: true,
			noScroll: true
		});
	}
</script>

<Page>
	<PageHeader title="Battle Boards" />
	<Typography>
		<p>Recent battles on the {regionState.label} server.</p>
	</Typography>

	<div class="mb-4">
		<AllianceSearchBar bind:value={searchQuery} onselect={updateUrl} />
	</div>

	<BattleTable />
</Page>