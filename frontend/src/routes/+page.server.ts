import { pb } from '$lib/pocketbase';
import { error } from '@sveltejs/kit';
import add from 'date-fns/add/index.js';
import format from 'date-fns/format/index.js';

export async function load() {
	const [newestVods, popularVods, clips] = await Promise.all([
		// new vods
		pb
			.collection('vod')
			.getList(1, 13, { sort: '-date', skipTotal: true, requestKey: 'newest_vods' })
			.catch((e) => {
				return e;
			}),
		// popular vods all time
		pb
			.collection('vod')
			.getList(1, 12, {
				sort: '-viewcount',
				skipTotal: true,
				requestKey: 'popular_vods'
			})
			.catch((e) => {
				return e;
			}),
		// clips last month
		pb
			.collection('clip')
			.getList(1, 12, {
				filter: `date >= '${format(add(new Date(), { months: -1 }), 'yyyy-MM-dd')}'`,
				sort: '-viewcount',
				skipTotal: true,
				requestKey: 'clips_last_month'
			})
			.catch((e) => {
				return e;
			})
	]);

	const data = {
		new: newestVods,
		popular: popularVods,
		clips: clips
	};

	if (data) {
		return data;
	}

	throw error(404, 'Not found');
}
