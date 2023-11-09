FROM golang:latest as base
WORKDIR /build
COPY cmd external internal go.mod go.sum ./
RUN go mod download && \
    apt-get update && \
    apt-get install libc-bin gcc && \
    go build -ldflags="-s -w" ./cmd/wubbl0rz-archiv

FROM debian:stable-slim
WORKDIR /app
COPY --from=base /build/wubbl0rz-archiv ./
RUN apt-get update && \
    apt-get install tzdata ffmpeg &&\
    apt-get clean autoclean &&\
    apt-get autoremove --yes &&\
    rm -rf /var/lib/{apt,dpkg,cache,log}/
HEALTHCHECK --interval=10s \
    CMD wget --no-verbose --tries=1 --spider "http://localhost:8090/api/health" || exit 1
ENTRYPOINT ["./wubbl0rz-archiv", "serve", "--http", "0.0.0.0:8090"]
