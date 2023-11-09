<script lang="ts">
	import { page } from '$app/stores';
	import { PUBLIC_API_URL } from '$env/static/public';
	import Button from '$lib/components/Button.svelte';
	import Card from '$lib/components/Card.svelte';
	import Player from '$lib/components/Player.svelte';
	import SEO from '$lib/components/SEO.svelte';
	import { formatBytes, toHHMMSS } from '$lib/functions';
	import { DefaultOpenGraph } from '$lib/types/opengraph';
	import Icon from '@iconify/svelte';
	import { format, formatDistance, parseISO } from 'date-fns';
	import de from 'date-fns/locale/de/index.js';
	import type { RecordModel } from 'pocketbase';
	import { onMount } from 'svelte';
	import type { MediaPlayerElement } from 'vidstack/elements.js';

	export let data;

	$: og = {
		...DefaultOpenGraph,
		title: data.vod?.title,
		image: `${PUBLIC_API_URL}/vods/${data.vod?.filename}-lg.webp`,
		updated_time: parseISO(data.vod?.date).toISOString()
	};

	let player: MediaPlayerElement;
	let currentTime: number;

	$: vod = data.vod as RecordModel;
	$: vodsCount = data.vodsCount;
	$: vodPosition = data.vodPosition;
	$: recommendations = data.recommendations;

	$: percentile = (vodPosition * 100) / vodsCount;
	$: percentileRounded = percentile < 1 ? percentile.toFixed(2) : Math.round(percentile);

	onMount(() => {
		if ($page.url.searchParams.has('t')) {
			player.currentTime = parseInt($page.url.searchParams.get('t') || '0');
		}
	});

	function copyLink(withTimestamp: boolean) {
		const url = new URL($page.url.origin + $page.url.pathname);
		if (withTimestamp) {
			url.searchParams.set('t', currentTime.toFixed(0));
		}
		navigator.clipboard.writeText(url.toString());
	}
</script>

<SEO bind:og />

<div
	class="absolute top-0 left-0 w-full h-full aspect-video bg-cover bg-center blur-2xl opacity-10 -z-10"
	style="background-image: url('{PUBLIC_API_URL}/vods/{vod.filename}-lg.webp');"
/>
<div class="max-w-[120rem] mx-auto flex flex-col xl:flex-row gap-8">
	<div class="xl:basis-4/5 flex flex-col gap-4">
		<Player bind:player bind:currentTime video={vod} />
		<h1 class="text-4xl font-bold">
			{vod.title}
		</h1>
		<div class="stats stats-vertical lg:stats-horizontal shadow bg-base-200 w-full">
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
				<div class="stat-value text-2xl flex items-center gap-1">
					#{vodPosition}
					{#if percentile <= 5}
						<div title="Top {percentileRounded}%">
							<Icon icon="solar:fire-bold-duotone" class="text-red-500 text-3xl" />
						</div>
					{:else if percentile <= 33}
						<div title="Top {percentileRounded}%">
							<Icon
								icon="solar:round-arrow-right-up-bold-duotone"
								class="text-green-500 text-3xl"
							/>
						</div>
					{:else if percentile <= 66}
						<div title="Top {percentileRounded}%">
							<Icon icon="solar:round-arrow-right-bold-duotone" class="text-slate-500" />
						</div>
					{:else}
						<div title="Top {percentileRounded}%">
							<Icon
								icon="solar:round-arrow-right-down-bold-duotone"
								class="text-red-500 text-3xl"
							/>
						</div>
					{/if}
				</div>
				<div class="stat-desc">
					von {vodsCount} Vods (Views &gt; 0)
				</div>
			</div>
			<div class="stat">
				<div class="stat-title text-lg">Auflösung</div>
				<div class="stat-value text-2xl">{vod.resolution}</div>
				<div class="stat-desc">{vod.fps} FPS</div>
			</div>
		</div>
		<div class="flex flex-col md:flex-row md:flex-wrap gap-4 justify-start">
			<div>
				<div
					title="Beim Download kann die genaue Dateigröße nicht vorhergesagt werden, weil das Video nicht als Ganzes existiert und die Videosegmente im Hintergrund zusammengesetzt werden, was nur eine grobe Schätzung ermöglicht."
				>
					<Button href="{PUBLIC_API_URL}/download/vods/{vod.id}" color="">
						<Icon icon="solar:download-square-bold-duotone" class="text-violet-500 text-2xl" /> Download
						(~ {formatBytes(vod.size)})
					</Button>
				</div>
			</div>
			<div class="md:ms-auto">
				<div class="dropdown md:dropdown-end">
					<label
						for="btn-share"
						tabindex="-1"
						class="btn shadow rounded-xl bg-gradient-to-r border-none hover:shadow-lg transition duration-200"
					>
						<Icon icon="solar:share-bold-duotone" class="text-violet-500 text-2xl" /> Teilen
					</label>
					<ul
						id="btn-share"
						tabindex="-1"
						class="dropdown-content z-[1] menu p-2 shadow bg-base-200 rounded-box"
					>
						<li>
							<button on:click={() => copyLink(false)}>Link kopieren</button>
						</li>
						<li>
							<button class="whitespace-nowrap" on:click={() => copyLink(true)}
								>Link bei {toHHMMSS(currentTime, false)} kopieren</button
							>
						</li>
					</ul>
				</div>
			</div>
		</div>
		{#if vod['expand']}
			{#if vod['expand']['clip(vod)'] && vod['expand']['clip(vod)'].length > 0}
				<h2 class="text-3xl font-bold">
					<span
						class="bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-violet-500 drop-shadow-md"
						>Clips für diesen Stream</span
					>
				</h2>
				<div class="grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					{#each vod['expand']['clip(vod)'] as video}
						<Card {video} />
					{/each}
				</div>
			{/if}
		{/if}
	</div>
	<div class="xl:basis-1/5 hidden xl:flex xl:flex-col xl:gap-4">
		<h2 class="text-3xl font-bold">
			<span
				class="bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-violet-500 drop-shadow-md"
				>Empfohlene Streams</span
			>
		</h2>
		<div class="grid grid-flow-row-dense grid-cols-1 gap-4">
			{#each recommendations as video}
				<Card {video} />
			{:else}
				Keine Ergebnisse
			{/each}
		</div>
	</div>
	<div class="xl:hidden flex flex-col gap-4">
		<h2 class="text-3xl font-bold">
			<span
				class="bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-violet-500 drop-shadow-md"
				>Empfohlene Streams</span
			>
		</h2>
		<div class="grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each recommendations as video}
				<Card {video} />
			{:else}
				Keine Ergebnisse
			{/each}
		</div>
	</div>
</div>
