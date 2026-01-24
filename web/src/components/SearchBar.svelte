<script>
	let {
		value = $bindable(),
		placeholder = '',
		label = '',
		oninput,
		onfocus,
		onkeydown,
		itemCount = 0,
		activeIndex = $bindable(0),
		onselectIndex,
		allowSubmitOnEnter = false,
		isSearching = false,
		showDropdown = $bindable(false),
		dropdownContent
	} = $props();

	import SearchIcon from './icons/SearchIcon.svelte';
	import { onMount } from 'svelte';

	let container = $state();

	const handleDocumentClick = (event) => {
		if (container && !container.contains(event.target)) {
			showDropdown = false;
		}
	};

	const handleKeydown = (event) => {
		onkeydown?.(event);
		if (event.defaultPrevented) return;
		if (event.key === 'Tab') {
			if (!itemCount) return;
			event.preventDefault();
			activeIndex = (activeIndex + 1) % itemCount;
			showDropdown = true;
			return;
		}

		if (event.key === 'Enter') {
			if (!itemCount) {
				event.preventDefault();
				return;
			}
			if (!allowSubmitOnEnter) {
				event.preventDefault();
			}
			onselectIndex?.(activeIndex);
			showDropdown = false;
		}
	};

	onMount(() => {
		document.addEventListener('click', handleDocumentClick);
		return () => document.removeEventListener('click', handleDocumentClick);
	});
</script>

<div class="flex flex-col gap-1" bind:this={container}>
	{#if label}
		<label class="text-xs font-medium text-gray-600 dark:text-gray-400">
			{label}
		</label>
	{/if}
	<div class="relative">
		<div class="relative flex-1">
			<input
				type="text"
				bind:value={value}
				{oninput}
				{onfocus}
				onkeydown={handleKeydown}
				{placeholder}
				class="w-full rounded border border-gray-300 bg-gray-50 px-3 py-2 pl-9 text-sm text-gray-900 placeholder-gray-500 focus:border-gray-400 focus:outline-none dark:border-neutral-800 dark:bg-neutral-900 dark:text-neutral-100 dark:placeholder-neutral-500 dark:focus:border-neutral-700"
			/>
			<div class="absolute top-[11px] left-3 flex items-center text-gray-500 dark:text-neutral-500">
				<SearchIcon size={16} />
			</div>
		</div>

		{#if showDropdown}
			<div
				class="absolute top-full right-0 left-0 z-10 mt-2 max-h-60 overflow-y-auto rounded border border-gray-200 bg-white shadow-lg dark:border-neutral-700 dark:bg-neutral-900"
			>
				{#if isSearching}
					<div class="px-3 py-2 text-sm text-gray-500 dark:text-neutral-400">Searching...</div>
				{:else}
					{@render dropdownContent?.()}
				{/if}
			</div>
		{/if}
	</div>
</div>
