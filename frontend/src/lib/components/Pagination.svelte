<script lang="ts">
	let {
		totalPages,
		currentPage,
		onPageChange
	}: { totalPages: number; currentPage: number; onPageChange: (page: number) => void } = $props();

	function range(start: number, end: number) {
		let a = Array(end - start + 1);
		let current = start;
		let index = 0;
		while (current <= end) a[index++] = current++;
		return a;
	}
</script>

<div class="join flex-wrap gap-y-1">
	{#if currentPage > 1}
		<button class="btn join-item" onclick={() => onPageChange(1)}>Erste</button>
	{/if}
	{#each range(1, totalPages) as i (i)}
		{#if i <= currentPage + 3 && i >= currentPage - 3}
			{#if currentPage == i}
				<button class="btn join-item btn-active" onclick={() => onPageChange(i)}>{i}</button>
			{:else}
				<button class="btn join-item" onclick={() => onPageChange(i)}>{i}</button>
			{/if}
		{/if}
	{/each}
	{#if currentPage < totalPages}
		<button class="btn join-item" onclick={() => onPageChange(totalPages)}>Letzte</button>
	{/if}
</div>
