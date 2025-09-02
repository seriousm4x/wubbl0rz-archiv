<script lang="ts">
	import { page } from '$app/state';
	import { PUBLIC_API_URL } from '$env/static/public';
	import Card from '$lib/components/Card.svelte';
	import Player from '$lib/components/Player.svelte';
	import SEO from '$lib/components/SEO.svelte';
	import { formatBytes, toHHMMSS } from '$lib/functions';
	import { DefaultOpenGraph } from '$lib/types/opengraph';
	import Icon from '@iconify/svelte';
	import { format, formatDistance, parseISO } from 'date-fns';
	import { de } from 'date-fns/locale';
	import type { RecordModel } from 'pocketbase';
	import type { MediaPlayerElement } from 'vidstack/elements';

	let { data } = $props();

	let og = $state({
		...DefaultOpenGraph,
		title: data.vod?.title,
		image: `${PUBLIC_API_URL}/vods/${data.vod?.filename}-lg.webp`,
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

	function copyLink(withTimestamp: boolean) {
		const url = new URL(page.url.origin + page.url.pathname);
		if (withTimestamp) {
			url.searchParams.set('t', currentTime.toFixed(0));
		}
		navigator.clipboard.writeText(url.toString());
	}
</script>

<SEO bind:og />

<div
	class="absolute left-0 top-0 -z-10 aspect-video h-full w-full bg-cover bg-center opacity-10 blur-2xl"
	style="background-image: url('{PUBLIC_API_URL}/vods/{vod.filename}-lg.webp');"
></div>
<div class="mx-auto flex max-w-[120rem] flex-col gap-8 xl:flex-row">
	<div class="flex flex-col gap-4 xl:basis-4/5">
		<Player bind:player bind:currentTime video={vod} />
		<h1 class="text-4xl font-bold">
			{vod.title}
		</h1>
		<div class="stats stats-vertical w-full bg-base-200 shadow lg:stats-horizontal">
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
							<Icon icon="solar:fire-bold-duotone" class="text-3xl text-red-500" />
						</div>
					{:else if percentile <= 33}
						<div title="Top {percentileRounded}%">
							<Icon
								icon="solar:round-arrow-right-up-bold-duotone"
								class="text-3xl text-green-500"
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
								class="text-3xl text-red-500"
							/>
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
		<div class="flex flex-col justify-start gap-4 md:flex-row md:flex-wrap">
			<div>
				<div
					title="Beim Download kann die genaue Dateigröße nicht vorhergesagt werden, weil das Video nicht als Ganzes existiert und die Videosegmente im Hintergrund zusammengesetzt werden, was nur eine grobe Schätzung ermöglicht."
				>
					<a
						href="{PUBLIC_API_URL}/download/vods/{vod.id}"
						class="btn rounded-xl bg-gradient-to-r shadow transition duration-200 hover:shadow-lg"
						aria-label="link"
					>
						<Icon icon="solar:download-square-bold-duotone" class="text-2xl text-violet-500" /> Download
						(~ {formatBytes(vod.size)})
					</a>
				</div>
			</div>
			<div class="md:ms-auto">
				<div class="dropdown md:dropdown-end">
					<label
						for="btn-share"
						tabindex="-1"
						class="btn rounded-xl border-none bg-gradient-to-r shadow transition duration-200 hover:shadow-lg"
					>
						<Icon icon="solar:share-bold-duotone" class="text-2xl text-violet-500" /> Teilen
					</label>
					<ul
						id="btn-share"
						tabindex="-1"
						class="menu dropdown-content z-[1] rounded-box bg-base-200 p-2 shadow"
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
		</div>
		{#if vod['expand']}
			{#if vod['expand']['clip_via_vod'] && vod['expand']['clip_via_vod'].length > 0}
				<h2 class="text-3xl font-bold">
					<span
						class="bg-gradient-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent drop-shadow-md"
						>Clips für diesen Stream</span
					>
				</h2>
				<div class="grid grid-flow-row-dense grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
					{#each vod['expand']['clip_via_vod'] as video}
						<Card {video} />
					{/each}
				</div>
			{/if}
		{/if}
	</div>
	<div class="hidden xl:flex xl:basis-1/5 xl:flex-col xl:gap-4">
		<h2 class="text-3xl font-bold">
			<span
				class="bg-gradient-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent drop-shadow-md"
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
	<div class="flex flex-col gap-4 xl:hidden">
		<h2 class="text-3xl font-bold">
			<span
				class="bg-gradient-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent drop-shadow-md"
				>Empfohlene Streams</span
			>
		</h2>
		<div class="grid grid-flow-row-dense grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each recommendations as video}
				<Card {video} />
			{:else}
				Keine Ergebnisse
			{/each}
		</div>
	</div>
</div>
