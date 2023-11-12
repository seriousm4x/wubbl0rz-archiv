<script>
	import Navbar from '$lib/components/Navbar.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import '../app.css';
	import { page } from '$app/stores';
</script>

<Sidebar>
	<div class="w-full h-full">
		<Navbar />
		<div class="p-4">
			<slot />
		</div>
	</div>
</Sidebar>

{#if $page.url.hostname === 'wubbl0rz.tv'}
	<script defer data-domain="wubbl0rz.tv" src="https://metrics.wubbl0rz.tv/js/script.js"></script>
{/if}
