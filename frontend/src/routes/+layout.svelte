<script>
	import { beforeNavigate } from '$app/navigation';
	import { page, updated } from '$app/state';
	import Favicons from '$lib/components/Favicons.svelte';
	import Navbar from '$lib/components/Navbar.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import '../app.css';

	beforeNavigate(({ willUnload, to }) => {
		if (updated.current && !willUnload && to?.url) {
			location.href = to.url.href;
		}
	});
</script>

<Favicons />

<Sidebar>
	<div class="h-full w-full">
		<Navbar />
		<div class="p-4">
			<slot />
		</div>
	</div>
</Sidebar>

{#if page.url.hostname === 'wubbl0rz.tv'}
	<script defer data-domain="wubbl0rz.tv" src="https://metrics.wubbl0rz.tv/js/script.js"></script>
{/if}
