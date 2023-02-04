#syntax=docker/dockerfile:latest

FROM golang:1.19-alpine AS build

WORKDIR /app

COPY ../ .

RUN apk add pkgconfig alsa-tools-dev

RUN go mod download
RUN go build -o /swears ./

FROM golang:1.19-alpine

RUN apk update && apk add ffmpeg

WORKDIR /

COPY --link --from=build /swears /swears
COPY --link --from=build /app/misc /misc

EXPOSE 3000

ENTRYPOINT ["/resonator"]