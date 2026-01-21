<script>
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import PlayerSearchBar from '$components/PlayerSearchBar.svelte';
	import Typography from '$components/Typography.svelte';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { regionState } from '$lib/regionState.svelte';

	function handlePlayerSelect(player) {
		goto(resolve(`/players/${regionState.value}/${encodeURIComponent(player.name)}`));
	}
</script>

<Page>
	<PageHeader title="Players" />
	<Typography>
		<p>Search for a player to view their stats.</p>
	</Typography>

	<div class="mt-8 max-w-xl">
		<PlayerSearchBar onSelect={handlePlayerSelect} />
	</div>
</Page>
