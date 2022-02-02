from django.db import models


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
