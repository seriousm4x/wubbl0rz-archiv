import argparse
import json
import os

from django.conf import settings
from main.tasks import Downloader


def parse_args():
    parser = argparse.ArgumentParser(
        description='Generate all media files for vods from an mp4.')
    parser.add_argument('-id', type=str, required=True,
                        help='Vod id, such as "v2849020238"')
    args = parser.parse_args()
    return args


def main(args):
    dl = Downloader()
    vod_dir = os.path.join(settings.MEDIA_ROOT, "vods")
    json_file = os.path.join(vod_dir, args.id + ".json")

    with open(json_file, "w", encoding="utf-8") as f:
        entry = json.load(f)

    duration, resolution, filesize = dl.get_metadata(vod_dir, args.id)
    print(args.id, entry["title"], duration,
          entry["timestamp"], resolution, entry["fps"], filesize)
    # dl.update_vod(args.id, entry["title"], duration,
    #               entry["timestamp"], resolution, entry["fps"], filesize)


if __name__ == "__main__":
    args = parse_args()
    main(args)
