<script lang="ts">
	import { applyAction, enhance } from '$app/forms';
	import { PUBLIC_API_URL } from '$env/static/public';
	import { toHHMMSS } from '$lib/functions';
	import { pb } from '$lib/pocketbase';
	import { currentUser } from '$lib/stores/user.js';
	import Icon from '@iconify/svelte';
	import { parseISO } from 'date-fns';
	import type { RecordModel } from 'pocketbase';
	import { onMount } from 'svelte';
	import type { PageData } from './$types';
	import Pagination from '$lib/components/Pagination.svelte';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';

	export let data: PageData;
	let currentPage = parseInt($page.url.searchParams.get('page') || `${data.vods.page}`);

	$: vods = data.vods.items as RecordModel[];
	$: if (browser) goto(`/admin?page=${currentPage}`);
	$: currentUser.set(data.user);

	onMount(() => {
		pb.collection('vod').subscribe('*', (data) => {
			const i = vods.findIndex((vod) => vod.id === data.record.id);
			if (i < 0) return;
			vods[i] = data.record;
		});
	});

	function upload(id: string) {
		fetch(`${PUBLIC_API_URL}/youtube/upload/${id}`, {
			headers: {
				Authorization: `Bearer ${pb.authStore.token}`
			}
		});
	}
</script>

<div class="container mx-auto flex flex-col gap-4">
	<h1 class="mb-4 text-4xl font-bold md:ms-3 md:mt-10">
		<span
			class="bg-gradient-to-r from-pink-500 to-violet-500 bg-clip-text text-transparent drop-shadow-md"
			>Admin</span
		>
	</h1>

	<form
		class="ms-auto"
		method="POST"
		action="/logout"
		use:enhance={() => {
			return async ({ result }) => {
				pb.authStore.clear();
				await applyAction(result);
			};
		}}
	>
		<button class="btn btn-outline btn-error"
			><Icon icon="solar:logout-3-bold-duotone" class="text-3xl" /> Log out</button
		>
	</form>

	{#each vods as vod}
		<div
			class="collapse {vod.youtube_upload === 'done'
				? 'bg-green-500/30'
				: vod.youtube_upload === 'pending'
				  ? 'bg-amber-300/40'
				  : 'bg-base-200'}"
		>
			<input type="checkbox" />
			<div class="collapse-title">
				<p class="break-all text-xl font-medium">{vod.title}</p>
				<p class="flex flex-row flex-wrap gap-1">
					<span class="badge badge-neutral"
						><span class="me-1 font-bold">Datum:</span>{parseISO(vod.date).toLocaleDateString(
							'de',
							{
								weekday: 'long',
								year: 'numeric',
								month: 'numeric',
								day: 'numeric'
							}
						)}</span
					>
					<span class="badge badge-neutral"
						><span class="me-1 font-bold">Länge:</span>{toHHMMSS(vod.duration, true)}</span
					>
					<span class="badge badge-neutral"
						><span class="me-1 font-bold">Views:</span>{vod.viewcount.toLocaleString('de-DE')}</span
					>
					<span class="badge badge-neutral"><span class="me-1 font-bold">ID:</span>{vod.id}</span>
				</p>
			</div>
			<div class="collapse-content bg-base-300/40">
				<dir class="flex flex-col gap-2 p-0">
					{#if vod.youtube_upload === ''}
						<button class="btn btn-outline btn-success" on:click={() => upload(vod.id)}
							>Export zu YouTube</button
						>
					{:else if vod.youtube_upload === 'done'}
						<button class="btn btn-outline btn-success" on:click={() => upload(vod.id)}
							>Bereits hochgeladen. Du kannst den Stream aber erneut hochladen</button
						>
					{:else}
						<button class="btn btn-disabled btn-outline btn-success">Upload läuft...</button>
					{/if}
				</dir>
			</div>
		</div>
	{/each}

	<Pagination totalPages={data.vods.totalPages} bind:currentPage />
</div>
