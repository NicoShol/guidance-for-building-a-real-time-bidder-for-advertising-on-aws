package app

import (
	"apps_minimalist/bidder/code/bidhandler"
	diagnosticServer "apps_minimalist/bidder/code/diagnostic_server"
	"apps_minimalist/bidder/code/server"
)

type Config struct {
	Server 			server.Config
	BidHandlerCfg 	bidhandler.Config
	DiagnosticServer 	diagnosticServer.Config
	LogLevel 		string `envconfig:"LOG_LEVEL" required:"true"`
}