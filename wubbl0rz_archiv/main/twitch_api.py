import datetime
import os

import requests
from discord_webhook import DiscordEmbed, DiscordWebhook
from django.utils import timezone

from main.models import ApiStorage


class TwitchApi:
    def __init__(self) -> None:
        self.apistorage = ApiStorage.objects.get()
        self.client_id = self.apistorage.ttv_client_id
        self.client_secret = self.apistorage.ttv_client_secret
        self.broadcaster_id = self.apistorage.broadcaster_id

    def is_expired(self) -> bool:
        # no expire date set
        if not self.apistorage.ttv_bearer_expire_date:
            return True

        expire_date = self.apistorage.ttv_bearer_expire_date

        # if bearer expires in less then 1 day
        if timezone.now() + datetime.timedelta(days=1) >= expire_date:
            return True
        return False

    def get_bearer(self):
        if self.is_expired():
            # refresh twitch credentials
            tokenurl = "https://id.twitch.tv/oauth2/token?client_id={}&client_secret={}&grant_type=client_credentials".format(
                self.client_id, self.client_secret)
            try:
                # get bearer token
                token_response = requests.post(tokenurl)
                token_response.raise_for_status()
                token_json_response = token_response.json()

                # write to database
                ApiStorage.objects.update_or_create(
                    broadcaster_id=self.broadcaster_id,
                    defaults={
                        "ttv_bearer_token": token_json_response["access_token"],
                        "ttv_bearer_expire_date": timezone.now() + datetime.timedelta(seconds=token_json_response["expires_in"])
                    }
                )
                return token_json_response["access_token"]
            except requests.exceptions.HTTPError as http_err:
                print("HTTP error occurred: {}".format(http_err))
            except Exception as err:
                print("Other error occurred: {}".format(err))

        return self.apistorage.ttv_bearer_token

    def check_live(self):
        bearer = self.get_bearer()
        helix_header = {
            "Client-ID": self.client_id,
            "Authorization": f"Bearer {bearer}",
        }
        stream_url = f"https://api.twitch.tv/helix/streams?user_id={self.broadcaster_id}"
        stream_resp = requests.get(
            stream_url, headers=helix_header).json()["data"]
        if len(stream_resp) > 0:
            is_live = True
            stream_resp = stream_resp[0]
        else:
            is_live = False

        if is_live != self.apistorage.is_live:
            ApiStorage.objects.update_or_create(
                broadcaster_id=self.broadcaster_id,
                defaults={
                    "is_live": is_live
                }
            )

            if is_live and os.getenv("DISCORD_WEBHOOK"):
                user_url = f"https://api.twitch.tv/helix/users?id={self.broadcaster_id}"
                user_resp = requests.get(
                    user_url, headers=helix_header).json()["data"]

                timestamp = timezone.make_aware(datetime.datetime.strptime(
                    stream_resp["started_at"], "%Y-%m-%dT%H:%M:%SZ"))
                webhook = DiscordWebhook(url=os.getenv(
                    "DISCORD_WEBHOOK"), content=f"Auf gehts @everyone, {stream_resp['user_name']} macht nun Streamstelz")
                embed = DiscordEmbed(
                    title=f"{stream_resp['title']}, https://twitch.tv/{stream_resp['user_login']}", description=stream_resp["game_name"], color="f28e2b")
                embed.set_author(
                    name=stream_resp["user_name"], url=f"https://twitch.tv/{stream_resp['user_login']}", icon_url=user_resp[0]["profile_image_url"])
                embed.set_image(url=stream_resp["thumbnail_url"].replace(
                    r"{width}", "1280").replace(r"{height}", "720"))
                embed.set_timestamp(datetime.datetime.timestamp(timestamp))
                webhook.add_embed(embed)
                webhook.execute()
