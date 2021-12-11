import re

from django.core.paginator import Paginator
from django.db.models import Count
from django.db.models.functions import TruncYear
from django.shortcuts import get_object_or_404, render

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
    search = request.GET.get("s")
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
