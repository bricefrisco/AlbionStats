<script>
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import PlayerSearchBar from '$components/PlayerSearchBar.svelte';
	import Typography from '$components/Typography.svelte';
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
	<PageHeader title="Players" />
	<Typography>
		<p>Search for a player to view their stats.</p>
	</Typography>

	<div class="mt-8">
		<PlayerSearchBar bind:value={searchQuery} onselect={updateUrl} />
	</div>
</Page>
