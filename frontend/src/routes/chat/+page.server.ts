import { createInstance } from '$lib/stores/pocketbase';
import { error } from '@sveltejs/kit';
import type { ListResult, RecordModel } from 'pocketbase';

export async function load() {
	const pb = createInstance();
	let chatmessages = {} as ListResult<RecordModel>;

	await pb
		.collection('chatmessage')
		.getList(1, 200, { sort: '-date', skipTotal: true, requestKey: 'chatmessages' })
		.then((data) => {
			chatmessages = data;
		})
		.catch((e) => {
			return e;
		});

	if (chatmessages.totalItems === 0) {
		throw error(404, 'Not found');
	}

	return structuredClone(chatmessages);
}
