from django.contrib import admin

from .models import ApiStorage, ChatMessage, Emote, Vod


class VodAdmin(admin.ModelAdmin):
    list_display = ["title", "uuid", "duration", "date", "filename"]
    search_fields = ["title", "uuid", "duration", "date", "filename"]
    readonly_fields = ["bitrate"]
    ordering = ["-date"]


class EmoteAdmin(admin.ModelAdmin):
    list_display = ["name", "id", "url", "provider", "outdated"]
    search_fields = ["name", "id", "url", "provider", "outdated"]


class ApiStorageAdmin(admin.ModelAdmin):
    list_display = ["broadcaster_id", "ttv_client_id",
                    "ttv_client_secret", "ttv_bearer_token"]
    search_fields = ["broadcaster_id", "ttv_client_id",
                     "ttv_client_secret", "ttv_bearer_token"]


class ChatMessageAdmin(admin.ModelAdmin):
    list_display = ["timestamp", "raw"]
    search_fields = ["timestamp", "raw"]


admin.site.register(Vod, VodAdmin)
admin.site.register(Emote, EmoteAdmin)
admin.site.register(ApiStorage, ApiStorageAdmin)
admin.site.register(ChatMessage, ChatMessageAdmin)
