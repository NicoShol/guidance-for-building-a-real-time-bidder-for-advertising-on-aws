package app

import (
	"apps_minimalist/bidder/code/bidhandler"
	bidserver "apps_minimalist/bidder/code/server"
	"os"
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

	bidHandler := bidhandler.New(cfg.BidHandlerCfg)
	server := bidserver.NewServer(cfg.Server, bidHandler)

	// Handle CTR+C 
	stop := make(chan os.Signal, 1)

	server.AsyncListenAndServe(func(err error) {
		errReturn = errors.Wrap(err, "error during server operation")
		stop <- os.Interrupt
	})

	<-stop
	return 
}