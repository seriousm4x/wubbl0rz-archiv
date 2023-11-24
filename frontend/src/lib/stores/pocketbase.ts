import { PUBLIC_API_URL } from '$env/static/public';
import PocketBase from 'pocketbase';
import { writable } from 'svelte/store';

export function createInstance() {
	return new PocketBase(PUBLIC_API_URL);
}

export const pb = writable(createInstance());
