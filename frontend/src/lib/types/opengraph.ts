export type OpenGraph = {
	title: string;
	description: string;
	image: string;
	image_width: number;
	image_height: number;
	updated_time: string;
};

export const DefaultOpenGraph: OpenGraph = {
	title: 'Wubbl0rz Twitch Archiv',
	description:
		'Twitch Archiv von m4xfps aka. wubbl0rz. Hier findest du alle VODs und Clips seit 2017. Egal ob Aufwachstreams am Sonntag oder Technik- und Programmierstreams, hier findest du alles. Die Volltextsuche hilft dir dabei das richtige Video zu finden.',
	image: '/og.webp',
	image_width: 1536,
	image_height: 860,
	updated_time: new Date().toISOString()
};
