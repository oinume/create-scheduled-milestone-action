# syntax=docker/dockerfile:1

FROM golang:1.22-alpine

WORKDIR /work
COPY . /work
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o app .
ENTRYPOINT ["/work/app"]
