<script lang="ts">
	import { resolve } from '$app/paths';
	import { PUBLIC_API_URL } from '$env/static/public';
	import IconAltArrowRightBoldDuotone from '@iconify-icons/solar/alt-arrow-right-bold-duotone';
	import Icon from '@iconify/svelte';
	import { formatDistance, parseISO } from 'date-fns';
	import { de } from 'date-fns/locale';
	import type { RecordModel } from 'pocketbase';
	import { fade } from 'svelte/transition';

	let { vod }: { vod: RecordModel } = $props();
	let hover = $state(false);
</script>

<section class="mx-auto max-w-6xl space-y-4" aria-labelledby="featured-vod-title">
	<a
		href={resolve(`/vods/${vod.id}`)}
		class="card group rounded-box relative block aspect-21/9 overflow-hidden shadow-sm"
		onmouseenter={() => (hover = true)}
		onmouseleave={() => (hover = false)}
	>
		{#if hover}
			<video
				muted
				loop
				autoplay
				playsinline
				transition:fade={{ duration: 200 }}
				class="absolute inset-0 z-10 h-full w-full object-cover"
			>
				<source
					src="{PUBLIC_API_URL}/vods/{vod.filename}/preview.webm"
					type="video/webm;codecs=vp9"
				/>
				<source
					src="{PUBLIC_API_URL}/vods/{vod.filename}/preview.mp4"
					type="video/mp4; codecs=hvc1"
				/>
				<track kind="captions" />
			</video>
		{/if}
		<img
			src="{PUBLIC_API_URL}/vods/{vod.filename}/thumb-lg.webp"
			alt="Vorschaubild: {vod.title}"
			class="h-full w-full object-cover transition duration-300 group-hover:scale-[1.01]"
			width="1536"
			height="860"
		/>
	</a>

	<div
		class="border-base-content/10 flex flex-col gap-4 border-b pb-6 sm:flex-row sm:items-end sm:justify-between"
	>
		<div class="space-y-2">
			<p class="text-base-content/60 text-sm font-medium tracking-wide uppercase">Neuster Stream</p>
			<h1
				id="featured-vod-title"
				class="text-2xl leading-tight font-semibold tracking-tight sm:text-4xl"
			>
				<a class="link link-hover" href={resolve(`/vods/${vod.id}`)}>{vod.title}</a>
			</h1>
			<p class="text-base-content/70">
				{vod.viewcount.toLocaleString('de-DE')} Aufrufe
				<span class="mx-1.5" aria-hidden="true">•</span>
				{formatDistance(parseISO(vod.date), Date.now(), {
					addSuffix: true,
					includeSeconds: true,
					locale: de
				})}
			</p>
		</div>
		<a
			class="btn border-base-content/15 bg-base-content text-base-100 hover:bg-base-content/85 sm:shrink-0"
			href={resolve(`/vods/${vod.id}`)}
		>
			Jetzt ansehen
			<Icon icon={IconAltArrowRightBoldDuotone} class="text-xl" />
		</a>
	</div>
</section>
