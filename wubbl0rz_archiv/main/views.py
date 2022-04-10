import datetime
import os
import queue
import re
import subprocess
from threading import Thread

from clips.models import Clip
from dateutil import relativedelta
from django.conf import settings
from django.db.models import Count, Sum
from django.db.models.functions.datetime import ExtractHour, ExtractWeekDay
from django.http.response import (HttpResponse, HttpResponseServerError,
                                  StreamingHttpResponse)
from django.shortcuts import get_object_or_404, render
from django.template.defaultfilters import date as _date
from django.utils import timezone
from django.utils.text import slugify
from vods.models import Vod

from main.models import ApiStorage, Emote


def index(request):
    all_vods = Vod.objects.filter(publish=True)
    vods = all_vods.order_by("-date")[:12]

    last_month = datetime.datetime.now(
        tz=timezone.get_current_timezone()) - datetime.timedelta(weeks=4)
    clips = Clip.objects.filter(
        date__gte=last_month).order_by("-view_count")[:12]

    api_obj = ApiStorage.objects.first()
    for v in vods:
        match_emotes(v)

    ctx = {
        "vods": vods,
        "clips": clips,
        "api_obj": api_obj
    }
    return render(request, "main_index.html", ctx)


def stats(request):
    all_vods = Vod.objects.filter(publish=True)
    all_clips = Clip.objects.filter()
    api_obj = ApiStorage.objects.first()
    total_duration = int(all_vods.aggregate(
        Sum("duration"))["duration__sum"]/3600)
    vod_size = int(all_vods.aggregate(Sum("size"))["size__sum"])
    clip_size = int(all_clips.aggregate(Sum("size"))["size__sum"])
    vods_per_weekday = list(Vod.objects.annotate(weekday=ExtractWeekDay("date")).order_by(
        "weekday").values("weekday").annotate(count=Count("uuid")).values_list("count", flat=True))

    # vods per month / weekday / hour
    vods_per_month_values = []
    vods_per_month_labels = []
    for i in range(11, -1, -1):
        first_day_of_month = timezone.now().replace(
            day=1) - relativedelta.relativedelta(months=i)
        month = _date(timezone.now() -
                      relativedelta.relativedelta(months=i), "M y")
        amount = Vod.objects.filter(
            date__range=[first_day_of_month, first_day_of_month + relativedelta.relativedelta(months=1)]).count()
        vods_per_month_labels.append(month)
        vods_per_month_values.append(amount)

    vods_per_hour_values = list(Vod.objects.annotate(hour=ExtractHour("date")).order_by(
        "hour").values("hour").annotate(count=Count("uuid")).values_list("count", flat=True))
    vods_per_hour_labels = list(Vod.objects.annotate(hour=ExtractHour("date")).order_by(
        "hour").values("hour").annotate(count=Count("uuid")).values_list("hour", flat=True))

    # clips per user
    most_clips_per_user = Clip.objects.values("creator__name").annotate(
        amount=Count("creator__name")).order_by("-amount")[:10]
    most_views_per_user = Clip.objects.values("creator__name").annotate(
        amount=Sum("view_count")).order_by("-amount")[:10]

    # emotes
    emote_count = Emote.objects.all().count()
    all_twitch_emotes = Emote.objects.filter(
        provider="twitch").order_by("name")
    all_bttv_emotes = Emote.objects.filter(provider="bttv").order_by("name")
    all_ffz_emotes = Emote.objects.filter(provider="ffz").order_by("name")

    ctx = {
        "all_vods": all_vods,
        "all_clips": all_clips,
        "api_obj": api_obj,
        "total_duration": total_duration,
        "total_size": vod_size+clip_size,
        "vods_per_weekday": vods_per_weekday,
        "vods_per_month_labels": vods_per_month_labels,
        "vods_per_month_values": vods_per_month_values,
        "vods_per_hour_values": vods_per_hour_values,
        "vods_per_hour_labels": vods_per_hour_labels,
        "most_clips_per_user": most_clips_per_user,
        "most_views_per_user": most_views_per_user,
        "emote_count": emote_count,
        "all_twitch_emotes": all_twitch_emotes,
        "all_bttv_emotes": all_bttv_emotes,
        "all_ffz_emotes": all_ffz_emotes,
    }
    return render(request, "stats.html", ctx)


def search(request):
    api_obj = ApiStorage.objects.first()
    search = request.GET.get("q")
    if not search:
        return render(request, "search.html", {"api_obj": api_obj})
    vods = Vod.objects.filter(title__icontains=search).order_by("-date")
    for v in vods:
        match_emotes(v)

    ctx = {
        "vods": vods,
        "searchquery": search,
        "api_obj": api_obj
    }
    return render(request, "search.html", ctx)


def match_emotes(item):
    all_emotes = map(
        str.lower, Emote.objects.all().values_list("name", flat=True))
    findemotes = re.compile(r'([A-Z]\w*)', re.IGNORECASE)
    for possible_emote in set(findemotes.findall(item.title)):
        if possible_emote.lower() in all_emotes:
            this_emote = Emote.objects.filter(
                name__icontains=possible_emote).first()
            item.emote_title = item.title.replace(
                possible_emote, f'<img src="{this_emote.url}" data-toggle="tooltip" title="{this_emote.name}" loading="lazy">')


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

    ff_queue = queue.Queue()
    cmd = ["ffmpeg", "-i", os.path.join(settings.MEDIA_ROOT, type, filename + "-segments", filename + ".m3u8"),
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
    response["Content-Disposition"] = f"attachment; filename={slugify(obj.date)}-{slugify(obj.title)}.mp4"
    return response
