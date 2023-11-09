import add from 'date-fns/add/index.js';
import format from 'date-fns/format/index.js';
import parseISO from 'date-fns/parseISO/index.js';
import PocketBase, { type RecordModel } from 'pocketbase';
import { PUBLIC_API_URL } from '$env/static/public';

export async function load({ params }) {
	const pb = new PocketBase(PUBLIC_API_URL);
	const [clip, allClips] = await Promise.all([
		pb
			.collection('clip')
			.getOne(params.id, {
				expand: 'vod,creator,game',
				requestKey: 'single_clip'
			})
			.catch((e) => {
				return e;
			}),
		pb
			.collection('clip')
			.getList(1, 1, {
				requestKey: 'clip_count'
			})
			.catch((e) => {
				return e;
			})
	]);
	const [clipPosition, recommendations] = await Promise.all([
		pb
			.collection('clip')
			.getList(1, 1, {
				sort: '-date',
				filter: `viewcount >= ${clip.viewcount}`,
				requestKey: 'clip_position'
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
			.catch((e) => {
				return e;
			})
	]);

	return {
		clip: clip,
		clipsCount: allClips.totalItems,
		clipPosition: clipPosition.totalItems,
		recommendations: recommendations.items.filter((v: RecordModel) => v.id !== clip.id)
	};
}
