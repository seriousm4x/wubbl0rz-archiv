# API documentation

The api is build with the [django rest framework](https://www.django-rest-framework.org/). Checkout the endpoint at [archiv.wubbl0rz.tv/api/](https://archiv.wubbl0rz.tv/api/) or continue reading.

## GET `/api/`

Returns all available routes.

```json
{
    "vods": "https://archiv.wubbl0rz.tv/api/vods/",
    "years": "https://archiv.wubbl0rz.tv/api/years/",
    "emotes": "https://archiv.wubbl0rz.tv/api/emotes/",
    "stats": "https://archiv.wubbl0rz.tv/api/stats/"
}
```

## GET `/api/vods/`

Returns vod infos.

```json
{
    "count": 620,
    "next": "https://archiv.wubbl0rz.tv/api/vods/?page=2",
    "previous": null,
    "results": [
        {
            "uuid": "c7e167",
            "title": "🚀 wir contributen zu open source projekt auf github 👀",
            "duration": 14050,
            "date": "2022-01-04T18:25:57+01:00",
            "filename": "v1252822335",
            "resolution": "1920x1080",
            "fps": 47.997,
            "size": 9182102360
        },
        {
            "uuid": "975406",
            "title": "😴 Aufwachen Stream 😴",
            "duration": 9517,
            "date": "2022-01-02T09:21:09+01:00",
            "filename": "v1250522854",
            "resolution": "1920x1080",
            "fps": 47.995,
            "size": 6212785240
        }
    ]
}
```
### Examples

* `/api/vods/` - returns all vods with default page size of 50
* `/api/vods/c7e167` - returns single vod by uuid
* `/api/vods/?year=2019` - returns vods filtered by year
* `/api/vods/?page_size=100` - extend results to 100 (max. 500)
* `/api/vods/?page=2` - returns page 2

## GET `/api/years/`

Returns years with vod count.

```json
[
    {
        "year": 2022,
        "count": 2
    },
    {
        "year": 2021,
        "count": 195
    },
    {
        "year": 2020,
        "count": 210
    }
]
```

## GET `/api/emotes/`

Returns twitch, bttv and ffz emotes.

```json
{
    "count": 151,
    "next": "https://archiv.wubbl0rz.tv/api/emotes/?page=2",
    "previous": null,
    "results": [
        {
            "id": "139407",
            "name": "LULW",
            "url": "https://cdn.frankerfacez.com/emote/139407/4",
            "provider": "ffz"
        },
        {
            "id": "14761",
            "name": "HotGrill",
            "url": "https://cdn.frankerfacez.com/emote/14761/4",
            "provider": "ffz"
        }
    ]
}
```

### Examples

* `/api/emotes/` - returns all vods with default page size of 50
* `/api/emotes/187256` - returns single emote by id
* `/api/emotes/?provider=ffz` - returns emotes filtered by provider. Possible providers are `twitch`, `ffz` and `bttv`.
* `/api/emotes/?name=KEKW` - returns emotes filtered by name
* `/api/emotes/?page_size=100` - extend results to 100 (max. 500)
* `/api/emotes/?page=2` - returns page 2

## GET `/api/stats`

Returns stats about the archive.

```json
{
    "count_vods_total": 620,
    "count_vods_1m": 21,
    "count_h_streamed": 1695,
    "archiv_size_bytes": 4268057583164,
    "vods_per_month": [
        {
            "month": "Feb 21",
            "count": 12
        },
        {
            "month": "Mär 21",
            "count": 15
        },
        {
            "month": "Apr 21",
            "count": 17
        }
    ],
    "vods_per_weekday": [
        {
            "weekday": "Sonntag",
            "count": 178
        },
        {
            "weekday": "Montag",
            "count": 105
        },
        {
            "weekday": "Dienstag",
            "count": 61
        }
    ],
    "start_by_time": [
        {
            "hour": 8,
            "count": 3
        },
        {
            "hour": 9,
            "count": 57
        },
        {
            "hour": 10,
            "count": 61
        }
    ]
}
```

## GET `/api/stats/db`

Returns info about last vod and emote update in UTC.

```json
{
    "last_vod_sync": "2022-01-06T02:00:00.038000Z",
    "last_emote_sync": "2022-01-06T01:00:00.101000Z"
}
```
