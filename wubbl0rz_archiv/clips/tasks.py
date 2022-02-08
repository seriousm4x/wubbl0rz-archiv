import datetime
import json
import os

import requests
from celery import shared_task
from django.conf import settings
from main.models import ApiStorage
from main.tasks import Downloader
from main.twitch_api import TwitchApi
import shutil

from clips.models import Clip, Game


class ClipDownloader:
    def __init__(self):
        self.clipdir = os.path.join(settings.MEDIA_ROOT, "clips")
        self.gamedir = os.path.join(settings.MEDIA_ROOT, "games")
        self.broadcaster_id = ApiStorage.objects.get().broadcaster_id
        self.downloader = Downloader()
        TwitchApi().update_bearer()
        self.helix_header = {
            "Client-ID": ApiStorage.objects.get().ttv_client_id,
            "Authorization": f"Bearer {ApiStorage.objects.get().ttv_bearer_token}",
        }

    def clips(self):
        channel_creation_date = requests.get(
            "https://api.twitch.tv/helix/users?id={}".format(
                self.broadcaster_id), headers=self.helix_header).json()["data"][0]["created_at"]

        week = 1
        started_at = (datetime.datetime.now() -
                      datetime.timedelta(weeks=week)).isoformat("T")+"Z"
        api_url = "https://api.twitch.tv/helix/clips?broadcaster_id={}&first=100&started_at={}".format(
            self.broadcaster_id, started_at)

        while True:
            req = requests.get(api_url, headers=self.helix_header)
            clips = req.json()

            for data in clips["data"]:
                # save clip
                if data["view_count"] < 3:
                    continue
                if Clip.objects.filter(clip_id=data["id"]).exists():
                    c = Clip.objects.filter(clip_id=data["id"]).first()
                    data["duration"] = c.duration
                    data["resolution"] = c.resolution
                    data["size"] = c.size
                    self.downloader.update_clip(data)
                    continue

                if not os.path.isfile(os.path.join(self.clipdir, data["id"] + ".json")):
                    with open(os.path.join(self.clipdir, data["id"] + ".json"), "w", encoding="utf-8") as f:
                        json.dump(data, f)

                self.downloader.download(self.clipdir, data["url"], data["id"])
                self.downloader.dl_post_processing(self.clipdir, data["id"])
                data["duration"], data["resolution"], data["size"] = self.downloader.get_metadata(
                    self.clipdir, data["id"])
                self.downloader.create_thumbnail(
                    self.clipdir, data["id"], data["duration"])
                self.downloader.update_clip(data)

            try:
                cursor = clips["pagination"]["cursor"]
                api_url = "https://api.twitch.tv/helix/clips?broadcaster_id={}&first=100&started_at={}&after={}".format(
                    self.broadcaster_id, (datetime.datetime.now() - datetime.timedelta(weeks=week)).isoformat("T")+"Z", cursor)
            except KeyError:
                week += 1
                cursor_week = (datetime.datetime.now() -
                               datetime.timedelta(weeks=week)).isoformat("T")+"Z"
                if cursor_week < channel_creation_date:
                    print("Channel creation date reached")
                    break

                print(cursor_week)
                api_url = "https://api.twitch.tv/helix/clips?broadcaster_id={}&first=100&started_at={}".format(
                    self.broadcaster_id, cursor_week)

    def game(self, obj):
        if not os.path.isdir(self.gamedir):
            os.mkdir(self.gamedir)
        game_req = requests.get(
            f"https://api.twitch.tv/helix/games?id={obj.game_id}", headers=self.helix_header)
        try:
            game_title = game_req.json()["data"][0]["name"]
            box_art = game_req.json()["data"][0]["box_art_url"].replace(
                r"{width}x{height}", "100x133")
            resp = requests.get(box_art, stream=True)
            with open(os.path.join(self.gamedir, str(obj.game_id) + ".jpg"), "wb") as f:
                shutil.copyfileobj(resp.raw, f)
            Game.objects.update_or_create(
                game_id=obj.game_id,
                defaults={
                    "name": game_title,
                    "box_art": box_art
                }
            )
        except IndexError:
            print(f"Game ID {obj.game_id} got deleted on Twitch")


@shared_task
def download_clips():
    cd = ClipDownloader()
    cd.clips()
    for g in Game.objects.all():
        cd.game(g)
