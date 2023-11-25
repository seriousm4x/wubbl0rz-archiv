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
	<div class="flex flex-row items-center justify-between">
		<div class="w-full text-4xl font-bold">
			<span
				class="bg-gradient-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent drop-shadow-md"
				>{title}</span
			>
		</div>
		<div>
			<button
				class="btn btn-sm flex-nowrap rounded-full"
				on:click={() => (showFilter = !showFilter)}
			>
				<Icon icon="solar:filter-bold-duotone" class="text-lg text-violet-500" /> Filter
			</button>
		</div>
	</div>
	{#if showFilter}
		<form on:submit|preventDefault={() => search()}>
			<div class="flex flex-col flex-wrap justify-end gap-2 md:flex-row">
				<input
					class="input w-full rounded-full border-base-content/20 bg-base-300/50 text-base-content/50 drop-shadow-md transition duration-200 hover:border-base-content/50 hover:bg-base-300/80 hover:text-base-content md:max-w-lg"
					{placeholder}
					bind:value={searchValue}
				/>
				<div class="join rounded-full">
					<span
						class="join-item flex items-center justify-center border-base-200/70 bg-base-100 px-4"
					>
						Sortieren
					</span>
					<select
						class="join-item select rounded-e-full border-base-200/70 bg-base-200"
						aria-label="Sortieren"
						bind:value={selectedSort}
					>
						{#each sorts as sort}
							<option value={sort} selected={sort === selectedSort}>{sort.text}</option>
						{/each}
					</select>
					<div title="Aufsteigend">
						<input
							class="btn join-item"
							type="radio"
							name="options"
							aria-label="&#9650;"
							value=""
							bind:group={ordering}
						/>
					</div>
					<div title="Absteigend">
						<input
							class="btn join-item"
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
						class="join-item flex items-center justify-center border-base-200/70 bg-base-100 px-4"
					>
						Von
					</span>
					<input
						type="date"
						class="input join-item border-base-200/70 bg-base-200"
						bind:value={dateFrom}
					/>
				</div>
				<div class="join rounded-full">
					<span
						class="join-item flex items-center justify-center border-base-200/70 bg-base-100 px-4"
					>
						Bis
					</span>
					<input
						type="date"
						class="input join-item border-base-200/70 bg-base-200"
						bind:value={dateTo}
					/>
				</div>
				<div class="flex gap-2">
					<button class="btn btn-primary w-fit rounded-full" type="submit">Suchen</button>
					<button class="btn btn-error w-fit rounded-full" type="button" on:click={reset}
						>Reset</button
					>
				</div>
			</div>
		</form>
	{/if}
	<div
		class="grid grid-flow-row-dense grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5"
	>
		<div
			class="transition- absolute left-0 top-0 -z-10 aspect-video h-auto w-full bg-cover bg-center opacity-40 blur-2xl duration-200"
		>
			<img
				src={data.items[0]?.filename
					? `${PUBLIC_API_URL}/${type}/${data.items[0]?.filename}-sm.webp`
					: ''}
				alt={data.items[0]?.title}
				class="h-auto w-full"
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
