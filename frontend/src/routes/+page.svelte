<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import Card from '$lib/components/Card.svelte';
	import Hero from '$lib/components/Hero.svelte';
	import SEO from '$lib/components/SEO.svelte';
	import { DefaultOpenGraph } from '$lib/types/opengraph';
	import { parseISO } from 'date-fns';

	export let data;
	let innerWidth: number;

	$: og = {
		...DefaultOpenGraph,
		updated_time: parseISO(data.new?.items?.[0]?.date).toISOString()
	};

	$: showHero = innerWidth > 767;
</script>

<svelte:window bind:innerWidth />

<SEO bind:og />

{#if showHero && data.new?.items?.length > 0}
	<Hero vod={data.new.items[0]} />
{/if}

<h1 class="mb-4 text-4xl font-bold md:ms-3 md:mt-10">
	<span
		class="bg-gradient-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent drop-shadow-md"
		>Neue Streams</span
	>
</h1>
<div
	class="grid grid-flow-row-dense grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5"
>
	<div
		class="absolute left-0 top-0 -z-10 aspect-video h-auto w-full bg-cover bg-center opacity-40 blur-2xl"
	>
		<img
			src={data.new?.items?.[0]?.filename
				? `${PUBLIC_API_URL}/vods/${data.new?.items?.[0]?.filename}-lg.webp`
				: ''}
			alt={data.new?.items?.[0]?.title}
			class="h-auto w-full"
			width="1536"
			height="860"
		/>
	</div>
	{#if showHero && data.new?.items?.length > 0}
		{#each data.new?.items?.filter((vod) => vod !== data.new?.items?.[0]) as video}
			<Card {video} />
		{/each}
	{:else}
		{#each data.new?.items as video}
			<Card {video} />
		{/each}
	{/if}
</div>

<h1 class="mb-4 ms-3 mt-10 text-4xl font-bold">
	<span
		class="bg-gradient-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent drop-shadow-md"
		>Beliebte Streams</span
	>
</h1>
<div
	class="grid grid-flow-row-dense grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5"
>
	{#each data.popular?.items as video}
		<Card {video} />
	{/each}
</div>

<h1 class="mb-4 ms-3 mt-10 text-4xl font-bold">
	<span
		class="bg-gradient-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent drop-shadow-md"
		>Top Clips des Monats</span
	>
</h1>
<div
	class="grid grid-flow-row-dense grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5"
>
	{#each data.clips?.items as video}
		<Card {video} />
	{/each}
</div>
