<script>
	export let content = '';
	export let offset = 12;

	let open = false;
	let x = 0;
	let y = 0;

	function handleMove(event) {
		x = event.clientX;
		y = event.clientY;
	}

	function handleEnter(event) {
		open = true;
		handleMove(event);
	}

	function handleLeave() {
		open = false;
	}
</script>

<span
	class="relative inline-flex"
	on:mouseenter={handleEnter}
	on:mousemove={handleMove}
	on:mouseleave={handleLeave}
	on:focusin={handleEnter}
	on:focusout={handleLeave}
>
	<slot />
	{#if open}
		<div
			class="pointer-events-none fixed z-50 max-w-xs rounded bg-neutral-900 px-2 py-1 text-xs text-white shadow-lg"
			style={`left: ${x}px; top: ${y - offset}px; transform: translate(-50%, -100%);`}
			role="tooltip"
		>
			{content}
		</div>
	{/if}
</span>
