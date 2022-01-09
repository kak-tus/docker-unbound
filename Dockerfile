FROM golang:1.13.2-alpine3.10 AS go-build

WORKDIR /go/docker-unbound

COPY *.go ./
COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0

RUN go test && go build -o /go/bin/check

FROM hashicorp/consul-template:0.27.2 AS build

FROM alpine:3.10

RUN \
  apk add --no-cache \
    unbound \
  \
  && echo 'include: "/etc/unbound/unbound.conf.d/local.conf"' >> /etc/unbound/unbound.conf \
  \
  # Update DNSSEC keys
  && ( /usr/sbin/unbound-anchor ; echo 'ok' )

COPY --from=build /bin/consul-template /usr/local/bin/consul-template
COPY templates /root/templates
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
COPY --from=go-build /go/bin/check /usr/local/bin/check
COPY etc /etc/

ENV \
  CHECK_PORT=9000 \
  UNBOUND_FORWARD_ZONE= \
  UNBOUND_LOCAL_DATA= \
  UNBOUND_STUB_ZONE=

EXPOSE 53 9000

CMD ["/usr/local/bin/entrypoint.sh"]
