<script lang="ts">
	import { page } from '$app/state';
	import { PUBLIC_API_URL } from '$env/static/public';
	import { connectPlayer, type ArchivePlayerElement } from '$lib/player';
	import '@videojs/html/video/player';
	import '@videojs/html/video/skin';
	import '@videojs/html/audio/player';
	import type { RecordModel } from 'pocketbase';
	import { untrack, type Snippet } from 'svelte';

	let {
		video = {} as RecordModel,
		class: className = '',
		player = $bindable(),
		// eslint-disable-next-line no-useless-assignment
		currentTime = $bindable(0),
		isAudio = $bindable(false),
		info,
		actions
	}: {
		video: RecordModel;
		class?: string;
		player?: ArchivePlayerElement;
		currentTime?: number;
		isAudio?: boolean;
		info?: Snippet;
		actions?: Snippet;
	} = $props();

	let media = $state<HTMLMediaElement>();
	const componentId = $props.id();
	let type = $derived(video.collectionName === 'vod' ? 'vods' : 'clips');
	let base = $derived(`${PUBLIC_API_URL}/${type}/${video.filename}`);
	let src = $derived(`${base}/${isAudio ? 'audio.ogg' : `${video.collectionName}.mp4`}`);

	$effect(() => {
		if (!player || !media) return;
		return untrack(() =>
			connectPlayer(
				player!,
				media!,
				src,
				() => Number(page.url.searchParams.get('t')),
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

{#key `${video.id}:${src}`}
	<svelte:element
		this={isAudio ? 'audio-player' : 'video-player'}
		bind:this={player}
		content-title={video.title}
		poster={`${base}/thumb-lg.webp`}
	>
		<media-container
			role="group"
			aria-label={video.title}
			class={[
				'player media-skin media-container video-skin block w-full',
				isAudio ? 'min-h-32' : 'aspect-video',
				className
			]}
			data-theme="default"
			data-preset="video"
			style:--media-accent-color="var(--color-primary)"
			style:--media-accent-text-color="var(--color-primary-content)"
			style:--media-border-radius="0px"
		>
			{#if isAudio}
				<audio bind:this={media} {src} crossorigin="anonymous" autoplay preload="metadata"></audio>
			{:else}
				<!-- Clips have no caption asset; VODs provide the existing German subtitle track. -->
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
			{/if}
			{#if !isAudio}
				<media-poster class="media-poster"
					><img alt="" decoding="async" class="media-poster-image" /></media-poster
				>
			{/if}
			<media-buffering-indicator class="media-buffering-indicator">
				<media-icon name="spinner" class="media-buffering-indicator-spinner-icon"></media-icon>
			</media-buffering-indicator>
			<media-error-dialog class="media-dialog-root">
				<media-dialog-backdrop class="media-dialog-backdrop"></media-dialog-backdrop>
				<media-dialog-popup class="media-dialog-popup">
					<div class="media-dialog-content">
						<media-dialog-title class="media-dialog-title"></media-dialog-title>
						<media-dialog-description class="media-dialog-description"></media-dialog-description>
					</div>
					<media-dialog-close class="media-button media-dialog-close"></media-dialog-close>
				</media-dialog-popup>
			</media-error-dialog>
			<media-controls visibility={isAudio ? 'always' : 'auto'}>
				<media-controls-backdrop class="video-controls-backdrop"></media-controls-backdrop>
				<media-controls-content class="video-controls video-controls-content backdrop-filter-none">
					<div
						class="pointer-events-none absolute inset-0 -z-10 rounded-[inherit] backdrop-blur-2xl backdrop-saturate-150"
						aria-hidden="true"
					></div>
					<media-controls-group class="video-controls-primary">
						<media-play-button class="media-button media-play-button">
							<media-icon name="restart" class="media-button-icon media-play-button-restart-icon"
							></media-icon>
							<media-icon name="play" class="media-button-icon media-play-button-play-icon"
							></media-icon>
							<media-icon name="pause" class="media-button-icon media-play-button-pause-icon"
							></media-icon>
						</media-play-button>
						<media-controls-group
							class="volume-control flex shrink-0 items-center gap-0"
							aria-label="Lautstärke"
						>
							<media-mute-button class="media-button media-mute-button">
								<media-icon name="volume-off" class="media-button-icon media-mute-button-off-icon"
								></media-icon>
								<media-icon name="volume-low" class="media-button-icon media-mute-button-low-icon"
								></media-icon>
								<media-icon name="volume-high" class="media-button-icon media-mute-button-high-icon"
								></media-icon>
							</media-mute-button>
							<div class="volume-panel">
								<media-volume-slider
									class="media-slider media-volume-slider"
									orientation="horizontal"
									style="width: 72px; flex: none; margin-inline: 8px"
									thumb-alignment="edge"
									aria-label="Lautstärke"
								>
									<media-slider-track class="media-slider-track"
										><media-slider-fill class="media-slider-fill"
										></media-slider-fill></media-slider-track
									>
									<media-slider-thumb class="media-slider-thumb media-volume-slider-thumb"
									></media-slider-thumb>
								</media-volume-slider>
							</div>
						</media-controls-group>
						<media-controls-group class="video-time-slider-group">
							<media-time class="media-time-value video-time-value" type="current"></media-time>
							<media-time-slider class="media-slider media-time-slider">
								<media-slider-track class="media-slider-track">
									<media-slider-buffer class="media-slider-buffer"></media-slider-buffer>
									<media-slider-fill class="media-slider-fill"></media-slider-fill>
								</media-slider-track>
								<media-slider-thumb class="media-slider-thumb media-time-slider-thumb"
								></media-slider-thumb>
								<media-slider-preview class="media-slider-preview" overflow="visible">
									<media-slider-thumbnail
										class="media-slider-preview-content media-popup-surface media-slider-thumbnail"
									>
										<img
											alt=""
											aria-hidden="true"
											decoding="async"
											class="media-slider-thumbnail-image"
										/>
									</media-slider-thumbnail>
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
						{#if !isAudio}
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
									<media-menu-item
										commandfor={`${componentId}-settings-speed`}
										class="media-menu-trigger-item"
									>
										<media-icon name="speed" class="media-menu-trigger-item-icon"></media-icon>
										<media-text token="menu.speed">Speed</media-text>
										<span class="media-menu-hint">
											<span data-part="value" class="media-menu-hint-label"></span>
											<media-icon name="chevron" class="media-menu-forward-chevron"></media-icon>
										</span>
									</media-menu-item>
									{#if actions}{@render actions()}{/if}
									<media-menu-content
										class="media-menu-content"
										id={`${componentId}-settings-speed`}
									>
										<media-menu-item class="media-menu-back-item">
											<media-icon name="chevron" class="media-menu-back-chevron"></media-icon>
											<media-text token="menu.speed">Speed</media-text>
										</media-menu-item>
										<media-menu-separator class="media-menu-separator"></media-menu-separator>
										<media-playback-rate-radio-group class="media-menu-radio-group">
											<template>
												<media-menu-radio-item class="media-menu-radio-item">
													<span data-part="label"></span>
													<media-menu-item-indicator force-mount class="media-menu-item-indicator">
														<media-icon name="check" class="media-menu-radio-item-icon"
														></media-icon>
													</media-menu-item-indicator>
												</media-menu-radio-item>
											</template>
										</media-playback-rate-radio-group>
									</media-menu-content>
								</media-menu-content>
							</media-menu>
							<media-airplay-button class="media-button media-airplay-button">
								<media-icon
									name="airplay-enter"
									class="media-button-icon media-airplay-button-enter-icon"
								></media-icon>
								<media-icon
									name="airplay-exit"
									class="media-button-icon media-airplay-button-exit-icon"
								></media-icon>
							</media-airplay-button>
							<media-pip-button class="media-button media-pip-button">
								<media-icon name="pip-enter" class="media-button-icon media-pip-button-enter-icon"
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
						{/if}
					</media-controls-group>
				</media-controls-content>
			</media-controls>
			{#if info && !isAudio}
				<media-controls>
					<media-controls-content class="player-overlay block transition-opacity">
						{#if info && !isAudio}
							<div
								class="bg-linear-to-b from-black/85 via-black/45 to-transparent px-4 pt-5 pb-16 text-white"
							>
								{@render info()}
							</div>
						{/if}
					</media-controls-content>
				</media-controls>
			{/if}
			<media-hotkey keys="Space" action="togglePaused"></media-hotkey>
			<media-hotkey keys="k" action="togglePaused"></media-hotkey>
			<media-hotkey keys="m" action="toggleMuted"></media-hotkey>
			<media-hotkey keys="ArrowRight" action="seekStep" value="10"></media-hotkey>
			<media-hotkey keys="ArrowLeft" action="seekStep" value="-10"></media-hotkey>
			<media-hotkey keys="Shift+ArrowRight" action="seekStep" value="20"></media-hotkey>
			<media-hotkey keys="Shift+ArrowLeft" action="seekStep" value="-20"></media-hotkey>
			<media-hotkey keys="l" action="seekStep" value="10"></media-hotkey>
			<media-hotkey keys="j" action="seekStep" value="-10"></media-hotkey>
			<media-hotkey keys="ArrowUp" action="volumeStep" value="0.05"></media-hotkey>
			<media-hotkey keys="ArrowDown" action="volumeStep" value="-0.05"></media-hotkey>
			<media-hotkey keys="0-9" action="seekToPercent"></media-hotkey>
			<media-hotkey keys=">" action="speedUp"></media-hotkey>
			<media-hotkey keys="<" action="speedDown"></media-hotkey>
			{#if !isAudio}
				<media-hotkey keys="f" action="toggleFullscreen"></media-hotkey>
				<media-hotkey keys="c" action="toggleSubtitles"></media-hotkey>
				<media-hotkey keys="i" action="togglePictureInPicture"></media-hotkey>
				<media-gesture type="tap" action="togglePaused" pointer="mouse" region="center"
				></media-gesture>
				<media-gesture type="tap" action="toggleControls" pointer="touch"></media-gesture>
				<media-gesture type="doubletap" action="seekStep" value="-10" region="left"></media-gesture>
				<media-gesture type="doubletap" action="toggleFullscreen" region="center"></media-gesture>
				<media-gesture type="doubletap" action="seekStep" value="10" region="right"></media-gesture>
			{/if}
			<media-status-announcer class="media-status-announcer"></media-status-announcer>
		</media-container>
	</svelte:element>
{/key}

<style>
	@container media-root (width < 32rem) {
		.video-controls-content {
			position: absolute;
			inset-inline: 8px;
			bottom: 8px;
			z-index: 30;
			display: flex;
			flex-wrap: wrap;
			align-items: center;
			background: var(--media-popover);
			border-radius: var(--media-controls-radius);
		}
		.video-controls-primary,
		.video-controls-secondary {
			position: static;
			padding: 0;
			background: none;
			box-shadow: none;
			backdrop-filter: none;
		}
		.video-controls-primary {
			width: 100%;
		}
		.video-controls-secondary {
			margin-inline-start: auto;
		}
		:global(.video-controls-content:not([data-visible]):not(:focus-within)) {
			opacity: 0;
			pointer-events: none;
		}
	}
	.volume-panel {
		display: flex;
		align-items: center;
		width: 0;
		overflow: hidden;
		opacity: 0;
		visibility: hidden;
		pointer-events: none;
		transition:
			width 0.2s ease,
			opacity 0.2s ease;
	}
	.volume-control:hover .volume-panel,
	.volume-control:focus-within .volume-panel {
		width: 88px;
		opacity: 1;
		visibility: visible;
		pointer-events: auto;
	}
	@media (prefers-reduced-motion: reduce) {
		.volume-panel {
			transition: none;
		}
	}
	:global(.volume-control:has(media-volume-slider[data-availability='unsupported']) .volume-panel) {
		display: none;
	}
	.player-overlay {
		position: absolute;
		inset: 0 0 auto;
		z-index: 2;
		pointer-events: none;
	}
	:global(.player video) {
		object-fit: contain;
	}
	:global(.player-overlay:not([data-visible]):not(:focus-within)) {
		opacity: 0;
		visibility: hidden;
	}
</style>
