<script>
	import { untrack } from 'svelte';
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import Typography from '$components/Typography.svelte';
	import { page } from '$app/state';

	let region = $derived(page.params.region);
	let battleIds = $derived(page.params.battleIds);

	let data = $state(null);
	let loading = $state(true);
	let error = $state(null);

	async function fetchBattles() {
		loading = true;
		error = null;
		try {
			const response = await fetch(
				`https://albionstats.bricefrisco.com/api/battles/${region}/${battleIds}`
			);
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			data = await response.json();
		} catch (err) {
			error = err.message;
			console.error('Failed to fetch battle data:', err);
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (region && battleIds) {
			untrack(() => fetchBattles());
		}
	});
</script>

<Page>
	<PageHeader title="Battle Boards" />
	<Typography>
		<p>Battle board results for the {region} server: {battleIds}</p>
		{#if loading}
			<p>Loading battle data...</p>
		{:else if error}
			<p class="text-red-600">{error}</p>
		{:else if data}
			<p>Successfully loaded battle data.</p>
		{/if}
	</Typography>
</Page>
