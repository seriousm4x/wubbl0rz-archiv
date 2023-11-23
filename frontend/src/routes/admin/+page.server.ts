import { PUBLIC_API_URL } from '$env/static/public';
import { pb } from '$lib/pocketbase.js';
import { error, redirect } from '@sveltejs/kit';
import type { ListResult, RecordModel } from 'pocketbase';

export async function load({ locals, url }) {
	// check if auth valid
	if (!locals.pb.authStore.isValid) {
		throw redirect(302, '/login');
	}

	// check valid youtube bearer token
	const resp = await fetch(`${PUBLIC_API_URL}/wubbl0rz/youtube/verify`, {
		headers: {
			Authorization: `Bearer ${locals.pb.authStore.token}`
		}
	});

	if (resp.status !== 200) {
		return structuredClone({
			tokenErr: resp.statusText,
			vods: {} as ListResult<RecordModel>,
			user: locals.user
		});
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
		tokenErr: null,
		vods: vods,
		user: locals.user
	});
}
