<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import { toHHMMSS } from '$lib/functions';
	import { watchHistory } from '$lib/stores/localstorage';
	import type { RecordModel } from 'pocketbase';
	import { fade } from 'svelte/transition';

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
			class="absolute left-0 top-0 z-10 h-auto w-full"
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
		<picture>
			<source
				type="image/webp"
				srcset="{PUBLIC_API_URL}/{type}/{video.filename}-sm.webp"
				media="(max-width: 545px)"
				class="h-auto w-full"
				width="512"
				height="286"
			/>
			<source
				type="image/webp"
				srcset="{PUBLIC_API_URL}/{type}/{video.filename}-md.webp"
				media="(max-width: 767px)"
				class="h-auto w-full"
				width="768"
				height="430"
			/>
			<source
				type="image/webp"
				srcset="{PUBLIC_API_URL}/{type}/{video.filename}-sm.webp"
				media="(min-width: 768px)"
				class="h-auto w-full"
				width="512"
				height="286"
			/>
			<img
				src="{PUBLIC_API_URL}/{type}/{video.filename}-md.webp"
				class="h-auto w-full"
				alt={video.title}
				loading="lazy"
				width="768"
				height="430"
			/>
		</picture>
		<div class="bg-base-300 absolute bottom-0 right-0 mx-2 my-3 rounded-md px-1 font-bold">
			{toHHMMSS(video.duration, false)}
		</div>
		{#if type in $watchHistory && video.id in $watchHistory[type]}
			<progress
				class="progress progress-primary bg-base-100 absolute bottom-0 w-full rounded-none"
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
