import PocketBase from 'pocketbase';
import { PUBLIC_API_URL } from '$env/static/public';

export type Emotes = {
	[name: string]: string;
};

type ttvEmote = {
	code: string;
	urls: [
		{
			url: string;
		}
	];
};

type bttvEmote = {
	id: string;
	code: string;
};

type ffzEmote = {
	id: string;
	name: string;
};

type seventvEmote = {
	id: string;
	name: string;
	data: {
		host: {
			files: [{ name: string }];
		};
	};
};

export async function getEmotes(): Promise<[Emotes, RegExp]> {
	const finalEmotes: Emotes = {};
	const pb = new PocketBase(PUBLIC_API_URL);

	await Promise.all([
		pb
			.collection('emote')
			.getFullList()
			.then((pbEmotes) => {
				pbEmotes.forEach((e) => {
					finalEmotes[e.name.toLowerCase()] = e.url;
				});
			}),

		fetch('https://emotes.adamcy.pl/v1/global/emotes/twitch')
			.then((response) => response.json())
			.then((data) => {
				data.forEach((e: ttvEmote) => {
					finalEmotes[e.code.toLowerCase()] = e.urls[0].url;
				});
			}),

		fetch('https://api.betterttv.net/3/cached/emotes/global')
			.then((response) => response.json())
			.then((data) => {
				data.forEach((e: bttvEmote) => {
					finalEmotes[e.code.toLowerCase()] = `https://cdn.betterttv.net/emote/${e.id}/1x`;
				});
			}),

		fetch('https://api.frankerfacez.com/v1/set/global')
			.then((response) => response.json())
			.then((data) => {
				data.sets[data.default_sets[0]].emoticons.forEach((e: ffzEmote) => {
					finalEmotes[e.name.toLowerCase()] = `https://cdn.frankerfacez.com/emoticon/${e.id}/1`;
				});
			}),

		fetch('https://7tv.io/v3/emote-sets/62cdd34e72a832540de95857')
			.then((response) => response.json())
			.then((data) => {
				data.emotes.forEach((e: seventvEmote) => {
					finalEmotes[
						e.name.toLowerCase()
					] = `https://cdn.7tv.app/emote/${e.id}/${e.data.host.files[0].name}`;
				});
			})
	]);

	// sort emotes by length to match long emote names before short ones
	let emoteKeys = Object.keys(finalEmotes).sort((a, b) => b.length - a.length);

	// escape all emotes for regex: https://stackoverflow.com/a/6969486/6574444
	emoteKeys = emoteKeys.map((emote) => emote.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'));

	// join all emotes in word boundaries and match case insensitive
	const re = new RegExp(`\\b(${emoteKeys.join('|')})\\b`, 'gi');

	return [finalEmotes, re];
}
