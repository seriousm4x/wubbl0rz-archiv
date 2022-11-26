FROM golang:alpine3.16 as base
WORKDIR /build
COPY main.go go.mod go.sum ./
COPY pkg ./pkg
COPY docs ./docs
RUN go mod download && \
    apk update && \
    apk add --no-cache libc-dev gcc vips-dev  && \
    go build -ldflags="-s -w" -o apiserver .

FROM alpine:3.16
WORKDIR /app
COPY --from=base /build/apiserver ./
RUN apk update && \
    apk add --no-cache tzdata ffmpeg vips curl
HEALTHCHECK --interval=10s \
    CMD curl -fs "http://localhost:5000/health" || exit 1
ENTRYPOINT ["./apiserver"]
