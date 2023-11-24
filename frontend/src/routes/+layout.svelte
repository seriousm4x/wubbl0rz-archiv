<script>
	import { page } from '$app/stores';
	import Favicons from '$lib/components/Favicons.svelte';
	import Navbar from '$lib/components/Navbar.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import { pb } from '$lib/stores/pocketbase';
	import { onMount } from 'svelte';
	import '../app.css';
	import { currentUser } from '$lib/stores/user';

	onMount(async () => {
		// load cookie from localstorage
		const pbCookie = localStorage.getItem('pocketbase_auth');
		if (!pbCookie) {
			return;
		}

		// set cookie as auth
		$pb.authStore.loadFromCookie('pb_auth=' + pbCookie);
		try {
			$pb.authStore.isValid && (await $pb.collection('users').authRefresh({ autoCancel: false }));
		} catch (_) {
			$pb.authStore.clear();
		}

		// update localstorage and user store on auth change
		$pb.authStore.onChange((_, model) => {
			currentUser.set(model);
		}, true);
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

{#if $page.url.hostname === 'wubbl0rz.tv'}
	<script defer data-domain="wubbl0rz.tv" src="https://metrics.wubbl0rz.tv/js/script.js"></script>
{/if}
