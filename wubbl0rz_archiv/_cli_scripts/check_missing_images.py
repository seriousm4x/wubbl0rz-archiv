import argparse
import logging
import os
import pathlib
import subprocess


def create_thumbnail(dir, id, duration):
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
    logger.info(f"{id}: Creating jpg sm thumbnail")
    cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab, "-i", m3u8,
           "-vframes", "1", "-vf", f"scale={sm_width}:-1", "-y", os.path.join(dir, id + "-sm.jpg")]
    proc = subprocess.Popen(
        cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    proc.communicate()

    # jpg md
    logger.info(f"{id}: Creating jpg md thumbnail")
    cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab, "-i", m3u8,
           "-vframes", "1", "-vf", f"scale={md_width}:-1", "-y", os.path.join(dir, id + "-md.jpg")]
    proc = subprocess.Popen(
        cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    proc.communicate()

    # jpg lg
    logger.info(f"{id}: Creating jpg lg thumbnail")
    cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab, "-i", m3u8,
           "-vframes", "1", "-vf", f"scale={lg_width}:-1", "-y", os.path.join(dir, id + "-lg.jpg")]
    proc = subprocess.Popen(
        cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    proc.communicate()

    # lossless source png for avif sm
    logger.info(f"{id}: Creating source png sm")
    cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab, "-i", m3u8,
           "-vframes", "1", "-vf", f"scale={sm_width}:-1", "-f", "image2", "-y", os.path.join(dir, id + ".png")]
    proc = subprocess.Popen(
        cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    proc.communicate()

    # avif sm final
    logger.info(f"{id}: Creating avif sm thumbnail")
    cmd = ["avifenc", os.path.join(
        dir, id + ".png"), os.path.join(dir, id + "-sm.avif")]
    proc = subprocess.Popen(
        cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    proc.communicate()

    # lossless source png for avif md
    logger.info(f"{id}: Creating source png md")
    cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab, "-i", m3u8,
           "-vframes", "1", "-vf", f"scale={md_width}:-1", "-f", "image2", "-y", os.path.join(dir, id + ".png")]
    proc = subprocess.Popen(
        cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    proc.communicate()

    # avif md final
    logger.info(f"{id}: Creating avif md thumbnail")
    cmd = ["avifenc", os.path.join(
        dir, id + ".png"), os.path.join(dir, id + "-md.avif")]
    proc = subprocess.Popen(
        cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    proc.communicate()

    # remove lossless image
    os.remove(os.path.join(dir, id + ".png"))

    # create .webp preview animation
    logger.info(f"{id}: Creating webm preview")
    cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", timecode_framegrab,
           "-i", m3u8, "-c:v", "libwebp", "-vf", "scale=260:-1,fps=fps=15", "-lossless",
           "0", "-compression_level", "3", "-q:v", "70", "-loop", "0", "-preset", "picture",
           "-an", "-vsync", "0", "-t", "4", "-y", os.path.join(dir, id + "-preview.webp")]
    proc = subprocess.Popen(cmd, stderr=subprocess.PIPE,
                            stdout=subprocess.PIPE)
    proc.communicate()

    # create sprites
    logger.info(f"{id}: Creating sprite thumbnails")
    sprite_dir = os.path.join(dir, id + "-sprites")
    if not os.path.isdir(sprite_dir):
        os.mkdir(sprite_dir)
    cmd = ["ffmpeg", "-i", m3u8, "-vf", "fps=1/20,scale=-1:90,tile",
           "-c:v", "libwebp", os.path.join(sprite_dir, id+r"_%03d.webp")]
    proc = subprocess.Popen(
        cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    proc.communicate()


def get_metadata(dir, id):
    m3u8 = os.path.join(dir, id + "-segments", id + ".m3u8")

    # duration
    cmd = ["ffprobe", "-v", "error", "-show_entries", "format=duration", "-of",
           "default=noprint_wrappers=1:nokey=1", m3u8]
    proc = subprocess.Popen(
        cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    out, _ = proc.communicate()
    duration = float(out.decode().strip())
    return duration


def set_logger():
    logger = logging.getLogger('generate_vod')
    logger.setLevel(logging.DEBUG)
    ch = logging.StreamHandler()
    ch.setLevel(logging.DEBUG)
    formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
    ch.setFormatter(formatter)
    logger.addHandler(ch)
    return logger


def parse_args():
    parser = argparse.ArgumentParser(
        description='Check media folder for missing images')
    parser.add_argument('-m', '--media', type=pathlib.Path,
                        required=True, help='Path to media directory')
    args = parser.parse_args()
    return args


def check_files(path):
    expected_files = [
        "-sm.jpg",
        "-md.jpg",
        "-lg.jpg",
        "-sm.avif",
        "-md.avif",
        "-preview.webp"
    ]

    for ext in expected_files:
        expected_file = path + ext
        if not os.path.isfile(expected_file):
            logger.info("Missing:", expected_file)
            return False
    return True


def gen_files(p, id):
    duration = get_metadata(p, id)
    create_thumbnail(p, id, duration)


def main(args):
    vods_path = os.path.join(args.media, "vods")
    clips_path = os.path.join(args.media, "clips")

    for p in [vods_path, clips_path]:
        for filename in os.listdir(p):
            id, ext = os.path.splitext(filename)
            if ext != ".json":
                continue
            if not check_files(os.path.join(p, id)):
                gen_files(p, id)


if __name__ == "__main__":
    logger = set_logger()
    args = parse_args()
    main(args)
