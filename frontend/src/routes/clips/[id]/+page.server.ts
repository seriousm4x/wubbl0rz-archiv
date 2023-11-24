import { createInstance } from '$lib/stores/pocketbase.js';
import add from 'date-fns/add/index.js';
import format from 'date-fns/format/index.js';
import parseISO from 'date-fns/parseISO/index.js';
import type { ListResult, RecordModel } from 'pocketbase';
import { error } from '@sveltejs/kit';

export async function load({ params }) {
	const pb = createInstance();
	let clip = {} as RecordModel;
	let allClips = {} as ListResult<RecordModel>;
	let clipPosition = {} as ListResult<RecordModel>;
	let recommendations = {} as ListResult<RecordModel>;

	await Promise.all([
		pb
			.collection('clip')
			.getOne(params.id, {
				expand: 'vod,creator,game',
				requestKey: 'single_clip'
			})
			.then((data) => {
				clip = data;
			})
			.catch((e) => {
				return e;
			}),
		pb
			.collection('clip')
			.getList(1, 1, {
				requestKey: 'clip_count'
			})
			.then((data) => {
				allClips = data;
			})
			.catch((e) => {
				return e;
			})
	]);

	if (!clip.id) {
		throw error(404, 'Not found');
	}

	await Promise.all([
		pb
			.collection('clip')
			.getList(1, 1, {
				sort: '-date',
				filter: `viewcount >= ${clip.viewcount}`,
				requestKey: 'clip_position'
			})
			.then((data) => {
				clipPosition = data;
			})
			.catch((e) => {
				return e;
			}),
		pb
			.collection('clip')
			.getList(1, 12, {
				sort: '-viewcount',
				filter: `date >= '${format(
					add(parseISO(clip.date), { months: -2 }),
					'yyyy-MM-dd'
				)}' && date < '${format(add(parseISO(clip.date), { months: +2 }), 'yyyy-MM-dd')}'`,
				requestKey: 'clip_recommendations'
			})
			.then((data) => {
				recommendations = data;
			})
			.catch((e) => {
				return e;
			})
	]);

	return structuredClone({
		clip: clip,
		clipsCount: allClips.totalItems,
		clipPosition: clipPosition.totalItems,
		recommendations: recommendations.items.filter((v: RecordModel) => v.id !== clip.id)
	});
}
