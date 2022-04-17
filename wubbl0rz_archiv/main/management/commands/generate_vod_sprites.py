import os
import subprocess

from django.conf import settings
from django.core.management.base import BaseCommand


class Command(BaseCommand):
    def handle(self, **options):
        self.vod_dir = os.path.join(settings.MEDIA_ROOT, "vods")
        self.gen_sprites()

    def gen_sprites(self):
        for f in os.listdir(self.vod_dir):
            if not f.endswith("-segments"):
                continue
            id = f.replace("-segments", "")
            print("create sprites for vod:", id)
            sprite_dir = os.path.join(self.vod_dir, id + "-sprites")
            if os.path.isdir(sprite_dir):
                return
            else:
                os.mkdir(sprite_dir)
            cmd = ["ffmpeg", "-i", os.path.join(self.vod_dir, f, id+".m3u8"), "-vf", "fps=1/20,scale=-1:90,tile",
                   "-c:v", "libwebp", os.path.join(sprite_dir, id+r"_%03d.webp")]
            proc = subprocess.Popen(
                cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
            proc.communicate()
