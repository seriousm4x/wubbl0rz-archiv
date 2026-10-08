<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { PUBLIC_API_URL } from '$app/env/public';
	import ChatReplay from '#lib/components/ChatReplay.svelte';
	import Player from '#lib/components/Player.svelte';
	import SEO from '#lib/components/SEO.svelte';
	import { formatBytes } from '#lib/functions';
	import { DefaultOpenGraph } from '#lib/types/opengraph';
	import IconDownloadMinimalisticBoldDuotone from '@iconify-icons/solar/download-minimalistic-bold-duotone';
	import IconShareBoldDuotone from '@iconify-icons/solar/share-bold-duotone';
	import Icon from '@iconify/svelte';
	import { parseISO } from 'date-fns';
	import type { RecordModel } from 'pocketbase';
	import type { ArchivePlayerElement } from '#lib/player';

	let { data } = $props();

	let og = $derived({
		...DefaultOpenGraph,
		title: data.vod?.title,
		image: `${PUBLIC_API_URL}/vods/${data.vod?.filename}/thumb-lg.webp`,
		updated_time: parseISO(data.vod?.date).toISOString()
	});

	let player = $state<ArchivePlayerElement>();
	let playerHeight = $state(0);
	let currentTime = $state(0);
	let isAudio = $state(false);
	let vod = $derived(data.vod as RecordModel);

	function copyLink(withTimecode = false) {
		const url = new URL(page.url.href);
		if (withTimecode) url.searchParams.set('t', Math.floor(currentTime).toString());
		else url.searchParams.delete('t');
		void navigator.clipboard.writeText(url.href);
	}

	function download(audio = false) {
		const url = resolve('/download/[type]/[filename]', {
			type: vod.collectionName,
			filename: vod.filename
		});
		window.location.href = audio ? `${url}?audio=true` : url;
	}
</script>

<SEO {og} />

<div
	class="-m-4 grid h-[calc(100dvh-4rem)] grid-rows-[auto_minmax(0,1fr)] overflow-hidden xl:grid-cols-[minmax(0,1fr)_22rem] xl:grid-rows-[minmax(0,1fr)]"
>
	<div
		class="min-h-0 min-w-0 {isAudio ? 'xl:self-start' : 'bg-black xl:h-full'}"
		bind:clientHeight={playerHeight}
	>
		<Player
			class="vod-player w-full {isAudio ? '' : 'h-full'}"
			bind:player
			bind:currentTime
			video={vod}
			{isAudio}
		>
			{#snippet quality()}
				<media-menu-item commandfor="vod-quality-menu" class="media-menu-trigger-item">
					<media-icon name="quality" class="media-menu-trigger-item-icon"></media-icon>
					<media-text>Qualität</media-text>
					<span class="media-menu-hint">
						{isAudio ? 'Nur Audio' : '1080p'}
						<media-icon name="chevron" class="media-menu-forward-chevron"></media-icon>
					</span>
				</media-menu-item>
				<media-menu-content id="vod-quality-menu" class="media-menu-content">
					<media-menu-item class="media-menu-back-item">
						<media-icon name="chevron" class="media-menu-back-chevron"></media-icon>
						<media-text>Qualität</media-text>
					</media-menu-item>
					<media-menu-radio-group
						value={isAudio ? 'audio' : '1080p'}
						onvalue-change={(event: CustomEvent<{ value: string }>) =>
							(isAudio = event.detail.value === 'audio')}
					>
						<media-menu-radio-item value="1080p" class="media-menu-radio-item">
							1080p
							<media-menu-item-indicator checked={!isAudio}>✓</media-menu-item-indicator>
						</media-menu-radio-item>
						<media-menu-radio-item value="audio" class="media-menu-radio-item">
							Nur Audio
							<media-menu-item-indicator checked={isAudio}>✓</media-menu-item-indicator>
						</media-menu-radio-item>
					</media-menu-radio-group>
				</media-menu-content>
			{/snippet}
			{#snippet actions()}
				<media-menu-item commandfor="vod-download-menu" class="media-menu-trigger-item">
					<Icon icon={IconDownloadMinimalisticBoldDuotone} class="media-menu-trigger-item-icon" />
					<media-text>Herunterladen</media-text>
					<span class="media-menu-hint">
						<media-icon name="chevron" class="media-menu-forward-chevron"></media-icon>
					</span>
				</media-menu-item>
				<media-menu-content id="vod-download-menu" class="media-menu-content">
					<media-menu-item class="media-menu-back-item">
						<media-icon name="chevron" class="media-menu-back-chevron"></media-icon>
						<media-text>Herunterladen</media-text>
					</media-menu-item>
					<media-menu-item class="media-menu-radio-item" onselect={() => download()}>
						Video <span class="media-menu-hint">{formatBytes(Number(vod.size))}</span>
					</media-menu-item>
					<media-menu-item class="media-menu-radio-item" onselect={() => download(true)}>
						Audio <span class="media-menu-hint">{formatBytes(Number(vod.size_audio))}</span>
					</media-menu-item>
				</media-menu-content>
				<media-menu-item commandfor="vod-share-menu" class="media-menu-trigger-item">
					<Icon icon={IconShareBoldDuotone} class="media-menu-trigger-item-icon" />
					<media-text>Teilen</media-text>
					<span class="media-menu-hint">
						<media-icon name="chevron" class="media-menu-forward-chevron"></media-icon>
					</span>
				</media-menu-item>
				<media-menu-content id="vod-share-menu" class="media-menu-content">
					<media-menu-item class="media-menu-back-item">
						<media-icon name="chevron" class="media-menu-back-chevron"></media-icon>
						<media-text>Teilen</media-text>
					</media-menu-item>
					<media-menu-item class="media-menu-radio-item" onselect={() => copyLink()}
						>Ohne Zeitstempel</media-menu-item
					>
					<media-menu-item class="media-menu-radio-item" onselect={() => copyLink(true)}
						>Mit Zeitstempel</media-menu-item
					>
				</media-menu-content>
			{/snippet}
		</Player>
		{#if isAudio}
			<button type="button" class="btn btn-sm mt-2" onclick={() => (isAudio = false)}>
				Video wiedergeben
			</button>
		{/if}
	</div>

	<aside class="min-h-0 p-4 xl:h-full xl:p-0">
		<ChatReplay {vod} {currentTime} playerHeight={isAudio ? undefined : playerHeight} />
	</aside>
</div>
