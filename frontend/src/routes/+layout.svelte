<script lang="ts">
	import { page } from '$app/state';
	import type { Snippet } from 'svelte';
	import Favicons from '$lib/components/Favicons.svelte';
	import Navbar from '$lib/components/Navbar.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import '../app.css';

	let { children }: { children: Snippet } = $props();
</script>

<Favicons />

<Sidebar>
	<div class="h-full w-full">
		<Navbar />
		<div class="p-4">
			{@render children()}
		</div>
	</div>
</Sidebar>

{#if page.url.hostname === 'wubbl0rz.tv'}
	<script defer data-domain="wubbl0rz.tv" src="https://metrics.wubbl0rz.tv/js/script.js"></script>
{/if}
