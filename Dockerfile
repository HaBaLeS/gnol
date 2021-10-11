# syntax=docker/dockerfile:1

FROM golang:1.16-alpine
RUN apk update
RUN apk upgrade
RUN apk add --update gcc g++

WORKDIR /gnol

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
COPY data ./data/
COPY server ./server
COPY container.cfg ./



RUN go generate
RUN go build

CMD [ "./gnol", "-c", "container.cfg" ]
