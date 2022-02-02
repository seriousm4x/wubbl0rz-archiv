from django.db import models


def gen_id():
    from uuid import uuid4
    while True:
        id = str(uuid4())[:6]
        if not Clip.objects.filter(uuid=id).exists():
            return id


class Clip(models.Model):
    uuid = models.SlugField(
        primary_key=True, default=gen_id, editable=False, unique=True)
    clip_id = models.SlugField(max_length=100)
    creator_name = models.SlugField(max_length=30)
    game_id = models.PositiveIntegerField(blank=True, null=True)
    title = models.CharField(max_length=150)
    view_count = models.PositiveIntegerField()
    created_at = models.DateTimeField()
    duration = models.PositiveIntegerField(null=True, blank=True)
    resolution = models.TextField(null=True, blank=True)
    size = models.PositiveBigIntegerField(null=True, blank=True)

    @property
    def bitrate(self):
        return int(self.size * 8 / self.duration)


class Game(models.Model):
    game_id = models.PositiveIntegerField(blank=True, null=True)
    name = models.CharField(max_length=150)
    box_art_url = models.URLField(max_length=150, default="0")
