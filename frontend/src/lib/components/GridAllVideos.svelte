<script lang="ts">
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { PUBLIC_API_URL } from '$env/static/public';
	import Card from '$lib/components/Card.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import Icon from '@iconify/svelte';
	import type { ListResult, RecordModel } from 'pocketbase';

	export let data: ListResult<RecordModel>;
	export let title: string;
	export let placeholder: string;

	const type = data.items[0]?.collectionName === 'vod' ? 'vods' : 'clips';

	const origin = $page.url.origin;
	const pathname = $page.url.pathname;
	let showFilter = false;

	// parse sort from url
	let paramSort = $page.url.searchParams.get('sort') || '';
	type sort = {
		value: string;
		text: string;
	};

	// parse page from url
	let currentPage = parseInt($page.url.searchParams.get('page') || `${data.page}`);
	$: if (browser) search(currentPage);

	// available sorts
	let sorts: sort[] = [
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
		},
		{
			value: 'size',
			text: 'Dateigröße'
		}
	];
	let ordering: '' | '-' = '-';

	// form elements
	let searchValue: string;
	let selectedSort =
		sorts.find((sort) => {
			return sort.value === paramSort;
		}) || sorts[0];
	let dateFrom: string;
	let dateTo: string;

	function reset() {
		searchValue = '';
		selectedSort = sorts[0];
		ordering = '-';
		dateFrom = '';
		dateTo = '';
		search(1);
	}

	// set url params and reload content
	function search(page = 1) {
		const url = new URL(origin + pathname);
		let filter: string[] = [];

		if (searchValue !== undefined && searchValue !== '') {
			filter = [...filter, `title ~ '${searchValue}'`];
		}
		if (dateFrom !== undefined && dateFrom !== '') {
			filter = [...filter, `date >= '${dateFrom}'`];
		}
		if (dateTo !== undefined && dateTo !== '') {
			filter = [...filter, `date <= '${dateTo}'`];
		}

		url.searchParams.append('filter', filter.join(' && ') || '');
		url.searchParams.append('sort', `${ordering}${selectedSort.value}`);
		url.searchParams.append('page', page.toString());
		goto(url);
		currentPage = page;
	}
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-row justify-between items-center">
		<div class="text-4xl font-bold w-full">
			<span
				class="bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-violet-500 drop-shadow-md"
				>{title}</span
			>
		</div>
		<div>
			<button
				class="btn btn-sm rounded-full flex-nowrap"
				on:click={() => (showFilter = !showFilter)}
			>
				<Icon icon="solar:filter-bold-duotone" class="text-violet-500 text-lg" /> Filter
			</button>
		</div>
	</div>
	{#if showFilter}
		<form on:submit|preventDefault={() => search()}>
			<div class="flex md:flex-row flex-col flex-wrap gap-2 justify-end">
				<input
					class="input input-bordered border-base-content/20 md:max-w-lg w-full bg-base-300/50 drop-shadow-md hover:bg-base-300/80 hover:border-base-content/50 cursor-pointer transition duration-200 rounded-full text-base-content/50 hover:text-base-content"
					{placeholder}
					bind:value={searchValue}
				/>
				<div class="join rounded-full">
					<span
						class="join-item flex items-center px-4 justify-center border-base-200/70 bg-base-100"
					>
						Sortieren
					</span>
					<select
						class="join-item select select-bordered rounded-e-full border-base-200/70 bg-base-200"
						aria-label="Sortieren"
						bind:value={selectedSort}
					>
						{#each sorts as sort}
							<option value={sort} selected={sort === selectedSort}>{sort.text}</option>
						{/each}
					</select>
					<div title="Aufsteigend">
						<input
							class="join-item btn"
							type="radio"
							name="options"
							aria-label="&#9650;"
							value=""
							bind:group={ordering}
						/>
					</div>
					<div title="Absteigend">
						<input
							class="join-item btn"
							type="radio"
							name="options"
							aria-label="&#9660;"
							value="-"
							bind:group={ordering}
						/>
					</div>
				</div>
				<div class="join rounded-full">
					<span
						class="join-item flex items-center px-4 justify-center border-base-200/70 bg-base-100"
					>
						Von
					</span>
					<input
						type="date"
						class="join-item input border-base-200/70 bg-base-200"
						bind:value={dateFrom}
					/>
				</div>
				<div class="join rounded-full">
					<span
						class="join-item flex items-center px-4 justify-center border-base-200/70 bg-base-100"
					>
						Bis
					</span>
					<input
						type="date"
						class="join-item input border-base-200/70 bg-base-200"
						bind:value={dateTo}
					/>
				</div>
				<div class="flex gap-2">
					<button class="btn btn-primary rounded-full w-fit" type="submit">Suchen</button>
					<button class="btn btn-error rounded-full w-fit" type="button" on:click={reset}
						>Reset</button
					>
				</div>
			</div>
		</form>
	{/if}
	<div
		class="grid grid-flow-row-dense grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4"
	>
		<div
			class="absolute top-0 left-0 w-full h-auto aspect-video bg-cover bg-center blur-2xl opacity-40 -z-10 transition- duration-200"
		>
			<img
				src={data.items[0]?.filename
					? `${PUBLIC_API_URL}/${type}/${data.items[0]?.filename}-sm.webp`
					: ''}
				alt={data.items[0]?.title}
				class="w-full h-auto"
			/>
		</div>
		{#each data.items as video}
			<Card {video} />
		{:else}
			Keine Ergebnisse
		{/each}
	</div>
	<Pagination bind:currentPage totalPages={data.totalPages} />
</div>
