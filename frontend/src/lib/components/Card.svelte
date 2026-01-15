<script lang="ts">
	import { resolve } from '$app/paths';
	import VideoThumbnail from '$lib/components/VideoThumbnail.svelte';
	import { format, formatDistance, parseISO } from 'date-fns';
	import { de } from 'date-fns/locale';
	import type { RecordModel } from 'pocketbase';

	let { video = {} as RecordModel, offset = 0 }: { video: RecordModel; offset?: number } = $props();
	const type = video.collectionName === 'vod' ? 'vods' : 'clips';
</script>

<div class="card bg-base-200 w-full overflow-hidden rounded-xl transition hover:shadow-lg">
	<VideoThumbnail {video} {offset} />
	<div class="card-body justify-between gap-1 p-3">
		<a href={resolve(`/${type}/${video.id}${offset > 0 ? `?t=${offset}` : ''}`)}>
			<h2 class="text-md font-bold">{video.title}</h2>
		</a>
		<div class="card-actions text-base-content/70 mt-4 justify-between text-sm font-semibold">
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
