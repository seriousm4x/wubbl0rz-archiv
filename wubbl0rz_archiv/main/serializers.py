from clips.models import Clip, Game
from dateutil import relativedelta
from django.db.models import CharField, Count, F, Func, Sum, Value
from django.db.models.functions.datetime import ExtractHour
from django.template.defaultfilters import date as _date
from django.utils import timezone
from django.utils.dateparse import parse_datetime
from rest_framework import filters, mixins, serializers, viewsets
from rest_framework.pagination import PageNumberPagination
from rest_framework.response import Response
from rest_framework.utils.urls import replace_query_param
from vods.models import Vod

from main.models import ApiStorage, Emote


class StandardResultsSetPagination(PageNumberPagination):
    page_size = 50
    max_page_size = 500
    page_size_query_param = "page_size"

    def get_first_page(self):
        url = self.request.build_absolute_uri()
        return replace_query_param(url, self.page_query_param, 1)

    def get_last_page(self):
        url = self.request.build_absolute_uri()
        final = self.page.paginator.num_pages
        return replace_query_param(url, self.page_query_param, final)

    def get_paginated_response(self, data):
        return Response({
            "links": {
                "next": self.get_next_link(),
                "previous": self.get_previous_link(),
                "first": self.get_first_page(),
                "last": self.get_last_page()
            },
            "count": self.page.paginator.count,
            "current_page": self.page.number,
            "total_pages": self.page.paginator.num_pages,
            "results": data
        })


class VodSerializer(serializers.HyperlinkedModelSerializer):
    clip_set = serializers.SlugRelatedField(
        many=True,
        read_only=True,
        slug_field="uuid"
    )
    class Meta:
        model = Vod
        fields = ["uuid", "title", "duration", "bitrate",
                  "date", "filename", "resolution", "fps", "size", "clip_set"]


class VodViewSet(viewsets.ReadOnlyModelViewSet):
    serializer_class = VodSerializer

    def get_queryset(self):
        queryset = Vod.objects.filter(publish=True)
        try:
            year = self.request.query_params.get("year")
            if year is not None and year is not "":
                queryset = queryset.filter(date__year=year)
        except:
            pass
        finally:
            return queryset
    filter_backends = [filters.SearchFilter, filters.OrderingFilter]
    search_fields = ["title"]
    ordering_fields = ["date"]
    ordering = ["-date"]
    pagination_class = StandardResultsSetPagination


class GameSerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Game
        fields = ["game_id", "name"]


class ClipSerializer(serializers.HyperlinkedModelSerializer):
    creator = serializers.SlugRelatedField(
        read_only=True,
        slug_field="name"
    )
    game = GameSerializer()
    vod = serializers.SlugRelatedField(
        read_only=True,
        slug_field="uuid"
    )

    class Meta:
        model = Clip
        fields = ["uuid", "title", "clip_id", "creator", "view_count", "date", "duration", "resolution", "size", "game", "vod", "bitrate"]


class ClipViewSet(viewsets.ReadOnlyModelViewSet):
    serializer_class = ClipSerializer

    def get_queryset(self):
        queryset = Clip.objects.all()
        try:
            date_from = self.request.query_params.get("date_from")
            date_to = self.request.query_params.get("date_to")
            if date_from is not None and date_from is not "":
                queryset = queryset.filter(date__gt=parse_datetime(date_from))
            if date_to is not None and date_to is not "":
                queryset = queryset.filter(date__lt=parse_datetime(date_to))
        except:
            pass
        finally:
            return queryset

    filter_backends = [filters.SearchFilter, filters.OrderingFilter]
    search_fields = ["title"]
    ordering_fields = ["view_count", "date"]
    ordering = ["-date"]
    pagination_class = StandardResultsSetPagination


class YearsSerializer(serializers.HyperlinkedModelSerializer):
    year = serializers.IntegerField()
    count = serializers.IntegerField()

    class Meta:
        model = Vod
        fields = ["year", "count"]


class YearsViewSet(mixins.ListModelMixin, viewsets.GenericViewSet):
    queryset = Vod.objects.annotate(year=Func(
        F("date"),
        Value("yyyy"),
        function="to_char",
        output_field=CharField()
    )).values("year").annotate(count=Count("year")).order_by("-year")
    serializer_class = YearsSerializer


class EmoteSerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Emote
        fields = ["id", "name", "url", "provider"]


class EmoteViewSet(viewsets.ReadOnlyModelViewSet):
    def get_queryset(self):
        queryset = Emote.objects.all()
        provider = self.request.query_params.get("provider")
        name = self.request.query_params.get("name")
        if provider is not None:
            queryset = queryset.filter(provider=provider)
        if name is not None:
            queryset = queryset.filter(name__icontains=name)
        return queryset
    serializer_class = EmoteSerializer
    pagination_class = StandardResultsSetPagination


class StatsViewSet(viewsets.ViewSet):
    def get_vods_per_month(self, i):
        first_day_of_month = timezone.now().replace(
            day=1) - relativedelta.relativedelta(months=i)
        month = _date(timezone.now() -
                      relativedelta.relativedelta(months=i), "M y")
        count = Vod.objects.filter(date__range=[
            first_day_of_month, first_day_of_month +
            relativedelta.relativedelta(months=1)]).count()
        return month, count

    def list(self, request):
        all_vods = Vod.objects.filter(publish=True)
        ctx = {}
        ctx["count_vods_total"] = all_vods.count()
        ctx["count_vods_1m"] = all_vods.filter(date__range=[timezone.now(
        ) - relativedelta.relativedelta(months=1), timezone.now()]).count()
        ctx["count_h_streamed"] = int(all_vods.aggregate(
            Sum("duration"))["duration__sum"]/3600)
        ctx["archiv_size_bytes"] = int(
            all_vods.aggregate(Sum("size"))["size__sum"])

        ctx["vods_per_month"] = []
        for i in range(11, -1, -1):
            month, count = self.get_vods_per_month(i)
            ctx["vods_per_month"].append({
                "month": month,
                "count": count
            })

        ctx["vods_per_weekday"] = []
        weekdays = [(1, "Sonntag"),
                    (2, "Montag"),
                    (3, "Dienstag"),
                    (4, "Mittwoch"),
                    (5, "Donnerstag"),
                    (6, "Freitag"),
                    (7, "Samstag")]
        for day in weekdays:
            ctx["vods_per_weekday"].append({
                "weekday": day[1],
                "count": Vod.objects.filter(date__week_day=day[0]).count()
            })

        ctx["start_by_time"] = Vod.objects.annotate(hour=ExtractHour("date")).order_by(
            "hour").values("hour").annotate(count=Count("uuid"))

        return Response(ctx)


class DBViewSet(viewsets.ViewSet):
    def list(self, request):
        ctx = {}
        ctx["last_vod_sync"] = ApiStorage.objects.first().date_vods_updated
        ctx["last_emote_sync"] = ApiStorage.objects.first().date_emotes_updated
        return Response(ctx)
