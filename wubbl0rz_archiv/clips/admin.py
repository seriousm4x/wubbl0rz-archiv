from django.contrib import admin

from clips.models import Clip, Game


class ClipAdmin(admin.ModelAdmin):
    list_display = ["title", "uuid", "clip_id",
                    "creator_name", "view_count", "created_at", "duration"]
    search_fields = ["title", "uuid", "clip_id", "creator_name", "created_at"]
    readonly_fields = ["bitrate"]
    ordering = ["-created_at"]


class GameAdmin(admin.ModelAdmin):
    list_display = ["name", "game_id", "box_art_url"]
    search_fields = ["name", "game_id", "box_art_url"]
    ordering = ["name"]


admin.site.register(Clip, ClipAdmin)
admin.site.register(Game, GameAdmin)
