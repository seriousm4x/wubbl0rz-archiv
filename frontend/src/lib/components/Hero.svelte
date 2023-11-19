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
		class="hero relative aspect-video max-h-[45vh] min-h-[20rem] max-w-[90rem] place-items-center justify-items-end overflow-hidden rounded-xl border-2 border-base-100/20 shadow-md transition duration-200 hover:shadow-2xl"
		role="none"
		on:mouseenter={() => (hover = true)}
		on:mouseleave={() => (hover = false)}
	>
		<a href="/vods/{vod.id}" class="h-full w-full">
			{#if hover}
				<video
					muted
					loop
					autoplay
					transition:fade={{ duration: 250 }}
					class="absolute left-0 top-0 z-10 h-auto w-full"
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
				class="h-full w-full bg-cover bg-center"
				role="none"
				style="background-image: url('{PUBLIC_API_URL}/vods/{vod.filename}-lg.webp');"
			></div>
		</a>
		<div class="hero-content p-0">
			<div class="card me-4 w-fit max-w-xl bg-base-100/80 shadow-xl backdrop-blur-sm xl:me-16">
				<div class="card-body p-5">
					<h1 class="card-title text-2xl font-extrabold xl:text-4xl">
						<span class="bg-gradient-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent"
							>Aktueller Stream</span
						>
					</h1>
					<h2 class="card-title line-clamp-3 text-lg font-bold xl:text-2xl">
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
						<a class="btn btn-primary btn-sm border-none hover:bg-purple-700" href="/vods/{vod.id}"
							>Jetzt ansehen
							<Icon icon="solar:alt-arrow-right-bold-duotone" class="text-3xl" />
						</a>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
