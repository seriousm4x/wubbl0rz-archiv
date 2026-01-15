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
	import { onMount } from 'svelte';
	import type { MediaPlayerElement } from 'vidstack/elements';

	let { data } = $props();

	let player: MediaPlayerElement = $state({} as MediaPlayerElement);
	let currentTime: number = $state(0);

	let og = $state({
		...DefaultOpenGraph,
		title: data.clip?.title,
		image: `${PUBLIC_API_URL}/clips/${data.clip.filename}-lg.webp`,
		updated_time: parseISO(data.clip?.date).toISOString()
	});
	let clip = $derived(data.clip as RecordModel);
	let clipsCount = $derived(data.clipsCount);
	let clipPosition = $derived(data.clipPosition);
	let recommendations = $derived(data.recommendations);
	let percentile = $derived((clipPosition * 100) / clipsCount);
	let percentileRounded = $derived(percentile < 1 ? percentile.toFixed(2) : Math.round(percentile));

	onMount(() => {
		if (page.url.searchParams.has('t')) {
			player.currentTime = parseInt(page.url.searchParams.get('t') || '0');
		}
	});

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
	class="absolute top-0 left-0 -z-10 aspect-video h-full w-full bg-cover bg-center opacity-10 blur-2xl"
	style="background-image: url('{PUBLIC_API_URL}/clips/{clip.filename}-lg.webp');"
></div>
<div class="mx-auto flex max-w-480 flex-col gap-8 xl:flex-row">
	<div class="flex flex-col gap-4 xl:basis-4/5">
		<Player bind:player bind:currentTime video={clip} />
		<h1 class="text-4xl font-bold">
			{clip.title}
		</h1>
		<div class="stats stats-vertical bg-base-200 lg:stats-horizontal w-full shadow">
			<div class="stat">
				<div class="stat-title text-lg">Geclippt am</div>
				<div class="stat-value text-2xl">
					{format(parseISO(clip.date), 'dd.MM.yyyy', { locale: de })}
				</div>
				<div class="stat-desc">
					{formatDistance(parseISO(clip.date), Date.now(), {
						addSuffix: true,
						includeSeconds: true,
						locale: de
					})} um {format(parseISO(clip.date), 'HH:mm', { locale: de })} Uhr
				</div>
			</div>
			{#if clip?.expand?.creator?.name !== undefined}
				<div class="stat">
					<div class="stat-title text-lg">Erstellt von</div>
					<div class="stat-value text-2xl">
						<a
							href="https://twitch.tv/{clip.expand.creator.name}"
							target="_blank"
							class="hover:underline">{clip.expand.creator.name}</a
						>
					</div>
				</div>
			{/if}
			<div class="stat">
				<div class="stat-title text-lg">Views</div>
				<div class="stat-value text-2xl">{clip.viewcount.toLocaleString('de-DE')}</div>
				<div class="stat-desc">
					Sync {formatDistance(parseISO(clip.updated), Date.now(), {
						addSuffix: true,
						includeSeconds: true,
						locale: de
					})}
				</div>
			</div>
			<div class="stat">
				<div class="stat-title text-lg">Platz</div>
				<div class="stat-value flex items-center gap-1 text-2xl">
					#{clipPosition}
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
					von {clipsCount} Clips
				</div>
			</div>
			<div class="stat">
				<div class="stat-title text-lg">Auflösung</div>
				<div class="stat-value text-2xl">{clip.resolution}</div>
				<div class="stat-desc">{clip.fps} FPS</div>
			</div>
		</div>
		<div class="flex flex-row flex-wrap gap-4 md:justify-between">
			<div
				title="Beim Download kann die genaue Dateigröße nicht vorhergesagt werden, weil das Video nicht als Ganzes existiert und die Videosegmente im Hintergrund zusammengesetzt werden, was nur eine grobe Schätzung ermöglicht."
			>
				<a
					href="{PUBLIC_API_URL}/download/clips/{clip.id}"
					class="btn rounded-xl bg-linear-to-r shadow transition duration-200 hover:shadow-lg"
					aria-label="link"
				>
					<Icon icon="solar:download-square-bold-duotone" class="text-2xl text-violet-500" /> Download
					(~ {formatBytes(clip.size)})
				</a>
			</div>
			<div>
				<div class="dropdown dropdown-left md:dropdown-bottom">
					<label
						for="btn-share"
						tabindex="-1"
						class="btn rounded-xl border-none bg-linear-to-r shadow transition duration-200 hover:shadow-lg"
					>
						<Icon icon="solar:share-bold-duotone" class="text-2xl text-violet-500" /> Teilen
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
		</div>
		{#if clip['expand']}
			{#if clip['expand']['vod']}
				<h2 class="text-3xl font-bold">
					<span
						class="bg-linear-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent drop-shadow-md"
						>Clip stammt aus diesem Stream</span
					>
				</h2>
				<div class="grid grid-flow-row-dense grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
					<Card video={clip['expand']['vod']} offset={clip.vod_offset} />
				</div>
			{/if}
		{/if}
	</div>
	<div class="hidden xl:flex xl:basis-1/5 xl:flex-col xl:gap-4">
		<h2 class="text-3xl font-bold">
			<span
				class="bg-linear-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent drop-shadow-md"
				>Empfohlene Clips</span
			>
		</h2>
		<div class="grid grid-flow-row-dense grid-cols-1 gap-4">
			{#each recommendations as video, index (index)}
				<Card {video} />
			{:else}
				Keine Ergebnisse
			{/each}
		</div>
	</div>
	<div class="flex flex-col gap-4 xl:hidden">
		<h2 class="text-3xl font-bold">
			<span
				class="bg-linear-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent drop-shadow-md"
				>Empfohlene Clips</span
			>
		</h2>
		<div class="grid grid-flow-row-dense grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each recommendations as video, index (index)}
				<Card {video} />
			{:else}
				Keine Ergebnisse
			{/each}
		</div>
	</div>
</div>
