from django.contrib import admin

from clips.models import Clip, Creator, Game


class ClipAdmin(admin.ModelAdmin):
    list_display = ["title", "uuid", "clip_id",
                    "creator", "view_count", "date", "duration", "vod"]
    search_fields = ["title", "uuid", "clip_id",
                     "creator", "date", "vod"]
    readonly_fields = ["bitrate"]
    ordering = ["-date"]


class GameAdmin(admin.ModelAdmin):
    list_display = ["name", "game_id", "box_art"]
    search_fields = ["name", "game_id"]
    ordering = ["name"]


class CreatorAdmin(admin.ModelAdmin):
    list_display = ["name", "creator_id"]
    search_fields = ["name", "creator_id"]
    ordering = ["name"]


class ClipInline(admin.TabularInline):
    model = Clip
    extra = 0


admin.site.register(Clip, ClipAdmin)
admin.site.register(Game, GameAdmin)
admin.site.register(Creator, CreatorAdmin)
