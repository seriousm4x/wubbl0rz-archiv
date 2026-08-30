<script lang="ts">
	import SEO from '$lib/components/SEO.svelte';
	import { formatBytes } from '$lib/functions';
	import { DefaultOpenGraph } from '$lib/types/opengraph';
	import IconClapperboardBoldDuotone from '@iconify-icons/solar/clapperboard-bold-duotone';
	import IconClockCircleBoldDuotone from '@iconify-icons/solar/clock-circle-bold-duotone';
	import IconPieChart2BoldDuotone from '@iconify-icons/solar/pie-chart-2-bold-duotone';
	import IconRoundArrowRightUpBoldDuotone from '@iconify-icons/solar/round-arrow-right-up-bold-duotone';
	import IconSmartphoneUpdateBoldDuotone from '@iconify-icons/solar/smartphone-update-bold-duotone';
	import IconTextBoldDuotone from '@iconify-icons/solar/text-bold-duotone';
	import IconVideocameraRecordBoldDuotone from '@iconify-icons/solar/videocamera-record-bold-duotone';
	import Icon from '@iconify/svelte';
	import { formatRelative, parseISO, formatDistanceToNow, format } from 'date-fns';
	import { de } from 'date-fns/locale';
	import type { RecordModel } from 'pocketbase';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let og = $derived({
		...DefaultOpenGraph,
		title: 'Statistik',
		description: 'Übersicht der Metadaten des Archivs.',
		updated_time: parseISO(data.stats.last_update).toISOString()
	});

	let sevenTv = $derived(
		data.emotes
			.filter((emote: RecordModel) => emote.provider === '7tv')
			.sort((a: RecordModel, b: RecordModel) =>
				a.name.localeCompare(b.name, 'de', { sensitivity: 'base', numeric: true })
			)
	);
	let bttv = $derived(
		data.emotes
			.filter((emote: RecordModel) => emote.provider === 'bttv')
			.sort((a: RecordModel, b: RecordModel) =>
				a.name.localeCompare(b.name, 'de', { sensitivity: 'base', numeric: true })
			)
	);
	let twitch = $derived(
		data.emotes
			.filter((emote: RecordModel) => emote.provider === 'twitch')
			.sort((a: RecordModel, b: RecordModel) =>
				a.name.localeCompare(b.name, 'de', { sensitivity: 'base', numeric: true })
			)
	);
	let ffz = $derived(
		data.emotes
			.filter((emote: RecordModel) => emote.provider === 'ffz')
			.sort((a: RecordModel, b: RecordModel) =>
				a.name.localeCompare(b.name, 'de', { sensitivity: 'base', numeric: true })
			)
	);
</script>

<SEO {og} />

