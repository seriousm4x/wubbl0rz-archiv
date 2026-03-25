<script lang="ts">
	import { page } from '$app/state';
	import { PUBLIC_API_URL } from '$env/static/public';
	import Card from '$lib/components/Card.svelte';
	import Player from '$lib/components/Player.svelte';
	import SEO from '$lib/components/SEO.svelte';
	import { formatBytes, toHHMMSS } from '$lib/functions';
	import { DefaultOpenGraph } from '$lib/types/opengraph';
	import IconArrowAutofitWidthDotted24Filled from '@iconify-icons/fluent/arrow-autofit-width-dotted-24-filled';
	import IconDownloadSquareBoldDuotone from '@iconify-icons/solar/download-square-bold-duotone';
	import IconFireBoldDuotone from '@iconify-icons/solar/fire-bold-duotone';
	import IconHeadphonesRoundSoundBold from '@iconify-icons/solar/headphones-round-sound-bold';
	import IconRoundArrowRightBoldDuotone from '@iconify-icons/solar/round-arrow-right-bold-duotone';
	import IconRoundArrowRightDownBoldDuotone from '@iconify-icons/solar/round-arrow-right-down-bold-duotone';
	import IconRoundArrowRightUpBoldDuotone from '@iconify-icons/solar/round-arrow-right-up-bold-duotone';
	import IconShareBoldDuotone from '@iconify-icons/solar/share-bold-duotone';
	import Icon from '@iconify/svelte';
	import { format, formatDistance, parseISO } from 'date-fns';
	import { de } from 'date-fns/locale';
	import type { RecordModel } from 'pocketbase';
	import type { MediaPlayerElement } from 'vidstack/elements';

	let { data } = $props();

	let og = $state({
		...DefaultOpenGraph,
		title: data.vod?.title,
		image: `${PUBLIC_API_URL}/vods/${data.vod?.filename}/thumb-lg.webp`,
		updated_time: parseISO(data.vod?.date).toISOString()
	});

	let player: MediaPlayerElement = $state({} as MediaPlayerElement);
	let currentTime: number = $state(0);

	let vod = $derived(data.vod as RecordModel);
	let vodsCount = $derived(data.vodsCount);
	let vodPosition = $derived(data.vodPosition);
	let recommendations = $derived(data.recommendations);
	let percentile = $derived((vodPosition * 100) / vodsCount);
	let percentileRounded = $derived(percentile < 1 ? percentile.toFixed(2) : Math.round(percentile));
	let isAudio = $state(false);
	let theaterEnabled = $state(false);

	function copyLink(withTimestamp: boolean) {
		const url = new URL(page.url.origin + page.url.pathname);
		if (withTimestamp) {
			url.searchParams.set('t', currentTime.toFixed(0));
		}
		navigator.clipboard.writeText(url.toString());
	}

	function onKeyDown(e: KeyboardEvent) {
		switch (e.key) {
			case 't':
				theaterEnabled = !theaterEnabled;
				break;
			default:
				break;
		}
	}
</script>

<svelte:window on:keydown={onKeyDown} />

<SEO bind:og />

<div
	class="absolute top-0 left-0 -z-10 aspect-video h-full w-full bg-cover bg-center opacity-10 blur-2xl"
	style="background-image: url('{PUBLIC_API_URL}/vods/{vod.filename}/thumb-lg.webp');"
