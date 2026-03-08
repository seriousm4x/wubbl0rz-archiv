<script lang="ts">
	import { resolve } from '$app/paths';
	import { PUBLIC_API_URL } from '$env/static/public';
	import { toHHMMSS } from '$lib/functions';
	import type { RecordModel } from 'pocketbase';
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';

	let {
		video = {} as RecordModel,
		offset = 0,
		isVod = false
	}: { video: RecordModel; offset: number; isVod?: boolean } = $props();

	const type = isVod || video.collectionName === 'vod' ? 'vods' : 'clips';
	let hover = $state(false);
	let progress = $state(0);

	onMount(() => {
		const filename = isVod || video.collectionName === 'vod' ? 'vod.mp4' : 'clip.mp4';
		const localStorageItem = localStorage.getItem(
			`${PUBLIC_API_URL}/${type}/${video.filename}/${filename}:0:0`
		);
		progress = localStorageItem ? parseFloat(localStorageItem) : 0;
	});
</script>

<a
	href={resolve(`/${type}/${video.id}${offset > 0 ? `?t=${offset}` : ''}`)}
	class="aspect-video overflow-hidden"
	onmouseenter={() => (hover = true)}
	onmouseleave={() => (hover = false)}
>
	{#if hover}
		<video
			muted
			loop
			autoplay
			transition:fade={{ duration: 250 }}
			class="absolute top-0 left-0 z-10 h-auto w-full"
		>
			<source
				src="{PUBLIC_API_URL}/{type}/{video.filename}/preview.webm"
				type="video/webm;codecs=vp9"
			/>
			<source
				src="{PUBLIC_API_URL}/{type}/{video.filename}/preview.mp4"
				type="video/mp4; codecs=hvc1"
			/>
			<track kind="captions" />
		</video>
	{/if}
	<div class="relative">
		<picture>
			<source
				type="image/webp"
				srcset="{PUBLIC_API_URL}/{type}/{video.filename}/thumb-sm.webp"
				media="(max-width: 545px)"
				class="h-auto w-full"
				width="512"
				height="286"
			/>
			<source
				type="image/webp"
				srcset="{PUBLIC_API_URL}/{type}/{video.filename}/thumb-md.webp"
				media="(max-width: 767px)"
				class="h-auto w-full"
				width="768"
				height="430"
			/>
			<source
				type="image/webp"
				srcset="{PUBLIC_API_URL}/{type}/{video.filename}/thumb-sm.webp"
				media="(min-width: 768px)"
				class="h-auto w-full"
				width="512"
				height="286"
			/>
			<img
				src="{PUBLIC_API_URL}/{type}/{video.filename}/thumb-md.webp"
				class="h-auto w-full"
				alt={video.title}
				loading="lazy"
				width="768"
				height="430"
			/>
		</picture>
		<div class="bg-base-300 absolute right-0 bottom-0 mx-2 my-3 rounded-md px-1 font-bold">
			{toHHMMSS(video.duration, false)}
		</div>
		{#if progress > 0}
			<progress
				class="progress progress-primary bg-base-100 absolute bottom-0 w-full rounded-none"
				value={progress}
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
