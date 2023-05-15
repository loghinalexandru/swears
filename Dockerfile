#syntax=docker/dockerfile:latest

FROM golang:alpine3.17 AS build

WORKDIR /app

COPY ../ .

RUN go mod download
RUN go build -o /swears ./

FROM golang:alpine3.17

RUN apk update && apk add ffmpeg

WORKDIR /

COPY --link --from=build /swears /swears
COPY --link --from=build /app/misc /misc

EXPOSE 3000

ENTRYPOINT ["/swears"]