import { pb } from '$lib/pocketbase';
import { error } from '@sveltejs/kit';

export async function load() {
	const chatmessages = await pb
		.collection('chatmessage')
		.getList(1, 200, { sort: '-date', skipTotal: true, requestKey: 'chatmessages' })
		.catch((e) => {
			return e;
		});

	if (chatmessages.items.length > 0) {
		return structuredClone(chatmessages);
	}

	throw error(404, 'Not found');
}
