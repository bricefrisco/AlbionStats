<script>
	import ChevronDownIcon from './icons/ChevronDownIcon.svelte';

	let { value = $bindable(), options = [], classes = '', ...restProps } = $props();

	function handleChange(event) {
		value = event.target.value;
	}
</script>

<div class="relative {classes}">
	<select
		{value}
		onchange={handleChange}
		class="cursor-pointer appearance-none rounded border border-gray-300 bg-gray-50 px-3 py-2 pr-8 text-sm text-gray-900 focus:border-gray-400 focus:outline-none dark:border-neutral-800 dark:bg-neutral-900 dark:text-neutral-100 dark:focus:border-neutral-700"
		{...restProps}
	>
		{#each options as option (option.value)}
			<option value={option.value}>
				{option.label}
			</option>
		{/each}
	</select>
	<div
		class="pointer-events-none absolute top-1/2 right-2 flex -translate-y-1/2 items-center text-gray-500 dark:text-neutral-500"
	>
		<ChevronDownIcon size={16} />
	</div>
</div>
