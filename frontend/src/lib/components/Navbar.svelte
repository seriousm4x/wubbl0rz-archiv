<script lang="ts">
	import { navigating } from '$app/state';
	import { PUBLIC_MEILI_SEARCH_KEY, PUBLIC_MEILI_URL } from '$env/static/public';
	import Pagination from '$lib/components/Pagination.svelte';
	import SearchResult from '$lib/components/SearchResult.svelte';
	import IconMenuDotsSquareBoldDuotone from '@iconify-icons/solar/menu-dots-square-bold-duotone';
	import Icon from '@iconify/svelte';
	import { Meilisearch } from 'meilisearch';

	let scrollY: number = $state(0);
	let modal: HTMLDialogElement;
	let searchText: string = $state('');
	let meiliIndex: 'transcripts' | 'vods' = $state('transcripts');
	let currentPage = $state(1);

	const client = new Meilisearch({
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
		if (navigating && modal) modal.close();
	});

	let searchConfig = $derived({
		page: currentPage,
		hitsPerPage: 30,
		attributesToHighlight: ['*'],
		highlightPreTag: '<span>',
		highlightPostTag: '</span>',
		sort: selectedSort.value === 'relevancy' ? [] : [`${selectedSort.value}:${ordering}`]
	});

	function showModal() {
		document.body.classList.add('overflow-hidden');
		modal.show();
	}

	function onModalClose() {
		document.body.classList.remove('overflow-hidden');
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

<svelte:window bind:scrollY onkeydown={onKeyDown} />

<div
	class="navbar sticky top-0 z-40 items-center justify-center gap-2 backdrop-blur-sm transition duration-200 {scrollY >
	0
		? 'bg-base-200/70'
		: ''}"
>
	<label
		for="menu-drawer"
		class="btn btn-ghost btn-square drawer-button flex shrink-0 items-center justify-center self-center p-0 md:hidden"
	>
		<Icon icon={IconMenuDotsSquareBoldDuotone} class="text-primary block text-3xl" />
	</label>
	<div class="aura aura-sm text-primary/30 w-full max-w-lg rounded-full">
		<button
			class="input border-base-content/10 bg-base-200 text-base-content/65 hover:border-base-content/20 hover:bg-base-200 w-full cursor-pointer rounded-full transition-colors"
			onclick={showModal}
		>
			<span>Suchen</span>
			<span class="ms-auto"><kbd class="kbd">/</kbd></span>
		</button>
	</div>
</div>
<dialog
	class="modal items-start overflow-hidden backdrop-blur-md transition duration-200"
	bind:this={modal}
	onclose={onModalClose}
>
	<div
		class="modal-box border-base-content/10 bg-base-200/95 flex max-h-[80dvh] w-full max-w-6xl flex-col items-center border p-0 backdrop-blur-md md:mt-20"
	>
		<div class="flex w-full flex-col gap-4 p-6">
			<input
				type="text"
				placeholder="Suchen"
				class="input border-base-content/10 bg-base-100 hover:border-base-content/20 hover:bg-base-100 w-full rounded-full transition-colors"
				bind:value={searchText}
				oninput={() => (currentPage = 1)}
			/>
			<div class="flex flex-row flex-wrap gap-4">
				<div class="join">
					<input
						class="btn join-item btn-sm border-base-content/10 md:btn-md rounded-s-full {meiliIndex ===
						'transcripts'
							? 'bg-primary! text-primary-content!'
							: 'bg-base-100 text-base-content'}"
						type="radio"
						name="searchIn"
						aria-label="Transcripts"
						value="transcripts"
						bind:group={meiliIndex}
						onchange={() => (currentPage = 1)}
					/>
					<input
						class="btn join-item btn-sm border-base-content/10 md:btn-md rounded-e-full {meiliIndex ===
						'vods'
							? 'bg-primary! text-primary-content!'
							: 'bg-base-100 text-base-content'}"
						type="radio"
						name="searchIn"
						aria-label="Streamtitel"
						value="vods"
						bind:group={meiliIndex}
						onchange={() => (currentPage = 1)}
					/>
				</div>
				<div class="join">
					<span
						class="join-item bg-base-300 flex items-center justify-center rounded-s-full border-0 px-4"
					>
						Sortieren
					</span>
					<select
						class="join-item select bg-base-100 border-0 {selectedSort.value !== 'relevancy'
							? ''
							: 'rounded-e-full'}"
						aria-label="Sortieren"
						bind:value={selectedSort}
						onchange={() => (currentPage = 1)}
					>
						{#each sorts as sort (sort.value)}
							<option value={sort}>{sort.text}</option>
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
								onchange={() => (currentPage = 1)}
							/>
						</div>
						<div title="Absteigend">
							<input
								class="btn join-item rounded-e-full"
								type="radio"
								name="options"
								aria-label="&#9660;"
								value="desc"
								bind:group={ordering}
								onchange={() => (currentPage = 1)}
							/>
						</div>
					{/if}
				</div>
			</div>
		</div>
		<hr class="border-base-content/20 w-full rounded" />
		<div class="flex max-h-[90dvh] w-full flex-col gap-4 overflow-y-auto p-6">
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
							<SearchResult {hit} searchIn={meiliIndex} onresultselect={() => modal.close()} />
						{/each}
					</div>
					<Pagination
						{currentPage}
						totalPages={result.totalPages}
						onPageChange={(page) => (currentPage = page)}
					/>
				{/await}
			{/if}
		</div>
	</div>
	<form method="dialog" class="modal-backdrop bg-base-100/60">
		<button>close</button>
	</form>
</dialog>
