import os

from celery import Celery
from celery.schedules import crontab

os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'settings.settings')

app = Celery('settings')

app.config_from_object('django.conf:settings', namespace='CELERY')

app.conf.timezone = "Europe/Berlin"

app.conf.beat_schedule = {
    "check_live": {
        "task": "main.tasks.check_live",
        "schedule": crontab(minute="*/3")
    },
    "update_emotes": {
        "task": "main.tasks.update_emotes",
        "schedule": crontab(hour=2, minute=0)
    },
    "download_vods": {
        "task": "vods.tasks.download_vods",
        "schedule": crontab(hour=3, minute=0)
    },
    "download_clips": {
        "task": "clips.tasks.download_clips",
        "schedule": crontab(hour=5, minute=0)
    }
}

app.autodiscover_tasks()
