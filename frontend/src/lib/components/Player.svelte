<script lang="ts">
	import { page } from '$app/state';
	import { PUBLIC_API_URL } from '$env/static/public';
	import type { RecordModel } from 'pocketbase';
	import { onMount } from 'svelte';
	import type { MediaTimeUpdateEventDetail } from 'vidstack';
	import { LocalMediaStorage } from 'vidstack';
	import type { MediaPlayerElement } from 'vidstack/elements';

	let {
		video = {} as RecordModel,
		player = $bindable(),
		currentTime = $bindable()
	}: { video: RecordModel; player: MediaPlayerElement; currentTime: number } = $props();

	class CustomLocalMediaStorage extends LocalMediaStorage {
		getTime(): Promise<number | null> {
			const t = parseInt(page.url.searchParams.get('t') || '0');
			if (t) {
				return Promise.resolve(t);
			}
			return Promise.resolve(currentTime);
		}
	}

	onMount(async () => {
		await import('vidstack/player/styles/default/theme.css');
		await import('vidstack/player/styles/default/layouts/video.css');
		await import('vidstack/player');
		await import('vidstack/player/layouts');
		await import('vidstack/player/ui');

		if (player) {
			player.storage = new CustomLocalMediaStorage();
			thumbnails = `${PUBLIC_API_URL}/${type}/${video.filename}/sprites/sprites.vtt`;

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

	let thumbnails: string = $state('');
	let type = video.collectionName === 'vod' ? 'vods' : 'clips';
</script>

<div class="aspect-video overflow-hidden rounded-xl">
	<media-player
		class="player h-full w-full"
		title={video.title}
		src="{PUBLIC_API_URL}/{type}/{video.filename}/{video.collectionName}.mp4"
		crossorigin
		bind:this={player}
	>
		<media-provider>
			<media-poster class="vds-poster" src="{PUBLIC_API_URL}/{type}/{video.filename}/thumb-lg.webp"
			></media-poster>
			{#if type === 'vods'}
				<track
					label="Deutsch"
					src="{PUBLIC_API_URL}/{type}/{video.filename}/subtitles.vtt"
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
