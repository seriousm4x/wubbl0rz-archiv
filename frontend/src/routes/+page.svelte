<script lang="ts">
	import Card from '$lib/components/Card.svelte';
	import SEO from '$lib/components/SEO.svelte';
	import { DefaultOpenGraph } from '$lib/types/opengraph';
	import { parseISO } from 'date-fns';
	import type { PageData } from './$types.js';

	let { data }: { data: PageData } = $props();

	let og = $derived({
		...DefaultOpenGraph,
		updated_time: parseISO(data.new?.items?.[0]?.date).toISOString()
	});
</script>

<SEO {og} />

<section class="space-y-4">
	<h1 class="text-2xl font-semibold tracking-tight">Neue Streams</h1>
	<div
		class="grid grid-cols-1 gap-x-3 gap-y-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 2xl:grid-cols-5"
	>
		{#each data.new?.items as video (video.id)}
			<Card {video} />
		{/each}
	</div>
</section>

<section class="mt-10 space-y-4">
	<h2 class="text-2xl font-semibold tracking-tight">Beliebte Streams</h2>
	<div
		class="grid grid-cols-1 gap-x-3 gap-y-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 2xl:grid-cols-5"
	>
		{#each data.popular?.items as video (video.id)}
			<Card {video} />
		{/each}
	</div>
</section>

<section class="mt-10 space-y-4">
	<h2 class="text-2xl font-semibold tracking-tight">Top Clips des Monats</h2>
	<div
		class="grid grid-cols-1 gap-x-3 gap-y-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 2xl:grid-cols-5"
	>
		{#each data.clips?.items as video (video.id)}
			<Card {video} />
		{:else}
			<p class="text-base-content/65 col-span-full py-8 text-sm">Keine Clips im letzten Monat.</p>
		{/each}
	</div>
</section>
