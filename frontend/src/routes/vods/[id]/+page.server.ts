import add from 'date-fns/add/index.js';
import format from 'date-fns/format/index.js';
import parseISO from 'date-fns/parseISO/index.js';
import PocketBase, { type RecordModel } from 'pocketbase';
import { PUBLIC_API_URL } from '$env/static/public';
import { error } from '@sveltejs/kit';

export async function load({ params }) {
	const pb = new PocketBase(PUBLIC_API_URL);
	const [vod, allVods] = await Promise.all([
		pb
			.collection('vod')
			.getOne(params.id, {
				expand: 'clip(vod)',
				requestKey: 'single_vod'
			})
			.catch((e) => {
				return e;
			}),
		pb
			.collection('vod')
			.getList(1, 1, {
				filter: 'viewcount > 0',
				requestKey: 'vod_count'
			})
			.catch((e) => {
				return e;
			})
	]);

	if (!vod.id) {
		throw error(404, 'Not found');
	}

	const [vodPosition, recommendations] = await Promise.all([
		pb
			.collection('vod')
			.getList(1, 1, {
				sort: '-date',
				filter: `viewcount >= ${vod.viewcount}`,
				requestKey: 'vod_position'
			})
			.catch((e) => {
				return e;
			}),
		pb
			.collection('vod')
			.getList(1, 12, {
				sort: '-viewcount',
				filter: `date >= '${format(
					add(parseISO(vod.date), { months: -2 }),
					'yyyy-MM-dd'
				)}' && date < '${format(add(parseISO(vod.date), { months: +2 }), 'yyyy-MM-dd')}'`,
				requestKey: 'vod_recommendations'
			})
			.catch((e) => {
				return e;
			})
	]);

	return {
		vod: vod,
		vodsCount: allVods.totalItems,
		vodPosition: vodPosition.totalItems,
		recommendations: recommendations.items.filter((v: RecordModel) => v.id !== vod.id)
	};
}
