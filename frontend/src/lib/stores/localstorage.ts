import { browser } from '$app/environment';
import { writable } from 'svelte/store';

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
	if (browser) {
		localStorage.watched = JSON.stringify(value);
	}
});
