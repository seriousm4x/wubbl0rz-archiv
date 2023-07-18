<div align="center" width="100%">
    <img src="https://raw.githubusercontent.com/AgileProggers/archiv-frontend/master/static/favicon.ico" width="128"/>
</div>

<div align="center" width="100%">
    <h2>wubbl0rz VOD Archiv</h2>
    <p>Stack: Go, Gin, Gorm, FFmpeg</p>
</div>

## 🐳 Deploy

Copy `.env.sample` to `.env` and replace the required variables.

### Example `docker-compose.yml`

```
version: '3'
services:
  api:
    container_name: archiv-api
    build: .
    restart: unless-stopped
    env_file: .env
    ports:
      - 127.0.0.1:5000:5000
    volumes:
      - /etc/timezone:/etc/timezone:ro
      - /etc/localtime:/etc/localtime:ro
      - /path/to/media/:/var/www/media/
    depends_on:
      - db
  db:
    container_name: archiv-db
    image: postgres:15-alpine
    restart: unless-stopped
    env_file: .env
    volumes:
      - /etc/timezone:/etc/timezone:ro
      - /etc/localtime:/etc/localtime:ro
      - /path/to/postgres/:/var/lib/postgresql/data
  meili:
    container_name: archiv-meili
    image: getmeili/meilisearch:v1.1
    env_file: .env
    ports:
      - 127.0.0.1:7700:7700
    volumes:
      - ./meili_data:/meili_data
    restart: unless-stopped
```

## 🚪 Reverse Proxy

The easiest way is to use caddy. Paste the following into a file called `Caddyfile`.

```
api.wubbl0rz.tv {
    reverse_proxy localhost:5000
    encode zstd gzip
    header Access-Control-Allow-Origin "*"
}
```

## Backup and restore database

Required to upgrade major postgres versions.

**Backup**
`docker exec -t archiv-db pg_dumpall -c -U YOUR_DB_USER > /path/to/backup/dump_$(date +%Y-%m-%d"_"%H_%M_%S).sql`

**Restore**
`docker exec -i archiv-db psql -d YOUR_DB_NAME -U YOUR_DB_USER < /path/to/backup/dump_<some-date>.sql`

## 🔎 Meilisearch

Meilisearch index is filled with [archiv-transcribe](https://github.com/AgileProggers/archiv-transcribe).

A custom config is requires for our indexes. [Use the api](https://docs.meilisearch.com/reference/api/settings.html#update-settings) to patch the index settings like so:

> PATCH `http://localhost:7700/indexes/transcripts/settings/`

```json
{
	"displayedAttributes": [
		"*"
	],
	"searchableAttributes": [
		"text"
	],
	"filterableAttributes": [],
	"sortableAttributes": [
		"date",
		"duration",
		"viewcount"
	],
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
	"displayedAttributes": [
		"*"
	],
	"searchableAttributes": [
		"title"
	],
	"filterableAttributes": [],
	"sortableAttributes": [
		"date",
		"duration",
		"viewcount"
	],
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
	"displayedAttributes": [
		"*"
	],
	"searchableAttributes": [
		"*"
	],
	"filterableAttributes": [],
	"sortableAttributes": [
		"date",
		"duration",
		"viewcount"
	]
}
```

> PATCH `http://localhost:7700/indexes/transcripts`

```json
{
	"primaryKey": "id"
}
```


## Documentation

Docs are generated with swagger. Install `swag` to build docs:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Then build:

```bash
swag i
```

Check out [the swagger page](https://api.wubbl0rz.tv/swagger/index.html) for the build documentation.
