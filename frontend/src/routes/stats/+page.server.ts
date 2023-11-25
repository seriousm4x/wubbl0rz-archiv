import { PRIVATE_MEILI_ADMIN_KEY } from '$env/static/private';
import { PUBLIC_API_URL, PUBLIC_MEILI_URL } from '$env/static/public';
import { createInstance } from '$lib/stores/pocketbase.js';
import type { Stats } from 'meilisearch';
import type { RecordModel } from 'pocketbase';

export async function load({ fetch }) {
	const pb = createInstance();
	let emotes = [] as RecordModel[];
	let stats = {
		last_update: '',
		count_vods: 0,
		trend_vods: 0,
		count_clips: 0,
		trend_clips: 0,
		count_hours: 0,
		trend_hours: 0,
		count_size: 0,
		chatters: [
			{
				name: '',
				msg_count: 0
			}
		]
	};
	let meili = {} as Stats;

	await Promise.all([
		// pocketbase stats
		fetch(`${PUBLIC_API_URL}/stats`)
			.then((response) => response.json())
			.then((data) => {
				stats = data;
			})
			.catch((e) => {
				return e;
			}),

		// meilisearch stats
		fetch(`${PUBLIC_MEILI_URL}/stats`, {
			headers: {
				Authorization: `Bearer ${PRIVATE_MEILI_ADMIN_KEY}`
			}
		})
			.then((response) => response.json())
			.then((data) => {
				meili = data;
			})
			.catch((e) => {
				return e;
			}),

		// all emotes
		pb
			.collection('emote')
			.getFullList({
				requestKey: 'all_emotes'
			})
			.then((data: RecordModel[]) => {
				emotes = data;
			})
			.catch((e: Error) => {
				return e;
			})
	]);

	const data = {
		stats: stats,
		meili: {
			transcripts: meili.indexes.transcripts.numberOfDocuments,
			title: meili.indexes.vods.numberOfDocuments,
			lastUpdate: meili.lastUpdate,
			databaseSize: meili.databaseSize
		},
		emotes: emotes
	};

	return structuredClone(data);
}
