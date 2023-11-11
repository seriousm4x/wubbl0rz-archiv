import { PUBLIC_API_URL } from '$env/static/public';
import PocketBase from 'pocketbase';

export function createInstance() {
	return new PocketBase(PUBLIC_API_URL);
}

export const pb = createInstance();
