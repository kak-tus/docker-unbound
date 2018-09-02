max_stale = "2m"

template {
  source = "/root/templates/local.conf.template"
  destination = "/etc/unbound/unbound.conf.d/local.conf"
}

exec {
  command = "unbound -d"
  splay = "60s"
}
