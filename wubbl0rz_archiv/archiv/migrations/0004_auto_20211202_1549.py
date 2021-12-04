# Generated by Django 3.2.9 on 2021-12-02 14:49

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('archiv', '0003_rename_format_note_vod_format_id'),
    ]

    operations = [
        migrations.RemoveField(
            model_name='vod',
            name='format_id',
        ),
        migrations.AddField(
            model_name='vod',
            name='size',
            field=models.PositiveBigIntegerField(blank=True, null=True),
        ),
    ]