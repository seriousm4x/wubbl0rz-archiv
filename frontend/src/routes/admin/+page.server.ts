import { pb } from '$lib/pocketbase.js';
import { error, redirect } from '@sveltejs/kit';
import type { ListResult, RecordModel } from 'pocketbase';

export async function load({ locals, url }) {
	// check if auth valid
	if (!locals.pb.authStore.isValid) {
		throw redirect(302, '/login');
	}

	// get vods
	let vods = {} as ListResult<RecordModel>;
	await pb
		.collection('vod')
		.getList(1, 20, {
			sort: '-date',
			page: parseInt(url.searchParams.get('page') || '1') || 1,
			requestKey: 'all_vods'
		})
		.then((data) => (vods = data))
		.catch((e) => {
			return e;
		});

	if (vods.totalItems === 0) {
		throw error(404, 'No vods found');
	}

	return structuredClone({
		vods: vods,
		user: locals.user
	});
}
