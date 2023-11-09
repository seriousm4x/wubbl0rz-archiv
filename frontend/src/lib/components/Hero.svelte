<script lang="ts">
	import { PUBLIC_API_URL } from '$env/static/public';
	import Icon from '@iconify/svelte';
	import { formatDistance, parseISO } from 'date-fns';
	import de from 'date-fns/locale/de/index.js';
	import type { RecordModel } from 'pocketbase';
	import { fade } from 'svelte/transition';

	export let vod: RecordModel;
	let hover = false;
</script>

<div class="flex items-center justify-center">
	<div
		class="hero aspect-video min-h-[20rem] max-h-[45vh] max-w-[90rem] relative overflow-hidden rounded-xl shadow-md place-items-center justify-items-end hover:shadow-2xl duration-200 transition border-base-100/20 border-2"
		role="none"
		on:mouseenter={() => (hover = true)}
		on:mouseleave={() => (hover = false)}
	>
		<a href="/vods/{vod.id}" class="w-full h-full">
			{#if hover}
				<video
					muted
					loop
					autoplay
					transition:fade={{ duration: 250 }}
					class="absolute top-0 left-0 w-full h-auto z-10"
				>
					<source
						src="{PUBLIC_API_URL}/vods/{vod.filename}-preview.webm"
						type="video/webm;codecs=vp9"
					/>
					<source
						src="{PUBLIC_API_URL}/vods/{vod.filename}-preview.mp4"
						type="video/mp4; codecs=hvc1"
					/>
					<track kind="captions" />
				</video>
			{/if}

			<div
				class="w-full h-full bg-cover bg-center"
				role="none"
				style="background-image: url('{PUBLIC_API_URL}/vods/{vod.filename}-lg.webp');"
			></div>
		</a>
		<div class="hero-content p-0">
			<div class="card w-fit max-w-xl bg-base-100/80 shadow-xl me-4 xl:me-16 backdrop-blur-sm">
				<div class="card-body p-5">
					<h1 class="card-title text-2xl xl:text-4xl font-extrabold">
						<span class="bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-violet-500"
							>Aktueller Stream</span
						>
					</h1>
					<h2 class="card-title text-lg xl:text-2xl font-bold line-clamp-3">
						{vod.title}
					</h2>
					<p>
						{vod.viewcount.toLocaleString('de-DE')} Aufrufe <span class="mx-1">•</span>
						{formatDistance(parseISO(vod.date), Date.now(), {
							addSuffix: true,
							includeSeconds: true,
							locale: de
						})}
					</p>
					<div class="card-actions justify-end">
						<a class="btn btn-primary btn-sm hover:bg-purple-700 border-none" href="/vods/{vod.id}"
							>Jetzt ansehen
							<Icon icon="solar:alt-arrow-right-bold-duotone" class="text-3xl" />
						</a>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
