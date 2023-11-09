<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import { toHHMMSS } from '$lib/functions';
	import type { RecordModel } from 'pocketbase';
	import { fade } from 'svelte/transition';
	import { watchHistory } from '$lib/stores/localstorage';

	export let video: RecordModel = {} as RecordModel;
	export let offset = 0;
	export let isVod = false;

	const type = isVod || video.collectionName === 'vod' ? 'vods' : 'clips';
	let hover = false;
</script>

<a
	href="/{type}/{video.id}{offset > 0 ? `?t=${offset}` : ''}"
	class="aspect-video overflow-hidden"
	on:mouseenter={() => (hover = true)}
	on:mouseleave={() => (hover = false)}
>
	{#if hover}
		<video
			muted
			loop
			autoplay
			transition:fade={{ duration: 250 }}
			class="absolute top-0 left-0 w-full h-auto z-10"
		>
			<source
				src="{PUBLIC_API_URL}/{type}/{video.filename}-preview.webm"
				type="video/webm;codecs=vp9"
			/>
			<source
				src="{PUBLIC_API_URL}/{type}/{video.filename}-preview.mp4"
				type="video/mp4; codecs=hvc1"
			/>
			<track kind="captions" />
		</video>
	{/if}
	<div class="relative">
		<picture role="img">
			<source
				type="image/webp"
				srcset="{PUBLIC_API_URL}/{type}/{video.filename}-sm.webp"
				media="(max-width: 545px)"
				class="w-full h-auto"
			/>
			<source
				type="image/webp"
				srcset="{PUBLIC_API_URL}/{type}/{video.filename}-md.webp"
				media="(max-width: 767px)"
				class="w-full h-auto"
			/>
			<source
				type="image/webp"
				srcset="{PUBLIC_API_URL}/{type}/{video.filename}-sm.webp"
				media="(min-width: 768px)"
				class="w-full h-auto"
			/>
			<img
				width="768"
				height="432"
				src="{PUBLIC_API_URL}/{type}/{video.filename}-md.webp"
				class="w-full h-auto"
				alt={video.title}
				loading="lazy"
			/>
		</picture>
		<div class="absolute bottom-0 right-0 mx-2 my-3 px-1 rounded-md bg-base-300 font-bold">
			{toHHMMSS(video.duration, false)}
		</div>
		{#if type in $watchHistory && video.id in $watchHistory[type]}
			<progress
				class="progress progress-primary bg-base-100 rounded-none w-full bottom-0 absolute"
				value={$watchHistory[type][video.id]}
				max={video.duration}
			></progress>
		{/if}
	</div>
</a>

<style>
	::-moz-progress-bar,
	::-webkit-progress-bar {
		border-radius: 0px;
	}
</style>
