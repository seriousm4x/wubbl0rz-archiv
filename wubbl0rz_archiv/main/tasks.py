from fileinput import filename
import os
import subprocess
from datetime import datetime

import requests
import yt_dlp
from celery import shared_task
from clips.models import Clip, Game, Creator
from django.utils import timezone
from django.utils.timezone import make_aware
from vods.models import Vod

from main.models import ApiStorage, Emote
from main.twitch_api import TwitchApi


class MyLogger:
    def debug(self, msg):
        pass

    def info(self, msg):
        pass

    def warning(self, msg):
        pass

    def error(self, msg):
        pass


class Downloader:
    def __init__(self) -> None:
        obj = ApiStorage.objects.first()
        obj.date_vods_updated = timezone.now()
        obj.save()
        TwitchApi().update_bearer()
        self.helix_header = {
            "Client-ID": ApiStorage.objects.get().ttv_client_id,
            "Authorization": f"Bearer {ApiStorage.objects.get().ttv_bearer_token}"
        }

    def get_vod_infos(self):
        ydl_opts = {
            "retries": 10,
            "logger": MyLogger(),
        }
        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            info_dict = ydl.extract_info(
                "https://www.twitch.tv/wubbl0rz/videos?filter=all&sort=time", download=False)
        return info_dict

    def download(self, dir, url, id):
        ydl_opts = {
            "format": "best",
            "concurrent-fragments": 8,
            "retries": 10,
            "outtmpl": os.path.join(dir, f"{id}.%(ext)s"),
            "logger": MyLogger()
        }
        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            ydl.download(url)

    def dl_post_processing(self, dir, id):
        mp4 = os.path.join(dir, id + ".mp4")
        if not os.path.isdir(os.path.join(dir, id + "-segments")):
            os.mkdir(os.path.join(dir, id + "-segments"))
        cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-stats", "-i", mp4, "-c", "copy",
               "-hls_playlist_type", "vod", "-hls_time", "10", "-hls_segment_filename", os.path.join(dir, id + "-segments", id + "_%04d.ts"), os.path.join(dir, id + "-segments", id + ".m3u8")]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        proc.communicate()
        os.remove(mp4)

    def create_thumbnail(self, dir, id, duration):
        m3u8 = os.path.join(dir, id + "-segments", id + ".m3u8")
        if duration <= 10:
            timecode_framegrab = "0"
        else:
            timecode_framegrab = str(round(duration/2))

        # thumbnail sizes
        sm_width = "260"
        md_width = "520"
        lg_width = "1592"

        # jpg sm
        cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab, "-i", m3u8,
               "-vframes", "1", "-vf", f"scale={sm_width}:-1", "-y", os.path.join(dir, id + "-sm.jpg")]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        proc.communicate()

        # jpg md
        cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab, "-i", m3u8,
               "-vframes", "1", "-vf", f"scale={md_width}:-1", "-y", os.path.join(dir, id + "-md.jpg")]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        proc.communicate()

        # jpg lg
        cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab, "-i", m3u8,
               "-vframes", "1", "-vf", f"scale={lg_width}:-1", "-y", os.path.join(dir, id + "-lg.jpg")]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        proc.communicate()

        # lossless source png for avif sm
        cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab, "-i", m3u8,
               "-vframes", "1", "-vf", f"scale={sm_width}:-1", "-f", "image2", "-y", os.path.join(dir, id + ".png")]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        proc.communicate()

        # avif sm final
        cmd = ["avifenc", os.path.join(
            dir, id + ".png"), os.path.join(dir, id + "-sm.avif")]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        proc.communicate()

        # lossless source png for avif md
        cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab, "-i", m3u8,
               "-vframes", "1", "-vf", f"scale={md_width}:-1", "-f", "image2", "-y", os.path.join(dir, id + ".png")]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        proc.communicate()

        # avif md final
        cmd = ["avifenc", os.path.join(
            dir, id + ".png"), os.path.join(dir, id + "-md.avif")]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        proc.communicate()

        # remove lossless image
        os.remove(os.path.join(dir, id + ".png"))

        # create .webp preview animation
        cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab,
               "-i", m3u8, "-c:v", "libwebp", "-vf", "scale=260:-1,fps=fps=15", "-lossless",
               "0", "-compression_level", "3", "-q:v", "70", "-loop", "0", "-preset", "picture",
               "-an", "-vsync", "0", "-t", "4", "-y", os.path.join(dir, id + "-preview.webp")]
        proc = subprocess.Popen(cmd, stderr=subprocess.PIPE,
                                stdout=subprocess.PIPE)
        proc.communicate()

    def get_metadata(self, dir, id):
        m3u8 = os.path.join(dir, id + "-segments", id + ".m3u8")

        # duration
        cmd = ["ffprobe", "-v", "error", "-show_entries", "format=duration", "-of",
               "default=noprint_wrappers=1:nokey=1", m3u8]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        out, _ = proc.communicate()
        duration = float(out.decode().strip())

        # resolution
        cmd = ["ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries",
               "stream=width,height", "-of", "csv=s=x:p=0", m3u8]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        out, _ = proc.communicate()
        resolution = out.decode().splitlines()[0].strip()

        # filesize
        filesize = 0
        for e in os.scandir(os.path.join(dir, id + "-segments")):
            filesize += os.path.getsize(e)

        return duration, resolution, filesize

    def update_vod(self, id, title, duration, timestamp, resolution, fps, filesize):
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

    def update_clip(self, data):
        if data["game_id"] == "":
            game_obj = None
        else:
            game_obj, created = Game.objects.get_or_create(game_id=data["game_id"])
            if created:
                game_req = requests.get(
                    f"https://api.twitch.tv/helix/games?id={data['game_id']}", headers=self.helix_header)
                try:
                    game_title = game_req.json()["data"][0]["name"]
                    Game.objects.update_or_create(
                        game_id=data["game_id"],
                        defaults={
                            "name": game_title
                        }
                    )
                except IndexError:
                    print(f"Game ID {data['game_id']} got deleted on Twitch")

        creator, _ = Creator.objects.get_or_create(
            creator_id=data["creator_id"], name=data["creator_name"])
        vod = Vod.objects.filter(filename=f"v{data['video_id']}").first()

        Clip.objects.update_or_create(
            clip_id=data["id"],
            defaults={
                "creator": creator,
                "game": game_obj,
                "title": data["title"],
                "view_count": data["view_count"],
                "date": data["created_at"],
                "duration": data["duration"],
                "resolution": data["resolution"],
                "size": data["size"],
                "vod": vod
            }
        )


