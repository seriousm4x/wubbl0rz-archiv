export type watched = {
	vods: {
		[key: string]: number;
	};
	clips: {
		[key: string]: number;
	};
};

export type bookmarks = {
	vods: [string];
	clips: [string];
};
