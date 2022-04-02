import os
import queue
import subprocess
from threading import Thread

from django.conf import settings
from django.core.paginator import Paginator
from django.db.models import Count
from django.db.models.functions import TruncYear
from django.http.response import StreamingHttpResponse
from django.shortcuts import get_object_or_404, render
from django.utils.text import slugify
from main.models import ApiStorage
from main.views import match_emotes

from vods.models import Vod


def vods(request):
    all_vods = Vod.objects.filter(publish=True)
    paginator = Paginator(all_vods.order_by("-date"), 36)
    page_number = request.GET.get("p")
    vods = paginator.get_page(page_number)
    api_obj = ApiStorage.objects.first()
    for v in vods:
        match_emotes(v)

    ctx = {
        "vods": vods,
        "api_obj": api_obj
    }
    return render(request, "vods.html", ctx)


def single_vod(request, uuid):
    vod = get_object_or_404(Vod, uuid=uuid)

    if request.GET.get("dl") == "1" and vod:
        ff_queue = queue.Queue()
        cmd = ["ffmpeg", "-i", os.path.join(settings.MEDIA_ROOT, "vods", vod.filename + "-segments", vod.filename + ".m3u8"),
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
        response["Content-Disposition"] = f"attachment; filename={slugify(vod.date)}-{slugify(vod.title)}.mp4"
        return response

    if vod:
        match_emotes(vod)

    api_obj = ApiStorage.objects.first()

    ctx = {
        "vod": vod,
        "api_obj": api_obj
    }
    return render(request, "single_vod.html", ctx)


def years(request):
    vods = Vod.objects.filter(publish=True).order_by("-date")
    grouped_years = Vod.objects.annotate(year=TruncYear("date")).values(
        "year").annotate(c=Count('uuid')).values('year', 'c').order_by("-year")
    api_obj = ApiStorage.objects.first()
    for v in vods:
        match_emotes(v)

    ctx = {
        "vods": vods,
        "grouped_years": grouped_years,
        "api_obj": api_obj
    }
    return render(request, "years.html", ctx)
