from django.urls import include, path
from clips import views

urlpatterns = [
    path("all", views.all, name="all"),
    path("top30", views.top30, name="top30"),
    path("watch/<slug:uuid>", views.single_clip, name="single_clip")
]
