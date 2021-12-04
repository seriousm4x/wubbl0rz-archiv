import json
import os
from pymediainfo import MediaInfo
from datetime import datetime

from archiv.models import Vod
from django.core.management.base import BaseCommand
from django.utils.timezone import make_aware


def get_duration(f):
    media_info = MediaInfo.parse(f)
    for track in media_info.tracks:
        if track.track_type == "Video":
            return track.duration

def get_resolution(f):
    media_info = MediaInfo.parse(f)
    for track in media_info.tracks:
        if track.track_type == "Video":
            return f"{track.width}x{track.height}"

def get_bitrate(f):
    media_info = MediaInfo.parse(f)
    return media_info.general_tracks[0].to_data()["overall_bit_rate"]

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
                        "duration": get_duration(os.path.join(self.vods_dir, name + ".ts")),
                        "date": make_aware(datetime.fromtimestamp(info["timestamp"])),
                        "resolution": get_resolution(os.path.join(self.vods_dir, name + ".ts")),
                        "bitrate": get_bitrate(os.path.join(self.vods_dir, name + ".ts")),
                        "fps": info["fps"],
                        "size": os.path.getsize(os.path.join(self.vods_dir, name + ".ts"))
                    })
