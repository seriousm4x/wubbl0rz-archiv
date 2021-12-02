from django.db.models import Count
from django.db.models.functions import TruncYear
from django.shortcuts import get_object_or_404, render
from django.core.paginator import Paginator
from .models import Vod


def index(request):
    all_vods = Vod.objects.all().order_by("-date")
    paginator = Paginator(all_vods, 12)
    page_number = request.GET.get("p")
    vods = paginator.get_page(page_number)
    ctx = {
        "vods": vods,
    }
    return render(request, "index.html", ctx)


def single_vod(request, uuid):
    ctx = {
        "vod": get_object_or_404(Vod, uuid=uuid)
    }
    return render(request, "single_vod.html", ctx)

def years(request):
    vods = Vod.objects.all().order_by("-date")
    grouped_years = Vod.objects.annotate(year=TruncYear("date")).values("year").annotate(c=Count('uuid')).values('year', 'c').order_by("-year")
    ctx = {
        "vods": vods,
        "grouped_years": grouped_years
    }
    return render(request, "years.html", ctx)

def search(request):
    search = request.GET.get("s")
    vods = Vod.objects.filter(title__icontains=search)
    ctx = {
        "vods": vods,
        "searchquery": search
    }
    return render(request, "search.html", ctx)
