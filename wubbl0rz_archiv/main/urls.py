from django.urls import include, path
from rest_framework import routers

from main import serializers
from main import views as main_views

router = routers.DefaultRouter()
router.register(r"vods", serializers.VodViewSet, basename="vods")
router.register(r"clips", serializers.ClipViewSet, basename="clips")
router.register(r"years", serializers.YearsViewSet, basename="years")
router.register(r"emotes", serializers.EmoteViewSet, basename="emotes")
router.register(r"stats", serializers.StatsViewSet, basename="stats")
router.register(r"stats/db", serializers.DBViewSet, basename="stats/db")


urlpatterns = [
    path("", include(router.urls)),
    path("health/", main_views.health, name="health"),
    path("dl/<slug:type>/<slug:uuid>/", main_views.download, name="download"),
]
