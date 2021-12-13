import re

from dateutil import relativedelta
from django.core.paginator import Paginator
from django.db.models import Count, Sum
from django.db.models.functions import TruncYear
from django.db.models.functions.datetime import ExtractHour, ExtractWeekDay
from django.http.response import HttpResponse, HttpResponseServerError
from django.shortcuts import get_object_or_404, render
from django.template.defaultfilters import date as _date
from django.utils import timezone

from .models import ApiStorage, Emote, Vod


def match_emotes(vod):
    all_emotes = Emote.objects.all().values_list("name", flat=True)
    findemotes = re.compile(r'([A-Z]\w*)')
    for possible_emote in set(findemotes.findall(vod.title)):
        if possible_emote in all_emotes:
            this_emote = Emote.objects.filter(
                name__iexact=possible_emote).first()
            vod.emote_title = vod.title.replace(
                possible_emote, f'<img src="{this_emote.url}" data-toggle="tooltip" title="{this_emote.name}">')


def index(request):
    all_vods = Vod.objects.all()
    paginator = Paginator(all_vods.order_by("-date"), 30)
    page_number = request.GET.get("p")
    vods = paginator.get_page(page_number)
    api_obj = ApiStorage.objects.first()
    for v in vods:
        match_emotes(v)

    ctx = {
        "vods": vods,
        "all_vod_titles": list(all_vods.values_list("title", flat=True)),
        "api_obj": api_obj
    }
    return render(request, "index.html", ctx)


def single_vod(request, uuid):
    all_vods = Vod.objects.all()
    api_obj = ApiStorage.objects.first()
    vod = get_object_or_404(Vod, uuid=uuid)
    match_emotes(vod)

    ctx = {
        "vod": vod,
        "all_vod_titles": list(all_vods.values_list("title", flat=True)),
        "api_obj": api_obj
    }
    return render(request, "single_vod.html", ctx)


def years(request):
    all_vods = Vod.objects.all()
    vods = Vod.objects.all().order_by("-date")
    grouped_years = Vod.objects.annotate(year=TruncYear("date")).values(
        "year").annotate(c=Count('uuid')).values('year', 'c').order_by("-year")
    api_obj = ApiStorage.objects.first()
    for v in vods:
        match_emotes(v)

    ctx = {
        "vods": vods,
        "all_vod_titles": list(all_vods.values_list("title", flat=True)),
        "grouped_years": grouped_years,
        "api_obj": api_obj
    }
    return render(request, "years.html", ctx)


def search(request):
    api_obj = ApiStorage.objects.first()
    search = request.GET.get("q")
    if not search:
        return render(request, "search.html", {"api_obj": api_obj})
    all_vods = Vod.objects.all()
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


def stats(request):
    all_vods = Vod.objects.all()
    all_vod_titles = list(all_vods.values_list("title", flat=True))
    api_obj = ApiStorage.objects.first()
    vods_this_month = all_vods.filter(date__range=[
        timezone.now() - relativedelta.relativedelta(months=1), timezone.now()])
    total_duration = int(all_vods.aggregate(
        Sum("duration"))["duration__sum"]/3600)
    total_size = int(all_vods.aggregate(Sum("size"))["size__sum"])
    vods_per_weekday = list(Vod.objects.annotate(weekday=ExtractWeekDay("date")).values(
        "weekday").annotate(count=Count("uuid")).values_list("count", flat=True))

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

    # emotes
    emote_count = Emote.objects.all().count()
    all_twitch_emotes = Emote.objects.filter(
        provider="twitch").order_by("name")
    all_bttv_emotes = Emote.objects.filter(provider="bttv").order_by("name")
    all_ffz_emotes = Emote.objects.filter(provider="ffz").order_by("name")

    ctx = {
        "all_vods": all_vods,
        "api_obj": api_obj,
        "all_vod_titles": all_vod_titles,
        "vods_this_month": vods_this_month,
        "total_duration": total_duration,
        "total_size": total_size,
        "vods_per_weekday": vods_per_weekday,
        "vods_per_month_labels": vods_per_month_labels,
        "vods_per_month_values": vods_per_month_values,
        "vods_per_hour_values": vods_per_hour_values,
        "vods_per_hour_labels": vods_per_hour_labels,
        "emote_count": emote_count,
        "all_twitch_emotes": all_twitch_emotes,
        "all_bttv_emotes": all_bttv_emotes,
        "all_ffz_emotes": all_ffz_emotes,
    }
    return render(request, "stats.html", ctx)


def health(request):
    try:
        ApiStorage.objects.all()
        return HttpResponse("Ok")
    except Exception:
        return HttpResponseServerError("db: cannot connect to database.")
