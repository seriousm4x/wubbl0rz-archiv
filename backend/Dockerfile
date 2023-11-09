FROM golang:latest as base
ARG DEBIAN_FRONTEND=noninteractive
WORKDIR /build
COPY . .
RUN go mod download && \
    go build -ldflags="-s -w" ./cmd/wubbl0rz-archiv

FROM debian:stable-slim
ARG DEBIAN_FRONTEND=noninteractive
WORKDIR /app
COPY --from=base /build/wubbl0rz-archiv ./
RUN apt-get update && \
    apt-get install tzdata ffmpeg -y &&\
    apt-get clean autoclean -y  &&\
    apt-get autoremove -y  &&\
    rm -rf /var/lib/{apt,dpkg,cache,log}/
HEALTHCHECK --interval=10s \
    CMD wget --no-verbose --tries=1 --spider "http://localhost:8090/api/health" || exit 1
ENTRYPOINT ["./wubbl0rz-archiv", "serve", "--http", "0.0.0.0:8090"]
