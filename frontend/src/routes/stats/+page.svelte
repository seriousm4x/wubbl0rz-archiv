<script lang="ts">
	import SEO from '$lib/components/SEO.svelte';
	import { formatBytes } from '$lib/functions';
	import { DefaultOpenGraph } from '$lib/types/opengraph';
	import Icon from '@iconify/svelte';
	import { formatRelative, parseISO } from 'date-fns';
	import de from 'date-fns/locale/de/index.js';
	import type { RecordModel } from 'pocketbase';

	export let data;

	$: og = {
		...DefaultOpenGraph,
		title: 'Statistik',
		description: 'Übersicht der Metadaten des Archivs.',
		updated_time: parseISO(data.stats.last_update).toISOString()
	};

	const sevenTv = data.emotes
		.filter((emote: RecordModel) => emote.provider === '7tv')
		.sort((a: RecordModel, b: RecordModel) =>
			a.name.localeCompare(b.name, 'de', { sensitivity: 'base', numeric: true })
		);
	const bttv = data.emotes
		.filter((emote: RecordModel) => emote.provider === 'bttv')
		.sort((a: RecordModel, b: RecordModel) =>
			a.name.localeCompare(b.name, 'de', { sensitivity: 'base', numeric: true })
		);
	const twitch = data.emotes
		.filter((emote: RecordModel) => emote.provider === 'twitch')
		.sort((a: RecordModel, b: RecordModel) =>
			a.name.localeCompare(b.name, 'de', { sensitivity: 'base', numeric: true })
		);
	const ffz = data.emotes
		.filter((emote: RecordModel) => emote.provider === 'ffz')
		.sort((a: RecordModel, b: RecordModel) =>
			a.name.localeCompare(b.name, 'de', { sensitivity: 'base', numeric: true })
		);
</script>

<SEO bind:og />

