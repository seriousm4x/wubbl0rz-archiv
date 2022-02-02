from django.conf import settings
from django.conf.urls.static import static
from django.urls import path

from vods import views

urlpatterns = [
    path("", views.vods, name="vods"),
    path("watch/<slug:uuid>/", views.single_vod, name="single_vod"),
    path("years/", views.years, name="years"),
] + static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
