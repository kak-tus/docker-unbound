FROM golang:1.10.3-alpine3.8 AS go-build

WORKDIR /go/src/github.com/kak-tus/docker-unbound/check

COPY check/vendor ./vendor
COPY check/main.go .

RUN go test && go build -o /go/bin/check

FROM alpine:3.8 AS build

ENV \
  CONSUL_TEMPLATE_VERSION=0.19.4 \
  CONSUL_TEMPLATE_SHA256=5f70a7fb626ea8c332487c491924e0a2d594637de709e5b430ecffc83088abc0

RUN \
  apk add --no-cache \
    curl \
    unzip \
  \
  && cd /usr/local/bin \
  && curl -L https://releases.hashicorp.com/consul-template/${CONSUL_TEMPLATE_VERSION}/consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip -o consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip \
  && echo -n "$CONSUL_TEMPLATE_SHA256  consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip" | sha256sum -c - \
  && unzip consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip \
  && rm consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip

FROM alpine:3.8

RUN \
  apk add --no-cache \
    unbound \
  \
  && echo 'include: "/etc/unbound/unbound.conf.d/local.conf"' >> /etc/unbound/unbound.conf \
  \
  # Update DNSSEC keys
  && ( /usr/sbin/unbound-anchor ; echo 'ok' )

COPY --from=build /usr/local/bin/consul-template /usr/local/bin/consul-template
COPY templates /root/templates
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
COPY --from=go-build /go/bin/check /usr/local/bin/check
COPY check/etc /etc/

ENV \
  CHECK_PORT=9000 \
  UNBOUND_FORWARD_ZONE= \
  UNBOUND_LOCAL_DATA= \
  UNBOUND_STUB_ZONE=

EXPOSE 53 9000

CMD ["/usr/local/bin/entrypoint.sh"]
