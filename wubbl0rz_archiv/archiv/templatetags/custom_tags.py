import datetime

from django import template

register = template.Library()



@register.filter
def duration(milliseconds):
    return str(datetime.timedelta(milliseconds=milliseconds))[:-7]