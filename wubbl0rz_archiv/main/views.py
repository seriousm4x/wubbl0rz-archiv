import datetime
import re

from clips.models import Clip
from dateutil import relativedelta
from django.db.models import Count, Sum
from django.db.models.functions.datetime import ExtractHour, ExtractWeekDay
from django.http.response import HttpResponse, HttpResponseServerError
from django.shortcuts import render
from django.template.defaultfilters import date as _date
from django.utils import timezone
from vods.models import Vod
from clips.models import Clip

from main.models import ApiStorage, Emote


def index(request):
    all_vods = Vod.objects.filter(publish=True)
    vods = all_vods.order_by("-date")[:12]

    last_month = datetime.datetime.now(
        tz=timezone.get_current_timezone()) - datetime.timedelta(weeks=4)
    clips = Clip.objects.filter(
        created_at__gte=last_month).order_by("-view_count")[:12]

    api_obj = ApiStorage.objects.first()
    for v in vods:
        match_emotes(v)

    ctx = {
        "vods": vods,
        "clips": clips,
        "all_vod_titles": list(all_vods.values_list("title", flat=True)),
        "api_obj": api_obj
    }
    return render(request, "main_index.html", ctx)


def stats(request):
    all_vods = Vod.objects.filter(publish=True)
    all_clips = Clip.objects.filter()
    all_vod_titles = list(all_vods.values_list("title", flat=True))
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
        "all_vod_titles": all_vod_titles,
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
    all_vods = Vod.objects.filter(publish=True)
    vods = Vod.objects.filter(title__icontains=search).order_by("-date")
    for v in vods:
        match_emotes(v)

    ctx = {
        "vods": vods,
        "all_vod_titles": list(all_vods.values_list("title", flat=True)),
        "searchquery": search,
        "api_obj": api_obj
    }
    return render(request, "search.html", ctx)


def match_emotes(item):
    all_emotes = Emote.objects.all().values_list("name", flat=True)
    findemotes = re.compile(r'([A-Z]\w*)')
    for possible_emote in set(findemotes.findall(item.title)):
        if possible_emote in all_emotes:
            this_emote = Emote.objects.filter(
                name__iexact=possible_emote).first()
            item.emote_title = item.title.replace(
                possible_emote, f'<img src="{this_emote.url}" data-toggle="tooltip" title="{this_emote.name}" loading="lazy">')


def health(request):
    try:
        ApiStorage.objects.all()
        return HttpResponse("Ok")
    except Exception:
        return HttpResponseServerError("db: cannot connect to database.")
