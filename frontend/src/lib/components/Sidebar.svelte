<script lang="ts">
	import { page } from '$app/state';
	import IconClapperboardBoldDuotone from '@iconify-icons/solar/clapperboard-bold-duotone';
	import IconKeyboardBoldDuotone from '@iconify-icons/solar/keyboard-bold-duotone';
	import IconPieChart2BoldDuotone from '@iconify-icons/solar/pie-chart-2-bold-duotone';
	import IconVideocameraRecordBoldDuotone from '@iconify-icons/solar/videocamera-record-bold-duotone';
	import Icon from '@iconify/svelte';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();

	const sidebarHome = { title: 'Home', href: '/', icon: '/home.webp' };
	const sidebarItems = [
		{ title: 'Vods', href: '/vods', icon: IconVideocameraRecordBoldDuotone },
		{ title: 'Clips', href: '/clips', icon: IconClapperboardBoldDuotone },
		{ title: 'Stats', href: '/stats', icon: IconPieChart2BoldDuotone },
		{ title: 'Chat', href: '/chat', icon: IconKeyboardBoldDuotone }
	];
</script>

<div class="drawer md:drawer-open">
	<input id="menu-drawer" type="checkbox" class="drawer-toggle" />
	<div class="drawer-content">
		{@render children()}
	</div>
	<div class="drawer-side border-base-content/10 z-50 border-r">
		<label for="menu-drawer" aria-label="close sidebar" class="drawer-overlay"></label>
		<ul class="menu bg-base-100 text-base-content min-h-full w-20 gap-1 p-2">
			<li>
				<a
					href={sidebarHome.href}
					class="flex flex-col items-center gap-1 rounded-lg px-1 py-2 text-[10px] font-medium {page
						.url.pathname === sidebarHome.href
						? 'bg-base-200 text-base-content'
						: 'text-base-content/65 hover:bg-base-200 hover:text-base-content'}"
				>
					<img src={sidebarHome.icon} alt="" class="h-7 w-7 rounded-full object-cover" />
					<span>{sidebarHome.title}</span>
				</a>
			</li>
			{#each sidebarItems as item (item.href)}
				<li>
					<a
						href={item.href}
						class="flex flex-col items-center gap-1 rounded-lg px-1 py-2 text-[10px] font-medium {page.url.pathname.startsWith(
							item.href
						)
							? 'bg-base-200 text-base-content'
							: 'text-base-content/65 hover:bg-base-200 hover:text-base-content'}"
					>
						<Icon icon={item.icon} class="text-2xl" />
						<span>{item.title}</span>
					</a>
				</li>
			{/each}
		</ul>
	</div>
</div>
