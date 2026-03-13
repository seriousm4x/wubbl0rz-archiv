import { PUBLIC_API_URL } from '$env/static/public';
import type { RequestHandler } from '@sveltejs/kit';

export const GET: RequestHandler = async ({ params, url, setHeaders }) => {
	if (params.type !== 'clip' && params.type !== 'vod') {
		return new Response('Invalid type', { status: 400 });
	}

	const fileName =
		params.type === 'vod'
			? url.searchParams.get('audio') === 'true'
				? 'audio.ogg'
				: 'vod.mp4'
			: 'clip.mp4';

	const res = await fetch(`${PUBLIC_API_URL}/${params.type}s/${params.filename}/${fileName}`);

	if (!res.ok) {
		return new Response('Not found', { status: 404 });
	}

	const contentType = res.headers.get('content-type') ?? 'application/octet-stream';
	const contentLength = res.headers.get('content-length');

	const headers: Record<string, string> = {
		'Content-Type': contentType,
		'Content-Disposition': `attachment; filename="${fileName}"`
	};

	// Nur setzen, wenn contentLength existiert
	if (contentLength !== null) {
		headers['Content-Length'] = contentLength;
	}

	setHeaders(headers);

	return new Response(res.body);
};
