import { PUBLIC_API_URL } from '$env/static/public';

export type TwitchBadgeVersion = {
	image_url_1x: string;
	title: string;
};

export type TwitchBadges = Record<string, TwitchBadgeVersion>;

const badgeRequests = new Map<string, Promise<TwitchBadges>>();

export function getTwitchBadges(roomId: string): Promise<TwitchBadges> {
	const existingRequest = badgeRequests.get(roomId);
	if (existingRequest) return existingRequest;

	const url = new URL('/twitch/badges', PUBLIC_API_URL);
	url.searchParams.set('broadcaster_id', roomId);

	const request = fetch(url)
		.then((response) => {
			if (!response.ok) throw new Error('Unable to load Twitch badges');
			return response.json() as Promise<TwitchBadges>;
		})
		.catch(() => ({}));

	badgeRequests.set(roomId, request);
	return request;
}
