package main

import (
	"apps_minimalist/bidder/app"
	"os"

	"github.com/rs/zerolog/log"
)

// Simple main function - logging & delegating to app.App() which will initialize and run the bidder.
func main() {
    log.Info().Msg("bidder starting...")
    if err := app.App(); err != nil {
        log.Error().Err(err).Msg("")
        log.Info().Msg("bidder exit")
        os.Exit(1)
    }
    log.Info().Msg("bidder exit")
    os.Exit(0)
}
