# syntax=docker/dockerfile:1


####
#### BUILD
####
FROM golang:1.25.1-alpine AS build
RUN apk update
RUN apk upgrade
RUN apk add --update gcc g++ make mupdf-dev build-base libc6-compat alpine-sdk

WORKDIR /gnol-build

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
COPY docs ./docs/
COPY data ./data/
COPY server ./server
COPY cmd ./cmd
COPY Makefile ./

RUN make build-server

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

COPY --from=build /gnol-build/bin/* ./
COPY container.cfg ./


CMD [ "./gnol", "-c", "container.cfg" ]
