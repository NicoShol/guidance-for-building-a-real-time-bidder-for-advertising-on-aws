package app

import (
	"apps_minimalist/bidder/code/auction"
	"apps_minimalist/bidder/code/bidhandler"
	diagnosticServer "apps_minimalist/bidder/code/diagnostic_server"
	bidserver "apps_minimalist/bidder/code/server"
	"apps_minimalist/bidder/code/stream"
	"os"
	"os/signal"
	"runtime"
	"time"

	"emperror.dev/errors"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func App() (errReturn error) {
	cfg := Config{}

	// Checking valid config 
	if err := envconfig.Process("", &cfg); err != nil {
		return errors.Wrap(err, "error during config initialization")
	}

	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return err 
	}
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = time.Stamp

	log.Debug().Msgf(
		"bidder will use %d of %d available threads for goroutine calls",
		runtime.GOMAXPROCS(0),
		runtime.NumCPU(),
	)

	auctionFn := auction.New()
	dataStream := stream.New(cfg.Stream)
	bidHandler := bidhandler.New(cfg.BidHandlerCfg, auctionFn, dataStream)
	server := bidserver.NewServer(cfg.Server, bidHandler)
	diagServer := diagnosticServer.New(cfg.DiagnosticServer)

	// Handle CTR+C 
	stop := make(chan os.Signal, 1)

	server.AsyncListenAndServe(func(err error) {
		errReturn = errors.Wrap(err, "error during server operation")
		stop <- os.Interrupt
	})

	diagServer.AsyncListenAndServe(func(err error) {
		errReturn = errors.Wrap(err, "error during diagnostic server operation")
		stop <- os.Interrupt
	})

	log.Info().Msg("bidder ready")
	signal.Notify(stop, os.Interrupt)
	<-stop

	err = shutdown(cfg, server, diagServer)
	if err != nil {
		if errReturn == nil {
			return err
		}
	}

	return
}