import { createInstance } from './stores/pocketbase';
const pb = createInstance();

export type Emotes = {
	[name: string]: string;
};

export async function getEmotes(): Promise<[Emotes, RegExp]> {
	const finalEmotes: Emotes = {};

	const pbEmotes = await pb.collection('emote').getFullList({ requestKey: null });
	pbEmotes.forEach((emote) => {
		finalEmotes[emote.name.toLowerCase()] = emote.url;
	});

	// sort emotes by length to match long emote names before short ones
	let emoteKeys = Object.keys(finalEmotes).sort((a, b) => b.length - a.length);

	// escape all emotes for regex: https://stackoverflow.com/a/6969486/6574444
	emoteKeys = emoteKeys.map((emote) => emote.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'));

	// Match emotes without matching normal names inside a larger word. Unlike \b,
	// this also supports punctuation-only emotes such as :) and <3.
	const re = emoteKeys.length
		? new RegExp(`(?<![\\p{L}\\p{N}_])(${emoteKeys.join('|')})(?![\\p{L}\\p{N}_])`, 'giu')
		: /$a/;

	return [finalEmotes, re];
}
