from django.conf import settings
from django.conf.urls.static import static
from django.urls import path

from archiv import views

urlpatterns = [
    path("", views.index, name="index"),
    path("v/<slug:uuid>/", views.single_vod, name="single_vod"),
    path("years/", views.years, name="years"),
    path("search/", views.search, name="search"),
] + static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
