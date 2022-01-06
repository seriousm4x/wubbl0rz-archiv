from django.conf import settings
from django.conf.urls.static import static
from django.urls import include, path
from rest_framework import routers

from archiv import serializers, views

router = routers.DefaultRouter()
router.register(r"vods", serializers.VodViewSet, basename="vods")
router.register(r"years", serializers.YearsList, basename="years")
router.register(r"emotes", serializers.EmoteViewSet, basename="emotes")
router.register(r"stats", serializers.StatsViewSet, basename="stats")

urlpatterns = [
    path("", views.index, name="index"),
    path("v/<slug:uuid>/", views.single_vod, name="single_vod"),
    path("years/", views.years, name="years"),
    path("search/", views.search, name="search"),
    path("stats/", views.stats, name="stats"),
    path("health/", views.health, name="health"),
    path("api/", include(router.urls))
] + static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
