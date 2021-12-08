import json
import os
import subprocess
from datetime import datetime

from archiv.models import Vod
from django.conf import settings
from django.core.management.base import BaseCommand
from django.utils.timezone import make_aware
from pymediainfo import MediaInfo


class Command(BaseCommand):
    def handle(self, **options):
        for f in os.listdir(settings.MEDIA_ROOT):
            if not f.endswith(".ts"):
                continue

            id, _ = os.path.splitext(f)

            print("importing:", id)

            ts = os.path.join(settings.MEDIA_ROOT, id + ".ts")

            cmd = ["ffprobe", "-v", "error", "-show_entries", "format=duration", "-of",
                "default=noprint_wrappers=1:nokey=1", os.path.join(settings.MEDIA_ROOT, ts)]
            proc = subprocess.Popen(
                cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
            out, _ = proc.communicate()
            duration = float(out.decode().strip())

            media_info = MediaInfo.parse(ts)
            for track in media_info.tracks:
                if track.track_type == "Video":
                    resolution = f"{track.width}x{track.height}"

            bitrate = media_info.general_tracks[0].to_data()[
                "overall_bit_rate"]

            filesize = os.path.getsize(ts)

            with open(os.path.join(settings.MEDIA_ROOT, id + ".json"), "r", encoding="utf-8") as fh:
                info_dict = json.load(fh)
            title = info_dict["title"]
            timestamp = info_dict["timestamp"]
            fps = info_dict["fps"]

            Vod.objects.update_or_create(
                filename=id,
                defaults={
                    "title": title,
                    "duration": duration,
                    "date": make_aware(datetime.fromtimestamp(timestamp)),
                    "resolution": resolution,
                    "bitrate": bitrate,
                    "fps": fps,
                    "size": filesize
                })
