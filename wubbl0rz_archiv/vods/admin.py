from clips.admin import ClipInline
from django.contrib import admin
from django.db import models

from vods.models import Vod


class VodAdmin(admin.ModelAdmin):
    list_display = ["title", "uuid", "duration", "date", "filename", "publish", "clip_count"]
    search_fields = ["title", "uuid", "duration",
                     "date", "filename", "clip_count"]
    readonly_fields = ["bitrate"]
    ordering = ["-date"]
    inlines = [ClipInline]

    def get_queryset(self, request):
        queryset = super().get_queryset(request)
        queryset = queryset.annotate(
            _clip_count=models.Count("clip", distinct=True)
        )
        return queryset

    def clip_count(self, obj):
        return obj.clip_set.count()

    clip_count.admin_order_field = "_clip_count"


admin.site.register(Vod, VodAdmin)
