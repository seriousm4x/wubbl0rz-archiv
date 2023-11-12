<script lang="ts">
	import { page } from '$app/stores';
	import Icon from '@iconify/svelte';

	const sidebarHome = {
		title: 'Home',
		href: '/',
		icon: '/home.webp',
		iconInactive: 'grayscale group-hover:grayscale-0',
		bgColor: 'from-yellow-800 to-orange-800'
	};
	const sidebarItems = [
		{
			title: 'Vods',
			href: '/vods',
			icon: 'solar:videocamera-record-bold-duotone',
			iconActive: 'text-[#9146FF] group-focus:text-[#772CE8] group-active:text-[#772CE8]',
			iconInactive:
				'group-hover:text-[#9146FF] group-focus:text-[#772CE8] group-active:text-[#772CE8]',
			bgColor: 'from-pink-800 to-purple-800'
		},
		{
			title: 'Clips',
			href: '/clips',
			icon: 'solar:clapperboard-bold-duotone',
			iconActive: 'text-red-500 group-focus:text-red-600 group-active:text-red-600',
			iconInactive: 'group-hover:text-red-500 group-focus:text-red-600 group-active:text-red-600',
			bgColor: 'from-rose-800 to-orange-800'
		},
		{
			title: 'Stats',
			href: '/stats',
			icon: 'solar:pie-chart-2-bold-duotone',
			iconActive: 'text-green-500 group-focus:text-green-600 group-active:text-green-600',
			iconInactive:
				'group-hover:text-green-500 group-focus:text-green-600 group-active:text-green-600',
			bgColor: 'from-lime-800 to-teal-800'
		},
		{
			title: 'Chat',
			href: '/chat',
			icon: 'solar:keyboard-bold-duotone',
			iconActive: 'text-sky-500 group-focus:text-sky-600 group-active:text-sky-600',
			iconInactive: 'group-hover:text-sky-500 group-focus:text-sky-600 group-active:text-sky-600',
			bgColor: 'from-emerald-800 to-cyan-800'
		}
	];
</script>

<div class="drawer md:drawer-open">
	<input id="menu-drawer" type="checkbox" class="drawer-toggle" />
	<div class="drawer-content">
		<slot />
	</div>
	<div class="drawer-side z-50 shadow-xl">
		<label for="menu-drawer" aria-label="close sidebar" class="drawer-overlay" />
		<ul
			class="menu w-20 xl:w-28 min-h-full backdrop-blur-md md:backdrop-blur-none bg-base-200/80 md:bg-base-200/40 text-base-content justify-items-center gap-2"
		>
			<li class="mask mask-squircle relative group md:my-4">
				<a
					href={sidebarHome.href}
					class="p-2 flex flex-col gap-1 justify-center text-center group hover:bg-transparent"
				>
					<div
						class="absolute -inset-0.5 bg-gradient-to-r opacity-0 group-hover:opacity-[.2] transition duration-200 group-hover:duration-200 {sidebarHome.bgColor}"
					/>
					<img
						src={sidebarHome.icon}
						alt="{sidebarHome.title} Icon"
						class="w-10 h-10 mask mask-squircle group-hover:-rotate-6 duration-200 group-hover:scale-110 contrast-[1.3] saturate-[.75] {$page
							.url.pathname === sidebarHome.href
							? ''
							: sidebarHome.iconInactive}"
					/>
					<span class="text-xs font-bold uppercase transition">{sidebarHome.title}</span>
				</a>
			</li>
			<li class="menu-title md:mb-4"><hr class="opacity-25 rounded border-base-content" /></li>
			{#each sidebarItems as item}
				<li class="mask mask-squircle relative group">
					<a
						href={item.href}
						class="p-2 flex flex-col gap-0 justify-center text-center group hover:bg-transparent"
					>
						<div
							class="absolute -inset-0.5 bg-gradient-to-r opacity-0 group-hover:opacity-[.2] transition duration-200 group-hover:duration-200 {item.bgColor}"
						/>
						<Icon
							icon={item.icon}
							class="text-5xl group-hover:-rotate-6 duration-200 group-hover:scale-110 {$page.url.pathname.startsWith(
								item.href
							)
								? item.iconActive
								: item.iconInactive}"
						/>
						<span class="text-xs font-bold uppercase">{item.title}</span>
					</a>
				</li>
			{/each}
		</ul>
	</div>
</div>
