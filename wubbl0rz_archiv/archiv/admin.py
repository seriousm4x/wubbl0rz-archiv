from django.contrib import admin
from .models import Vod

class VodAdmin(admin.ModelAdmin):
    list_display = ["title", "uuid", "duration", "date", "filename"]
    search_fields = ["title", "uuid", "duration", "date", "filename"]

admin.site.register(Vod, VodAdmin)