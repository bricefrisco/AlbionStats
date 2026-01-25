<script>
	import { resolve } from '$app/paths';
	import Page from '$components/Page.svelte';
	import PageHeader from '$components/PageHeader.svelte';
	import Typography from '$components/Typography.svelte';

	let { error, status } = $props();
	let statusCode = $derived(status ?? error?.status ?? error?.statusCode);
	let isNotFound = $derived(statusCode === 404);
</script>

<Page>
	<PageHeader title={isNotFound ? 'Page Not Found' : 'Something went wrong'} />
	<Typography>
		{#if isNotFound}
			<p>The page you are looking for does not exist.</p>
			<p class="mt-2">
				Go back to the
				<a href={resolve('/')} class="underline">home page</a>.
			</p>
		{:else}
			<p>
				{#if statusCode}
					{statusCode} - {error?.message || 'An unexpected error occurred.'}
				{:else}
					{error?.message || 'An unexpected error occurred.'}
				{/if}
			</p>
		{/if}
	</Typography>
</Page>
