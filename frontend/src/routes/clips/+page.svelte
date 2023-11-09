<script lang="ts">
	import GridAllVideos from '$lib/components/GridAllVideos.svelte';
	import SEO from '$lib/components/SEO.svelte';
	import { DefaultOpenGraph } from '$lib/types/opengraph';
	import { parseISO } from 'date-fns';
	import type { ListResult, RecordModel } from 'pocketbase';

	export let data: ListResult<RecordModel>;

	$: og = {
		...DefaultOpenGraph,
		title: 'Alle Clips',
		updated_time: parseISO(data.items[0]?.date).toISOString()
	};
</script>

<SEO bind:og />

<GridAllVideos {data} title="Alle Clips" placeholder="Cliptitel" />
