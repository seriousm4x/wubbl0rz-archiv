from django.contrib import admin
from .models import Vod, Emote, ApiStorage

class VodAdmin(admin.ModelAdmin):
    list_display = ["title", "uuid", "duration", "date", "filename"]
    search_fields = ["title", "uuid", "duration", "date", "filename"]

class EmoteAdmin(admin.ModelAdmin):
    list_display = ["name", "id", "url", "provider", "outdated"]
    search_fields = ["name", "id", "url", "provider", "outdated"]

class ApiStorageAdmin(admin.ModelAdmin):
    list_display = ["broadcaster_id", "ttv_client_id", "ttv_client_secret", "ttv_bearer_token"]
    search_fields = ["broadcaster_id", "ttv_client_id", "ttv_client_secret", "ttv_bearer_token"]


admin.site.register(Vod, VodAdmin)
admin.site.register(Emote, EmoteAdmin)
admin.site.register(ApiStorage, ApiStorageAdmin)
