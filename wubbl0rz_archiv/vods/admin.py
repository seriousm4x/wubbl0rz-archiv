from django.contrib import admin

from vods.models import Vod


class VodAdmin(admin.ModelAdmin):
    list_display = ["title", "uuid", "duration", "date", "filename", "publish"]
    search_fields = ["title", "uuid", "duration", "date", "filename"]
    readonly_fields = ["bitrate"]
    ordering = ["-date"]


admin.site.register(Vod, VodAdmin)
