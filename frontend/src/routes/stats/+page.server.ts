import { PRIVATE_MEILI_ADMIN_KEY } from '$env/static/private';
import { PUBLIC_API_URL, PUBLIC_MEILI_URL } from '$env/static/public';
import { pb } from '$lib/pocketbase.js';
import { error } from '@sveltejs/kit';

export async function load({ fetch }) {
	const [stats, meili, emotes] = await Promise.all([
		// pocketbase stats
		fetch(`${PUBLIC_API_URL}/stats`)
			.then((response) => response.json())
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
			.catch((e) => {
				return e;
			}),

		// all emotes
		pb
			.collection('emote')
			.getFullList({
				requestKey: 'all_emotes'
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

	if (data) {
		return data;
	}

	throw error(404, 'Not found');
}
