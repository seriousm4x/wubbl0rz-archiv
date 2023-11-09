import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export type WatchHistory = {
	vods: {
		[id: string]: number;
	};
	clips: {
		[id: string]: number;
	};
};

const defaultWatchHistory = {
	vods: {},
	clips: {}
} as WatchHistory;

const lsWatchHistory = browser
	? (JSON.parse(localStorage.watched || JSON.stringify(defaultWatchHistory)) as WatchHistory) ||
	  defaultWatchHistory
	: defaultWatchHistory;
export const watchHistory = writable(lsWatchHistory);

// update localstorage on change
watchHistory.subscribe((value) => {
	browser ? (localStorage.watched = JSON.stringify(value)) : '';
});
