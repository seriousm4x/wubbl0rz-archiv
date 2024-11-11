<script lang="ts">
	import GridAllVideos from '$lib/components/GridAllVideos.svelte';
	import SEO from '$lib/components/SEO.svelte';
	import { DefaultOpenGraph } from '$lib/types/opengraph';
	import { parseISO } from 'date-fns';
	import type { ListResult, RecordModel } from 'pocketbase';

	let { data }: { data: ListResult<RecordModel> } = $props();

	let og = $state({
		...DefaultOpenGraph,
		title: 'Alle Streams',
		updated_time:
			data.items.length > 0 ? parseISO(data.items[0].date).toISOString() : new Date().toISOString()
	});
</script>

<SEO bind:og />

<GridAllVideos {data} title="Alle Streams" placeholder="Streamtitel" />
