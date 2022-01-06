from dateutil import relativedelta
from django.db.models import CharField, Count, F, Func, Sum, Value
from django.db.models.functions.datetime import ExtractHour
from django.template.defaultfilters import date as _date
from django.utils import timezone
from rest_framework import mixins, serializers, viewsets
from rest_framework.pagination import PageNumberPagination
from rest_framework.response import Response

from archiv.models import Emote, Vod


class StandardResultsSetPagination(PageNumberPagination):
    page_size = 50
    max_page_size = 500
    page_size_query_param = "page_size"


class VodSerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Vod
        fields = ["uuid", "title", "duration",
                  "date", "filename", "resolution", "fps", "size"]


class VodViewSet(viewsets.ReadOnlyModelViewSet):
    def get_queryset(self):
        queryset = Vod.objects.all().order_by("-date")
        year = self.request.query_params.get("year")
        if year is not None:
            queryset = queryset.filter(date__year=year)
        return queryset
    serializer_class = VodSerializer
    pagination_class = StandardResultsSetPagination


class YearsSerializer(serializers.HyperlinkedModelSerializer):
    year = serializers.IntegerField()
    count = serializers.IntegerField()

    class Meta:
        model = Vod
        fields = ["year", "count"]


class YearsList(mixins.ListModelMixin, viewsets.GenericViewSet):
    queryset = Vod.objects.annotate(year=Func(
        F("date"),
        Value("yyyy"),
        function="to_char",
        output_field=CharField()
    )).values("year").annotate(count=Count("year")).order_by("-year")
    serializer_class = YearsSerializer
    http_method_names = ["get"]


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
        amount = Vod.objects.filter(date__range=[
                                    first_day_of_month, first_day_of_month +
                                    relativedelta.relativedelta(months=1)]).count()
        return month, amount

    def list(self, request):
        all_vods = Vod.objects.all()
        ctx = {}
        ctx["count_vods_total"] = all_vods.count()
        ctx["count_vods_1m"] = all_vods.filter(date__range=[timezone.now(
        ) - relativedelta.relativedelta(months=1), timezone.now()]).count()
        ctx["count_h_streamed"] = int(all_vods.aggregate(
            Sum("duration"))["duration__sum"]/3600)
        ctx["archiv_size_bytes"] = int(
            all_vods.aggregate(Sum("size"))["size__sum"])

        ctx["vods_per_month"] = {}
        for i in range(11, -1, -1):
            month, amount = self.get_vods_per_month(i)
            ctx["vods_per_month"][month] = amount

        ctx["vods_per_weekday"] = {}
        weekdays = [(1, "Sonntag"),
                    (2, "Montag"),
                    (3, "Dienstag"),
                    (4, "Mittwoch"),
                    (5, "Donnerstag"),
                    (6, "Freitag"),
                    (7, "Samstag")]
        for day in weekdays:
            ctx["vods_per_weekday"][day[1]] = Vod.objects.filter(
                date__week_day=day[0]).count()

        ctx["start_by_time"] = Vod.objects.annotate(hour=ExtractHour("date")).order_by(
            "hour").values("hour").annotate(count=Count("uuid"))

        return Response(ctx)
