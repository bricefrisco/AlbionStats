<script>
	// Props
	let { tabs = [], activeTab = $bindable(null), ontabChange } = $props();

	// Set default active tab if none provided
	$effect(() => {
		if (!activeTab && tabs.length > 0) {
			activeTab = tabs[0].id;
		}
	});

	function selectTab(tabId) {
		activeTab = tabId;
		if (ontabChange) {
			ontabChange({ tabId });
		}
	}
</script>

<div class="border-b border-gray-200 dark:border-gray-700">
	<nav class="flex space-x-8">
		{#each tabs as tab (tab.id)}
			<button
				class="border-b-2 px-1 py-2 text-sm font-medium transition-colors"
				class:border-gray-600={activeTab === tab.id}
				class:text-gray-900={activeTab === tab.id}
				class:dark:text-white={activeTab === tab.id}
				class:border-transparent={activeTab !== tab.id}
				class:text-gray-500={activeTab !== tab.id}
				class:hover:text-gray-700={activeTab !== tab.id}
				class:dark:text-gray-400={activeTab !== tab.id}
				class:dark:hover:text-gray-300={activeTab !== tab.id}
				on:click={() => selectTab(tab.id)}
			>
				{tab.label}
			</button>
		{/each}
	</nav>
</div>
