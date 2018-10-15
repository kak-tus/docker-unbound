#!/usr/bin/env sh

/usr/local/bin/consul-template -config /root/templates/service.hcl &
child1=$!

/usr/local/bin/check &
child2=$!

trap "kill -s SIGINT $child1 ; kill $child2" SIGINT SIGTERM

while true; do
  kill -0 "$child1" >/dev/null 2>&1
  if [ "$?" = "0" ]; then
    sleep 5
  else
    echo "Exited consul-template"
    exit
  fi

  kill -0 "$child2" >/dev/null 2>&1
  if [ "$?" = "0" ]; then
    sleep 5
  else
    echo "Exited check"
    exit
  fi
done
