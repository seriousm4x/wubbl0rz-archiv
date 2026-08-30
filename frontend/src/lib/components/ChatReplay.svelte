<script lang="ts">
	import { getEmotes, type Emotes } from '$lib/emotes';
	import { replaceEmotesInString, toHHMMSS } from '$lib/functions';
	import { pb } from '$lib/stores/pocketbase';
	import { getTwitchBadges, type TwitchBadges } from '$lib/twitch-badges';
	import type { RecordModel } from 'pocketbase';
	import { onMount, tick } from 'svelte';
	import { SvelteMap, SvelteSet } from 'svelte/reactivity';

	const WINDOW_SECONDS = 60 * 2;
	const HISTORY_SECONDS = 10 * 60;
	const RETAINED_WINDOWS = 100;
	const EMPTY_BADGES: TwitchBadges = {};
	const pendingEmotes = new Promise<[Emotes, RegExp]>(() => {});

	let {
		vod,
		currentTime,
		playerHeight = 0
	}: { vod: RecordModel; currentTime: number; playerHeight?: number } = $props();
	let messages = $state.raw<RecordModel[]>([]);
	let isLoading = $state(false);
	let hasLoadError = $state(false);
	let shouldAutoScroll = $state(true);
	let messageContainer: HTMLDivElement;
	let emotesPromise = $state<Promise<[Emotes, RegExp]>>(pendingEmotes);
	let badgesPromise = $state<Promise<TwitchBadges>>(Promise.resolve(EMPTY_BADGES));

	const loadedWindows = new SvelteSet<number>();
	let loadedVodId = '';

	let startTime = $derived(Date.parse(vod.date));
	let duration = $derived(Number(vod.duration) || 0);
	let playbackTimestamp = $derived(startTime + Math.max(0, currentTime) * 1000);
	let currentWindow = $derived(Math.floor(Math.max(0, currentTime) / WINDOW_SECONDS));
	let visibleMessages = $derived(
		messages.filter((message) => Date.parse(message.date) <= playbackTimestamp)
	);
	let latestVisibleMessageId = $derived(visibleMessages.at(-1)?.id);
	let roomId = $derived(String(messages.at(0)?.tags?.['room-id'] || ''));

	onMount(() => {
		emotesPromise = getEmotes();
	});

	function getBadgeKeys(message: RecordModel): string[] {
		const badges = message.tags?.badges;
		return typeof badges === 'string' && badges ? badges.split(',') : [];
	}

	function toPocketBaseDate(timestamp: number) {
		return new Date(timestamp)
			.toISOString()
			.replace('T', ' ')
			.replace(/\.\d{3}Z$/, '');
	}

	function handleChatScroll() {
		const remainingScroll =
			messageContainer.scrollHeight - messageContainer.scrollTop - messageContainer.clientHeight;
		shouldAutoScroll = remainingScroll < 24;
	}

	$effect(() => {
		if (roomId) badgesPromise = getTwitchBadges(roomId);
	});

	$effect(() => {
		const requestedVodId = vod.id;
		const requestedWindow = currentWindow;

		if (!Number.isFinite(startTime) || duration <= 0) return;

		if (loadedVodId !== requestedVodId) {
			loadedVodId = requestedVodId;
			loadedWindows.clear();
			messages = [];
		}

		if (loadedWindows.has(requestedWindow)) return;

		const windowStart = Math.max(0, requestedWindow * WINDOW_SECONDS - HISTORY_SECONDS);
		const windowEnd = Math.min(duration, (requestedWindow + 1) * WINDOW_SECONDS);
		if (windowStart >= windowEnd) return;

		loadedWindows.add(requestedWindow);
		isLoading = true;
		hasLoadError = false;

		const start = toPocketBaseDate(startTime + windowStart * 1000);
		const end = toPocketBaseDate(startTime + windowEnd * 1000);

		void $pb
			.collection('chatmessage')
			.getFullList({
				batch: 500,
				filter: $pb.filter('date >= {:start} && date < {:end}', { start, end }),
				sort: 'date'
			})
			.then((receivedMessages) => {
				if (vod.id !== requestedVodId) return;

				const messagesById = new SvelteMap(messages.map((message) => [message.id, message]));
				receivedMessages.forEach((message) => messagesById.set(message.id, message));
				const oldestRetainedTimestamp =
					startTime +
					Math.max(0, (requestedWindow - RETAINED_WINDOWS) * WINDOW_SECONDS - HISTORY_SECONDS) *
						1000;
				messages = [...messagesById.values()]
					.filter((message) => Date.parse(message.date) >= oldestRetainedTimestamp)
					.sort((a, b) => Date.parse(a.date) - Date.parse(b.date));

				for (const loadedWindow of loadedWindows) {
					if (loadedWindow < requestedWindow - RETAINED_WINDOWS) loadedWindows.delete(loadedWindow);
				}
			})
			.catch(() => {
				loadedWindows.delete(requestedWindow);
				hasLoadError = true;
			})
			.finally(() => {
				if (vod.id === requestedVodId) isLoading = false;
			});
	});

	$effect(() => {
		if (!messageContainer || !shouldAutoScroll || !latestVisibleMessageId) return;

		void tick().then(() => {
			if (shouldAutoScroll) messageContainer.scrollTop = messageContainer.scrollHeight;
		});
	});
