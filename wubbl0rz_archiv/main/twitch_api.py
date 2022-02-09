from main.models import ApiStorage
import requests


class TwitchApi:
    def update_bearer(self):
        ttv_client_id = ApiStorage.objects.get().ttv_client_id
        ttv_client_secret = ApiStorage.objects.get().ttv_client_secret
        broadcaster_id = ApiStorage.objects.get().broadcaster_id

        # refresh twitch credentials
        tokenurl = "https://id.twitch.tv/oauth2/token?client_id={}&client_secret={}&grant_type=client_credentials".format(
            ttv_client_id, ttv_client_secret)
        try:
            # get bearer token
            token_response = requests.post(tokenurl)
            token_response.raise_for_status()
            token_json_response = token_response.json()
            bearer = token_json_response["access_token"]

            # write to database
            ApiStorage.objects.update_or_create(
                broadcaster_id=broadcaster_id,
                defaults={
                    "ttv_bearer_token": bearer
                }
            )
        except requests.exceptions.HTTPError as http_err:
            print("HTTP error occurred: {}".format(http_err))
        except Exception as err:
            print("Other error occurred: {}".format(err))
