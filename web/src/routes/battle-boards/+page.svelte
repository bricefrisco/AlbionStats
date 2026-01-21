<script>
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import BattleTable from '$components/BattleTable.svelte';
	import Typography from '$components/Typography.svelte';
	import { regionState } from '$lib/regionState.svelte';
</script>

<Page>
	<PageHeader title="Battle Boards" />
	<Typography>
		<p>Recent battles on the {regionState.label} server.</p>
	</Typography>
	<BattleTable />
</Page>