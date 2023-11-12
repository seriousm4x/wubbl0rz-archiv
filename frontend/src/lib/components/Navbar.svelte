<script lang="ts">
	import { navigating } from '$app/stores';
	import { PUBLIC_MEILI_SEARCH_KEY, PUBLIC_MEILI_URL } from '$env/static/public';
	import Pagination from '$lib/components/Pagination.svelte';
	import SearchResult from '$lib/components/SearchResult.svelte';
	import Icon from '@iconify/svelte';
	import { MeiliSearch } from 'meilisearch';
	import { onMount } from 'svelte';

	let scrollY: number;
	let modal: HTMLDialogElement;
	let searchText: string = '';
	let meiliIndex: 'transcripts' | 'vods' = 'transcripts';
	let currentPage = 1;
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
	let ordering: 'asc' | 'desc' = 'desc';
	let selectedSort = sorts[0];

	$: if (meiliIndex || searchText) currentPage = 1;
	$: if ($navigating && modal) modal.close();
	$: searchConfig = {
		page: currentPage,
		hitsPerPage: 30,
		attributesToHighlight: ['*'],
		highlightPreTag: '<span>',
		highlightPostTag: '</span>',
		sort: selectedSort.value === 'relevancy' ? [] : [`${selectedSort.value}:${ordering}`]
	};

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
	class="navbar backdrop-blur-sm sticky top-0 z-40 justify-center transition duration-200 {scrollY >
	0
		? 'bg-base-200/70'
		: ''}"
>
	<label
		for="menu-drawer"
		class="btn btn-ghost p-0 mask mask-squircle drawer-button md:hidden flex-nowrap"
	>
		<Icon icon="solar:menu-dots-square-bold-duotone" class="text-5xl text-primary" />
	</label>
	<button
		class="input input-bordered border-base-content/20 md:max-w-lg w-full bg-base-300/50 drop-shadow-md hover:bg-base-300/80 hover:border-base-content/50 cursor-pointer transition duration-200 rounded-full text-base-content/50 hover:text-base-content"
		on:click={showModal}
	>
		<span>Suchen</span>
		<span class="ms-auto"><kbd class="kbd">/</kbd></span>
	</button>
</div>
<dialog
	class="modal items-start backdrop-blur-md transition duration-200 overflow-hidden"
	bind:this={modal}
	on:close={onModalClose}
>
	<div
		class="modal-box p-0 md:mt-20 max-w-6xl h-full md:h-4/5 max-h-full flex flex-col items-center bg-base-300/90 border border-base-content/20 backdrop-blur-md"
	>
		<div class="flex flex-col w-full gap-4 p-6">
			<input
				type="text"
				placeholder="Suchen"
				class="input input-bordered border-base-content/20 w-full bg-base-300/50 drop-shadow-md hover:bg-base-300/80 hover:border-base-content/50 cursor-pointer transition duration-200 rounded-full"
				bind:value={searchText}
			/>
			<div class="flex flex-row flex-wrap gap-4">
				<div class="join rounded-full">
					<input
						class="join-item btn btn-sm md:btn-md bg-base-100"
						type="radio"
						name="searchIn"
						aria-label="Transcripts"
						value="transcripts"
						bind:group={meiliIndex}
					/>
					<input
						class="join-item btn btn-sm md:btn-md bg-base-100"
						type="radio"
						name="searchIn"
						aria-label="Streamtitel"
						value="vods"
						bind:group={meiliIndex}
					/>
				</div>
				<div class="join rounded-full">
					<span
						class="join-item flex items-center px-4 justify-center border-base-200/70 bg-base-100"
					>
						Sortieren
					</span>
					<select
						class="join-item select select-bordered rounded-e-full border-base-200/70 bg-base-200"
						bind:value={selectedSort}
					>
						{#each sorts as sort}
							<option value={sort} selected={sort === selectedSort}>{sort.text}</option>
						{/each}
					</select>
					{#if selectedSort.value !== 'relevancy'}
						<div title="Aufsteigend">
							<input
								class="join-item btn"
								type="radio"
								name="options"
								aria-label="&#9650;"
								value="asc"
								bind:group={ordering}
							/>
						</div>
						<div title="Absteigend">
							<input
								class="join-item btn"
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
		<hr class="rounded border-base-content/20 w-full" />
		<div class="flex flex-col w-full h-full gap-4 p-6 overflow-y-scroll">
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
					<div class="grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
						{#each result.hits as hit}
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
