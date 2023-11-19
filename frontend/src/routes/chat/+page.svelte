<script lang="ts">
	import SEO from '$lib/components/SEO.svelte';
	import { getEmotes } from '$lib/emotes';
	import { replaceEmotesInString } from '$lib/functions';
	import { pb } from '$lib/pocketbase';
	import { DefaultOpenGraph } from '$lib/types/opengraph';
	import Icon from '@iconify/svelte';
	import { format, formatDistance, parseISO } from 'date-fns';
	import de from 'date-fns/locale/de/index.js';
	import type { ListResult, RecordModel } from 'pocketbase';
	import { onMount } from 'svelte';

	export let data: ListResult<RecordModel>;
	let sliderValue = 200;

	$: og = {
		...DefaultOpenGraph,
		title: 'Livechat'
	};

	onMount(async () => {
		pb.collection('chatmessage').subscribe('*', function (e) {
			data.items = [e.record, ...data.items];
			if (sliderValue === 500) return;
			data.items = data.items.slice(0, sliderValue);
		});
	});
</script>

<SEO bind:og />

<div class="container mx-auto">
	<h1 class="text-4xl font-bold mb-4">
		<span
			class="bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-violet-500 drop-shadow-md"
			>Livechat</span
		>
	</h1>
	<p class="text-sm text-base-content/80 mb-8">Neue Nachrichten werden automatisch geladen...</p>
	<h2 class="text-xl">Zu behaltende Nachrichten</h2>
	<div class="flex flex-row flex-wrap gap-4">
		<input
			type="range"
			class="range"
			min="100"
			max="500"
			step="100"
			bind:value={sliderValue}
			on:change={() => (data.items = data.items.slice(0, sliderValue))}
		/>
		<div class="w-full flex justify-between text-xs px-2">
			<span>100</span>
			<span>200</span>
			<span>300</span>
			<span>400</span>
			<span>Unendlich</span>
		</div>
	</div>
	<div class="flex flex-col gap-2 mt-8">
		{#await getEmotes()}
			<div class="w-full text-center">
				<span class="loading loading-spinner loading-lg"></span>
			</div>
		{:then [emotes, re]}
			{#each data.items as message}
				<div class="chat chat-start transition duration-100 rounded-sm">
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
							class="chat-bubble bg-gray-200 dark:bg-base-300 text-base-content transition duration-100 hover:shadow-lg scroll-mt-24 flex flex-col gap-1 [overflow-wrap:anywhere]"
						>
							<a
								href={`#${message.tags['reply-parent-msg-id']}`}
								class="text-xs font-bold bg-gray-100 dark:bg-base-200 text-base-content/80 px-2 py-1 rounded-lg w-fit flex flex-row items-center gap-1"
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
							class="chat-bubble bg-gray-200 dark:bg-base-300 text-base-content transition duration-100 hover:shadow-lg scroll-mt-24 flex flex-row flex-wrap items-center gap-1 [overflow-wrap:anywhere]"
						>
							<!-- eslint-disable-next-line svelte/no-at-html-tags -->
							{@html replaceEmotesInString(message.message, emotes, re)}
						</div>
					{/if}
				</div>
			{/each}
		{:catch}
			{#each data.items as message}
				<div class="chat chat-start transition duration-100 rounded-sm">
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
						class="chat-bubble bg-gray-200 dark:bg-base-300 text-base-content transition duration-100 hover:shadow-lg scroll-mt-24 flex flex-row flex-wrap items-center gap-1 [overflow-wrap:anywhere]"
					>
						{message.message}
					</div>
				</div>
			{/each}
		{/await}
	</div>
</div>
