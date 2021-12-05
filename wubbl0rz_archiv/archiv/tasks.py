import json
import os
import subprocess
from datetime import datetime

import yt_dlp
from celery import shared_task
from django.utils.timezone import make_aware
from pymediainfo import MediaInfo

from .models import Vod


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
        cmd = ["ffmpeg", "-ss", str(duration/1000/2), "-i",
               ts, "-vframes", "1", os.path.join(vod_dir, id + ".jpg")]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        proc.communicate()

    def get_metadata(self, vod_dir, entry):
        ts = os.path.join(vod_dir, entry["id"] + ".ts")
        media_info = MediaInfo.parse(ts)
        for track in media_info.tracks:
            if track.track_type == "Video":
                duration = track.duration
                resolution = f"{track.width}x{track.height}"
        bitrate = media_info.general_tracks[0].to_data()[
            "overall_bit_rate"]
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


@shared_task
def main_dl():
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
