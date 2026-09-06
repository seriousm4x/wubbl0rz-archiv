<script lang="ts">
	import { page } from '$app/state';
	import { PUBLIC_API_URL } from '$env/static/public';
	import type { RecordModel } from 'pocketbase';
	import { onMount } from 'svelte';
	import type { Snippet } from 'svelte';
	import type { MediaTimeUpdateEventDetail } from 'vidstack';
	import { LocalMediaStorage } from 'vidstack';
	import 'vidstack/bundle';
	import type { MediaPlayerElement } from 'vidstack/elements';

	let {
		video = {} as RecordModel,
		class: className = '',
		player = $bindable(),
		// eslint-disable-next-line no-useless-assignment
		currentTime = $bindable(),
		isAudio = $bindable(false),
		info,
		actions
	}: {
		video: RecordModel;
		class?: string;
		player: MediaPlayerElement;
		currentTime: number;
		isAudio: boolean;
		info?: Snippet;
		actions?: Snippet;
	} = $props();

	class CustomLocalMediaStorage extends LocalMediaStorage {
		async getTime(): Promise<number | null> {
			const tParam = page.url.searchParams.get('t');
			const t = tParam ? Number(tParam) : 0;
			if (Number.isFinite(t) && t > 0) return t;
			return super.getTime();
		}
	}

	onMount(() => {
		if (!player) return;

		player.storage = new CustomLocalMediaStorage();
		const onTimeUpdate = (event: Event) => {
			const e = event as CustomEvent<MediaTimeUpdateEventDetail>;
			currentTime = e.detail.currentTime;
		};
		player.addEventListener('time-update', onTimeUpdate);

		return () => player.removeEventListener('time-update', onTimeUpdate);
	});

	$effect(() => {
		if (!player) return;

		const t = Number(page.url.searchParams.get('t'));
		if (Number.isFinite(t) && t > 0) player.currentTime = t;
	});

	let type = $derived(video.collectionName === 'vod' ? 'vods' : 'clips');
</script>

{#key video}
	<media-player
		class={['player', className]}
		title={video.title}
		viewType={isAudio ? 'audio' : 'video'}
		streamType="on-demand"
		bind:this={player}
		crossorigin
		playsInline
		autoplay
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
		{#if (info || actions) && !isAudio}
			<media-controls class="vds-controls pointer-events-none">
				{#if info}
					<div
						class="absolute inset-x-0 top-0 bg-linear-to-b from-black/85 via-black/45 to-transparent px-4 pt-5 pb-16 text-white"
					>
						{@render info()}
					</div>
				{/if}
				{#if actions}
					<div class="pointer-events-auto absolute top-3 right-3 flex items-center gap-1">
						{@render actions()}
					</div>
				{/if}
			</media-controls>
		{/if}
		<media-buffering-indicator></media-buffering-indicator>
	</media-player>
{/key}
