import argparse
import json
import logging
import os
import pathlib
import subprocess
import time


class Gen:
    def m3u8(self, mp4, outdir, id):
        logger.info(f"{id}: Creating m3u8 segments")
        if not os.path.isdir(os.path.join(outdir, id + "-segments")):
            os.mkdir(os.path.join(outdir, id + "-segments"))
        cmd = ["ffmpeg", "-hide_banner", "-loglevel", "error", "-stats", "-i", mp4, "-c", "copy",
               "-hls_playlist_type", "vod", "-hls_time", "10", "-hls_segment_filename", os.path.join(outdir, id + "-segments", id + "_%04d.ts"), os.path.join(outdir, id + "-segments", id + ".m3u8")]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        proc.communicate()

    def create_thumbnail(self, dir, id, duration):
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


    def get_metadata(self, dir, id):
        m3u8 = os.path.join(dir, id + "-segments", id + ".m3u8")

        # duration
        cmd = ["ffprobe", "-v", "error", "-show_entries", "format=duration", "-of",
               "default=noprint_wrappers=1:nokey=1", m3u8]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        out, _ = proc.communicate()
        duration = float(out.decode().strip())

        # fps
        cmd = ["ffprobe", "-v", "error", "-select_streams", "v", "-of",
               "default=noprint_wrappers=1:nokey=1", "-show_entries", "stream=r_frame_rate", m3u8]
        proc = subprocess.Popen(
            cmd, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        out, _ = proc.communicate()
        fps = float(out.decode().splitlines()[0].strip().split("/")[0])

        return duration, fps


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
        description='Generate all media files for vods from an mp4.')
    parser.add_argument('-i', '--input', type=pathlib.Path, required=True,
                        help='Input file path. Only use files, that already have h264 video with aac audio.')
    parser.add_argument('-id', type=str, required=True,
                        help='Vod id, such as "v2849020238"')
    parser.add_argument('-o', '--output', type=pathlib.Path, required=True,
                        help='Output folder path.')
    args = parser.parse_args()
    if not os.path.isfile(args.input):
        logger.error("Input file path is not a file")
        exit()
    if not os.path.isdir(args.output):
        logger.error("Output path is not a directory")
        exit()
    return args


def main(args):
    gen = Gen()
    gen.m3u8(args.input, args.output, args.id)
    duration, fps = gen.get_metadata(args.output, args.id)
    gen.create_thumbnail(args.output, args.id, duration)

    if not os.path.isfile(os.path.join(args.output, args.id + ".json")):
        json_file = os.path.join(args.output, args.id + ".json")
        with open(json_file, "w", encoding="utf-8") as f:
            json.dump({
                "title": args.id,
                "fps": fps,
                "timestamp": time.time()
            }, f)

    logger.info(
        f"Files have been created. Please edit {json_file} now, to set the correct metadata.")


if __name__ == "__main__":
    logger = set_logger()
    args = parse_args()
    main(args)
