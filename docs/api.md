# API documentation

The api is build with the [django rest framework](https://www.django-rest-framework.org/). Checkout the endpoint at [archiv.wubbl0rz.tv/api/](https://archiv.wubbl0rz.tv/api/) or continue reading.

## GET `/api/`

Returns all available routes.

```json
{
    "vods": "https://archiv.wubbl0rz.tv/api/vods/",
    "clips": "https://archiv.wubbl0rz.tv/api/clips/",
    "years": "https://archiv.wubbl0rz.tv/api/years/",
    "emotes": "https://archiv.wubbl0rz.tv/api/emotes/",
    "stats": "https://archiv.wubbl0rz.tv/api/stats/",
    "stats/db": "https://archiv.wubbl0rz.tv/api/stats/db/"
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
            "uuid": "176658",
            "title": "😴 Aufwachen Stream 😴",
            "duration": 12584,
            "date": "2022-01-30T08:58:22+01:00",
            "filename": "v1281016761",
            "resolution": "1920x1080",
            "fps": 47.992,
            "size": 8219765328,
            "clip_set": [
                "5e8869",
                "215b35",
                "437d5b"
            ]
        },
        {
            "uuid": "99a11f",
            "title": "😊 Svelte Autocomplete + PostgreSQL Text Search 😎🔍 | C# / JS",
            "duration": 14959,
            "date": "2022-01-24T13:44:39+01:00",
            "filename": "v1274355715",
            "resolution": "1920x1080",
            "fps": 47.997,
            "size": 9776643524,
            "clip_set": [
                "66118c"
            ]
        }
    ]
}
```

## GET `/api/clips/`

Returns clip infos.

```json
{
    "count": 2084,
    "next": "https://archiv.wubbl0rz.tv/api/clips/?page=2",
    "previous": null,
    "results": [
        {
            "uuid": "437d5b",
            "title": "KEKW",
            "clip_id": "AthleticObedientCrocodileCclamChamp-Sc-YFtu08q0BNApT",
            "creator": "maracujaxd",
            "view_count": 10,
            "date": "2022-01-30T11:43:58+01:00",
            "duration": 27,
            "resolution": "1920x1080",
            "size": 12264180,
            "game": "Just Chatting",
            "vod": "176658"
        },
        {
            "uuid": "5e8869",
            "title": "😴 Aufwachen Stream 😴",
            "clip_id": "DignifiedPolitePelicanCclamChamp-HtIEVCKNQId0lpq7",
            "creator": "niceiiiiii",
            "view_count": 3,
            "date": "2022-01-30T10:25:17+01:00",
            "duration": 29,
            "resolution": "1920x1080",
            "size": 19584336,
            "game": "Just Chatting",
            "vod": "176658"
        }
    ]
}
```
### Examples

* `/api/clips/` - returns all clips with default page size of 50
* `/api/clips/437d5b` - returns single clip by uuid
* `/api/clips/?title=kekw` - returns vods with title containing parameter
* `/api/clips/?year=2019` - returns vods filtered by year
* `/api/clips/?page_size=100` - extend results to 100 (max. 500)
* `/api/clips/?page=2` - returns page 2

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
