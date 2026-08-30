<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import ChatReplay from '$lib/components/ChatReplay.svelte';
	import Player from '$lib/components/Player.svelte';
	import SEO from '$lib/components/SEO.svelte';
	import { DefaultOpenGraph } from '$lib/types/opengraph';
	import { format, parseISO } from 'date-fns';
	import { de } from 'date-fns/locale';
	import type { RecordModel } from 'pocketbase';
	import type { MediaPlayerElement } from 'vidstack/elements';

	let { data } = $props();

	let og = $derived({
		...DefaultOpenGraph,
		title: data.vod?.title,
		image: `${PUBLIC_API_URL}/vods/${data.vod?.filename}/thumb-lg.webp`,
		updated_time: parseISO(data.vod?.date).toISOString()
	});

	let player: MediaPlayerElement = $state({} as MediaPlayerElement);
	let playerHeight = $state(0);
	let currentTime = $state(0);
	let vod = $derived(data.vod as RecordModel);
</script>

<SEO {og} />

<div
	class="-m-4 grid h-[calc(100dvh-4rem)] grid-rows-[auto_minmax(0,1fr)] overflow-hidden xl:grid-cols-[minmax(0,1fr)_22rem] xl:grid-rows-[minmax(0,1fr)]"
>
	<div class="group relative min-h-0 min-w-0 bg-black xl:h-full" bind:clientHeight={playerHeight}>
		<Player
			class="vod-player h-full w-full"
			bind:player
			bind:currentTime
			video={vod}
			isAudio={false}
		/>

		<div
			class="pointer-events-none absolute inset-x-0 top-0 bg-linear-to-b from-black/85 via-black/45 to-transparent px-4 pt-5 pb-16 text-white transition-opacity duration-200 xl:opacity-0 xl:group-hover:opacity-100"
		>
			<h1 class="max-w-4xl text-xl leading-tight font-semibold tracking-tight sm:text-2xl">
				{vod.title}
			</h1>
			<div class="mt-3 flex flex-wrap gap-x-2 gap-y-1 text-sm text-white/75">
				<span>
					{format(parseISO(vod.date), 'dd.MM.yyyy, HH:mm', { locale: de })} Uhr
				</span>
				|
				<span>{vod.viewcount.toLocaleString('de-DE')} Views</span> |
				<span>{vod.resolution}@{vod.fps} FPS</span>
			</div>
		</div>
	</div>

	<aside class="min-h-0 p-4 xl:h-full xl:p-0">
		<ChatReplay {vod} {currentTime} {playerHeight} />
	</aside>
</div>

<style>
	:global(.vod-player [data-media-provider] video) {
		width: 100%;
		height: 100%;
		object-fit: contain;
	}
</style>
