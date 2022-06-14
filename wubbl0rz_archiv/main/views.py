import os
import queue
import subprocess
from threading import Thread

from clips.models import Clip
from django.conf import settings
from django.http.response import (HttpResponse, HttpResponseServerError,
                                  StreamingHttpResponse)
from django.shortcuts import get_object_or_404
from django.utils.text import slugify
from vods.models import Vod

from main.models import ApiStorage


def health(request):
    try:
        ApiStorage.objects.all()
        return HttpResponse("Ok")
    except Exception:
        return HttpResponseServerError("db: cannot connect to database.")


def download(request, type, uuid):
    if type == "vods":
        obj = get_object_or_404(Vod, uuid=uuid)
        filename = obj.filename
    elif type == "clips":
        obj = get_object_or_404(Clip, uuid=uuid)
        filename = obj.clip_id
    else:
        return

    ff_queue = queue.Queue(maxsize=1024*8)
    cmd = ["ffmpeg", "-i", os.path.join(settings.MEDIA_ROOT, type, filename + "-segments", filename + ".m3u8"),
           "-c", "copy", "-bsf:a", "aac_adtstoasc", "-movflags", "frag_keyframe+empty_moov", "-f", "mp4", "-"]

    def read_output(proc):
        while True:
            data = proc.stdout.read(1024*4)
            if not data:
                break
            # try to put new data in queue. kill after timeout if no data has been pulled => client disconnected
            # a client must have downloaded data size within timeout, else the download quits
            try:
                ff_queue.put(data, timeout=60*5)
            except queue.Full:
                proc.kill()
                with ff_queue.mutex:
                    ff_queue.queue.clear()
                break

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

    proc = subprocess.Popen(
        cmd, bufsize=1024*4, stdout=subprocess.PIPE, stderr=subprocess.DEVNULL)
    t = Thread(target=read_output, args=(proc,))

    response = StreamingHttpResponse(iterator(), content_type="video/mp4")
    response["Content-Disposition"] = f"attachment; filename={slugify(obj.date)}-{slugify(obj.title)}.mp4"
    return response
