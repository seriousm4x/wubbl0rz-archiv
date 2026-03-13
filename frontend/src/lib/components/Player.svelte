<script lang="ts">
	import { page } from '$app/state';
	import { PUBLIC_API_URL } from '$env/static/public';
	import type { RecordModel } from 'pocketbase';
	import { onMount } from 'svelte';
	import type { MediaTimeUpdateEventDetail } from 'vidstack';
	import { LocalMediaStorage } from 'vidstack';
	import 'vidstack/bundle';
	import type { MediaPlayerElement } from 'vidstack/elements';

	let {
		video = {} as RecordModel,
		player = $bindable(),
		// eslint-disable-next-line no-useless-assignment
		currentTime = $bindable(),
		isAudio = $bindable(false)
	}: {
		video: RecordModel;
		player: MediaPlayerElement;
		currentTime: number;
		isAudio: boolean;
	} = $props();

	class CustomLocalMediaStorage extends LocalMediaStorage {
		async getTime(): Promise<number | null> {
			const tParam = page.url.searchParams.get('t');
			const t = tParam ? parseInt(tParam, 10) : 0;
			if (t !== 0) return t;
			return super.getTime();
		}
	}

	onMount(() => {
		if (player) {
			player.storage = new CustomLocalMediaStorage();

			player.addEventListener('time-update', (event: Event) => {
				const e = event as CustomEvent<MediaTimeUpdateEventDetail>;
				currentTime = e.detail.currentTime;
			});
		}
	});

	$effect(() => {
		if (player) {
			player.currentTime = parseInt(page.url.searchParams.get('t') || '0');
		}
	});

	let type = video.collectionName === 'vod' ? 'vods' : 'clips';
</script>

<div class="overflow-hidden rounded-xl">
	<media-player
		class="player"
		title={video.title}
		viewType={isAudio ? 'audio' : 'video'}
		streamType="on-demand"
		bind:this={player}
		crossorigin
		playsInline
	>
		<media-provider>
			<media-poster class="vds-poster" src="{PUBLIC_API_URL}/{type}/{video.filename}/thumb-lg.webp"
			></media-poster>
			{#if isAudio}
				<source src="{PUBLIC_API_URL}/{type}/{video.filename}/audio.ogg" type="audio/ogg" />
			{:else}
				<source
					src="{PUBLIC_API_URL}/{type}/{video.filename}/{video.collectionName}.mp4"
					type="video/mp4"
				/>
			{/if}
			{#if type === 'vods'}
				<track
					label="Deutsch"
					src="{PUBLIC_API_URL}/{type}/{video.filename}/subtitles.vtt"
					kind="subtitles"
					srclang="de"
				/>
			{/if}
		</media-provider>
		<media-audio-layout
			playbackRates={[
				0.25, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2, 2.25, 2.5, 2.75, 3, 3.25, 3.5, 3.75, 4
			]}
		></media-audio-layout>
		<media-video-layout
			thumbnails="{PUBLIC_API_URL}/{type}/{video.filename}/sprites/sprites.vtt"
			playbackRates={[
				0.25, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2, 2.25, 2.5, 2.75, 3, 3.25, 3.5, 3.75, 4
			]}
		></media-video-layout>
		<media-buffering-indicator></media-buffering-indicator>
	</media-player>
</div>
