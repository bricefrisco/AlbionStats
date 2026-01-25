<script>
	export let value = '';
	export let placeholder = 'Filter...';
</script>

<div class="mb-4">
	<input
		type="text"
		bind:value
		placeholder={placeholder}
		class="w-full rounded border border-gray-300 bg-gray-50 px-3 py-2 text-sm text-gray-900 placeholder-gray-500 focus:border-gray-400 focus:outline-none dark:border-neutral-800 dark:bg-neutral-900 dark:text-neutral-100 dark:placeholder-neutral-500 dark:focus:border-neutral-700"
	/>
</div>
