import { createInstance } from '$lib/stores/pocketbase.js';
import type { ListResult, RecordModel } from 'pocketbase';

export async function load({ url }) {
	const pb = createInstance();
	let allVods = {} as ListResult<RecordModel>;

	await pb
		.collection('vod')
		.getList(1, 36, {
			sort: url.searchParams.get('sort') || '-date',
			filter: url.searchParams.get('filter') || '',
			page: parseInt(url.searchParams.get('page') || '1') || 1,
			requestKey: 'vod_grid'
		})
		.then((data) => {
			allVods = data;
		})
		.catch((e) => {
			return e;
		});

	return structuredClone(allVods);
}
