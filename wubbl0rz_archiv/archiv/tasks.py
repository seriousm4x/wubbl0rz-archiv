import json
import os
import subprocess
from datetime import datetime

import pillow_avif
import requests
import yt_dlp
from celery import shared_task
from django.utils import timezone
from django.utils.timezone import make_aware
from PIL import Image
from pymediainfo import MediaInfo

from .models import ApiStorage, Emote, Vod


class MyLogger:
    def debug(self, msg):
        pass

    def info(self, msg):
        pass

    def warning(self, msg):
        pass

    def error(self, msg):
        print(msg)


class VODDownloader:
    def __init__(self) -> None:
        obj = ApiStorage.objects.first()
        obj.date_vods_updated = timezone.now()
        obj.save()

    def get_info_dict(self):
        ydl_opts = {
            'logger': MyLogger(),
        }
        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            info_dict = ydl.extract_info(
                "https://www.twitch.tv/wubbl0rz/videos?filter=all&sort=time", download=False)
        return info_dict

    def download_vod(self, vod_dir, entry):
        ydl_opts = {
            "format": "best",
            "concurrent-fragments": 8,
            "outtmpl": os.path.join(vod_dir, "%(id)s.%(ext)s"),
            'logger': MyLogger()
        }
        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            ydl.download(entry["webpage_url"])

    def dl_post_processing(self, vod_dir, entry):
        mp4 = os.path.join(vod_dir, entry["id"] + ".mp4")
        m3u8 = os.path.join(vod_dir, entry["id"] + ".m3u8")
        cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-stats", "-i", mp4, "-c", "copy",
               "-hls_playlist_type", "vod", "-hls_time", "10", "-hls_flags", "single_file", m3u8]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        proc.communicate()
        os.remove(mp4)

    def create_thumbnail(self, vod_dir, id, duration):
        ts = os.path.join(vod_dir, id + ".ts")

        # create lossless image to use as source for smaller thumbs
        cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", str(round(duration/2)), "-i",
               ts, "-vframes", "1", "-f", "image2", "-y", os.path.join(vod_dir, id + ".png")]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        proc.communicate()

        # jpg lg
        lossless = Image.open(os.path.join(vod_dir, id + ".png"))
        lossless.save(os.path.join(vod_dir, id + "-lg.jpg"), quality=90)

        # jpg sm
        lossless = Image.open(os.path.join(vod_dir, id + ".png"))
        img = lossless.resize((480, 270), Image.ANTIALIAS)
        img.save(os.path.join(vod_dir, id + "-sm.jpg"), quality=90)

        # avif sm
        lossless = Image.open(os.path.join(vod_dir, id + ".png"))
        img = lossless.resize((480, 270), Image.ANTIALIAS)
        img.save(os.path.join(vod_dir, id + "-sm.avif"), quality=90)

        os.remove(os.path.join(vod_dir, id + ".png"))

    def get_metadata(self, vod_dir, entry):
        ts = os.path.join(vod_dir, entry["id"] + ".ts")

        # duration
        cmd = ["ffprobe", "-v", "error", "-show_entries", "format=duration", "-of",
               "default=noprint_wrappers=1:nokey=1", os.path.join(vod_dir, id + ".ts")]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        out, _ = proc.communicate()
        duration = float(out.decode().strip())

        # resolution
        media_info = MediaInfo.parse(ts)
        for track in media_info.tracks:
            if track.track_type == "Video":
                resolution = f"{track.width}x{track.height}"

        # bitrate
        bitrate = media_info.general_tracks[0].to_data()[
            "overall_bit_rate"]

        # filesize
        filesize = os.path.getsize(ts)

        return duration, resolution, bitrate, filesize

    def update_db(self, id, title, duration, timestamp, resolution, bitrate, fps, filesize):
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
        ttv_client_id = ApiStorage.objects.get().ttv_client_id
        ttv_client_secret = ApiStorage.objects.get().ttv_client_secret

        # refresh twitch credentials
        tokenurl = "https://id.twitch.tv/oauth2/token?client_id={}&client_secret={}&grant_type=client_credentials".format(
            ttv_client_id, ttv_client_secret)
        try:
            # get bearer token
            token_response = requests.post(tokenurl)
            token_response.raise_for_status()
            token_jsonResponse = token_response.json()
            bearer = token_jsonResponse["access_token"]

            helix_header = {
                "Client-ID": ttv_client_id,
                "Authorization": "Bearer {}".format(bearer),
            }

            # write to database
            ApiStorage.objects.update_or_create(
                broadcaster_id=self.broadcaster_id,
                defaults={
                    "ttv_bearer_token": bearer
                }
            )
        except requests.exceptions.HTTPError as http_err:
            print("HTTP error occurred: {}".format(http_err))
        except Exception as err:
            print("Other error occurred: {}".format(err))

        # get emotes
        emote_url = f"https://api.twitch.tv/helix/chat/emotes?broadcaster_id={self.broadcaster_id}"
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
def download_vods():
    vod_dir = "/mnt/nas/Archiv/wubbl0rz-twitch-vods/media/"
    vodd = VODDownloader()
    print("getting info dict")
    info_dict = vodd.get_info_dict()
    for entry in info_dict["entries"]:
        if entry["live_status"] == "is_live" or Vod.objects.filter(filename=entry["id"]).exists() or os.path.isfile(os.path.join(vod_dir, entry["id"] + ".ts")):
            continue

        if not os.path.isfile(os.path.join(vod_dir, entry["id"] + ".json")):
            with open(os.path.join(vod_dir, entry["id"] + ".json"), "w", encoding="utf-8") as f:
                json.dump(entry, f)

        print(f"download vod: {entry['id']}")
        vodd.download_vod(vod_dir, entry)
        print("post processing")
        vodd.dl_post_processing(vod_dir, entry)
        print("get metadata")
        duration, resolution, bitrate, filesize = vodd.get_metadata(
            vod_dir, entry)
        print("create thumbnail")
        vodd.create_thumbnail(vod_dir, entry["id"], duration)
        filesize = os.path.getsize(os.path.join(vod_dir, entry["id"] + ".ts"))
        print("update db")
        vodd.update_db(entry["id"], entry["title"], duration,
                       entry["timestamp"], resolution, bitrate, entry["fps"], filesize)


@shared_task
def update_emotes():
    eu = EmoteUpdater()
    eu.update_all()


@shared_task
def check_live():
    try:
        ydl_opts = {
            'logger': MyLogger()
        }
        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            ydl.extract_info("https://www.twitch.tv/wubbl0rz/", download=False)
        live = True
    except yt_dlp.DownloadError:
        live = False
    finally:
        obj = ApiStorage.objects.first()
        obj.is_live = live
        obj.save()
