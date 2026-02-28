package app

import (
	"apps_minimalist/bidder/code/server"
)

type Config struct {
	Server 		server.Config
	LogLevel 	string `envconfig:"LOG_LEVEL" required:"true"`
}