# syntax=docker/dockerfile:1


####
#### BUILD
####
FROM golang:1.16-alpine AS build
RUN apk update
RUN apk upgrade
RUN apk add --update gcc g++

WORKDIR /gnol-build

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
COPY data ./data/
COPY server ./server

RUN go get github.com/shurcooL/vfsgen
RUN go mod download
RUN go generate
RUN go build

###
### RUN
###
FROM alpine:latest


WORKDIR /gnol-app

ENV USER=gnoluser
ENV UID=12345
ENV GID=23456

RUN addgroup --gid $GID $USER
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --ingroup "$USER" \
    --no-create-home \
    --uid "$UID" \
    "$USER"
USER $USER:$USER



COPY --from=build /gnol-build/gnol ./
COPY container.cfg ./


CMD [ "./gnol", "-c", "container.cfg" ]
