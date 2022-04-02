import datetime
import os
import queue
import subprocess
from threading import Thread

from django.conf import settings
from django.core.paginator import Paginator
from django.http.response import StreamingHttpResponse
from django.shortcuts import get_object_or_404, render
from django.utils import timezone
from django.utils.text import slugify
from main.models import ApiStorage
from main.views import match_emotes

from clips.models import Clip


def all(request):
    all_clips = Clip.objects.all().order_by("-view_count")
    paginator = Paginator(all_clips, 36)
    page_number = request.GET.get("p")
    clips = paginator.get_page(page_number)
    api_obj = ApiStorage.objects.first()
    for c in clips:
        match_emotes(c)

    ctx = {
        "title": "Alle Clips",
        "clips": clips,
        "api_obj": api_obj
    }

    return render(request, "clips.html", ctx)


def top30(request):
    last_month = datetime.datetime.now(
        tz=timezone.get_current_timezone()) - datetime.timedelta(days=30)
    all_clips = Clip.objects.filter(date__gte=last_month)
    paginator = Paginator(all_clips.order_by("-view_count"), 36)
    page_number = request.GET.get("p")
    clips = paginator.get_page(page_number)
    api_obj = ApiStorage.objects.first()
    for c in clips:
        match_emotes(c)

    ctx = {
        "title": "Top Clips 30 Tage",
        "clips": clips,
        "api_obj": api_obj
    }

    return render(request, "clips.html", ctx)


def single_clip(request, uuid):
    clip = get_object_or_404(Clip, uuid=uuid)

    if request.GET.get("dl") == "1" and clip:
        ff_queue = queue.Queue()
        cmd = ["ffmpeg", "-i", os.path.join(settings.MEDIA_ROOT, "clips", clip.clip_id + "-segments", clip.clip_id + ".m3u8"),
               "-c", "copy", "-bsf:a", "aac_adtstoasc", "-movflags", "frag_keyframe+empty_moov", "-f", "mp4", "-"]

        def read_output(proc):
            while True:
                data = proc.stdout.read(4096)
                if not data:
                    break
                ff_queue.put(data)

        proc = subprocess.Popen(
            cmd, stdout=subprocess.PIPE, stderr=subprocess.DEVNULL)
        t = Thread(target=read_output, args=(proc,))

        def iterator():
            t.start()
            while True:
                if proc.poll() is not None and ff_queue.empty():
                    proc.kill()
                    break
                try:
                    data = ff_queue.get()
                    ff_queue.task_done()
                    yield data
                except queue.Empty:
                    pass

        response = StreamingHttpResponse(iterator(), content_type="video/mp4")
        response["Content-Disposition"] = f"attachment; filename={slugify(clip.date)}-{slugify(clip.title)}.mp4"
        return response

    if clip:
        match_emotes(clip)

    api_obj = ApiStorage.objects.first()

    ctx = {
        "clip": clip,
        "api_obj": api_obj
    }

    return render(request, "single_clip.html", ctx)
