<div align="center" width="100%">
    <img src="wubbl0rz_archiv/archiv/static/img/wubPog.png" width="128"/>
</div>

<div align="center" width="100%">
    <h2>wubbl0rz VOD Archiv</h2>
    <p>Stack: Django, Celery, Redis, yt-dlp, Pillow, FFmpeg, Bootstrap</p>
</div>

## 🚀 Features

* Komplett autonome Webanwendung mit background tasks die sich um runterladen, files generieren und import kümmern
* VODs zurück bis 2017
* Ansichten: Kürzliche Streams, Suche, Jahresansicht
* Live Suche beim tippen
* Dark/Light Mode mit toggle/preferes-color-scheme/LocalStorage
* 100% Cross Platform kompatibel (.ts/h264 Videos mit .avif Thumbnails und.jpg Fallback)
* Live status check mit Twitchstream embed
* Statistiken
* Docker Image

## 🕒 Coming soon

* ❌ ~~Rechat Fenster neben dem VOD um den Chat mitlesen zu können~~ ([#3](https://github.com/AgileProggers/wubbl0rz-archiv/issues/3))

## 📸 Screenshots

| Dark | Light |
| ---- | ----- |
| ![](https://i.imgur.com/zYMnfly.png) | ![](https://i.imgur.com/j0DBl0T.png) |
| ![](https://i.imgur.com/Ln6yJwZ.png) | ![](https://i.imgur.com/mxP330u.png) |
| ![](https://i.imgur.com/DUTdbBY.png) | ![](https://i.imgur.com/1G0KKjq.png) |
| ![](https://i.imgur.com/t9iv9sM.png) | ![](https://i.imgur.com/xRfL6sh.png) |
| ![](https://i.imgur.com/0Pzx7UF.png) | ![](https://i.imgur.com/lRvGmqc.png) |

## 🐳 Deploy

### IMPORTANT VARIABLES TO CHANGE

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

* DB_*

  If you change them, make sure to change both entries, in web and db.

* DB_BACKUP_DIR

  Used for storing database backups every 24h in .json format.

* TWITCH_CLIENT_ID and TWITCH_CLIENT_SECRET

  Create your api keys at https://dev.twitch.tv/console/apps

* WEB_STATIC

  Path where .js, .css, .jpg files live.

* WEB_MEDIA

  Path where vod files, .m3u8, thumbnails etc. live. Make sure to use a larger drive for WEB_MEDIA, as the dir size will grow over time.

Here is a example docker-compose.yml file.

```
version: "3"
services:
  web:
    container_name: wub-web
    image: ghcr.io/AgileProggers/wub-archiv:latest
    restart: unless-stopped
    ports:
      - 127.0.0.1:8000:${DJANGO_PORT}
    volumes:
      - ${WEB_STATIC}:/var/www/static/
      - ${WEB_MEDIA}:/var/www/media/
      - ${DB_BACKUP_DIR}:${DB_BACKUP_DIR}
    environment:
      - DJANGO_SUPERUSER_USER=<user>
      - DJANGO_SUPERUSER_PASSWORD=<password>
      - DJANGO_SECRET_KEY=<secret>
      - DJANGO_DEBUG=False
      - DJANGO_ALLOWED_HOSTS=localhost
      - DJANGO_LANGUAGE_CODE=de
      - DJANGO_TIME_ZONE=Europe/Berlin
      - DJANGO_PORT=8000
      - DB_HOST=db
      - DB_NAME=wub
      - DB_USER=wub
      - DB_PASSWORD=wub
      - DB_PORT=5432
      - DB_BACKUP_DIR=/path/to/backup_dir/
      - TWITCH_CLIENT_ID=<client-id>
      - TWITCH_CLIENT_SECRET=<client-secret>
      - WEB_STATIC=/path/to/static/
      - WEB_MEDIA=/path/to/media/
    depends_on:
      - db
      - redis
  db:
    container_name: wub-db
    image: postgres:14-alpine
    restart: always
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
    healthcheck:
    test: redis-cli ping
    interval: 10s

volumes:
  wub_db:
```

## 🚪 Reverse Proxy

The Django app won't serve static and media files. A reverse proxy is needed. The easiest way is to use caddy. Here is an example config. Change the root path to the same as WEB_MEDIA in the `docker-compose.yml`. Then run `caddy run` from the same directory.

```
:8001 {
  root * /path/to/media/
  @notStatic {
    not path /static/* /media/*
  }
  reverse_proxy @notStatic :8000
  file_server
  encode gzip
}
```
