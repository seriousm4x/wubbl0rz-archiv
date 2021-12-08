from django.db.models import Count
from django.db.models.functions import TruncYear
from django.shortcuts import get_object_or_404, render
from django.core.paginator import Paginator
from .models import Emote, Vod, ApiStorage
import re


def match_emotes(vod):
    all_emotes = Emote.objects.all().values_list("name", flat=True)
    findemotes = re.compile(r'([A-Z]\w*)')
    for possible_emote in set(findemotes.findall(vod.title)):
        if possible_emote in all_emotes:
            vod.title = vod.title.replace(possible_emote, f'<img src="{Emote.objects.filter(name__iexact=possible_emote).first().url}" class="inline-emote">')


def index(request):
    all_vods = Vod.objects.all().order_by("-date")
    paginator = Paginator(all_vods, 30)
    page_number = request.GET.get("p")
    vods = paginator.get_page(page_number)
    api_obj = ApiStorage.objects.first()
    for v in vods:
        match_emotes(v)

    ctx = {
        "vods": vods,
        "api_obj": api_obj
    }
    return render(request, "index.html", ctx)


def single_vod(request, uuid):
    api_obj = ApiStorage.objects.first()
    vod = get_object_or_404(Vod, uuid=uuid)
    match_emotes(vod)

    ctx = {
        "vod": vod,
        "api_obj": api_obj
    }
    return render(request, "single_vod.html", ctx)

def years(request):
    vods = Vod.objects.all().order_by("-date")
    grouped_years = Vod.objects.annotate(year=TruncYear("date")).values("year").annotate(c=Count('uuid')).values('year', 'c').order_by("-year")
    api_obj = ApiStorage.objects.first()
    for v in vods:
        match_emotes(v)

    ctx = {
        "vods": vods,
        "grouped_years": grouped_years,
        "api_obj": api_obj
    }
    return render(request, "years.html", ctx)

def search(request):
    search = request.GET.get("s")
    vods = Vod.objects.filter(title__icontains=search)
    api_obj = ApiStorage.objects.first()
    for v in vods:
        match_emotes(v)

    ctx = {
        "vods": vods,
        "searchquery": search,
        "api_obj": api_obj
    }
    return render(request, "search.html", ctx)
