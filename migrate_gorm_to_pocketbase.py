import copy
import datetime
import math
from urllib.parse import quote

import requests

GORM_API = ""
PB_API = ""
PB_USER = ""
PB_PASS = ""


class Gorm:
    def get_all_vods(self):
        req = requests.get(GORM_API + "/vods?limit=-1")
        return req.json()["result"]

    def get_all_clips(self):
        req = requests.get(GORM_API + "/clips?limit=-1")
        return req.json()["result"]

    def get_all_creators(self):
        req = requests.get(GORM_API + "/creators?limit=-1")
        return req.json()["result"]

    def get_all_games(self):
        req = requests.get(GORM_API + "/games?limit=-1")
        return req.json()["result"]

    def get_all_emotes(self):
        req = requests.get(GORM_API + "/emotes?limit=-1")
        return req.json()["result"]

    def get_all_messages(self):
        req = requests.get(GORM_API + "/chat?from=946681200&to=" +
                           str(math.floor(datetime.datetime.now().timestamp())))
        return req.json()["result"]

    def get_game(self, uuid):
        req = requests.get(GORM_API + "/games/" + uuid)
        return req.json()["result"]

    def get_vod(self, uuid):
        req = requests.get(GORM_API + "/vods/" + uuid)
        return req.json()["result"]


class PocketBase:
    def __init__(self):
        data = {
            "identity": PB_USER,
            "password": PB_PASS
        }
        req = requests.post(PB_API + "/admins/auth-with-password", json=data)
        self.headers = {
            "Content-Type": "application/json",
            "Authorization": req.json()["token"]
        }

    def create_vod(self, data):
        data = copy.copy(data)
        data.pop("uuid")
        try:
            req = requests.post(PB_API + "/collections/vod/records",
                                json=data, headers=self.headers)
        except Exception as e:
            print(e)
            print("-"*20)
            print(req.text)
            exit(1)

    def create_clip(self, data):
        data = copy.copy(data)
        data.pop("uuid")
        data.pop("creator_uuid")
        data.pop("game_uuid")
        data.pop("vod_uuid")
        try:
            req = requests.post(PB_API + "/collections/clip/records",
                                json=data, headers=self.headers)
        except Exception as e:
            print(e)
            print("-"*20)
            print(req.text)
            exit(1)

    def create_creator(self, data):
        data = copy.copy(data)
        data["ttv_id"] = data["uuid"]
        data.pop("uuid")
        data.pop("clips")
        try:
            req = requests.post(PB_API + "/collections/creator/records",
                                json=data, headers=self.headers)
        except Exception as e:
            print(e)
            print("-"*20)
            print(req.text)
            exit(1)

    def create_game(self, data):
        data = copy.copy(data)
        data["ttv_id"] = data["uuid"]
        data.pop("uuid")
        data.pop("clips")
        try:
            req = requests.post(PB_API + "/collections/game/records",
                                json=data, headers=self.headers)
        except Exception as e:
            print(e)
            print("-"*20)
            print(req.text)
            exit(1)

    def create_emote(self, data):
        data = copy.copy(data)
        data.pop("id")
        try:
            req = requests.post(PB_API + "/collections/emote/records",
                                json=data, headers=self.headers)
        except Exception as e:
            print(e)
            print("-"*20)
            print(req.text)
            exit(1)

    def create_message(self, data):
        data = copy.copy(data)
        data["date"] = data["created_at"]
        data.pop("ID")
        data.pop("created_at")
        try:
            req = requests.post(PB_API + "/collections/chatmessage/records",
                                json=data, headers=self.headers)
        except Exception as e:
            print(e)
            print("-"*20)
            print(req.text)
            exit(1)

    def get_creator(self, name):
        try:
            req = requests.get(
                PB_API + f"/collections/creator/records?page=1&perPage=1&filter=name=\"{quote(name)}\"&fields=id")
            return req.json()["items"]
        except Exception as e:
            print(name)
            print("-"*20)
            print(req.text)
            exit(1)

    def get_game(self, name):
        try:
            req = requests.get(
                PB_API + f"/collections/game/records?page=1&perPage=1&filter=name=\"{quote(name)}\"&fields=id")
            return req.json()["items"]
        except Exception as e:
            print(name)
            print("-"*20)
            print(req.text)
            exit(1)

    def get_vod(self, filename):
        try:
            req = requests.get(
                PB_API + f"/collections/vod/records?page=1&perPage=1&filter=filename=\"{quote(filename)}\"&fields=id")
            return req.json()["items"]
        except Exception as e:
            print(filename)
            print("-"*20)
            print(req.text)
            exit(1)


def main():
    pb = PocketBase()
    gorm = Gorm()

    print("creating vods")
    vods = gorm.get_all_vods()
    for vod in vods:
        pb.create_vod(vod)

    print("creating creators")
    creators = gorm.get_all_creators()
    for creator in creators:
        pb.create_creator(creator)

    print("creating games")
    games = gorm.get_all_games()
    games = [i for i in games if i["name"] != "" and i.get("uuid")]
    for game in games:
        pb.create_game(game)

    print("creating emotes")
    emotes = gorm.get_all_emotes()
    for emote in emotes:
        pb.create_emote(emote)

    print("creating chatmessages")
    messages = gorm.get_all_messages()
    for message in messages:
        pb.create_message(message)

    print("creating clips")
    clips = gorm.get_all_clips()
    for clip in clips:
        # creator relation
        creator_req = pb.get_creator(clip["creator"]["name"])
        if len(creator_req) > 0:
            clip["creator"] = creator_req[0]["id"]
        else:
            clip["creator"] = ""

        # game relation
        for game in games:
            if game["uuid"] == clip["game_uuid"]:
                gorm_game = gorm.get_game(clip["game_uuid"])
                pb_vod = pb.get_game(gorm_game["name"])
                if len(pb_vod) > 0:
                    clip["game"] = pb_vod[0]["id"]
                else:
                    clip["game"] = ""
                break

        # vod relation
        for vod in vods:
            if vod["uuid"] == clip["vod_uuid"]:
                gorm_vod = gorm.get_vod(clip["vod_uuid"])
                pb_vod = pb.get_vod(gorm_vod["filename"])
                if len(pb_vod) > 0:
                    clip["vod"] = pb_vod[0]["id"]
                else:
                    clip["vod"] = ""
                break

        pb.create_clip(clip)


if __name__ == "__main__":
    main()
