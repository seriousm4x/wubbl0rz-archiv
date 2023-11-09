import { error } from '@sveltejs/kit';
import PocketBase, { type RecordModel } from 'pocketbase';
import type { ListResult } from 'pocketbase';
import { PUBLIC_API_URL } from '$env/static/public';

export async function load({ url }) {
	const pb = new PocketBase(PUBLIC_API_URL);
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
		return allVods;
	}

	throw error(404, 'Not found');
}
