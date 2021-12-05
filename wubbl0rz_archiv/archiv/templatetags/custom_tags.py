import datetime

from django import template

register = template.Library()


@register.filter
def duration(milliseconds):
    dur = str(datetime.timedelta(milliseconds=milliseconds))
    if "." in dur:
        dur = dur.split(".")[0]
    return dur
