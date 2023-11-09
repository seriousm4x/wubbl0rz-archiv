<div align="center" width="100%">
    <img src="https://raw.githubusercontent.com/AgileProggers/archiv-frontend/master/static/favicon.ico" width="128"/>
</div>

<div align="center" width="100%">
    <h2>wubbl0rz VOD Archiv</h2>
    <p>Backend Stack: Go, PocketBase, Echo, FFmpeg</p>
    <p>Frontend Stack: SvelteKit, Tailwind (DaisyUI), VidStack, PNPM</p>
</div>

## 🐳 Deploy

-   Check the `.env.sample` files in [frontend/](frontend/) and [backend/](backend/) and copy them to `*/.env`.
-   `docker-compose up`

## 🚪 Reverse Proxy

The easiest way is to use caddy. Paste the following into a file called `Caddyfile`.

```
wubbl0rz.tv {
    reverse_proxy localhost:8090
    encode zstd gzip
    header Cache-Control "max-age=31536000"
}
api.wubbl0rz.tv {
    reverse_proxy localhost:8090
    encode zstd gzip
    header /media/* Cache-Control "max-age=31536000"
    header Access-Control-Allow-Origin "*"
}
meili.wubbl0rz.tv {
    reverse_proxy localhost:7700
}
```

## 🔎 Meilisearch

Meilisearch index is filled with [wubbl0rz-archiv-transcribe](https://github.com/seriousm4x/wubbl0rz-archiv-transcribe).

A custom config is required for our indexes. [Use the api](https://docs.meilisearch.com/reference/api/settings.html#update-settings) to patch the index settings like so:

> PATCH `http://localhost:7700/indexes/transcripts/settings/`

```json
{
    "displayedAttributes": ["*"],
    "searchableAttributes": ["text"],
    "filterableAttributes": [],
    "sortableAttributes": ["date", "duration", "viewcount"],
    "rankingRules": [
        "sort",
        "words",
        "typo",
        "proximity",
        "attribute",
        "exactness"
    ]
}
```

> PATCH `http://localhost:7700/indexes/vods/settings/`

```json
{
    "displayedAttributes": ["*"],
    "searchableAttributes": ["title"],
    "filterableAttributes": [],
    "sortableAttributes": ["date", "duration", "viewcount"],
    "rankingRules": [
        "sort",
        "words",
        "typo",
        "proximity",
        "attribute",
        "exactness"
    ]
}
```

> PATCH `http://localhost:7700/indexes/clips/settings/`

```json
{
    "displayedAttributes": ["*"],
    "searchableAttributes": ["*"],
    "filterableAttributes": [],
    "sortableAttributes": ["date", "duration", "viewcount"]
}
```

> PATCH `http://localhost:7700/indexes/transcripts`

```json
{
    "primaryKey": "id"
}
```
