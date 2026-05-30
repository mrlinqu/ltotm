package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrlinqu/ltotm/internal/storage"
	"github.com/mrlinqu/ltotm/internal/webserver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Level(zerolog.DebugLevel)

	ctx := context.Background()

	cfg := parseArgs()

	stor := storage.NewFileStorage(cfg.storageDir)

	srv := webserver.New(
		cfg.iface+":"+cfg.port,
		cfg.tlsCert,
		cfg.tlsKey,
		stor,
	)

	webDoneCh := srv.Run(ctx)

	log.Debug().Msg("started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-webDoneCh:
	case <-sigChan:
		log.Debug().Msg("received shutdown signal")
	}

	<-srv.Shutdown(ctx)

	log.Debug().Msg("stoped")
}
