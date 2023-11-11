import { pb } from '$lib/pocketbase.js';
import { error } from '@sveltejs/kit';
import type { ListResult, RecordModel } from 'pocketbase';

export async function load({ url }) {
	const allVods = await pb
		.collection('vod')
		.getList(1, 36, {
			sort: url.searchParams.get('sort') || '-date',
			filter: url.searchParams.get('filter') || '',
			page: parseInt(url.searchParams.get('page') || '1') || 1,
			requestKey: 'vod_grid'
		})
		.catch((e) => {
			return e;
		});

	if ((allVods as ListResult<RecordModel>).totalItems > 0) {
		return structuredClone(allVods);
	}

	throw error(404, 'Not found');
}
