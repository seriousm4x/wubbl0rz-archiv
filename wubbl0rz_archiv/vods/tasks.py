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
    info_dict = dl.get_vod_infos()

    vod_count_total = 0

    for entry in info_dict["entries"]:
        if entry["live_status"] == "is_live" or Vod.objects.filter(filename=entry["id"]).exists():
            continue

        vod_count_total += 1

        if not os.path.isfile(os.path.join(vod_dir, entry["id"] + ".json")):
            with open(os.path.join(vod_dir, entry["id"] + ".json"), "w", encoding="utf-8") as f:
                json.dump(entry, f)

        dl.download(vod_dir, entry["webpage_url"], entry['id'])
        dl.dl_post_processing(vod_dir, entry["id"])
        duration, resolution, filesize = dl.get_metadata(
            vod_dir, entry["id"])
        dl.create_thumbnail(vod_dir, entry["id"], duration)
        dl.update_vod(entry["id"], entry["title"], duration,
                       entry["timestamp"], resolution, entry["fps"], filesize)

    print("Vods downloaded:", vod_count_total)
