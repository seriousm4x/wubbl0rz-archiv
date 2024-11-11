<script lang="ts">
	import SEO from '$lib/components/SEO.svelte';
	import { getEmotes } from '$lib/emotes';
	import { replaceEmotesInString } from '$lib/functions';
	import { pb } from '$lib/stores/pocketbase';
	import { DefaultOpenGraph } from '$lib/types/opengraph';
	import Icon from '@iconify/svelte';
	import { format, formatDistance, parseISO } from 'date-fns';
	import { de } from 'date-fns/locale';
	import type { ListResult, RecordModel } from 'pocketbase';
	import { onMount } from 'svelte';

	let { data }: { data: ListResult<RecordModel> } = $props();
	let sliderValue = $state(200);

	let og = $state({
		...DefaultOpenGraph,
		title: 'Livechat'
	});

	onMount(async () => {
		$pb.collection('chatmessage').subscribe('*', function (e) {
			data.items = [e.record, ...data.items];
			if (sliderValue === 500) return;
			data.items = data.items.slice(0, sliderValue);
		});
	});
</script>

<SEO bind:og />

<div class="container mx-auto">
	<h1 class="mb-4 text-4xl font-bold">
		<span
			class="bg-gradient-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent drop-shadow-md"
			>Livechat</span
		>
	</h1>
	<p class="mb-8 text-sm text-base-content/80">Neue Nachrichten werden automatisch geladen...</p>
	<h2 class="text-xl">Zu behaltende Nachrichten</h2>
	<div class="flex flex-row flex-wrap gap-4">
		<input
			type="range"
			class="range"
			min="100"
			max="500"
			step="100"
			bind:value={sliderValue}
			onchange={() => (data.items = data.items.slice(0, sliderValue))}
		/>
		<div class="flex w-full justify-between px-2 text-xs">
			<span>100</span>
			<span>200</span>
			<span>300</span>
			<span>400</span>
			<span>Unendlich</span>
		</div>
	</div>
	<div class="mt-8 flex flex-col gap-2">
		{#await getEmotes()}
			<div class="w-full text-center">
				<span class="loading loading-spinner loading-lg"></span>
			</div>
		{:then [emotes, re]}
			{#each data.items as message}
				<div class="chat chat-start rounded-sm transition duration-100">
					<div class="chat-header">
						<span
							class="font-semibold drop-shadow-md"
							style="color: {message?.tags['color'] || 'inherit'}">{message.user_display_name}</span
						>
						<div
							class="inline-block"
							title={format(parseISO(message.date), "dd.MM.yyyy 'um' HH:mm:ss")}
						>
							<time class="text-xs opacity-50">
								{formatDistance(parseISO(message.date), Date.now(), {
									addSuffix: true,
									includeSeconds: true,
									locale: de
								})}
							</time>
						</div>
					</div>
					{#if 'tags' in message && 'reply-parent-msg-body' in message.tags}
						<div
							id={message.tags['id']}
							class="chat-bubble flex scroll-mt-24 flex-col gap-1 transition duration-100 [overflow-wrap:anywhere] hover:shadow-lg"
						>
							<a
								href={`#${message.tags['reply-parent-msg-id']}`}
								class="flex w-fit flex-row items-center gap-1 rounded-lg bg-gray-100 px-2 py-1 text-xs font-bold text-base-content/80 dark:bg-base-200"
							>
								<div>
									<Icon icon="solar:chat-square-arrow-bold-duotone" class="text-lg text-primary" />
								</div>
								{message.tags['reply-parent-display-name']}: {message.tags['reply-parent-msg-body']}
							</a>
							<div class="flex flex-row flex-wrap items-center gap-1 [overflow-wrap:anywhere]">
								<!-- eslint-disable-next-line svelte/no-at-html-tags -->
								{@html replaceEmotesInString(message.message, emotes, re)}
							</div>
						</div>
					{:else}
						<div
							id={message.tags['id']}
							class="chat-bubble flex scroll-mt-24 flex-row flex-wrap items-center gap-1 text-slate-100 transition duration-100 [overflow-wrap:anywhere] hover:shadow-lg"
						>
							<!-- eslint-disable-next-line svelte/no-at-html-tags -->
							{@html replaceEmotesInString(message.message, emotes, re)}
						</div>
					{/if}
				</div>
			{/each}
		{:catch}
			{#each data.items as message}
				<div class="chat chat-start rounded-sm transition duration-100">
					<div class="chat-header">
						<span
							class="font-semibold drop-shadow-md"
							style="color: {message?.tags['color'] || 'inherit'}">{message.user_display_name}</span
						>
						<div
							class="inline-block"
							title={format(parseISO(message.date), "dd.MM.yyyy 'um' HH:mm:ss")}
						>
							<time class="text-xs opacity-50">
								{formatDistance(parseISO(message.date), Date.now(), {
									addSuffix: true,
									includeSeconds: true,
									locale: de
								})}
							</time>
						</div>
					</div>
					<div
						id={message.tags['id']}
						class="chat-bubble flex scroll-mt-24 flex-row flex-wrap items-center gap-1 text-slate-100 transition duration-100 [overflow-wrap:anywhere] hover:shadow-lg"
					>
						{message.message}
					</div>
				</div>
			{/each}
		{/await}
	</div>
</div>
