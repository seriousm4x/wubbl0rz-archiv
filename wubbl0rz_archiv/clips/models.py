from django.db import models
from vods.models import Vod


def gen_id():
    from uuid import uuid4
    while True:
        id = str(uuid4())[:6]
        if not Clip.objects.filter(uuid=id).exists():
            return id


class Game(models.Model):
    game_id = models.PositiveIntegerField()
    name = models.CharField(max_length=150, default="Unknown")

    def __str__(self):
        return self.name


class Creator(models.Model):
    creator_id = models.PositiveIntegerField()
    name = models.SlugField(max_length=30)

    def __str__(self):
        return self.name


class Clip(models.Model):
    uuid = models.SlugField(
        primary_key=True, default=gen_id, editable=False, unique=True)
    clip_id = models.SlugField(max_length=100)
    creator = models.ForeignKey(
        Creator, on_delete=models.CASCADE, null=True, blank=True)
    game = models.ForeignKey(
        Game, on_delete=models.CASCADE, null=True, blank=True)
    title = models.CharField(max_length=150)
    view_count = models.PositiveIntegerField()
    date = models.DateTimeField()
    duration = models.PositiveIntegerField(null=True, blank=True)
    resolution = models.TextField(null=True, blank=True)
    size = models.PositiveBigIntegerField(null=True, blank=True)
    vod = models.ForeignKey(
        Vod, on_delete=models.CASCADE, null=True, blank=True)

    @property
    def bitrate(self):
        return int(self.size * 8 / self.duration)

    def __str__(self):
        return self.clip_id
