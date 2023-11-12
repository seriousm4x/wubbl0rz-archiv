import { pb } from '$lib/pocketbase';
import { error } from '@sveltejs/kit';
import type { ListResult, RecordModel } from 'pocketbase';

export async function load() {
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

	if (chatmessages.items.length > 0) {
		return structuredClone(chatmessages);
	}

	throw error(404, 'Not found');
}
