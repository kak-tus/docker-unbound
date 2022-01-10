package main

import (
	"net"
	"net/http"
	"strings"

	"github.com/kak-tus/healthcheck"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	log := logger.Sugar()

	var cnf struct {
		Listen string
	}

	err = envconfig.Process("CHECK", &cnf)
	if err != nil {
		log.Panic(err)
	}

	hlth := healthcheck.NewHandler()

	hlth.AddReq("/dig/", func(r *http.Request) (healthcheck.State, string) {
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

	err = http.ListenAndServe(cnf.Listen, hlth)
	if err != nil && err != http.ErrServerClosed {
		log.Panic(err)
	}

	_ = log.Sync()
}
