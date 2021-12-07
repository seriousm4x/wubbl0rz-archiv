import os

from celery import Celery
from celery.schedules import crontab

os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'wubbl0rz_archiv.settings')

app = Celery('wubbl0rz_archiv')

app.config_from_object('django.conf:settings', namespace='CELERY')

app.conf.timezone = "Europe/Berlin"

app.conf.beat_schedule = {
    "download_vods": {
        "task": "archiv.tasks.download_vods",
        "schedule": crontab(hour=3, minute=0)
    },
    "update_emotes": {
        "task": "archiv.tasks.update_emotes",
        "schedule": crontab(hour=2, minute=0)
    },
    "check_live": {
        "task": "archiv.tasks.check_live",
        "schedule": crontab(minute="*/3")
    }
}

app.autodiscover_tasks()
