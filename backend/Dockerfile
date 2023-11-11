FROM golang:latest as base
ARG DEBIAN_FRONTEND=noninteractive
WORKDIR /build
COPY . .
RUN go mod download && \
    go build -ldflags="-s -w" ./cmd/wubbl0rz-archiv

FROM linuxserver/ffmpeg:latest
ARG DEBIAN_FRONTEND=noninteractive
USER 1000:1000
WORKDIR /app
COPY --from=base /build/wubbl0rz-archiv ./
HEALTHCHECK --interval=10s \
    CMD wget --no-verbose --tries=1 --spider "http://localhost:8090/api/health" || exit 1
ENTRYPOINT ["./wubbl0rz-archiv", "serve", "--http", "0.0.0.0:8090"]
