from django.db import models
from django.utils import timezone


def gen_id():
    from uuid import uuid4
    while True:
        id = str(uuid4())[:6]
        if not Vod.objects.filter(uuid=id).exists():
            return id


class Vod(models.Model):
    uuid = models.SlugField(
        primary_key=True, default=gen_id, editable=False, unique=True)
    title = models.TextField(null=False, blank=False)
    duration = models.PositiveIntegerField()
    date = models.DateTimeField()
    filename = models.SlugField(null=False, blank=False)
    resolution = models.TextField(null=True, blank=True)
    fps = models.FloatField(null=True, blank=True)
    size = models.PositiveBigIntegerField(null=True, blank=True)
    publish = models.BooleanField(default=True)

    @property
    def bitrate(self):
        return int(self.size * 8 / self.duration)

    readonly_fields = ["bitrate"]


class Emote(models.Model):
    id = models.SlugField(primary_key=True, blank=False, null=False)
    name = models.SlugField(blank=False, null=False)
    url = models.URLField(blank=False, null=False)
    provider = models.SlugField(blank=False, null=False)
    outdated = models.BooleanField(default=False)


class ApiStorage(models.Model):
    broadcaster_id = models.SlugField(
        blank=False, null=False, default="108776574")
    ttv_client_id = models.SlugField(blank=False, null=False)
    ttv_client_secret = models.SlugField(blank=False, null=False)
    ttv_bearer_token = models.SlugField(blank=True, null=True)
    date_vods_updated = models.DateTimeField(blank=True, null=True)
    date_emotes_updated = models.DateTimeField(blank=True, null=True)
    is_live = models.BooleanField(blank=True, null=True, default=False)
