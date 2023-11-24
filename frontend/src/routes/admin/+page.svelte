<script lang="ts">
	import { browser } from '$app/environment';
	import { applyAction, enhance } from '$app/forms';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { PUBLIC_API_URL } from '$env/static/public';
	import Pagination from '$lib/components/Pagination.svelte';
	import { toHHMMSS } from '$lib/functions';
	import { pb } from '$lib/pocketbase';
	import { currentUser } from '$lib/stores/user.js';
	import Icon from '@iconify/svelte';
	import { parseISO } from 'date-fns';
	import type { RecordModel } from 'pocketbase';
	import { onMount } from 'svelte';
	import type { PageData } from './$types';

	export let data: PageData;
	let currentPage = parseInt($page.url.searchParams.get('page') || `${data.vods.page || '1'}`);

	$: vods = data.vods.items as RecordModel[];
	$: if (browser) goto(`/admin?page=${currentPage}`);
	$: currentUser.set(data.user);

	onMount(() => {
		if (data.tokenErr) return;
		pb.collection('vod').subscribe('*', (data) => {
			const i = vods.findIndex((vod) => vod.id === data.record.id);
			if (i < 0) return;
			vods[i] = data.record;
		});
	});

	function youtubeRevoke() {
		fetch(`${PUBLIC_API_URL}/wubbl0rz/youtube/revoke`, {
			headers: {
				Authorization: `Bearer ${pb.authStore.token}`
			}
		}).then(() => goto('/admin', { invalidateAll: true }));
	}

	function youtubeLogin() {
		fetch(`${PUBLIC_API_URL}/wubbl0rz/youtube/login`, {
			headers: {
				Authorization: `Bearer ${pb.authStore.token}`
			}
		})
			.then((resp) => resp.text())
			.then((url) => {
				goto(url);
			});
	}

	function upload(id: string) {
		fetch(`${PUBLIC_API_URL}/wubbl0rz/youtube/upload/${id}`, {
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
	{#if data.tokenErr}
		<div role="alert" class="alert alert-error shadow-lg">
			<Icon icon="solar:danger-triangle-bold-duotone" class="text-4xl" />
			<div>
				<h3 class="text-lg font-bold">YouTube Token ist nicht valide!</h3>
				<div>
					Beim Überprüfen des YouTube Bearer Token ist ein Fehler aufgetreten. Klicke den Button
					rechts um dich mit deinem YouTube Account einzuloggen und den Token zu aktualisieren.
				</div>
			</div>
			<button class="btn rounded-full" on:click={youtubeLogin}>
				<Icon icon="solar:login-3-bold-duotone" class="text-2xl" /> Login
			</button>
		</div>
		<div class="ms-auto flex flex-row flex-wrap gap-2">
			<form
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
		</div>
	{:else}
		<div class="ms-auto flex flex-row flex-wrap gap-2">
			<button class="btn btn-outline btn-error" on:click={youtubeRevoke}
				><Icon icon="solar:trash-bin-2-bold-duotone" class="text-3xl" /> YouTube Token Löschen</button
			>
			<form
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
		</div>
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
							><span class="me-1 font-bold">Views:</span>{vod.viewcount.toLocaleString(
								'de-DE'
							)}</span
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
	{/if}
</div>
