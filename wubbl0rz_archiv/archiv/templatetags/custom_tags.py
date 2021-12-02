import datetime

from archiv.models import Vod
from django import template
from django.db.models import Count
from django.db.models.functions import TruncYear

register = template.Library()


@register.filter
def duration(seconds):
    return str(datetime.timedelta(seconds=seconds))
