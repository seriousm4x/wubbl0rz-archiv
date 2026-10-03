<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { PUBLIC_API_URL } from '$env/static/public';
	import ChatReplay from '$lib/components/ChatReplay.svelte';
	import Player from '$lib/components/Player.svelte';
	import SEO from '$lib/components/SEO.svelte';
	import { formatBytes, toHHMMSS } from '$lib/functions';
	import { DefaultOpenGraph } from '$lib/types/opengraph';
	import IconDownloadMinimalisticBoldDuotone from '@iconify-icons/solar/download-minimalistic-bold-duotone';
	import IconShareBoldDuotone from '@iconify-icons/solar/share-bold-duotone';
	import Icon from '@iconify/svelte';
	import { format, parseISO } from 'date-fns';
	import { de } from 'date-fns/locale';
	import type { RecordModel } from 'pocketbase';
	import type { ArchivePlayerElement } from '$lib/player';

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
	let vod = $derived(data.vod as RecordModel);

	function copyLink(withTimestamp: boolean) {
		const url = new URL(page.url.origin + page.url.pathname);
		if (withTimestamp) {
			url.searchParams.set('t', currentTime.toFixed(0));
		}
		navigator.clipboard.writeText(url.toString());
	}
</script>

<SEO {og} />

<div
	class="-m-4 grid h-[calc(100dvh-4rem)] grid-rows-[auto_minmax(0,1fr)] overflow-hidden xl:grid-cols-[minmax(0,1fr)_22rem] xl:grid-rows-[minmax(0,1fr)]"
>
	<div class="min-h-0 min-w-0 bg-black xl:h-full" bind:clientHeight={playerHeight}>
		<Player
			class="vod-player h-full w-full"
			bind:player
			bind:currentTime
			video={vod}
			isAudio={false}
		>
			{#snippet info()}
				<h1
					class="max-w-[calc(100%-5rem)] text-xl leading-tight font-semibold tracking-tight sm:text-2xl"
				>
					{vod.title}
				</h1>
				<div class="mt-3 flex flex-wrap gap-x-2 gap-y-1 text-sm text-white/75">
					<span>
						{format(parseISO(vod.date), 'dd.MM.yyyy, HH:mm', { locale: de })} Uhr
					</span>
					|
					<span>{vod.viewcount.toLocaleString('de-DE')} Views</span> |
					<span>{vod.resolution}@{vod.fps} FPS</span>
				</div>
			{/snippet}

			{#snippet actions()}
				<media-menu-item commandfor="vod-download" class="media-menu-trigger-item">
					<Icon icon={IconDownloadMinimalisticBoldDuotone} class="media-menu-trigger-item-icon" />
					<span>Herunterladen</span>
					<span class="media-menu-hint">
						<media-icon name="chevron" class="media-menu-forward-chevron"></media-icon>
					</span>
				</media-menu-item>
				<media-menu-content class="media-menu-content" id="vod-download">
					<media-menu-item class="media-menu-back-item">
						<media-icon name="chevron" class="media-menu-back-chevron"></media-icon>
						<span>Herunterladen</span>
					</media-menu-item>
					<media-menu-separator class="media-menu-separator"></media-menu-separator>
					<media-menu-item class="media-menu-radio-item">
						<a
							class="block w-full"
							download={`${vod.filename}.mp4`}
							href={resolve('/download/[type]/[filename]', {
								type: vod.collectionName,
								filename: vod.filename
							})}
						>
							Video ({formatBytes(vod.size)})
						</a>
					</media-menu-item>
					<media-menu-item class="media-menu-radio-item">
						<a
							class="block w-full"
							download={`${vod.filename}.ogg`}
							href={resolve('/download/[type]/[filename]?audio=true', {
								type: vod.collectionName,
								filename: vod.filename
							})}
						>
							Audio ({formatBytes(vod.size_audio)})
						</a>
					</media-menu-item>
				</media-menu-content>

				<media-menu-item commandfor="vod-share" class="media-menu-trigger-item">
					<Icon icon={IconShareBoldDuotone} class="media-menu-trigger-item-icon" />
					<span>Teilen</span>
					<span class="media-menu-hint">
						<media-icon name="chevron" class="media-menu-forward-chevron"></media-icon>
					</span>
				</media-menu-item>
				<media-menu-content class="media-menu-content" id="vod-share">
					<media-menu-item class="media-menu-back-item">
						<media-icon name="chevron" class="media-menu-back-chevron"></media-icon>
						<span>Teilen</span>
					</media-menu-item>
					<media-menu-separator class="media-menu-separator"></media-menu-separator>
					<media-menu-item class="media-menu-radio-item" onselect={() => copyLink(false)}>
						Link kopieren
					</media-menu-item>
					<media-menu-item class="media-menu-radio-item" onselect={() => copyLink(true)}>
						Link bei {toHHMMSS(currentTime, false)} kopieren
					</media-menu-item>
				</media-menu-content>
			{/snippet}
		</Player>
	</div>

	<aside class="min-h-0 p-4 xl:h-full xl:p-0">
		<ChatReplay {vod} {currentTime} {playerHeight} />
	</aside>
</div>

<style>
	:global(.vod-player video) {
		width: 100%;
		height: 100%;
		object-fit: contain;
	}
</style>
