<script>
	import Navbar from '$lib/components/Navbar.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import '../app.css';
</script>

<Sidebar>
	<div class="w-full h-full">
		<Navbar />
		<div class="p-4">
			<slot />
		</div>
	</div>
</Sidebar>
