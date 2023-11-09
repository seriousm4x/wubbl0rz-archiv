import PocketBase, { type RecordModel } from 'pocketbase';
import { error } from '@sveltejs/kit';
import { PUBLIC_API_URL } from '$env/static/public';
import { PRIVATE_ALLOW_SEARCH_INDEXING } from '$env/static/private';
import { PUBLIC_FRONTEND_URL } from '$env/static/public';
import parseISO from 'date-fns/parseISO/index.js';

export async function GET() {
    if (PRIVATE_ALLOW_SEARCH_INDEXING !== 'true') {
        throw error(404, 'not found');
    }

    const pb = new PocketBase(PUBLIC_API_URL);
    const [vods, clips, stats] = await Promise.all([
        // all vods
        pb
            .collection('vod')
            .getFullList({
                requestKey: 'all_vods'
            })
            .catch((e) => {
                return e;
            }),

        // all clips
        pb
            .collection('clip')
            .getFullList({
                requestKey: 'all_clips'
            })
            .catch((e) => {
                return e;
            }),

        // stats
        fetch(`${PUBLIC_API_URL}/stats`)
            .then((response) => response.json())
            .catch((e) => {
                return e;
            }),
    ]);

    const head = `
    <?xml version="1.0" encoding="UTF-8" ?>
    <urlset
      xmlns="https://www.sitemaps.org/schemas/sitemap/0.9"
      xmlns:xhtml="https://www.w3.org/1999/xhtml"
      xmlns:mobile="https://www.google.com/schemas/sitemap-mobile/1.0"
      xmlns:news="https://www.google.com/schemas/sitemap-news/0.9"
      xmlns:image="https://www.google.com/schemas/sitemap-image/1.1"
      xmlns:video="https://www.google.com/schemas/sitemap-video/1.1"
    >
    <url>
        <loc>${PUBLIC_FRONTEND_URL}</loc>
        <lastmod>${stats.last_update}</lastmod>
        <changefreq>hourly</changefreq>
        <priority>1.0</priority>
    </url>
    <url>
        <loc>${PUBLIC_FRONTEND_URL}/vods</loc>
        <lastmod>${stats.last_update}</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.99</priority>
    </url>
    <url>
        <loc>${PUBLIC_FRONTEND_URL}/clips</loc>
        <lastmod>${stats.last_update}</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.99</priority>
    </url>
    <url>
        <loc>${PUBLIC_FRONTEND_URL}/stats</loc>
        <lastmod>${stats.last_update}</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.98</priority>
    </url>`;
    let body = '';

    let vod_priority = 98;
    vods.forEach((vod: RecordModel) => {
        vod_priority -= 1;
        vod_priority = vod_priority > 0 ? vod_priority : 0.0;
        body =
            body +
            `
        <url>
            <loc>${PUBLIC_FRONTEND_URL}/vods/${vod.id}</loc>
            <lastmod>${parseISO(vod.date).toISOString()}</lastmod>
            <changefreq>monthly</changefreq>
            <priority>${vod_priority / 100}</priority>
        </url>
        `;
    });

    let clip_priority = 98;
    clips.forEach((clip: RecordModel) => {
        clip_priority -= 1;
        clip_priority = clip_priority > 0 ? clip_priority : 0.0;
        body =
            body +
            `
        <url>
            <loc>${PUBLIC_FRONTEND_URL}/clips/${clip.id}</loc>
            <lastmod>${parseISO(clip.date).toISOString()}</lastmod>
            <changefreq>monthly</changefreq>
            <priority>${clip_priority / 100}</priority>
        </url>
        `;
    });
    body = body + `</urlset>`;

    return new Response(head.trim() + body.trim(), {
        headers: {
            'Content-Type': 'application/xml'
        }
    });
}
