#!/usr/bin/env sh

/usr/local/bin/consul-template -config /root/templates/service.hcl &
child=$!

trap "kill -s SIGINT $child" SIGINT SIGTERM
wait "$child"
trap - SIGINT SIGTERM
wait "$child"