</script>

<section
	class="bg-base-200 flex h-full min-h-0 flex-col overflow-hidden shadow xl:h-(--player-height)"
	style:--player-height={playerHeight ? `${playerHeight}px` : undefined}
	aria-labelledby="chat-replay-title"
>
	<div class="bg-base-300/50 flex items-center justify-between px-4 py-3">
		<h2 id="chat-replay-title" class="text-base font-semibold">Chat-Wiederholung</h2>
		<span class="text-base-content/60 text-xs">{toHHMMSS(currentTime, false)}</span>
	</div>
	<div
		class="min-h-0 flex-1 overflow-x-hidden overflow-y-auto p-3"
		bind:this={messageContainer}
		onscroll={handleChatScroll}
		aria-live="polite"
	>
		{#await emotesPromise}
			<div class="flex h-full items-center justify-center">
				<span class="loading loading-spinner loading-sm"></span>
			</div>
		{:then [emotes, regex]}
			{#if visibleMessages.length > 0}
				{#await badgesPromise then badges}
					{#each visibleMessages as message (message.id)}
						<p class="mb-1 text-sm leading-5 wrap-break-word hyphens-auto">
							<span class="text-base-content/45 me-1 font-mono text-xs">
								{toHHMMSS((Date.parse(message.date) - startTime) / 1000, false)}
							</span>
							<span class="me-1 inline-flex align-middle">
								{#each getBadgeKeys(message) as badgeKey (badgeKey)}
									{#if badges[badgeKey]}
										<img
											src={badges[badgeKey].image_url_1x}
											alt={badges[badgeKey].title}
											title={badges[badgeKey].title}
											class="inline-block h-4.5 w-4.5"
											loading="lazy"
										/>
									{/if}
								{/each}
							</span>
							<span class="me-1 font-semibold" style="color: {message.tags?.color || 'inherit'}">
								{message.user_display_name}:
							</span>
							<!-- eslint-disable-next-line svelte/no-at-html-tags -->
							{@html replaceEmotesInString(message.message, emotes, regex)}
						</p>
					{/each}
				{/await}
			{:else if isLoading}
				<div class="flex h-full items-center justify-center">
					<span class="loading loading-spinner loading-sm"></span>
				</div>
			{:else if hasLoadError}
				<p class="text-base-content/60 text-center text-sm">Chat konnte nicht geladen werden.</p>
			{:else}
				<p class="text-base-content/60 text-center text-sm">Noch keine Nachrichten.</p>
			{/if}
		{/await}
	</div>
</section>
