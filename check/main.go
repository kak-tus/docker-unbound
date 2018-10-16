package main

import (
	"net"
	"net/http"
	"strings"

	"git.aqq.me/go/app/appconf"
	"git.aqq.me/go/app/launcher"
	"github.com/iph0/conf/envconf"
	"github.com/iph0/conf/fileconf"
	"github.com/kak-tus/healthcheck"
)

func init() {
	fileLdr := fileconf.NewLoader("etc", "/etc")
	envLdr := envconf.NewLoader()
	appconf.RegisterLoader("file", fileLdr)
	appconf.RegisterLoader("env", envLdr)
	appconf.Require("file:check.yml")
	appconf.Require("env:^CHECK_")
}

func main() {
	launcher.Run(func() error {
		healthcheck.AddReq("/dig/", func(r *http.Request) (healthcheck.State, string) {
			addr := r.FormValue("addr")
			if len(addr) == 0 {
				return healthcheck.StateCritical, "fail"
			}

			ips, err := net.LookupIP(addr)
			if err != nil {
				return healthcheck.StateCritical, "fail"
			}

			if len(ips) == 0 {
				return healthcheck.StateCritical, "fail"
			}

			var conv []string
			for _, ip := range ips {
				conv = append(conv, ip.String())
			}

			res := strings.Join(conv, ",")

			return healthcheck.StatePassing, "Resolved: " + res
		})

		return nil
	})
}
