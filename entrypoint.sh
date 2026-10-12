#!/usr/bin/env sh

/usr/sbin/unbound-anchor

/usr/local/bin/check &

unbound -d &
child=$!

trap "kill -s INT $child" SIGINT SIGTERM
wait "$child"
trap - SIGINT SIGTERM
wait "$child"
