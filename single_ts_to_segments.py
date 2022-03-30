import os
import subprocess
import sys


def main(media_dir):
    for media_type in ["vods", "clips"]:
        media_path = os.path.join(media_dir, media_type)
        for f in os.listdir(media_path):
            if f.endswith(".ts"):
                print("------------------------------------")
                print(f)
                full_path = os.path.join(media_path, f)
                filename, _ = os.path.splitext(f)
                dest_path = os.path.join(media_path, filename + "-segments")
                if not os.path.isdir(dest_path):
                    os.mkdir(dest_path)
                cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-stats", "-i", full_path, "-c", "copy",
                       "-hls_playlist_type", "vod", "-hls_time", "10", "-hls_segment_filename", os.path.join(dest_path, filename + "_%04d.ts"), os.path.join(dest_path, filename + ".m3u8")]
                proc = subprocess.Popen(
                    cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
                proc.communicate()
                os.remove(full_path)
                os.remove(os.path.join(media_path, filename + ".m3u8"))
        print(media_type, "done")


if __name__ == "__main__":
    if not len(sys.argv) > 1:
        exit()
    media_dir = sys.argv[1]
    main(media_dir)
