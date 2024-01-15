import { createInstance } from '$lib/stores/pocketbase.js';
import { error } from '@sveltejs/kit';
import { add, format, parseISO } from 'date-fns';
import type { ListResult, RecordModel } from 'pocketbase';

export async function load({ params }) {
	const pb = createInstance();
	let vod = {} as RecordModel;
	let allVods = {} as ListResult<RecordModel>;
	let vodPosition = {} as ListResult<RecordModel>;
	let recommendations = {} as ListResult<RecordModel>;

	await Promise.all([
		pb
			.collection('vod')
			.getOne(params.id, {
				expand: 'clip(vod)',
				requestKey: 'single_vod'
			})
			.then((data) => {
				vod = data;
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
			.then((data) => {
				allVods = data;
			})
			.catch((e) => {
				return e;
			})
	]);

	if (!vod.id) {
		throw error(404, 'Not found');
	}

	await Promise.all([
		pb
			.collection('vod')
			.getList(1, 1, {
				sort: '-date',
				filter: `viewcount >= ${vod.viewcount}`,
				requestKey: 'vod_position'
			})
			.then((data) => {
				vodPosition = data;
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
			.then((data) => {
				recommendations = data;
			})
			.catch((e) => {
				return e;
			})
	]);

	return structuredClone({
		vod: vod,
		vodsCount: allVods.totalItems,
		vodPosition: vodPosition.totalItems,
		recommendations: recommendations.items.filter((v: RecordModel) => v.id !== vod.id)
	});
}
