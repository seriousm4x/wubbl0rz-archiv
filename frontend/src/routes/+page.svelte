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

<h1 class="text-4xl font-bold md:mt-10 md:ms-3 mb-4">
	<span
		class="bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-violet-500 drop-shadow-md"
		>Neue Streams</span
	>
</h1>
<div
	class="grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4"
>
	<div
		class="absolute top-0 left-0 w-full h-auto aspect-video bg-cover bg-center blur-2xl opacity-40 -z-10"
	>
		<img
			src={data.new?.items?.[0]?.filename
				? `${PUBLIC_API_URL}/vods/${data.new?.items?.[0]?.filename}-lg.webp`
				: ''}
			alt={data.new?.items?.[0]?.title}
			class="w-full h-auto"
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

<h1 class="text-4xl font-bold mt-10 mb-4 ms-3">
	<span
		class="bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-violet-500 drop-shadow-md"
		>Beliebte Streams</span
	>
</h1>
<div
	class="grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4"
>
	{#each data.popular?.items as video}
		<Card {video} />
	{/each}
</div>

<h1 class="text-4xl font-bold mt-10 mb-4 ms-3">
	<span
		class="bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-violet-500 drop-shadow-md"
		>Top Clips des Monats</span
	>
</h1>
<div
	class="grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4"
>
	{#each data.clips?.items as video}
		<Card {video} />
	{/each}
</div>
