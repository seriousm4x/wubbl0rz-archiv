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
      - /path/to/media/:/var/www/
    depends_on:
      - db
  db:
    container_name: archiv-db
    image: postgres:14-alpine
    restart: unless-stopped
    env_file: .env
    volumes:
      - /path/to/postgres/:/var/lib/postgresql/data
```
## 🚪 Reverse Proxy

The easiest way is to use caddy. Paste the following into a file called `Caddyfile`. Change the `root` path to the parent directory of media (/path/to/) and run caddys system service.

```
api.wubbl0rz.tv {
    root * /path/to/
    @excludePaths {
        not path /media/*
    }
    reverse_proxy @excludePaths localhost:5000
    file_server
    encode gzip
    header * Cache-Control max-age=1
    header Access-Control-Allow-Origin "*"
    header Cross-Origin-Resource-Policy "*"
}
```

## Documentation

Check out [the swagger page](https://api.wubbl0rz.tv/swagger/index.html) for more documentation.
