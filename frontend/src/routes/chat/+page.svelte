<script lang="ts">
	import SEO from '$lib/components/SEO.svelte';
	import ChatterHistory from '$lib/components/ChatterHistory.svelte';
	import { getEmotes } from '$lib/emotes';
	import { replaceEmotesInString } from '$lib/functions';
	import { pb } from '$lib/stores/pocketbase';
	import { getTwitchBadges, type TwitchBadges } from '$lib/twitch-badges';
	import { DefaultOpenGraph } from '$lib/types/opengraph';
	import { format, parseISO } from 'date-fns';
	import type { ListResult, RecordModel } from 'pocketbase';
	import { onMount, tick } from 'svelte';

	let { data }: { data: ListResult<RecordModel> } = $props();
	let realtimeMessages = $state.raw<RecordModel[]>([]);
	let messages = $derived(realtimeMessages.length > 0 ? realtimeMessages : data.items.toReversed());
	let emotesPromise = getEmotes();
	let badgesPromise = $state<Promise<TwitchBadges>>(Promise.resolve({}));
	let roomId = $derived(String(messages.at(0)?.tags?.['room-id'] || ''));
	let lastMessageId = $derived(messages.at(-1)?.id);
	let shouldAutoScroll = $state(true);
	let chatterHistory: { show: (userName: string) => void };

	let og = $state({
		...DefaultOpenGraph,
		title: 'Livechat'
	});

	onMount(() => {
		$pb.collection('chatmessage').subscribe('*', (e) => {
			realtimeMessages = [...messages, e.record].slice(-1000);
		});
	});

	$effect(() => {
		if (roomId) badgesPromise = getTwitchBadges(roomId);
	});

	function getBadgeKeys(message: RecordModel): string[] {
		const badges = message.tags?.badges;
		return typeof badges === 'string' && badges ? badges.split(',') : [];
	}

	function handleScroll() {
		shouldAutoScroll =
			window.innerHeight + window.scrollY >= document.documentElement.scrollHeight - 24;
	}

	$effect(() => {
		if (!shouldAutoScroll || !lastMessageId) return;

		void Promise.allSettled([emotesPromise, badgesPromise])
			.then(() => tick())
			.then(() => {
				if (shouldAutoScroll) window.scrollTo({ top: document.documentElement.scrollHeight });
			});
	});
</script>

<SEO {og} />
<svelte:window onscroll={handleScroll} />
<ChatterHistory bind:this={chatterHistory} />

<div class="container mx-auto">
	<h1 class="mb-4 text-2xl font-semibold tracking-tight">
		<span class="text-base-content tracking-tight">Livechat</span>
	</h1>
	<p class="text-base-content/80 mb-8 text-sm">Neue Nachrichten werden automatisch geladen...</p>
	<div class="message-list">
		{#await emotesPromise}
			<div class="w-full text-center">
				<span class="loading loading-spinner loading-lg"></span>
			</div>
		{:then [emotes, re]}
			{#await badgesPromise then badges}
				{#each messages as message (message.id)}
					<p
						id={message.tags?.id}
						class={[
							'mb-1 text-sm leading-5 wrap-break-word hyphens-auto',
							message.tags?.['first-msg'] === '1' &&
								'relative border-r-4 border-fuchsia-500 bg-[#422342] px-2 py-1 pr-36'
						]}
					>
						{#if message.tags?.['first-msg'] === '1'}
							<span class="absolute top-1 right-2 text-xs font-semibold text-fuchsia-500">
								FIRST MESSAGE
							</span>
						{/if}
						<time
							class="text-base-content/45 me-1 font-mono text-xs"
							title={format(parseISO(message.date), "dd.MM.yyyy 'um' HH:mm:ss")}
						>
							{format(parseISO(message.date), 'HH:mm:ss')}
						</time>
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
						<button
							class="me-1 cursor-pointer border-0 bg-transparent p-0 font-semibold hover:underline"
							style="color: {message.tags?.color || 'inherit'}"
							onclick={() => chatterHistory.show(message.user_name)}
						>
							{message.user_display_name}
						</button>:
						<!-- eslint-disable-next-line svelte/no-at-html-tags -->
						{@html replaceEmotesInString(message.message, emotes, re)}
					</p>
				{/each}
			{/await}
		{:catch}
			{#each messages as message (message.id)}
				<p
					id={message.tags?.id}
					class={[
						'mb-1 text-sm leading-5 wrap-break-word hyphens-auto',
						message.tags?.['first-msg'] === '1' &&
							'relative border-r-4 border-fuchsia-500 bg-[#422342] px-2 py-1 pr-36'
					]}
				>
					{#if message.tags?.['first-msg'] === '1'}
						<span class="absolute top-1 right-2 text-xs font-semibold text-fuchsia-500">
							FIRST MESSAGE
						</span>
					{/if}
					<time
						class="text-base-content/45 me-1 font-mono text-xs"
						title={format(parseISO(message.date), "dd.MM.yyyy 'um' HH:mm:ss")}
					>
						{format(parseISO(message.date), 'HH:mm:ss')}
					</time>
					<button
						class="me-1 cursor-pointer border-0 bg-transparent p-0 font-semibold hover:underline"
						style="color: {message.tags?.color || 'inherit'}"
						onclick={() => chatterHistory.show(message.user_name)}
					>
						{message.user_display_name}
					</button>:
					{message.message}
				</p>
			{/each}
		{/await}
	</div>
</div>
