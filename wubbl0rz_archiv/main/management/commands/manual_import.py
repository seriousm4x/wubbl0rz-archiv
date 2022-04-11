import json
import os

from django.conf import settings
from django.core.management.base import BaseCommand
from main.tasks import Downloader


class Command(BaseCommand):
    def handle(self, **options):
        self.dl = Downloader()
        self.import_vods()
        self.import_clips()

    def import_vods(self):
        vod_dir = os.path.join(settings.MEDIA_ROOT, "vods")
        for f in os.listdir(vod_dir):
            if not f.endswith("-segments"):
                continue

            id, _ = f.split("-")
            print("importing vod:", id)
            with open(os.path.join(vod_dir, id + ".json"), "r", encoding="utf-8") as info:
                data = json.load(info)
            title = data["title"]
            timestamp = data["timestamp"]
            fps = data["fps"]
            duration, resolution, size = self.dl.get_metadata(vod_dir, id)
            self.dl.update_vod(id, title, duration, timestamp,
                               resolution, fps, size)

    def import_clips(self):
        clip_dir = os.path.join(settings.MEDIA_ROOT, "clips")
        for f in os.listdir(clip_dir):
            if not f.endswith("-segments"):
                continue
            id = f.replace("-segments", "")
            print("importing clip:", id)
            with open(os.path.join(clip_dir, id + ".json"), "r", encoding="utf-8") as info:
                data = json.load(info)
            duration, resolution, size = self.dl.get_metadata(clip_dir, id)
            data["duration"] = duration
            data["resolution"] = resolution
            data["size"] = size
            self.dl.update_clip(data)
