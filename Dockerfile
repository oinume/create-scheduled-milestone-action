FROM golang:1.22

WORKDIR /work
COPY . /work
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o app .
ENTRYPOINT ["/work/app"]
