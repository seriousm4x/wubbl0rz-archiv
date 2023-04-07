FROM golang:alpine3.17 as base
WORKDIR /build
COPY main.go go.mod go.sum ./
COPY pkg ./pkg
COPY docs ./docs
RUN go mod download && \
    apk update && \
    apk add --no-cache libc-dev gcc && \
    go build -ldflags="-s -w" -o apiserver .

FROM alpine:3.17
WORKDIR /app
COPY --from=base /build/apiserver ./
RUN apk update && \
    apk add --no-cache tzdata ffmpeg
HEALTHCHECK --interval=10s \
    CMD wget --no-verbose --tries=1 --spider "http://localhost:5000/health" || exit 1
ENTRYPOINT ["./apiserver"]
