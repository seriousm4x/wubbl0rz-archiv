from django.db import models
from django.shortcuts import get_object_or_404


def gen_id():
    from uuid import uuid4
    while True:
        id = str(uuid4())[:6]
        if not Vod.objects.filter(uuid=id).exists():
            return id

class Vod(models.Model):
    uuid = models.SlugField(primary_key=True, default=gen_id, editable=False, unique=True)
    title = models.TextField(null=False, blank=False)
    duration = models.PositiveIntegerField()
    date = models.DateTimeField()
    filename = models.SlugField(null=False, blank=False)
    resolution = models.TextField(null=True, blank=True)
    bitrate = models.PositiveIntegerField(null=True, blank=True)
    fps = models.PositiveSmallIntegerField(null=True, blank=True)
    size = models.PositiveBigIntegerField(null=True, blank=True)