import { pb } from '$lib/pocketbase.js';
import { error } from '@sveltejs/kit';
import type { ListResult, RecordModel } from 'pocketbase';

export async function load({ url }) {
	let allClips = {} as ListResult<RecordModel>;

	await pb
		.collection('clip')
		.getList(1, 36, {
			sort: url.searchParams.get('sort') || '-date',
			filter: url.searchParams.get('filter') || '',
			page: parseInt(url.searchParams.get('page') || '1') || 1,
			requestKey: 'clip_grid'
		})
		.then((data) => {
			allClips = data;
		})
		.catch((e) => {
			return e;
		});

	if (allClips.items.length > 0) {
		return structuredClone(allClips);
	}

	throw error(404, 'Not found');
}
