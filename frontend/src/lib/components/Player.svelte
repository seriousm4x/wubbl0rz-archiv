<script lang="ts">
	import { page } from '$app/stores';
	import { PUBLIC_API_URL } from '$env/static/public';
	import { watchHistory, type WatchHistory } from '$lib/stores/localstorage';
	import Hls from 'hls.js';
	import type { RecordModel } from 'pocketbase';
	import { onMount } from 'svelte';
	import type { HLSProvider, MediaProviderChangeEvent, MediaTimeUpdateEvent } from 'vidstack';
	import type { MediaPlayerElement } from 'vidstack/elements';

	let {
		video = {} as RecordModel,
		player = $bindable(),
		currentTime = $bindable()
	}: { video: RecordModel; player: MediaPlayerElement; currentTime: number } = $props();

	onMount(async () => {
		await import('vidstack/player/styles/default/theme.css');
		await import('vidstack/player/styles/default/layouts/video.css');
		await import('vidstack/player');
		await import('vidstack/player/layouts');
		await import('vidstack/player/ui');

		if (!$watchHistory[type as keyof WatchHistory]) {
			$watchHistory[type as keyof WatchHistory] = {};
		}
	});

	let thumbnails: string = $state('');
	let type = video.collectionName === 'vod' ? 'vods' : 'clips';

	function onCanPlay() {
		thumbnails = `${PUBLIC_API_URL}/${type}/${video.filename}-sprites/${video.filename}.vtt`;

		// set player time to url param
		if ($page.url.searchParams.has('t')) {
			player.currentTime = parseInt($page.url.searchParams.get('t') || '0');
			return;
		}

		// set player time to localstorage history
		if (type in $watchHistory) {
			const lsTime = $watchHistory[type as keyof WatchHistory][video.id];
			if (lsTime !== undefined) {
				player.currentTime = lsTime;
			}
		}
	}

	function onProviderChange(event: MediaProviderChangeEvent) {
		const provider = event.detail;
		if (provider && provider.type === 'hls') {
			(provider as HLSProvider).library = Hls;
		}
	}

	function onTimeUpdate(event: MediaTimeUpdateEvent) {
		currentTime = event.detail.currentTime;
		$watchHistory[type as keyof WatchHistory][video.id] = currentTime;
	}
</script>

<div class="aspect-video overflow-hidden rounded-xl">
	<media-player
		class="player h-full w-full"
		title={video.title}
		src="{PUBLIC_API_URL}/{type}/{video.filename}-segments/{video.filename}.m3u8"
		crossorigin
		bind:this={player}
		oncan-play={onCanPlay}
		onprovider-change={onProviderChange}
		ontime-update={onTimeUpdate}
	>
		<media-provider>
			<media-poster class="vds-poster" src="{PUBLIC_API_URL}/{type}/{video.filename}-lg.webp"
			></media-poster>
			{#if type === 'vods'}
				<track
					label="Deutsch"
					src="{PUBLIC_API_URL}/{type}/{video.filename}.vtt"
					kind="subtitles"
					srclang="de"
				/>
			{/if}
		</media-provider>
		<media-audio-layout></media-audio-layout>
		<media-video-layout {thumbnails}></media-video-layout>
		<media-buffering-indicator></media-buffering-indicator>
	</media-player>
</div>