></div>
<div class="mx-auto flex max-w-480 flex-col gap-8 {theaterEnabled ? '' : 'xl:flex-row'}">
	<div class="flex flex-col gap-4 {theaterEnabled ? '' : 'xl:w-4/5'}">
		<Player bind:player bind:currentTime video={vod} {isAudio} />
		<h1 class="text-4xl font-bold">
			{vod.title}
		</h1>
		<div class="stats stats-vertical bg-base-200 lg:stats-horizontal w-full shadow">
			<div class="stat">
				<div class="stat-title text-lg">Gestreamt am</div>
				<div class="stat-value text-2xl">
					{format(parseISO(vod.date), 'dd.MM.yyyy', { locale: de })}
				</div>
				<div class="stat-desc">
					{formatDistance(parseISO(vod.date), Date.now(), {
						addSuffix: true,
						includeSeconds: true,
						locale: de
					})} um {format(parseISO(vod.date), 'HH:mm', { locale: de })} Uhr
				</div>
			</div>
			<div class="stat">
				<div class="stat-title text-lg">Views</div>
				<div class="stat-value text-2xl">{vod.viewcount.toLocaleString('de-DE')}</div>
				<div class="stat-desc">
					Sync {formatDistance(parseISO(vod.updated), Date.now(), {
						addSuffix: true,
						includeSeconds: true,
						locale: de
					})}
				</div>
			</div>
			<div class="stat">
				<div class="stat-title text-lg">Platz</div>
				<div class="stat-value flex items-center gap-1 text-2xl">
					#{vodPosition}
					{#if percentile <= 5}
						<div title="Top {percentileRounded}%">
							<Icon icon={IconFireBoldDuotone} class="text-3xl text-red-500" />
						</div>
					{:else if percentile <= 33}
						<div title="Top {percentileRounded}%">
							<Icon icon={IconRoundArrowRightUpBoldDuotone} class="text-3xl text-green-500" />
						</div>
					{:else if percentile <= 66}
						<div title="Top {percentileRounded}%">
							<Icon icon={IconRoundArrowRightBoldDuotone} class="text-slate-500" />
						</div>
					{:else}
						<div title="Top {percentileRounded}%">
							<Icon icon={IconRoundArrowRightDownBoldDuotone} class="text-3xl text-red-500" />
						</div>
					{/if}
				</div>
				<div class="stat-desc">
					von {vodsCount} Vods
				</div>
			</div>
			<div class="stat">
				<div class="stat-title text-lg">Auflösung</div>
				<div class="stat-value text-2xl">{vod.resolution}</div>
				<div class="stat-desc">{vod.fps} FPS</div>
			</div>
		</div>
		<div class="flex flex-row flex-wrap justify-start gap-2 md:ms-auto">
			<button
				class="btn w-fit rounded-xl bg-linear-to-r shadow transition duration-200 hover:shadow-lg"
				onclick={() => (theaterEnabled = !theaterEnabled)}
			>
				<Icon icon={IconArrowAutofitWidthDotted24Filled} class="text-2xl text-violet-500" />
				{theaterEnabled ? 'Standardansicht' : 'Kinomodus'}
			</button>
			<button
				class="btn w-fit rounded-xl bg-linear-to-r shadow transition duration-200 hover:shadow-lg"
				onclick={() => (isAudio = !isAudio)}
			>
				<Icon icon={IconHeadphonesRoundSoundBold} class="text-2xl text-violet-500" />
				{isAudio ? 'Video' : 'Audio only'}
			</button>
			<div class="dropdown md:dropdown-end">
				<label
					for="btn-download"
					tabindex="-1"
					class="btn rounded-xl border-none bg-linear-to-r shadow transition duration-200 hover:shadow-lg"
				>
					<Icon icon={IconDownloadSquareBoldDuotone} class="text-2xl text-violet-500" /> Download
				</label>
				<ul
					id="btn-download"
					tabindex="-1"
					class="menu dropdown-content rounded-box bg-base-200 z-1 w-48 p-2 shadow"
				>
					<li>
						<a href="/download/{vod.collectionName}/{vod.filename}"
							>Video ({formatBytes(vod.size)})</a
						>
					</li>
					<li>
						<a href="/download/{vod.collectionName}/{vod.filename}?audio=true"
							>Audio ({formatBytes(vod.size_audio)})</a
						>
					</li>
				</ul>
			</div>
			<div class="dropdown md:dropdown-end">
				<label
					for="btn-share"
					tabindex="-1"
					class="btn rounded-xl border-none bg-linear-to-r shadow transition duration-200 hover:shadow-lg"
				>
					<Icon icon={IconShareBoldDuotone} class="text-2xl text-violet-500" /> Teilen
				</label>
				<ul
					id="btn-share"
					tabindex="-1"
					class="menu dropdown-content rounded-box bg-base-200 z-1 p-2 shadow"
				>
					<li>
						<button onclick={() => copyLink(false)}>Link kopieren</button>
					</li>
					<li>
						<button class="whitespace-nowrap" onclick={() => copyLink(true)}
							>Link bei {toHHMMSS(currentTime, false)} kopieren</button
						>
					</li>
				</ul>
			</div>
		</div>
		{#if vod['expand']}
			{#if vod['expand']['clip_via_vod'] && vod['expand']['clip_via_vod'].length > 0}
				<h2 class="text-3xl font-bold">
					<span
						class="bg-linear-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent drop-shadow-md"
						>Clips für diesen Stream</span
					>
				</h2>
				<div
					class="grid grid-flow-row-dense grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4"
				>
					{#each vod['expand']['clip_via_vod'] as video, index (index)}
						<Card {video} />
					{/each}
				</div>
			{/if}
		{/if}
	</div>
	<div class="flex flex-col gap-4 {theaterEnabled ? '' : 'xl:w-1/5'}">
		<h2 class="text-3xl font-bold">
			<span
				class="bg-linear-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent drop-shadow-md"
				>Empfohlene Streams</span
			>
		</h2>
		<div
			class="grid grid-flow-row-dense gap-4 {theaterEnabled
				? 'grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4'
				: 'grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-1'}"
		>
			{#each recommendations as video, index (index)}
				<Card {video} />
			{:else}
				Keine Ergebnisse
			{/each}
		</div>
	</div>
</div>
