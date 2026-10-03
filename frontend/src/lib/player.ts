import type { VideoPlayerElement } from '@videojs/html/video';
import type { AudioPlayerElement } from '@videojs/html/audio';
import { selectTextTrack } from '@videojs/html';

export type ArchivePlayerElement = VideoPlayerElement | AudioPlayerElement;

type Preferences = {
	volume?: number;
	muted?: boolean;
	rate?: number;
	lang?: string;
	captions?: boolean;
};

function readStorage(key: string): string | null {
	try {
		return localStorage.getItem(key);
	} catch {
		return null;
	}
}

function writeStorage(key: string, value: string) {
	try {
		localStorage.setItem(key, value);
	} catch {
		// Playback remains available when browser storage is blocked.
	}
}

// Keep Vidstack's keys so existing preferences and resume positions survive migration.
export function connectPlayer(
	player: ArchivePlayerElement,
	media: HTMLMediaElement,
	src: string,
	getTimestamp: () => number,
	onTime: (time: number) => void,
	onRate?: (rate: number) => void
) {
	let saved: Preferences = {};
	try {
		saved = JSON.parse(readStorage('vds-player') ?? '{}') ?? {};
	} catch {
		// Ignore invalid saved preferences.
	}
	const timeKey = `${src}:0:0`;
	const savedTime = Number(readStorage(timeKey));
	const store = player.store;
	let restored = false;
	let positionRestored = false;
	let restoringPosition = false;
	let captionsRestored = false;
	let lastCaptions = '';
	let lastPreferences = '';
	let lastSavedSecond = -1;

	const saveTime = () => {
		if (!positionRestored || media.readyState === 0) return;
		writeStorage(timeKey, String(media.ended ? 0 : media.currentTime));
	};
	const updateTime = () => {
		onTime(media.currentTime);
		const second = Math.floor(media.currentTime);
		if (second !== lastSavedSecond) {
			lastSavedSecond = second;
			saveTime();
		}
	};
	const sync = () => {
		if (!store.target) return;
		if (!restored) {
			restored = true;
			if (Number.isFinite(saved.volume) && saved.volume! >= 0 && saved.volume! <= 1)
				store.setVolume(saved.volume!);
			if (typeof saved.muted === 'boolean') store.setMuted(saved.muted);
			if (Number.isFinite(saved.rate) && saved.rate! >= 0.25 && saved.rate! <= 4)
				store.setPlaybackRate(saved.rate!);
		}
		if (!positionRestored && !restoringPosition && media.readyState >= 1) {
			restoringPosition = true;
			const timestamp = getTimestamp();
			const time = timestamp > 0 ? timestamp : savedTime;
			if (Number.isFinite(time) && time > 0) {
				void store
					.seek(time)
					.catch(() => {})
					.finally(() => {
						positionRestored = true;
						updateTime();
					});
			} else positionRestored = true;
		}

		// Audio has no caption feature. Preserve the video caption preference in audio mode.
		const textTracks = selectTextTrack(store.state);
		if (textTracks) {
			const tracks = textTracks.textTrackList.filter(
				(track) => track.kind === 'subtitles' || track.kind === 'captions'
			);
			if (tracks.length) {
				if (!captionsRestored) {
					captionsRestored = true;
					const track = tracks.find((track) => track.language === saved.lang) ?? tracks[0];
					if (typeof saved.captions === 'boolean')
						textTracks.selectSubtitlesTrack(saved.captions ? track.id : null);
					lastCaptions = JSON.stringify([saved.captions, saved.lang]);
					return;
				}
				const selected = tracks.find((track) => track.mode === 'showing');
				const selection = JSON.stringify([!!selected, selected?.language ?? saved.lang]);
				if (selection !== lastCaptions) {
					lastCaptions = selection;
					saved.captions = !!selected;
					if (selected) saved.lang = selected.language;
				}
			}
		}
		onRate?.(store.playbackRate);
		saved = { ...saved, volume: store.volume, muted: store.muted, rate: store.playbackRate };
		const preferences = JSON.stringify(saved);
		if (preferences !== lastPreferences) {
			lastPreferences = preferences;
			writeStorage('vds-player', preferences);
		}
	};

	const unsubscribe = store.subscribe(sync);
	sync();
	media.addEventListener('timeupdate', updateTime);
	media.addEventListener('loadedmetadata', sync);
	media.addEventListener('pause', saveTime);
	media.addEventListener('seeked', saveTime);
	media.addEventListener('ended', saveTime);
	window.addEventListener('pagehide', saveTime);
	return () => {
		saveTime();
		unsubscribe();
		media.removeEventListener('timeupdate', updateTime);
		media.removeEventListener('loadedmetadata', sync);
		media.removeEventListener('pause', saveTime);
		media.removeEventListener('seeked', saveTime);
		media.removeEventListener('ended', saveTime);
		window.removeEventListener('pagehide', saveTime);
	};
}
