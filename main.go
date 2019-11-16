package main

import (
	"net"
	"net/http"
	"strings"

	"github.com/iph0/conf"
	"github.com/iph0/conf/envconf"
	"github.com/iph0/conf/fileconf"
	"github.com/kak-tus/healthcheck"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	log := logger.Sugar()

	fileLdr := fileconf.NewLoader("etc", "/etc")
	envLdr := envconf.NewLoader()

	configProc := conf.NewProcessor(
		conf.ProcessorConfig{
			Loaders: map[string]conf.Loader{
				"file": fileLdr,
				"env":  envLdr,
			},
		},
	)

	configRaw, err := configProc.Load(
		"file:check.yml",
		"env:^CHECK_",
	)

	if err != nil {
		log.Panic(err)
	}

	var cnf struct {
		Listen string
	}

	if err := conf.Decode(configRaw["healthcheck"], &cnf); err != nil {
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
