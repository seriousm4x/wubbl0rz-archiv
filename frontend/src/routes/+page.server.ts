import { error } from '@sveltejs/kit';
import add from 'date-fns/add/index.js';
import format from 'date-fns/format/index.js';
import type { ListResult, RecordModel } from 'pocketbase';
import { createInstance } from '$lib/stores/pocketbase';

export async function load() {
	const pb = createInstance();
	let newestVods = {} as ListResult<RecordModel>;
	let popularVods = {} as ListResult<RecordModel>;
	let clips = {} as ListResult<RecordModel>;

	await Promise.all([
		// new vods
		pb
			.collection('vod')
			.getList(1, 13, { sort: '-date', skipTotal: true, requestKey: 'newest_vods' })
			.then((data) => {
				newestVods = data;
			})
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
			.then((data) => {
				popularVods = data;
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
			.then((data) => {
				clips = data;
			})
			.catch((e) => {
				return e;
			})
	]);

	if (!newestVods.items || !popularVods.items || !clips.items) {
		throw error(404, 'Not found');
	}

	const data = {
		new: newestVods,
		popular: popularVods,
		clips: clips
	};

	return structuredClone(data);
}
