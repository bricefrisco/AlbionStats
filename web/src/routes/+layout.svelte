<script>
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { afterNavigate, beforeNavigate } from '$app/navigation';
	import { onMount } from 'svelte';
	import NProgress from 'nprogress';

	let { children } = $props();
	let startTimer;

	onMount(() => {
		NProgress.configure({
			minimum: 0.12,
			trickleSpeed: 120,
			showSpinner: false
		});
	});

	beforeNavigate(() => {
		clearTimeout(startTimer);
		startTimer = setTimeout(() => {
			NProgress.start();
		}, 120);
	});

	afterNavigate(() => {
		clearTimeout(startTimer);
		NProgress.done();
	});
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
{@render children()}
