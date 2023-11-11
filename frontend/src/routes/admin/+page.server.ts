import { pb } from '$lib/pocketbase.js';
import { redirect } from '@sveltejs/kit';

export async function load({ locals }) {
	// check if auth valid
	if (!locals.pb.authStore.isValid) {
		throw redirect(302, '/login');
	}

	// get vods
	const vods = await pb
		.collection('vod')
		.getFullList({
			sort: '-date',
			requestKey: 'all_vods'
		})
		.catch((e) => {
			return e;
		});

	return {
		vods: structuredClone(vods),
		user: structuredClone(locals.user)
	};
}
