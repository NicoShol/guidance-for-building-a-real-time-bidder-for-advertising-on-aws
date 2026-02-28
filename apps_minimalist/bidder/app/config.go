package app

import (
	"apps_minimalist/bidder/code/bidhandler"
	"apps_minimalist/bidder/code/server"
)

type Config struct {
	Server 			server.Config
	BidHandlerCfg 	bidhandler.Config
	LogLevel 		string `envconfig:"LOG_LEVEL" required:"true"`
}