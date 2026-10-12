package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kak-tus/healthcheck"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	var cnf struct {
		Listen string
	}

	err := envconfig.Process("CHECK", &cnf)
	if err != nil {
		log.Panic().Err(err).Msg("config read fail")
	}

	log.Debug().Msg("start update certs")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	cmd := exec.CommandContext(ctx, "/usr/sbin/unbound-anchor")

	if err := cmd.Run(); err != nil {
		log.Error().Err(err).Msg("update certs failed")
	}

	cancel()

	log.Debug().Msg("done update certs")

	ticker := time.NewTicker(time.Hour * 24 * 7)

	go func() {
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

			cmd := exec.CommandContext(ctx, "/usr/sbin/unbound-anchor")

			if err := cmd.Run(); err != nil {
				log.Error().Err(err).Msg("update certs failed")
			}

			cancel()
		}
	}()

	hlth := healthcheck.NewHandler()

	hlth.AddReq("/dig/", func(r *http.Request) (healthcheck.State, string) {
		addr := r.FormValue("addr")
		if len(addr) == 0 {
			log.Error().Msg("no address")

			return healthcheck.StateCritical, "fail"
		}

		ips, err := net.LookupIP(addr)
		if err != nil {
			log.Error().Err(err).Msg("lookup failed")

			return healthcheck.StateCritical, "fail"
		}

		if len(ips) == 0 {
			log.Error().Msg("no ips")

			return healthcheck.StateCritical, "fail"
		}

		var conv []string
		for _, ip := range ips {
			conv = append(conv, ip.String())
		}

		res := strings.Join(conv, ",")

		log.Debug().Str("res", res).Msg("resolved")

		return healthcheck.StatePassing, "Resolved: " + res
	})

	err = http.ListenAndServe(cnf.Listen, hlth)
	if err != nil && err != http.ErrServerClosed {
		log.Panic().Err(err).Msg("listen fail")
	}
}
