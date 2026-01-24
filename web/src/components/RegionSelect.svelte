<script>
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import Select from './Select.svelte';
	import { regionRoutes } from '$lib/regionRoutes';
	import { regionState } from '$lib/regionState.svelte';

	const regions = [
		{ label: 'Americas', value: 'americas' },
		{ label: 'Europe', value: 'europe' },
		{ label: 'Asia', value: 'asia' }
	];

	const regionRouteSet = new Set(regionRoutes);
	let lastRegion = $state(regionState.value);

	$effect(() => {
		if (!browser) return;
		const nextRegion = regionState.value;
		if (nextRegion === lastRegion) return;
		lastRegion = nextRegion;

		const segments = page.url.pathname.split('/').filter(Boolean);
		if (!segments.length || !regionRouteSet.has(segments[0])) return;

		if (segments.length === 1) {
			segments.push(nextRegion);
		} else {
			segments[1] = nextRegion;
		}

		const nextPath = `/${segments.join('/')}`;
		const nextUrl = `${nextPath}${page.url.search}${page.url.hash}`;
		const currentUrl = `${page.url.pathname}${page.url.search}${page.url.hash}`;
		if (nextUrl === currentUrl) return;

		goto(nextUrl, { keepFocus: true, noScroll: true, replaceState: true });
	});
</script>

<Select bind:value={regionState.value} options={regions} />
