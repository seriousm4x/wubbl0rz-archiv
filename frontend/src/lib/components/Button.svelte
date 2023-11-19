<script lang="ts">
	export let href = '';
	export let color = 'from-pink-400 to-purple-400 dark:from-pink-800 dark:to-purple-800';
</script>

{#if href !== ''}
	<a
		{href}
		class="btn rounded-xl bg-gradient-to-r shadow transition duration-200 hover:shadow-lg {color}"
	>
		<slot />
	</a>
{:else}
	<button
		class="btn rounded-xl bg-gradient-to-r shadow transition duration-200 hover:shadow-lg {color}"
	>
		<slot />
	</button>
{/if}
