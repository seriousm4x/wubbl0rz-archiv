<script lang="ts">
	import { resolve } from '$app/paths';
	import VideoThumbnail from '$lib/components/VideoThumbnail.svelte';
	import { toHHMMSS } from '$lib/functions';
	import { formatDistance, fromUnixTime } from 'date-fns';
	import { de } from 'date-fns/locale';
	import type { Hit } from 'meilisearch';
	import type { RecordModel } from 'pocketbase';

	let {
		hit = {} as Hit,
		searchIn = 'transcripts'
	}: { hit: Hit; searchIn: 'transcripts' | 'vods' } = $props();

	const video = hit as RecordModel;
	const offset = searchIn === 'transcripts' ? hit.start : 0;
</script>

<div class="card bg-base-200 w-full overflow-hidden rounded-xl transition hover:shadow-lg">
	<VideoThumbnail {video} isVod={true} {offset} />
	<div class="card-body gap-1 p-4">
		<a href={resolve(`/vods/${video.id}${offset > 0 ? `?t=${offset}` : ''}`)}>
			{#if searchIn === 'vods'}
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				<h2 class="text-md matched-string font-bold">{@html hit?._formatted?.title}</h2>
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
				class="matched-string bg-base-300 h-full rounded-lg p-2 font-mono text-sm font-medium hover:shadow-md"
				href={resolve(`/vods/${video.id}${offset > 0 ? `?t=${offset}` : ''}`)}
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
		color: var(--color-warning-content);
		background-color: var(--color-warning);
	}
</style>
