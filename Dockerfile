FROM golang:1.14 as builder

WORKDIR /app
COPY . /app

RUN make build

# A distroless container image with some basics like SSL certificates
# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static

COPY --from=builder /app/create-scheduled-milestone /create-scheduled-milestone

ENTRYPOINT ["/create-scheduled-milestone"]