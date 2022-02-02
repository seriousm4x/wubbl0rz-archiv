from django.contrib import admin

from main.models import ApiStorage, Emote


class EmoteAdmin(admin.ModelAdmin):
    list_display = ["name", "id", "url", "provider", "outdated"]
    search_fields = ["name", "id", "url", "provider", "outdated"]


class ApiStorageAdmin(admin.ModelAdmin):
    list_display = ["broadcaster_id", "ttv_client_id",
                    "ttv_client_secret", "ttv_bearer_token"]
    search_fields = ["broadcaster_id", "ttv_client_id",
                     "ttv_client_secret", "ttv_bearer_token"]


admin.site.register(Emote, EmoteAdmin)
admin.site.register(ApiStorage, ApiStorageAdmin)
