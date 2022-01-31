#!/bin/sh

# wait for db and redis
/usr/bin/env bash ./wait-for-it.sh wub-db:"${DB_PORT}" -t 300 -s
/usr/bin/env bash ./wait-for-it.sh wub-redis:6379 -t 300 -s

# init django
python manage.py makemigrations
python manage.py migrate
python manage.py collectstatic --noinput

# import existing db if file exists
if [ -f "db.json" ]; then
    python manage.py loaddata db.json
    echo "DB imported from file"
fi

# create django super user
if [ -n "$DJANGO_SUPERUSER_USER" ] && [ -n "$DJANGO_SUPERUSER_PASSWORD" ]; then
    python -u manage.py shell -c "from django.contrib.auth.models import User; User.objects.create_superuser('$DJANGO_SUPERUSER_USER', password='$DJANGO_SUPERUSER_PASSWORD') if not User.objects.filter(username='$DJANGO_SUPERUSER_USER').exists() else print('Django user exists')"
fi

# Create twitch settings if none exists
python -u manage.py shell -c "from vods.models import ApiStorage; ApiStorage(ttv_client_id='$TWITCH_CLIENT_ID', ttv_client_secret='$TWITCH_CLIENT_SECRET').save() if ApiStorage.objects.filter().count() == 0 else print('ApiStorage exists')"

# run webserver and celery tasks
gunicorn --bind 0.0.0.0:8000 --workers $(($(nproc) + 1)) -k gevent settings.wsgi:application &
celery -A settings worker &
celery -A settings beat &

# create backups every 24 hours
while :; do
    sleep $(( 24 * 60 * 60 )) # 24 hours
    python -u manage.py dumpdata > "/backups/dump_$(date +%Y-%m-%d"_"%H-%M-%S).json"
done
