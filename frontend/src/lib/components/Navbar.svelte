<script lang="ts">
	import { navigating } from '$app/stores';
	import { PUBLIC_MEILI_SEARCH_KEY, PUBLIC_MEILI_URL } from '$env/static/public';
	import Pagination from '$lib/components/Pagination.svelte';
	import SearchResult from '$lib/components/SearchResult.svelte';
	import Icon from '@iconify/svelte';
	import { MeiliSearch } from 'meilisearch';
	import { onMount } from 'svelte';

	let scrollY: number = $state(0);
	let modal: HTMLDialogElement;
	let searchText: string = $state('');
	let meiliIndex: 'transcripts' | 'vods' = $state('transcripts');
	let currentPage = $state(1);
	let body: HTMLBodyElement;

	const client = new MeiliSearch({
		host: PUBLIC_MEILI_URL,
		apiKey: PUBLIC_MEILI_SEARCH_KEY
	});

	// available sorts
	type sort = {
		value: string;
		text: string;
	};
	let sorts: sort[] = [
		{
			value: 'relevancy',
			text: 'Relevanz'
		},
		{
			value: 'date',
			text: 'Datum'
		},
		{
			value: 'viewcount',
			text: 'Views'
		},
		{
			value: 'duration',
			text: 'Dauer'
		}
	];
	let ordering: 'asc' | 'desc' = $state('desc');
	let selectedSort = $state(sorts[0]);

	$effect(() => {
		if (meiliIndex || searchText) currentPage = 1;
	});
	$effect(() => {
		if ($navigating && modal) modal.close();
	});

	let searchConfig = $derived({
		page: currentPage,
		hitsPerPage: 30,
		attributesToHighlight: ['*'],
		highlightPreTag: '<span>',
		highlightPostTag: '</span>',
		sort: selectedSort.value === 'relevancy' ? [] : [`${selectedSort.value}:${ordering}`]
	});

	onMount(() => {
		body = document.body as HTMLBodyElement;
	});

	function showModal() {
		body.classList.add('overflow-hidden');
		modal.show();
	}

	function onModalClose() {
		body.classList.remove('overflow-hidden');
	}

	function onKeyDown(e: KeyboardEvent) {
		switch (e.key) {
			case '/':
				e.preventDefault();
				showModal();
				break;
			case 'Escape':
				e.preventDefault();
				modal.close();
				break;
			default:
				break;
		}
	}
</script>

<svelte:window bind:scrollY on:keydown={onKeyDown} />

<div
	class="navbar sticky top-0 z-40 justify-center backdrop-blur-sm transition duration-200 {scrollY >
	0
		? 'bg-base-200/70'
		: ''}"
>
	<label
		for="menu-drawer"
		class="btn btn-ghost mask drawer-button mask-squircle flex-nowrap p-0 md:hidden"
	>
		<Icon icon="solar:menu-dots-square-bold-duotone" class="text-primary text-5xl" />
	</label>
	<button
		class="input border-base-content/20 bg-base-300/50 text-base-content/50 hover:border-base-content/50 hover:bg-base-300/80 hover:text-base-content w-full cursor-pointer rounded-full drop-shadow-md transition duration-200 md:max-w-lg"
		onclick={showModal}
	>
		<span>Suchen</span>
		<span class="ms-auto"><kbd class="kbd">/</kbd></span>
	</button>
</div>
<dialog
	class="modal items-start overflow-hidden backdrop-blur-md transition duration-200"
	bind:this={modal}
	onclose={onModalClose}
>
	<div
		class="modal-box border-base-content/20 bg-base-300/90 flex h-full max-h-full max-w-6xl flex-col items-center border p-0 backdrop-blur-md md:mt-20 md:h-4/5"
	>
		<div class="flex w-full flex-col gap-4 p-6">
			<input
				type="text"
				placeholder="Suchen"
				class="input border-base-content/20 bg-base-300/50 hover:border-base-content/50 hover:bg-base-300/80 w-full rounded-full drop-shadow-md transition duration-200"
				bind:value={searchText}
			/>
			<div class="flex flex-row flex-wrap gap-4">
				<div class="join">
					<input
						class="btn join-item btn-sm bg-base-100 md:btn-md rounded-s-full"
						type="radio"
						name="searchIn"
						aria-label="Transcripts"
						value="transcripts"
						bind:group={meiliIndex}
					/>
					<input
						class="btn join-item btn-sm bg-base-100 md:btn-md rounded-e-full"
						type="radio"
						name="searchIn"
						aria-label="Streamtitel"
						value="vods"
						bind:group={meiliIndex}
					/>
				</div>
				<div class="join">
					<span
						class="join-item border-base-200/70 bg-base-100 flex items-center justify-center rounded-s-full px-4"
					>
						Sortieren
					</span>
					<select
						class="join-item select border-base-200/70 bg-base-200 rounded-e-full"
						aria-label="Sortieren"
						bind:value={selectedSort}
					>
						{#each sorts as sort, index (index)}
							<option value={sort} selected={sort === selectedSort}>{sort.text}</option>
						{/each}
					</select>
					{#if selectedSort.value !== 'relevancy'}
						<div title="Aufsteigend">
							<input
								class="btn join-item"
								type="radio"
								name="options"
								aria-label="&#9650;"
								value="asc"
								bind:group={ordering}
							/>
						</div>
						<div title="Absteigend">
							<input
								class="btn join-item"
								type="radio"
								name="options"
								aria-label="&#9660;"
								value="desc"
								bind:group={ordering}
							/>
						</div>
					{/if}
				</div>
			</div>
		</div>
		<hr class="border-base-content/20 w-full rounded" />
		<div class="flex h-full w-full flex-col gap-4 overflow-y-scroll p-6">
			{#if searchText}
				{#await client.index(meiliIndex).search(searchText, searchConfig)}
					<div class="w-full text-center">
						<span class="loading loading-spinner loading-lg"></span>
					</div>
				{:then result}
					<p class="text-sm font-semibold">
						{#if result.totalHits === 1000}
							&gt;&equals;
						{/if}
						{result.totalHits} Ergebnisse in {result.processingTimeMs}ms
					</p>
					<div class="grid grid-flow-row-dense grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
						{#each result.hits as hit, index (index)}
							<SearchResult {hit} searchIn={meiliIndex} />
						{/each}
					</div>
					<Pagination bind:currentPage totalPages={result.totalPages} />
				{/await}
			{/if}
		</div>
	</div>
	<form method="dialog" class="modal-backdrop bg-base-300/30">
		<button>close</button>
	</form>
</dialog>
