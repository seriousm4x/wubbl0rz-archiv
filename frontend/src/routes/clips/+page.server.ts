import { pb } from '$lib/pocketbase.js';
import { error } from '@sveltejs/kit';

export async function load({ url }) {
	const allClips = await pb
		.collection('clip')
		.getList(1, 36, {
			sort: url.searchParams.get('sort') || '-date',
			filter: url.searchParams.get('filter') || '',
			page: parseInt(url.searchParams.get('page') || '1') || 1,
			requestKey: 'clip_grid'
		})
		.catch((e) => {
			return e;
		});

	if (allClips) {
		return allClips;
	}

	throw error(404, 'Not found');
}
