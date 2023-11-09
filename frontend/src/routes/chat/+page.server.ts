import { error } from '@sveltejs/kit';
import PocketBase from 'pocketbase';
import { PUBLIC_API_URL } from '$env/static/public';

export async function load() {
	const pb = new PocketBase(PUBLIC_API_URL);
	const chatmessages = await pb
		.collection('chatmessage')
		.getList(1, 100, { sort: '-date', skipTotal: true, requestKey: 'chatmessages' })
		.catch((e) => {
			return e;
		});

	if (chatmessages.items.length > 0) {
		return chatmessages;
	}

	throw error(404, 'Not found');
}
