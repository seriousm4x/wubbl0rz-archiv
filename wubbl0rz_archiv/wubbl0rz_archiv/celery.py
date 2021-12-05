import os

from celery import Celery
from celery.schedules import crontab

os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'wubbl0rz_archiv.settings')

app = Celery('wubbl0rz_archiv')

app.config_from_object('django.conf:settings', namespace='CELERY')

app.conf.timezone = "Europe/Berlin"

app.conf.beat_schedule = {
    "download_vods": {
        "task": "archiv.tasks.main_dl",
        "schedule": crontab(hour=3, minute=0)
    }
}

app.autodiscover_tasks()
