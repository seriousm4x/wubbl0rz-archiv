<script lang="ts">
	import { getEmotes } from '$lib/emotes';
	import { replaceEmotesInString } from '$lib/functions';
	import { pb } from '$lib/stores/pocketbase';
	import { getTwitchBadges, type TwitchBadges } from '$lib/twitch-badges';
	import { format, parseISO } from 'date-fns';
	import type { RecordModel } from 'pocketbase';

	let dialog: HTMLDialogElement;
	let messageContainer: HTMLDivElement;
	let messages = $state.raw<RecordModel[]>([]);
	let userName = $state('');
	let isLoading = $state(false);
	let isLoadingMore = $state(false);
	let hasLoadError = $state(false);
	let hasMore = $state(true);
	let page = $state(0);
	let emotesPromise = getEmotes();
	let badgesPromise = $state<Promise<TwitchBadges>>(Promise.resolve({}));

	export function show(requestedUserName: string) {
		userName = requestedUserName;
		messages = [];
		hasLoadError = false;
		hasMore = true;
		page = 0;
		badgesPromise = Promise.resolve({});
		dialog.showModal();
		loadMore(requestedUserName);
	}

	function loadMore(requestedUserName = userName) {
		if (isLoading || isLoadingMore || !hasMore) return;

		const requestedPage = page + 1;
		if (requestedPage === 1) isLoading = true;
		else isLoadingMore = true;

		void $pb
			.collection('chatmessage')
			.getList(requestedPage, 100, {
				filter: $pb.filter('user_name = {:userName}', { userName: requestedUserName }),
				sort: '-date'
			})
			.then((result) => {
				if (userName !== requestedUserName) return;

				page = requestedPage;
				hasMore = result.page < result.totalPages;
				messages = [...messages, ...result.items];
				const roomId = String(messages.at(0)?.tags?.['room-id'] || '');
				if (roomId) badgesPromise = getTwitchBadges(roomId);
			})
			.catch(() => {
				if (userName === requestedUserName) hasLoadError = true;
			})
			.finally(() => {
				if (userName !== requestedUserName) return;

				isLoading = false;
				isLoadingMore = false;
			});
	}

	function handleScroll() {
		const remainingScroll =
			messageContainer.scrollHeight - messageContainer.scrollTop - messageContainer.clientHeight;
		if (remainingScroll < 24) loadMore();
	}

	function getBadgeKeys(message: RecordModel): string[] {
		const badges = message.tags?.badges;
		return typeof badges === 'string' && badges ? badges.split(',') : [];
	}

	function getDate(message: RecordModel) {
		return format(parseISO(message.date), 'dd.MM.yyyy');
	}
</script>

<dialog class="modal backdrop-blur-md" bind:this={dialog} aria-labelledby="chatter-history-title">
	<div class="modal-box bg-base-200 max-h-[80dvh] max-w-2xl p-0">
		<div class="bg-base-300/50 flex items-center justify-between px-4 py-3">
			<h2 id="chatter-history-title" class="text-base font-semibold">Nachrichten von {userName}</h2>
			<form method="dialog">
				<button class="btn btn-ghost btn-sm btn-square" aria-label="Schließen">×</button>
			</form>
		</div>
		<div
			class="max-h-[calc(80dvh-3.5rem)] overflow-y-auto p-3"
			bind:this={messageContainer}
			onscroll={handleScroll}
		>
			{#if isLoading}
				<div class="flex justify-center py-8">
					<span class="loading loading-spinner loading-sm"></span>
				</div>
			{:else if hasLoadError}
				<p class="text-base-content/60 text-center text-sm">
					Nachrichten konnten nicht geladen werden.
				</p>
			{:else if messages.length === 0}
				<p class="text-base-content/60 text-center text-sm">Keine Nachrichten gefunden.</p>
			{:else}
				{#await emotesPromise then [emotes, regex]}
					{#await badgesPromise then badges}
						{#each messages as message, index (message.id)}
							{#if index === 0 || getDate(message) !== getDate(messages[index - 1])}
								<div
									class="border-base-content/15 my-2 border-t pt-2 text-center font-mono text-xs"
								>
									{getDate(message)}
								</div>
							{/if}
							<p
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
								<span class="me-1 font-semibold" style="color: {message.tags?.color || 'inherit'}">
									{message.user_display_name}
								</span>:
								<!-- eslint-disable-next-line svelte/no-at-html-tags -->
								{@html replaceEmotesInString(message.message, emotes, regex)}
							</p>
						{/each}
					{/await}
				{/await}
				{#if isLoadingMore}
					<div class="flex justify-center py-3">
						<span class="loading loading-spinner loading-sm"></span>
					</div>
				{/if}
			{/if}
		</div>
	</div>
	<form method="dialog" class="modal-backdrop"><button>Schließen</button></form>
</dialog>
