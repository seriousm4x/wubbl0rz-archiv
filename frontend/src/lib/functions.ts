import type { Emotes } from './emotes';

export function formatBytes(bytes: number, decimals = 2) {
	if (bytes === 0) return '0 Bytes';
	const k = 1024;
	const dm = decimals < 0 ? 0 : decimals;
	const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
	const i = Math.floor(Math.log(bytes) / Math.log(k));
	return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)).toLocaleString('de-DE') + ' ' + sizes[i];
}

export function toHHMMSS(t: number, showUnit: boolean): string {
	const hours = Math.floor(t / 3600);
	const minutes = Math.floor((t - hours * 3600) / 60);
	const seconds = Math.round(t - hours * 3600 - minutes * 60);
	const h = hours < 10 ? '0' + hours : hours;
	const m = minutes < 10 ? '0' + minutes : minutes;
	const s = seconds < 10 ? '0' + seconds : seconds;
	const mmss = m + ':' + s;
	if (hours == 0) {
		return showUnit ? mmss + 'm' : mmss;
	} else {
		return showUnit ? h + ':' + mmss + 'h' : h + ':' + mmss;
	}
}

export function replaceEmotesInString(
	message: string,
	emotes: Emotes,
	regexPattern: RegExp
): string {
	// decode url encoded string
	let decodedmsg: string;
	try {
		decodedmsg = decodeURIComponent(message);
	} catch {
		decodedmsg = message;
	}

	// and finally replace our emotes with html
	const modifiedString = decodedmsg.replace(regexPattern, (match) => {
		const replacement = `<div class="tooltip" data-tip="${match}"><img src="${
			emotes[match.toLowerCase()]
		}" alt="${match}" loading="lazy" style="height: 2em;" /></div>`;
		return replacement || match;
	});

	return modifiedString;
}
