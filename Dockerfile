FROM golang:1.17.6-alpine3.15 AS go-build

WORKDIR /go/docker-unbound

COPY *.go ./
COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0

RUN go build -o /go/bin/check

FROM alpine:3.15

RUN \
  apk add --no-cache \
    unbound \
  \
  && echo 'include: "/etc/unbound/unbound.conf.d/local.conf"' >> /etc/unbound/unbound.conf \
  \
  # Update DNSSEC keys
  && ( /usr/sbin/unbound-anchor ; echo 'ok' )

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
COPY --from=go-build /go/bin/check /usr/local/bin/check

ENV \
  CHECK_LISTEN=0.0.0.0:9000

EXPOSE 53 9000

CMD ["/usr/local/bin/entrypoint.sh"]
