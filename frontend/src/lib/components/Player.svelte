<script lang="ts">
	import { page } from '$app/state';
	import { PUBLIC_API_URL } from '$app/env/public';
	import { connectPlayer, type ArchivePlayerElement } from '#lib/player';
	import '@videojs/html/video/player';
	import '@videojs/html/video/skin';
	import '@videojs/html/audio/player';
	import '@videojs/html/audio/skin';
	import '@videojs/html/ui/menu-radio-group';
	import { qualityIcon, registerIcons } from '@videojs/html/icons';
	import type { RecordModel } from 'pocketbase';
	import { onMount, untrack, type Snippet } from 'svelte';

	registerIcons('default', { quality: qualityIcon });

	let {
		video = {} as RecordModel,
		class: className = '',
		player = $bindable(),
		 
		currentTime = $bindable(0),
		isAudio = $bindable(false),
		quality,
		actions
	}: {
		video: RecordModel;
		class?: string;
		player?: ArchivePlayerElement;
		currentTime?: number;
		isAudio?: boolean;
		quality?: Snippet;
		actions?: Snippet;
	} = $props();

	let media = $state<HTMLMediaElement>();
	let mounted = $state(false);
	const componentId = $props.id();
	let type = $derived(video.collectionName === 'vod' ? 'vods' : 'clips');
	let base = $derived(`${PUBLIC_API_URL}/${type}/${video.filename}`);
	let src = $derived(`${base}/${isAudio ? 'audio.ogg' : `${video.collectionName}.mp4`}`);

	function setPlaybackRate(rate: number) {
		if (media) media.playbackRate = rate;
	}

	onMount(() => {
		mounted = true;
	});

	$effect(() => {
		if (!player || !media) return;
		return untrack(() =>
			connectPlayer(
				player!,
				media!,
				src,
				() => Math.max(Number(page.url.searchParams.get('t')), currentTime),
				(time) => (currentTime = time)
			)
		);
	});

	$effect(() => {
		const time = Number(page.url.searchParams.get('t'));
		if (player?.store.target && media && media.readyState >= 1 && Number.isFinite(time) && time > 0)
			void player.store.seek(time).catch(() => {});
	});
</script>

