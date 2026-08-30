<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	import Card from '$lib/components/Card.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import IconFilterBoldDuotone from '@iconify-icons/solar/filter-bold-duotone';
	import Icon from '@iconify/svelte';
	import type { ListResult, RecordModel } from 'pocketbase';

	let {
		data,
		title,
		placeholder
	}: { data: ListResult<RecordModel>; title: string; placeholder: string } = $props();

	const origin = page.url.origin;
	const pathname = page.url.pathname;
	let showFilter = $state(false);

	// parse sort from url
	let paramSort = page.url.searchParams.get('sort') || '';
	type sort = {
		value: string;
		text: string;
	};

	// parse page from url
	let currentPage = $derived(parseInt(page.url.searchParams.get('page') || '1', 10));

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
	let ordering: '' | '-' = $state('-');

	// form elements
	let searchValue: string = $state('');
	let selectedSort: sort = $state(
		sorts.find((sort) => {
			return sort.value === paramSort;
		}) || sorts[0]
	);
	let dateFrom: string = $state('');
	let dateTo: string = $state('');

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
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto(url);
	}
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-row items-center justify-between">
		<div class="w-full text-2xl font-semibold tracking-tight">
			<span class="text-base-content tracking-tight">{title}</span>
		</div>
		<div>
			<button
				class="btn btn-sm flex-nowrap rounded-full"
				onclick={() => (showFilter = !showFilter)}
			>
				<Icon icon={IconFilterBoldDuotone} class="text-lg text-violet-500" /> Filter
			</button>
		</div>
	</div>
	{#if showFilter}
		<form
			onsubmit={(e) => {
				e.preventDefault();
				search();
			}}
		>
			<div class="flex flex-col flex-wrap justify-end gap-2 p-3 md:flex-row">
				<input
					class="input border-base-content/10 bg-base-100 text-base-content/70 hover:border-base-content/20 hover:bg-base-100 w-full rounded-full transition-colors md:max-w-lg"
					{placeholder}
					bind:value={searchValue}
				/>
				<div class="join rounded-full">
					<span
						class="join-item bg-base-300 flex items-center justify-center rounded-s-full border-0 px-4"
					>
						Sortieren
					</span>
					<select
						class="join-item select bg-base-100 border-0"
						aria-label="Sortieren"
						bind:value={selectedSort}
					>
						{#each sorts as sort (sort.value)}
							<option value={sort}>{sort.text}</option>
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
							class="btn join-item rounded-e-full"
							type="radio"
							name="options"
							aria-label="&#9660;"
							value="-"
							bind:group={ordering}
						/>
					</div>
				</div>
				<div class="join">
					<span
						class="join-item bg-base-300 flex items-center justify-center rounded-s-full border-0 px-4"
					>
						Von
					</span>
					<input
						type="date"
						class="input join-item bg-base-100 rounded-e-full border-0"
						bind:value={dateFrom}
					/>
				</div>
				<div class="join rounded-full">
					<span
						class="join-item bg-base-300 flex items-center justify-center rounded-s-full border-0 px-4"
					>
						Bis
					</span>
					<input
						type="date"
						class="input join-item bg-base-100 rounded-e-full border-0"
						bind:value={dateTo}
					/>
				</div>
				<div class="flex gap-2">
					<button class="btn btn-primary w-fit rounded-full" type="submit">Suchen</button>
					<button class="btn btn-error w-fit rounded-full" type="button" onclick={reset}
						>Reset</button
					>
				</div>
			</div>
		</form>
	{/if}
	<div
		class="grid grid-cols-1 gap-x-3 gap-y-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 2xl:grid-cols-5"
	>
		{#each data.items as video (video.id)}
			<Card {video} />
		{:else}
			Keine Ergebnisse
		{/each}
	</div>
	<Pagination {currentPage} totalPages={data.totalPages} onPageChange={search} />
</div>