<div class="container mx-auto">
	<div class="flex flex-col gap-4">
		<h1 class="flex flex-row gap-4 text-2xl font-semibold tracking-tight">
			<img src="/pocketbase.svg" alt="pocketbase" class="h-10" />
			<span>
				Pocket<span class="font-bold">Base</span>
			</span>
		</h1>
		<div class="stats stats-vertical bg-base-200 xl:stats-horizontal w-full shadow">
			<div class="stat">
				<div class="stat-figure text-primary">
					<Icon icon={IconVideocameraRecordBoldDuotone} class="text-5xl" />
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
						{data.stats.trend_vods} Vod{data.stats.trend_vods === -1 ? '' : 's'}
					{/if}
				</div>
			</div>
			<div class="stat">
				<div class="stat-figure text-red-500">
					<Icon icon={IconClapperboardBoldDuotone} class="text-5xl" />
				</div>
				<div class="stat-title">Clips</div>
				<div class="stat-value text-red-500">{data.stats.count_clips.toLocaleString('de-DE')}</div>
				<div class="stat-desc">
					30 Tage Trend:
					{#if data.stats.trend_clips === 0}
						{data.stats.trend_clips} Clips
					{:else if data.stats.trend_clips > 0}
						+{data.stats.trend_clips} Clip{data.stats.trend_clips === 1 ? '' : 's'}
					{:else}
						{data.stats.trend_clips} Clip{data.stats.trend_clips === 1 ? '' : 's'}
					{/if}
				</div>
			</div>
			<div class="stat">
				<div class="stat-figure text-amber-500">
					<Icon icon={IconClockCircleBoldDuotone} class="text-5xl" />
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
						+{data.stats.trend_hours} Stunde{data.stats.trend_hours === 1 ? '' : 'n'}
					{:else}
						{data.stats.trend_hours} Stunde{data.stats.trend_hours === 1 ? '' : 'n'}
					{/if}
				</div>
			</div>
			<div class="stat">
				<div class="stat-figure text-green-500">
					<Icon icon={IconPieChart2BoldDuotone} class="text-5xl" />
				</div>
				<div class="stat-title">Archivgröße</div>
				<div class="stat-value text-green-500">{formatBytes(data.stats.count_size)}</div>
				<div class="stat-desc">Vods und Clips gemeinsam</div>
			</div>
		</div>
		<h2 class="mt-8 flex flex-row gap-4 text-2xl font-semibold tracking-tight">
			<img src="/meilisearch.svg" alt="meilisearch" class="h-10" />
			<span>
				meili<span class="font-light">search</span>
			</span>
		</h2>
		<div class="stats stats-vertical bg-base-200 xl:stats-horizontal w-full shadow">
			<div class="stat">
				<div class="stat-figure text-yellow-500">
					<Icon icon={IconTextBoldDuotone} class="text-5xl" />
				</div>
				<div class="stat-title">Streamtitel</div>
				<div class="stat-value text-yellow-500">
					{data.meili.title.toLocaleString('de-DE')}
				</div>
			</div>
			<div class="stat">
				<div class="stat-figure text-emerald-500">
					<Icon icon={IconRoundArrowRightUpBoldDuotone} class="text-5xl" />
				</div>
				<div class="stat-title">Transkripte</div>
				<div class="stat-value text-emerald-500">
					{data.meili.transcripts.toLocaleString('de-DE')}
				</div>
				<div class="stat-desc">Einzelne Textzeilen</div>
			</div>
			<div class="stat">
				<div class="stat-figure text-red-400">
					<Icon icon={IconSmartphoneUpdateBoldDuotone} class="text-5xl" />
				</div>
				<div class="stat-title">Letztes Update</div>
				<div class="stat-value whitespace-break-spaces text-red-400">
					{formatRelative(parseISO(data.meili.lastUpdate), Date.now(), { locale: de })}
				</div>
			</div>
			<div class="stat">
				<div class="stat-figure text-violet-500">
					<Icon icon={IconPieChart2BoldDuotone} class="text-5xl" />
				</div>
				<div class="stat-title">Datenbankgröße</div>
				<div class="stat-value text-violet-500">{formatBytes(data.meili.databaseSize)}</div>
				<div class="stat-desc">Transkripte und Vods</div>
			</div>
		</div>
		<h2 class="mt-8 text-2xl font-semibold tracking-tight">
			<span class="text-base-content">Top Chatter</span>
		</h2>
		<div class="overflow-x-auto">
			<table class="table-sm table">
				<thead>
					<tr>
						<th>#</th>
						<th>Name</th>
						<th>Nachrichten</th>
						<th>Letzte Nachricht</th>
					</tr>
				</thead>
				<tbody>
					{#each data.stats.chatters as chatter, index (index)}
						<tr class="hover">
							<td>{index + 1}.</td>
							<td>{chatter.name}</td>
							<td>{chatter.msg_count.toLocaleString('de-DE')}</td>
							<td title={format(chatter.updated, "dd.MM.yyyy 'um' HH:mm:ss")}
								>vor {formatDistanceToNow(chatter.updated, {
									locale: de
								})}</td
							>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
		<h2
			class="mt-8 flex flex-row flex-wrap items-center gap-3 text-2xl font-semibold tracking-tight"
		>
			<span class="text-base-content">Emotes</span>
			<span class="badge badge-neutral">{data.emotes.length}</span>
		</h2>
		<h3 class="mt-4 flex flex-row flex-wrap items-center gap-3 text-lg font-semibold">
			<span>7tv</span>
			<span class="badge badge-neutral">{sevenTv.length}</span>
		</h3>
		<div class="flex flex-row flex-wrap gap-2">
			{#each sevenTv as emote, index (index)}
				<div title={emote.name}>
					<img
						src={emote.url}
						alt={emote.name}
						class="h-10 sm:h-14"
						loading="lazy"
						title={emote.name}
					/>
				</div>
			{:else}
				Keine Emotes
			{/each}
		</div>
		<h3 class="mt-4 flex flex-row flex-wrap items-center gap-3 text-lg font-semibold">
			<span>BetterTTV</span>
			<span class="badge badge-neutral">{bttv.length}</span>
		</h3>
		<div class="flex flex-row flex-wrap gap-2">
			{#each bttv as emote, index (index)}
				<div title={emote.name}>
					<img src={emote.url} alt={emote.name} class="h-10 sm:h-14" loading="lazy" />
				</div>
			{:else}
				Keine Emotes
			{/each}
		</div>
		<h3 class="mt-4 flex flex-row flex-wrap items-center gap-3 text-lg font-semibold">
			<span>FrankerFaceZ</span>
			<span class="badge badge-neutral">{ffz.length}</span>
		</h3>
		<div class="flex flex-row flex-wrap gap-2">
			{#each ffz as emote, index (index)}
				<div title={emote.name}>
					<img src={emote.url} alt={emote.name} class="h-10 sm:h-14" loading="lazy" />
				</div>
			{:else}
				Keine Emotes
			{/each}
		</div>
		<h3 class="mt-4 flex flex-row flex-wrap items-center gap-3 text-lg font-semibold">
			<span>Twitch</span>
			<span class="badge badge-neutral">{twitch.length}</span>
		</h3>
		<div class="flex flex-row flex-wrap gap-2">
			{#each twitch as emote, index (index)}
				<div title={emote.name}>
					<img src={emote.url} alt={emote.name} class="h-10 sm:h-14" loading="lazy" />
				</div>
			{:else}
				Keine Emotes
			{/each}
		</div>
	</div>
</div>
