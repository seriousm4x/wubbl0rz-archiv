from django.core.paginator import Paginator
from django.shortcuts import render, get_object_or_404
from main.models import ApiStorage
import os
from clips.models import Clip
import subprocess
from main.views import match_emotes
from django.conf import settings

from django.http.response import StreamingHttpResponse

def all(request):
    all_clips = Clip.objects.all()
    paginator = Paginator(all_clips.order_by("-created_at"), 36)
    page_number = request.GET.get("p")
    clips = paginator.get_page(page_number)
    api_obj = ApiStorage.objects.first()
    for c in clips:
        match_emotes(c)

    ctx = {
        "clips": clips,
        "api_obj": api_obj
    }

    return render(request, "clips.html", ctx)


def top30(request):
    return render(request, "clips.html", {})


def single_clip(request, uuid):
    clip = get_object_or_404(Clip, uuid=uuid)

    if request.GET.get("dl") == "1" and clip:
        cmd = ["ffmpeg", "-i", os.path.join(settings.MEDIA_ROOT, "clips", clip.clip_id + ".ts"),
               "-c", "copy", "-bsf:a", "aac_adtstoasc", "-movflags", "frag_keyframe+empty_moov", "-f", "mp4", "-"]
        proc = subprocess.Popen(
            cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)

        def iterator():
            while True:
                data = proc.stdout.read(4096)
                if not data:
                    proc.stdout.close()
                    break
                yield data

        response = StreamingHttpResponse(iterator(), content_type="video/mp4")
        response["Content-Disposition"] = f"attachment; filename={uuid}.mp4"
        return response

    api_obj = ApiStorage.objects.first()
    match_emotes(clip)

    ctx = {
        "clip": clip,
        "api_obj": api_obj
    }

    return render(request, "single_clip.html", ctx)