<div class="container mx-auto">
	<div class="flex flex-col gap-4">
		<h1 class="text-4xl flex flex-row gap-4 drop-shadow-md">
			<img src="/pocketbase.svg" alt="pocketbase" class="h-10" />
			<span>
				Pocket<span class="font-bold">Base</span>
			</span>
		</h1>
		<div class="stats stats-vertical lg:stats-horizontal shadow bg-base-200 w-full">
			<div class="stat">
				<div class="stat-figure text-primary">
					<Icon icon="solar:videocamera-record-bold-duotone" class="text-5xl" />
				</div>
				<div class="stat-title">Vods</div>
				<div class="stat-value text-primary">{data.stats.count_vods.toLocaleString('de-DE')}</div>
				<div class="stat-desc">
					30 Tage Trend:
					{#if data.stats.trend_vods === 0}
						{data.stats.trend_vods} Vods
					{:else if data.stats.trend_vods > 0}
						+{data.stats.trend_vods} Vod{data.stats.trend_vods === 1 ? '' : 's'}
					{:else}
						-{data.stats.trend_vods} Vod{data.stats.trend_vods === -1 ? '' : 's'}
					{/if}
				</div>
			</div>
			<div class="stat">
				<div class="stat-figure text-red-500">
					<Icon icon="solar:clapperboard-bold-duotone" class="text-5xl" />
				</div>
				<div class="stat-title">Clips</div>
				<div class="stat-value text-red-500">{data.stats.count_clips.toLocaleString('de-DE')}</div>
				<div class="stat-desc">
					30 Tage Trend:
					{#if data.stats.trens_clips === 0}
						{data.stats.trend_clips} Clips
					{:else if data.stats.trend_clips > 0}
						+{data.stats.trend_clips} Clip{data.stats.trend_clips === 1 ? '' : 's'}
					{:else}
						-{data.stats.trend_clips} Clip{data.stats.trend_clips === 1 ? '' : 's'}
					{/if}
				</div>
			</div>
			<div class="stat">
				<div class="stat-figure text-amber-500">
					<Icon icon="solar:clock-circle-bold-duotone" class="text-5xl" />
				</div>
				<div class="stat-title">Stunden gestreamt</div>
				<div class="stat-value text-amber-500">
					{data.stats.count_hours.toLocaleString('de-DE')}
				</div>
				<div class="stat-desc">
					30 Tage Trend:
					{#if data.stats.trend_hours === 0}
						{data.stats.trend_hours} Stunden
					{:else if data.stats.trend_hours > 0}
						+ {data.stats.trend_hours} Stunde{data.stats.trend_hours === 1 ? '' : 'n'}
					{:else}
						- {data.stats.trend_hours} Stunde{data.stats.trend_hours === 1 ? '' : 'n'}
					{/if}
				</div>
			</div>
			<div class="stat">
				<div class="stat-figure text-green-500">
					<Icon icon="solar:pie-chart-2-bold-duotone" class="text-5xl" />
				</div>
				<div class="stat-title">Archivgröße</div>
				<div class="stat-value text-green-500">{formatBytes(data.stats.count_size)}</div>
				<div class="stat-desc">Vods und Clips gemeinsam</div>
			</div>
		</div>
		<h1 class="text-4xl font-bold flex flex-row gap-4 mt-8 drop-shadow-md">
			<img src="/meilisearch.svg" alt="meilisearch" class="h-10" />
			<span>
				meili<span class="font-light">search</span>
			</span>
		</h1>
		<div class="stats stats-vertical lg:stats-horizontal shadow bg-base-200 w-full">
			<div class="stat">
				<div class="stat-figure text-emerald-500">
					<Icon icon="solar:subtitles-bold-duotone" class="text-5xl" />
				</div>
				<div class="stat-title">Transkripte</div>
				<div class="stat-value text-emerald-500">
					{data.meili.transcripts.toLocaleString('de-DE')}
				</div>
				<div class="stat-desc">Einzelne Textzeilen</div>
			</div>
			<div class="stat">
				<div class="stat-figure text-yellow-500">
					<Icon icon="solar:text-bold-duotone" class="text-5xl" />
				</div>
				<div class="stat-title">Streamtitel</div>
				<div class="stat-value text-yellow-500">
					{data.meili.title.toLocaleString('de-DE')}
				</div>
			</div>
			<div class="stat">
				<div class="stat-figure text-red-400">
					<Icon icon="solar:smartphone-update-bold-duotone" class="text-5xl" />
				</div>
				<div class="stat-title">Letztes Update</div>
				<div class="stat-value text-red-400 whitespace-break-spaces">
					{formatRelative(parseISO(data.meili.lastUpdate), Date.now(), { locale: de })}
				</div>
			</div>
			<div class="stat">
				<div class="stat-figure text-violet-500">
					<Icon icon="solar:pie-chart-2-bold-duotone" class="text-5xl" />
				</div>
				<div class="stat-title">Datenbankgröße</div>
				<div class="stat-value text-violet-500">{formatBytes(data.meili.databaseSize)}</div>
				<div class="stat-desc">Transkripte und Vods</div>
			</div>
		</div>
		<h1 class="text-4xl font-bold mt-8">
			<span
				class="bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-violet-500 drop-shadow-md"
				>Top Chatter</span
			>
		</h1>
		<div class="overflow-x-auto">
			<table class="table table-sm">
				<!-- head -->
				<thead>
					<tr>
						<th>Name</th>
						<th>Nachrichten</th>
					</tr>
				</thead>
				<tbody>
					{#each data.stats.chatters as chatter}
						<tr class="hover">
							<td>{chatter.name}</td>
							<td>{chatter.msg_count.toLocaleString('de-DE')}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
		<h1 class="text-4xl font-bold mt-8 flex flex-row flex-wrap items-center gap-4 drop-shadow-md">
			<span class="bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-violet-500">
				Emotes
			</span>
			<span class="badge badge-primary">{data.emotes.length}</span>
		</h1>
		<h1 class="text-2xl font-bold mt-4 flex flex-row flex-wrap items-center gap-4">
			<span>7tv</span>
			<span class="badge badge-neutral">{sevenTv.length}</span>
		</h1>
		<div class="flex flex-row flex-wrap gap-2">
			{#each sevenTv as emote}
				<div title="{emote.name} (7tv)">
					<img
						src={emote.url}
						alt={emote.name}
						class="h-14"
						loading="lazy"
						title="{emote.name} (7tv)"
					/>
				</div>
			{:else}
				Keine Emotes
			{/each}
		</div>
		<h1 class="text-2xl font-bold mt-4 flex flex-row flex-wrap items-center gap-4">
			<span>BetterTTV</span>
			<span class="badge badge-neutral">{bttv.length}</span>
		</h1>
		<div class="flex flex-row flex-wrap gap-2">
			{#each bttv as emote}
				<div title={emote.name}>
					<img src={emote.url} alt={emote.name} class="h-14" loading="lazy" />
				</div>
			{:else}
				Keine Emotes
			{/each}
		</div>
		<h1 class="text-2xl font-bold mt-4 flex flex-row flex-wrap items-center gap-4">
			<span>FrankerFaceZ</span>
			<span class="badge badge-neutral">{ffz.length}</span>
		</h1>
		<div class="flex flex-row flex-wrap gap-2">
			{#each ffz as emote}
				<div title={emote.name}>
					<img src={emote.url} alt={emote.name} class="h-14" loading="lazy" />
				</div>
			{:else}
				Keine Emotes
			{/each}
		</div>
		<h1 class="text-2xl font-bold mt-4 flex flex-row flex-wrap items-center gap-4">
			<span>Twitch</span>
			<span class="badge badge-neutral">{twitch.length}</span>
		</h1>
		<div class="flex flex-row flex-wrap gap-2">
			{#each twitch as emote}
				<div title={emote.name}>
					<img src={emote.url} alt={emote.name} class="h-14" loading="lazy" />
				</div>
			{:else}
				Keine Emotes
			{/each}
		</div>
	</div>
</div>
