import json
import os
from datetime import datetime

from archiv.models import Vod
from django.core.management.base import BaseCommand
from django.utils.timezone import make_aware


class Command(BaseCommand):
    def __init__(self):
        self.vods_dir = "/mnt/nas/archiv-dev/media/"

    def handle(self, **options):
        for f in os.listdir(self.vods_dir):
            name, ext = os.path.splitext(f)
            if ext == ".json":
                with open(os.path.join(self.vods_dir, f), "r", encoding="utf-8") as info_file:
                    info = json.load(info_file)
                print(info["id"])
                Vod.objects.update_or_create(
                    filename=name,
                    defaults={
                        "title": info["title"],
                        "duration": info["duration"],
                        "date": make_aware(datetime.fromtimestamp(info["timestamp"])),
                        "resolution": f"{info['width']}x{info['height']}",
                        "bitrate": info["tbr"],
                        "fps": info["fps"],
                        "size": os.path.getsize(os.path.join(self.vods_dir, name + ".ts"))
                    })
