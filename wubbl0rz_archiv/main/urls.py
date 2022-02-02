from django.urls import include, path
from rest_framework import routers
from vods import serializers as vods_serializers

from main import views as main_views

router = routers.DefaultRouter()
router.register(r"vods", vods_serializers.VodViewSet, basename="vods")
router.register(r"years", vods_serializers.YearsViewSet, basename="years")
router.register(r"emotes", vods_serializers.EmoteViewSet, basename="emotes")
router.register(r"stats", vods_serializers.StatsViewSet, basename="stats")
router.register(r"stats/db", vods_serializers.DBViewSet, basename="stats/db")


urlpatterns = [
    path("", main_views.index, name="index"),
    path("search/", main_views.search, name="search"),
    path("stats/", main_views.stats, name="stats"),
    path("health/", main_views.health, name="health"),
    path("api/", include(router.urls)),
]
