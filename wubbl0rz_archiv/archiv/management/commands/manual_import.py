import json
import os
import subprocess
from datetime import datetime

from archiv.models import Vod
from django.conf import settings
from django.core.management.base import BaseCommand
from django.utils.timezone import make_aware


class Command(BaseCommand):
    def handle(self, **options):
        for f in os.listdir(settings.MEDIA_ROOT):
            if not f.endswith(".ts"):
                continue

            id, _ = os.path.splitext(f)
            print("importing:", id)
            ts = os.path.join(settings.MEDIA_ROOT, id + ".ts")
            with open(os.path.join(settings.MEDIA_ROOT, id + ".json"), "r", encoding="utf-8") as info:
                data = json.load(info)
                title = data["title"]
                timestamp = data["timestamp"]
                fps = data["fps"]

            duration, resolution, filesize = self.get_metadata(ts)
            self.update_db(id, title, duration, timestamp,
                           resolution, fps, filesize)

    def get_metadata(self, ts):
        # duration
        cmd = ["ffprobe", "-v", "error", "-show_entries", "format=duration", "-of",
               "default=noprint_wrappers=1:nokey=1", ts]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        out, _ = proc.communicate()
        duration = float(out.decode().strip())

        # resolution
        cmd = ["ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries",
               "stream=width,height", "-of", "csv=s=x:p=0", ts]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        out, _ = proc.communicate()
        resolution = out.decode().splitlines()[0].strip()

        # filesize
        filesize = os.path.getsize(ts)

        return duration, resolution, filesize

    def update_db(self, id, title, duration, timestamp, resolution, fps, filesize):
        Vod.objects.update_or_create(
            filename=id,
            defaults={
                "title": title,
                "duration": duration,
                "date": make_aware(datetime.fromtimestamp(timestamp)),
                "resolution": resolution,
                "fps": fps,
                "size": filesize
            })
