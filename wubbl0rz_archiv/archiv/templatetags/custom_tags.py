import datetime

from django import template

register = template.Library()


@register.filter
def duration(seconds):
    dur = str(datetime.timedelta(seconds=seconds))
    if "." in dur:
        dur = dur.split(".")[0]
    return dur
