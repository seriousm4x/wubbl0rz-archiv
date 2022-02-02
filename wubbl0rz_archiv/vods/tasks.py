import json
import os

from celery import shared_task
from django.conf import settings
from main.tasks import Downloader

from vods.models import Vod


@shared_task
def download_vods():
    vod_dir = os.path.join(settings.MEDIA_ROOT, "vods")
    dl = Downloader()
    print("getting info dict")
    info_dict = dl.get_vod_infos()
    for entry in info_dict["entries"]:
        if entry["live_status"] == "is_live" or Vod.objects.filter(filename=entry["id"]).exists() or os.path.isfile(os.path.join(vod_dir, entry["id"] + ".ts")):
            continue

        if not os.path.isfile(os.path.join(vod_dir, entry["id"] + ".json")):
            with open(os.path.join(vod_dir, entry["id"] + ".json"), "w", encoding="utf-8") as f:
                json.dump(entry, f)

        print(f"download vod: {entry['id']}")
        dl.download(vod_dir, entry["webpage_url"], entry['id'])
        print("post processing")
        dl.dl_post_processing(vod_dir, entry["id"])
        print("get metadata")
        duration, resolution, filesize = dl.get_metadata(
            vod_dir, entry["id"])
        print("create thumbnail")
        dl.create_thumbnail(vod_dir, entry["id"], duration)
        print("update db")
        dl.update_vod(entry["id"], entry["title"], duration,
                       entry["timestamp"], resolution, entry["fps"], filesize)
