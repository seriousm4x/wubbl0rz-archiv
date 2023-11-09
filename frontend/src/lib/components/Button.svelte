<script lang="ts">
	export let href = '';
	export let color = 'from-pink-400 to-purple-400 dark:from-pink-800 dark:to-purple-800';
</script>

{#if href !== ''}
	<a
		{href}
		class="btn shadow rounded-xl bg-gradient-to-r hover:shadow-lg transition duration-200 {color}"
	>
		<slot />
	</a>
{:else}
	<button
		class="btn shadow rounded-xl bg-gradient-to-r hover:shadow-lg transition duration-200 {color}"
	>
		<slot />
	</button>
{/if}
