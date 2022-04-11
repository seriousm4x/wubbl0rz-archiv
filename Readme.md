<div align="center" width="100%">
    <img src="https://archiv.wubbl0rz.tv/favicon.ico" width="128"/>
</div>

<div align="center" width="100%">
    <h2>wubbl0rz VOD Archiv</h2>
    <p>Stack: Django, Celery, Redis, yt-dlp, FFmpeg, Bootstrap</p>
</div>

## 🐳 Deploy

#### IMPORTANT VARIABLES TO CHANGE

#### Volumes

* static: Path where .js, .css, .jpg files live (small files).

* media: Path where vod files, .m3u8, thumbnails etc. live. Make sure to use a larger drive, as the dir size will grow over time.

* backups: Used for storing database backups every 24h in .json format.

#### Environment

* DJANGO_SUPERUSER_USER

  Defines the django user name. Can be used to log into `/admin/`.

* DJANGO_SUPERUSER_PASSWORD

  Defines the django user password. Can be used to log into `/admin/`.

* DJANGO_SECRET_KEY

  The prefered way is to generate it with:

  ```
  from django.core.management.utils import get_random_secret_key
  get_random_secret_key()
  ```

  but you can also use [https://djecrety.ir/](https://djecrety.ir/).

* DJANGO_DOMAIN

  Set this to your domain.

* DB_*

  If you change them, make sure to change both entries, in web and db.

* TWITCH_CLIENT_ID and TWITCH_CLIENT_SECRET

  Create your api keys at https://dev.twitch.tv/console/apps

### Example `docker-compose.yml`

```
version: "3"
services:
  web:
    container_name: wub-api
    image: ghcr.io/agileproggers/archiv-frontend:latest
    restart: unless-stopped
    ports:
      - 127.0.0.1:8000:8000
    volumes:
      - /path/to/static/:/var/www/static/
      - /path/to/media/:/var/www/media/
      - /path/to/backups/:/backups/
    environment:
      - DJANGO_SUPERUSER_USER=<user>
      - DJANGO_SUPERUSER_PASSWORD=<password>
      - DJANGO_SECRET_KEY=<secret>
      - DJANGO_DOMAIN=your-domain.com
      - DJANGO_DEBUG=False
      - DJANGO_LANGUAGE_CODE=de
      - DJANGO_TIME_ZONE=Europe/Berlin
      - DB_HOST=db
      - DB_NAME=wub
      - DB_USER=wub
      - DB_PASSWORD=wub
      - DB_PORT=5432
      - TWITCH_CLIENT_ID=<client-id>
      - TWITCH_CLIENT_SECRET=<client-secret>
    depends_on:
      - db
      - redis
  db:
    container_name: wub-db
    image: postgres:14-alpine
    restart: unless-stopped
    environment:
      - "POSTGRES_USER=wub"
      - "POSTGRES_PASSWORD=wub"
      - "POSTGRES_DB=wub"
    volumes:
      - wub_db:/var/lib/postgresql/data
    healthcheck:
      test: pg_isready -U wub
      interval: 10s
  redis:
    container_name: wub-redis
    image: redis:6-alpine
    restart: unless-stopped
    command: redis-server --loglevel warning
    healthcheck:
      test: redis-cli ping
      interval: 10s

volumes:
  wub_db:
```

### Tags

Available tags are

* `latest`
* `v{MAJOR}`
* `v{MAJOR}.{MINOR}.{PATCH}`

## 🚪 Reverse Proxy

The Django app won't serve static and media files. A reverse proxy is needed. The easiest way is to use caddy. Paste the following into a file called `Caddyfile`. Change the root path to the parent directory of media and static files from `docker-compose.yml` (/path/to/). Then run `caddy run` from the same directory.

```
:8001 {
  root * /path/to/
  @notStatic {
    not path /static/* /media/*
  }
  reverse_proxy @notStatic :8000
  file_server
  encode gzip
}
```

## Documentation

Check out [the docs folder](https://github.com/AgileProggers/archiv-frontend/tree/master/docs) for more documentation.

## Contributing

Any help is always appreciated. Especially if you know Django and Javascript.