{#if mounted}
	{#key `${video.id}:${src}`}
		<svelte:element
			this={isAudio ? 'audio-player' : 'video-player'}
			bind:this={player}
			content-title={video.title}
			poster={`${base}/thumb-lg.webp`}
		>
			{#if isAudio}
				<audio-skin class={className}>
					<audio bind:this={media} {src} crossorigin="anonymous" autoplay preload="metadata"
					></audio>
				</audio-skin>
			{:else}
				<media-container
					class={['media-skin media-container video-skin block aspect-video w-full', className]}
					data-theme="default"
					data-preset="video"
				>
					<!-- svelte-ignore a11y_media_has_caption -->
					<video
						bind:this={media}
						{src}
						crossorigin="anonymous"
						playsinline
						autoplay
						preload="metadata"
					>
						{#if type === 'vods'}
							<track label="Deutsch" src={`${base}/subtitles.vtt`} kind="subtitles" srclang="de" />
						{/if}
						<track label="thumbnails" src={`${base}/sprites/sprites.vtt`} kind="metadata" default />
					</video>
					<media-controls visibility="auto">
						<media-controls-backdrop class="video-controls-backdrop"></media-controls-backdrop>
						<media-controls-content class="video-controls video-controls-content">
							<media-tooltip-group>
								<media-controls-group class="video-controls-primary">
									<media-play-button class="media-button media-play-button">
										<media-icon name="play" class="media-button-icon media-play-button-play-icon"
										></media-icon>
										<media-icon name="pause" class="media-button-icon media-play-button-pause-icon"
										></media-icon>
									</media-play-button>
									<media-mute-button
										commandfor={`${componentId}-volume`}
										class="media-button media-mute-button video-controls-volume-button"
									>
										<media-icon
											name="volume-off"
											class="media-button-icon media-mute-button-off-icon"
										></media-icon>
										<media-icon
											name="volume-low"
											class="media-button-icon media-mute-button-low-icon"
										></media-icon>
										<media-icon
											name="volume-high"
											class="media-button-icon media-mute-button-high-icon"
										></media-icon>
									</media-mute-button>
									<media-volume-popover
										id={`${componentId}-volume`}
										open-on-hover
										delay="200"
										close-delay="100"
										side="top"
										class="media-popup media-popup-safe-area media-popup-transition media-popup-surface media-volume-popover"
									>
										<media-volume-slider
											class="media-slider media-volume-slider"
											thumb-alignment="edge"
											orientation="vertical"
										>
											<media-slider-track class="media-slider-track"
												><media-slider-fill class="media-slider-fill"
												></media-slider-fill></media-slider-track
											>
											<media-slider-thumb class="media-slider-thumb media-volume-slider-thumb"
											></media-slider-thumb>
										</media-volume-slider>
									</media-volume-popover>
									<media-controls-group class="video-time-slider-group">
										<media-time class="media-time-value video-time-value" type="current"
										></media-time>
										<media-time-slider class="media-slider media-time-slider">
											<media-slider-track class="media-slider-track"
												><media-slider-buffer class="media-slider-buffer"
												></media-slider-buffer><media-slider-fill class="media-slider-fill"
												></media-slider-fill></media-slider-track
											>
											<media-slider-thumb class="media-slider-thumb media-time-slider-thumb"
											></media-slider-thumb>
											<media-slider-preview class="media-slider-preview" overflow="visible">
												<media-slider-thumbnail
													class="media-slider-preview-content media-popup-surface media-slider-thumbnail"
													><img
														alt=""
														aria-hidden="true"
														decoding="async"
														class="media-slider-thumbnail-image"
													/></media-slider-thumbnail
												>
												<div class="media-slider-preview-content media-time-slider-preview-content">
													<media-slider-value class="media-time-slider-value" type="pointer"
													></media-slider-value>
												</div>
											</media-slider-preview>
										</media-time-slider>
										<media-time class="media-time-toggle video-time-value" type="remaining" toggle
										></media-time>
									</media-controls-group>
								</media-controls-group>
								<media-controls-group class="video-controls-secondary">
									<media-captions-button class="media-button media-captions-button">
										<media-icon
											name="captions-off"
											class="media-button-icon media-captions-button-off-icon"
										></media-icon>
										<media-icon
											name="captions-on"
											class="media-button-icon media-captions-button-on-icon"
										></media-icon>
									</media-captions-button>
									<button
										type="button"
										commandfor={`${componentId}-settings`}
										class="media-button media-settings-menu-trigger video-controls-settings-button"
									>
										<media-icon
											name="gear"
											class="media-button-icon-base media-settings-menu-trigger-icon"
										></media-icon>
										<media-text class="media-settings-menu-trigger-label" token="menu.settings"
											>Settings</media-text
										>
									</button>
									<media-menu
										id={`${componentId}-settings`}
										side="top"
										align="center"
										class="media-popup media-popup-surface media-menu-popup media-menu-resizable-popup"
									>
										<media-menu-content class="media-menu-content">
											{#if quality}{@render quality()}{/if}
											<media-menu-item
												commandfor={`${componentId}-speed`}
												class="media-menu-trigger-item"
											>
												<media-icon name="speed" class="media-menu-trigger-item-icon"></media-icon>
												<media-text token="menu.speed">Speed</media-text>
												<span class="media-menu-hint">
													<span data-part="value" class="media-menu-hint-label"></span>
													<media-icon name="chevron" class="media-menu-forward-chevron"
													></media-icon>
												</span>
											</media-menu-item>
											<media-menu-content id={`${componentId}-speed`} class="media-menu-content">
												<media-menu-item class="media-menu-back-item">
													<media-icon name="chevron" class="media-menu-back-chevron"></media-icon>
													<media-text token="menu.speed">Speed</media-text>
												</media-menu-item>
												<media-menu-separator class="media-menu-separator"></media-menu-separator>
												<media-menu-item
													class="media-menu-radio-item"
													onselect={() => setPlaybackRate(0.5)}>0.5x</media-menu-item
												>
												<media-menu-item
													class="media-menu-radio-item"
													onselect={() => setPlaybackRate(0.75)}>0.75x</media-menu-item
												>
												<media-menu-item
													class="media-menu-radio-item"
													onselect={() => setPlaybackRate(1)}>1x</media-menu-item
												>
												<media-menu-item
													class="media-menu-radio-item"
													onselect={() => setPlaybackRate(1.25)}>1.25x</media-menu-item
												>
												<media-menu-item
													class="media-menu-radio-item"
													onselect={() => setPlaybackRate(1.5)}>1.5x</media-menu-item
												>
												<media-menu-item
													class="media-menu-radio-item"
													onselect={() => setPlaybackRate(2)}>2x</media-menu-item
												>
											</media-menu-content>
											{#if actions}{@render actions()}{/if}
										</media-menu-content>
									</media-menu>
									<media-pip-button class="media-button media-pip-button">
										<media-icon
											name="pip-enter"
											class="media-button-icon media-pip-button-enter-icon"
										></media-icon>
										<media-icon name="pip-exit" class="media-button-icon media-pip-button-exit-icon"
										></media-icon>
									</media-pip-button>
									<media-fullscreen-button class="media-button media-fullscreen-button">
										<media-icon
											name="fullscreen-enter"
											class="media-button-icon media-fullscreen-button-enter-icon"
										></media-icon>
										<media-icon
											name="fullscreen-exit"
											class="media-button-icon media-fullscreen-button-exit-icon"
										></media-icon>
									</media-fullscreen-button>
								</media-controls-group>
							</media-tooltip-group>
						</media-controls-content>
					</media-controls>
					<media-hotkey keys="Space" action="togglePaused"></media-hotkey>
					<media-hotkey keys="k" action="togglePaused"></media-hotkey>
					<media-hotkey keys="m" action="toggleMuted"></media-hotkey>
					<media-hotkey keys="ArrowRight" action="seekStep"></media-hotkey>
					<media-hotkey keys="ArrowLeft" action="seekStep"></media-hotkey>
					<media-hotkey keys="l" action="seekStep"></media-hotkey>
					<media-hotkey keys="j" action="seekStep"></media-hotkey>
					<media-hotkey keys="ArrowUp" action="volumeStep"></media-hotkey>
					<media-hotkey keys="ArrowDown" action="volumeStep"></media-hotkey>
					<media-hotkey keys="0-9" action="seekToPercent"></media-hotkey>
					<media-hotkey keys="Home" action="seekToPercent" value="0"></media-hotkey>
					<media-hotkey keys="End" action="seekToPercent" value="100"></media-hotkey>
					<media-hotkey keys="&gt;" action="speedUp"></media-hotkey>
					<media-hotkey keys="&lt;" action="speedDown"></media-hotkey>
					<media-hotkey keys="f" action="toggleFullscreen"></media-hotkey>
					<media-hotkey keys="c" action="toggleSubtitles"></media-hotkey>
					<media-hotkey keys="i" action="togglePictureInPicture"></media-hotkey>
					<media-gesture type="tap" action="togglePaused" pointer="mouse" region="center"
					></media-gesture>
					<media-gesture type="tap" action="toggleControls" pointer="touch"></media-gesture>
					<media-gesture type="doubletap" action="seekStep" region="left"></media-gesture>
					<media-gesture type="doubletap" action="toggleFullscreen" region="center"></media-gesture>
					<media-gesture type="doubletap" action="seekStep" region="right"></media-gesture>
					<media-status-announcer class="media-status-announcer"></media-status-announcer>
					<div class="video-status-indicators">
						<media-volume-indicator class="media-indicator media-volume-indicator">
							<media-volume-indicator-fill
								class="media-indicator-content media-volume-indicator-fill"
							>
								<media-icon name="volume-high" class="media-volume-indicator-high-icon"
								></media-icon>
								<media-icon name="volume-low" class="media-volume-indicator-low-icon"></media-icon>
								<media-icon name="volume-off" class="media-volume-indicator-off-icon"></media-icon>
								<media-volume-indicator-value class="media-volume-indicator-value"
								></media-volume-indicator-value>
							</media-volume-indicator-fill>
						</media-volume-indicator>
						<media-status-indicator
							actions="toggleSubtitles,toggleFullscreen,togglePictureInPicture"
							class="media-indicator media-status-indicator"
						>
							<div class="media-indicator-content media-status-indicator-content">
								<media-icon name="captions-on" class="media-status-indicator-captions-on-icon"
								></media-icon>
								<media-icon name="captions-off" class="media-status-indicator-captions-off-icon"
								></media-icon>
								<media-icon
									name="fullscreen-enter"
									class="media-status-indicator-fullscreen-enter-icon"
								></media-icon>
								<media-icon
									name="fullscreen-exit"
									class="media-status-indicator-fullscreen-exit-icon"
								></media-icon>
								<media-icon name="pip-enter" class="media-status-indicator-pip-enter-icon"
								></media-icon>
								<media-icon name="pip-exit" class="media-status-indicator-pip-exit-icon"
								></media-icon>
								<media-status-indicator-value class="media-status-indicator-value"
								></media-status-indicator-value>
							</div>
						</media-status-indicator>
						<media-seek-indicator class="media-seek-indicator">
							<media-icon name="chevron" class="media-seek-indicator-icon"></media-icon>
							<media-seek-indicator-value class="media-seek-indicator-value"
							></media-seek-indicator-value>
						</media-seek-indicator>
						<media-status-indicator actions="togglePaused" class="media-playback-status-indicator">
							<media-icon name="play" class="media-playback-status-indicator-play-icon"
							></media-icon>
							<media-icon name="pause" class="media-playback-status-indicator-pause-icon"
							></media-icon>
						</media-status-indicator>
					</div>
				</media-container>
			{/if}
		</svelte:element>
	{/key}
{/if}

<style>
	:root {
		--media-border-radius: 0px;
		--media-border-color: transparent;
	}
</style>
