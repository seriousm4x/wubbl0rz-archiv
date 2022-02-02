from django.db import models


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
