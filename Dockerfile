FROM python:3.10.1-slim-bullseye as base


FROM base as build-pip

ENV PYTHONUNBUFFERED 1
WORKDIR /build-pip
COPY requirements.txt .
RUN python -m pip install --no-cache-dir --upgrade pip && \
    pip install --prefix=/build-pip --no-cache-dir -r requirements.txt


FROM base as build-avif

WORKDIR /build-avif
RUN apt-get update && apt-get -y install git build-essential zlib1g-dev libpng-dev libjpeg-dev cmake ninja-build yasm &&\
    rm -rf /var/lib/apt/lists/* &&\
    git clone -b v0.10.1 https://github.com/AOMediaCodec/libavif.git &&\
    cd libavif/ext/ &&\
    ./aom.cmd &&\
    cd .. && mkdir build && cd build &&\
    cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=0 -DAVIF_CODEC_AOM=1 -DAVIF_LOCAL_AOM=1 -DAVIF_BUILD_APPS=1 .. &&\
    make


FROM base

COPY --from=build-pip /build-pip /usr/local
COPY --from=build-avif /build-avif/libavif/build/avifenc /usr/bin/avifenc
COPY wubbl0rz_archiv /app
WORKDIR /app
RUN apt-get update && \
    apt-get -y install ffmpeg mediainfo curl && \
    rm -rf /var/lib/apt/lists/*

HEALTHCHECK --interval=10s \
    CMD curl -fs "http://localhost:8000/health/" || exit 1

CMD ["./run.sh"]
