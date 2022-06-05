import datetime

import requests
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
        url = f"https://api.twitch.tv/helix/streams?user_id={self.broadcaster_id}"
        helix_header = {
            "Client-ID": self.client_id,
            "Authorization": f"Bearer {bearer}",
        }
        resp = requests.get(url, headers=helix_header).json()
        if resp["data"]:
            is_live = True
        else:
            is_live = False

        ApiStorage.objects.update_or_create(
            broadcaster_id=self.broadcaster_id,
            defaults={
                "is_live": is_live
            }
        )
