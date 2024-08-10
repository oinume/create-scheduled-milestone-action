# syntax=docker/dockerfile:1

FROM golang:1.22 AS build
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -v -o /bin/app

FROM debian:bookworm-slim AS final
RUN set -eux \
  apt-get install -qyy --no-install-recommends --no-install-suggests ca-certificates openssl
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
RUN set -eux \
  update-ca-certificates --fresh
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser
COPY --from=build /bin/app /bin/
ENTRYPOINT ["/bin/app"]
