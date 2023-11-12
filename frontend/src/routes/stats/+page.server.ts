import { PRIVATE_MEILI_ADMIN_KEY } from '$env/static/private';
import { PUBLIC_API_URL, PUBLIC_MEILI_URL } from '$env/static/public';
import { pb } from '$lib/pocketbase.js';
import type { RecordModel } from 'pocketbase';

export async function load({ fetch }) {
	let emotes = [] as RecordModel[];
	let stats = {};
	let meili = {
		indexes: {
			transcripts: {
				numberOfDocuments: 0
			},
			vods: {
				numberOfDocuments: 0
			}
		},
		lastUpdate: Date.now(),
		databaseSize: 0
	};

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
			.then((data) => {
				emotes = data;
			})
			.catch((e) => {
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
