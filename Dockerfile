FROM golang:alpine as base
WORKDIR /build
COPY main.go go.mod go.sum ./
COPY pkg ./pkg
COPY docs ./docs
RUN go mod download && \
    apk update && \
    apk add --no-cache libc-dev gcc vips-dev  && \
    go build -ldflags="-s -w" -o apiserver .

FROM alpine:latest
WORKDIR /app
COPY --from=base /build/apiserver ./
RUN apk update && \
    apk add --no-cache ffmpeg vips curl
HEALTHCHECK --interval=10s \
    CMD curl -fs "http://localhost:5000/health" || exit 1
ENTRYPOINT ["./apiserver"]
