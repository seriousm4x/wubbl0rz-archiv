<script lang="ts">
	import { applyAction, enhance } from '$app/forms';
	import { pb } from '$lib/pocketbase';

	let error = '';
</script>

<div class="container mx-auto flex flex-col items-center justify-center gap-4">
	<form
		method="POST"
		use:enhance={() => {
			return async ({ result }) => {
				pb.authStore.loadFromCookie(document.cookie);
				await applyAction(result);
			};
		}}
	>
		<div class="form-control">
			<label class="label" for="username">
				<span class="label-text">Username</span>
			</label>
			<input
				id="username"
				name="username"
				type="text"
				placeholder="Username"
				class="input input-bordered w-full max-w-xs rounded-full"
				required
			/>
		</div>
		<div class="form-control w-full max-w-xs">
			<label class="label" for="password">
				<span class="label-text">Password</span>
			</label>
			<input
				id="password"
				name="password"
				type="password"
				placeholder="Password"
				class="input input-bordered w-full max-w-xs rounded-full"
				min="8"
				required
			/>
		</div>
		<div class="form-control mt-4 w-full max-w-xs">
			<button class="btn" type="submit">Login</button>
		</div>
	</form>
	{#if error !== ''}
		<div class="alert alert-error">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="h-6 w-6 shrink-0 stroke-current"
				fill="none"
				viewBox="0 0 24 24"
				><path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
				/></svg
			>
			<span>{error}</span>
		</div>
	{/if}
</div>
