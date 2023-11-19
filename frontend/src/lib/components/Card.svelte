<script lang="ts">
	import VideoThumbnail from '$lib/components/VideoThumbnail.svelte';
	import { format, formatDistance, parseISO } from 'date-fns';
	import de from 'date-fns/locale/de/index.js';
	import type { RecordModel } from 'pocketbase';

	export let video: RecordModel = {} as RecordModel;
	export let offset = 0;
	const type = video.collectionName === 'vod' ? 'vods' : 'clips';
</script>

<div class="card w-full overflow-hidden rounded-xl bg-base-200 transition hover:shadow-lg">
	<VideoThumbnail {video} {offset} />
	<div class="card-body justify-between gap-1 p-3">
		<a href="/{type}/{video.id}{offset > 0 ? `?t=${offset}` : ''}">
			<h2 class="text-md font-bold">{video.title}</h2>
		</a>
		<div class="card-actions mt-4 justify-between text-sm font-semibold text-base-content/70">
			{video.viewcount.toLocaleString('de-DE')} Aufrufe
			<span title={format(parseISO(video.date), "dd.MM.yyyy 'um' HH:mm:ss")}>
				{formatDistance(parseISO(video.date), Date.now(), {
					addSuffix: true,
					includeSeconds: true,
					locale: de
				})}
			</span>
		</div>
	</div>
</div>
