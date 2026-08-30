<script lang="ts">
	import { resolve } from '$app/paths';
	import VideoThumbnail from '$lib/components/VideoThumbnail.svelte';
	import { format, formatDistance, parseISO } from 'date-fns';
	import { de } from 'date-fns/locale';
	import type { RecordModel } from 'pocketbase';

	let { video = {} as RecordModel, offset = 0 }: { video: RecordModel; offset?: number } = $props();
	let type = $derived<'vods' | 'clips'>(video.collectionName === 'vod' ? 'vods' : 'clips');
</script>

<article class="group min-w-0">
	<VideoThumbnail {video} {offset} />
	<div class="space-y-1 px-1 pt-3">
		<a class="block" href={resolve(`/${type}/${video.id}${offset > 0 ? `?t=${offset}` : ''}`)}>
			<h2 class="line-clamp-2 text-sm leading-5 font-semibold">
				{video.title}
			</h2>
		</a>
		<div class="text-base-content/65 flex flex-wrap gap-x-2 text-xs">
			<span>{video.viewcount.toLocaleString('de-DE')} Aufrufe</span>
			<span title={format(parseISO(video.date), "dd.MM.yyyy 'um' HH:mm:ss")}>
				{formatDistance(parseISO(video.date), Date.now(), {
					addSuffix: true,
					includeSeconds: true,
					locale: de
				})}
			</span>
		</div>
	</div>
</article>
