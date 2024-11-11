<script lang="ts">
	import { page } from '$app/stores';
	import type { OpenGraph } from '$lib/types/opengraph';

	let { og = $bindable() }: { og: OpenGraph } = $props();
</script>

<svelte:head>
	<meta name="description" content={og.description} />
	<meta property="og:title" content={og.title} />
	<meta property="og:description" content={og.description} />
	<meta property="og:url" content={$page.url.toString()} />
	<meta property="og:updated_time" content={og.updated_time} />
	<meta property="og:image" content={og.image} />
	<meta property="og:image:url" content={og.image} />
	<meta property="og:image:secure_url" content={og.image} />
	<meta property="og:image:width" content="1200" />
	<meta property="og:image:height" content="630" />
	<meta property="og:image:alt" content={og.title} />
	<meta name="twitter:title" content={og.title} />
	<meta name="twitter:description" content={og.description} />
	<meta name="twitter:image" content={og.image} />
	<meta name="twitter:image:alt" content={og.title} />
	<title>{og.title}</title>
</svelte:head>