class EmoteUpdater:
    def __init__(self) -> None:
        obj = ApiStorage.objects.first()
        obj.date_emotes_updated = timezone.now()
        obj.save()
        self.broadcaster_id = ApiStorage.objects.get().broadcaster_id

    def mark_outdated(self):
        for emote in Emote.objects.all():
            emote.outdated = True
            emote.save()

    def twitch(self):
        TwitchApi().update_bearer()
        client_id = ApiStorage.objects.get().ttv_client_id
        bearer = ApiStorage.objects.get().ttv_bearer_token

        # get emotes
        emote_url = f"https://api.twitch.tv/helix/chat/emotes?broadcaster_id={self.broadcaster_id}"
        helix_header = {
            "Client-ID": client_id,
            "Authorization": "Bearer {}".format(bearer),
        }
        emote_resp = requests.get(emote_url, headers=helix_header)
        emote_resp.raise_for_status()
        emote_json_resp = emote_resp.json()
        for emote in emote_json_resp["data"]:
            if "animated" in emote["format"]:
                image = emote["images"]["url_4x"].replace(
                    "/static/", "/animated/")
            else:
                image = emote["images"]["url_4x"]
            Emote.objects.update_or_create(
                id=emote["id"],
                provider="twitch",
                defaults={
                    "name": emote["name"],
                    "url": image,
                    "outdated": False
                }
            )

    def bttv(self):
        emote_url = f"https://api.betterttv.net/3/cached/users/twitch/{self.broadcaster_id}"
        emote_resp = requests.get(emote_url)
        emote_resp.raise_for_status()
        emote_json_resp = emote_resp.json()
        for emote in emote_json_resp["sharedEmotes"]:
            Emote.objects.update_or_create(
                id=emote["id"],
                provider="bttv",
                defaults={
                    "name": emote["code"],
                    "url": f"https://cdn.betterttv.net/emote/{emote['id']}/3x",
                    "outdated": False
                }
            )

    def ffz(self):
        emote_url = f"https://api.frankerfacez.com/v1/room/id/{self.broadcaster_id}"
        emote_resp = requests.get(emote_url)
        emote_resp.raise_for_status()
        emote_json_resp = emote_resp.json()
        for _, value in emote_json_resp["sets"].items():
            for emote in value["emoticons"]:
                Emote.objects.update_or_create(
                    id=emote["id"],
                    provider="ffz",
                    defaults={
                        "name": emote["name"],
                        "url": f"https://cdn.frankerfacez.com/emote/{emote['id']}/4",
                        "outdated": False
                    }
                )

    def delete_outdated(self):
        for emote in Emote.objects.all():
            if emote.outdated == True:
                emote.delete()

    def update_all(self):
        self.mark_outdated()
        self.twitch()
        self.bttv()
        self.ffz()
        self.delete_outdated()


@shared_task
def update_emotes():
    eu = EmoteUpdater()
    eu.update_all()


@shared_task
def check_live():
    live = False

    try:
        ydl_opts = {
            "retries": 10,
            "logger": MyLogger()
        }
        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            ydl.extract_info("https://www.twitch.tv/wubbl0rz/", download=False)
        live = True
    except yt_dlp.DownloadError:
        pass
    finally:
        obj = ApiStorage.objects.first()
        obj.is_live = live
        obj.save()
