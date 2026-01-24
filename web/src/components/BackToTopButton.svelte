<script>
	let { class: className = '' } = $props();
</script>

<button
	type="button"
	onclick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
	class="cursor-pointer rounded border border-gray-300 bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-200 focus:border-gray-400 focus:outline-none dark:border-neutral-800 dark:bg-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-700 dark:focus:border-neutral-700 {className}"
>
	Back to top
</button>
