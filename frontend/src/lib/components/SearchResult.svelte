<script lang="ts">
	import VideoThumbnail from '$lib/components/VideoThumbnail.svelte';
	import { toHHMMSS } from '$lib/functions';
	import { formatDistance, fromUnixTime } from 'date-fns';
	import de from 'date-fns/locale/de/index.js';
	import type { Hit } from 'meilisearch';
	import type { RecordModel } from 'pocketbase';

	export let hit: Hit = {} as Hit;
	export let searchIn: 'transcripts' | 'vods' = 'transcripts';
	const video = hit as RecordModel;
	const offset = searchIn === 'transcripts' ? hit.start : 0;
</script>

<div class="card w-full rounded-xl bg-base-200 hover:shadow-lg transition overflow-hidden">
	<VideoThumbnail {video} isVod={true} {offset} />
	<div class="card-body p-4 gap-1">
		<a href="/{searchIn}/{video.id}{offset > 0 ? `?t=${offset}` : ''}">
			{#if searchIn === 'vods'}
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				<h2 class="text-md font-bold matched-string">{@html hit?._formatted?.title}</h2>
			{:else}
				<h2 class="text-md font-bold">{hit.title}</h2>
			{/if}
		</a>
		<p class="text-sm font-normal">
			{hit.viewcount.toLocaleString('de-DE')} Aufrufe <span class="mx-1">•</span>
			{formatDistance(fromUnixTime(hit.date), Date.now(), {
				addSuffix: true,
				includeSeconds: true,
				locale: de
			})}
		</p>
		{#if searchIn === 'transcripts'}
			<a
				class="text-sm font-medium font-mono bg-base-300 h-full p-2 matched-string rounded-lg hover:shadow-md"
				href="/vods/{hit.id}{offset > 0 ? `?t=${offset}` : ''}"
			>
				<span class="underline">{toHHMMSS(hit.start, false)}</span>:
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				{@html hit?._formatted?.text}
			</a>
		{/if}
	</div>
</div>

<style global>
	.matched-string :global(span:not(.underline)) {
		font-weight: bold;
		border-radius: 5px;
		padding: 0.1rem;
		color: white;
		background-color: rgb(139 92 246 / var(--tw-bg-opacity));
	}
</style>
